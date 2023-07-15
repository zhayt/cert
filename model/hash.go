package model

import "time"

type CertHash struct {
	ID           uint64    `json:"id" db:"id"`
	InputStr     string    `json:"data" db:"input_str"`
	Hash         string    `json:"hash" db:"hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	CalculatedAt time.Time `json:"calculated_at" db:"calculated_at"`
}
