package service

import (
	"github.com/stretchr/testify/mock"
	"github.com/zhayt/cert-tz/model"
	"github.com/zhayt/cert-tz/storage"
	"github.com/zhayt/cert-tz/storage/mocks"
	"go.uber.org/zap"
	"testing"
)

func TestHashService_CalculateHash(t *testing.T) {
	type args struct {
		certHash model.CertHash
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"success", args{certHash: model.CertHash{InputStr: "Hello"}}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashStorage := mocks.NewIHashStorage(t)

			if tt.args.certHash.InputStr != "" {
				hashStorage.On("CreateHash", mock.Anything, mock.AnythingOfType("string")).
					Return(func(hash model.CertHash) (uint64, error) {
						return 1, nil
					})
			} else {
				hashStorage.On("CreateHash", mock.Anything, mock.AnythingOfType("string")).
					Return(func(hash model.CertHash) (uint64, error) {
						return 0, ErrInvalidData
					})
			}

			stor := &storage.Storage{Hash: hashStorage}

			s := NewHashService(stor, NewValidateService(), zap.NewNop())

			got, err := s.CalculateHash(tt.args.certHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}
