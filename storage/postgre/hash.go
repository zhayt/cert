package postgre

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/cert-tz/model"
)

type HashStorage struct {
	db *sqlx.DB
}

func (r *HashStorage) CreateHash(hash model.CertHash) (uint64, error) {
	qr := `INSERT INTO cert_hash (input_str, hash) VALUES ($1, $2) RETURNING id`

	var hashID uint64
	if err := r.db.Get(&hashID, qr, hash.InputStr, hash.Hash); err != nil {
		return 0, fmt.Errorf("cannot create hash: %w", err)
	}

	return hashID, nil
}

func (r *HashStorage) GetHash(hashID uint64) (model.CertHash, error) {
	qr := `SELECT id, input_str, hash, created_at, coalesce(calculated_at, '0001-01-01') AS calculated_at FROM cert_hash WHERE id = $1`

	var certHash model.CertHash

	if err := r.db.Get(&certHash, qr, hashID); err != nil {
		return model.CertHash{}, fmt.Errorf("cannot get hash: %w", err)
	}

	return certHash, nil
}

func (r *HashStorage) UpdateHash(hash model.CertHash) error {
	qr := `UPDATE cert_hash SET hash = $1, calculated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING id`

	var hashID uint64

	fmt.Println(hash)
	if err := r.db.Get(&hashID, qr, hash.Hash, hash.ID); err != nil {
		return fmt.Errorf("cannot update hash: %w", err)
	}

	return nil
}

func NewHashStorage(db *sqlx.DB) *HashStorage {
	return &HashStorage{db: db}
}
