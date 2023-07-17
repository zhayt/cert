package model

import "time"

type CertHash struct {
	ID           uint64    `json:"id" db:"id"`
	InputStr     string    `json:"input_str" db:"input_str" validate:"required,min=1,max=255"`
	Hash         string    `json:"hash" db:"hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	CalculatedAt time.Time `json:"calculated_at" db:"calculated_at"`
}
