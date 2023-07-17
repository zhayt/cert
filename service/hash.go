package service

import (
	"errors"
	"github.com/zhayt/cert-tz/model"
	"github.com/zhayt/cert-tz/storage"
	"go.uber.org/zap"
	"hash/crc64"
	"math/bits"
	"strconv"
	"sync"
	"time"
)

var ErrWorkersPool = errors.New("the workers' pool is full")

const N = 3

//go:generate mockery --name IHashService
type IHashService interface {
	CalculateHash(certHash model.CertHash) (uint64, error)
	GetCalculatedHash(hashID uint64) (model.CertHash, error)
}

type HashService struct {
	workersLimit uint64
	workersCount uint64
	m            sync.Mutex
	storage      *storage.Storage
	validate     *ValidateService
	crcTable     *crc64.Table
	l            *zap.Logger
}

func (s *HashService) GetCalculatedHash(hashID uint64) (model.CertHash, error) {
	return s.storage.Hash.GetHash(hashID)
}

func (s *HashService) CalculateHash(certHash model.CertHash) (uint64, error) {
	if err := s.validate.validateStruct(certHash); err != nil {
		return 0, ErrInvalidData
	}

	err := s.addWorker()
	if err != nil {
		return 0, err
	}

	certHash.Hash = "PENDING"

	certHash.ID, err = s.storage.Hash.CreateHash(certHash)
	if err != nil {
		s.deleteWorker()
		return 0, err
	}

	go func(model.CertHash) {
		defer s.deleteWorker()
		certHash.Hash = strconv.Itoa(s.certHash(certHash.InputStr))

		if err := s.storage.Hash.UpdateHash(certHash); err != nil {
			s.l.Error("UpdateHash error", zap.Error(err))
			return
		}

		s.l.Error("Hash calculated successfully", zap.Uint64("id", certHash.ID), zap.String("hash", certHash.Hash))
	}(certHash)

	return certHash.ID, nil
}

func (s *HashService) certHash(str string) int {
	// take crc64 hash
	crcHash := int64(crc64.Checksum([]byte(str), s.crcTable))

	// crcHash & current timestamp every 5 second
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan struct{}, 1)

	go func(int64) {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				crcHash = crcHash & time.Now().UnixNano()
			}
		}
	}(crcHash)

	// waiting 1 min
	time.Sleep(1*time.Minute + 5*time.Second)
	ticker.Stop()
	done <- struct{}{}

	// count set of bits and return
	return bits.OnesCount64(uint64(crcHash))
}

func (s *HashService) deleteWorker() {
	s.m.Lock()
	defer s.m.Unlock()
	s.workersCount--
}

func (s *HashService) addWorker() error {
	s.m.Lock()
	defer s.m.Unlock()

	if s.workersCount+1 > s.workersLimit {
		return ErrWorkersPool
	}

	s.workersCount++

	return nil
}

func NewHashService(storage *storage.Storage, validate *ValidateService, l *zap.Logger) *HashService {
	return &HashService{
		workersLimit: N,
		crcTable:     crc64.MakeTable(crc64.ISO),
		storage:      storage,
		validate:     validate,
		l:            l,
	}
}
