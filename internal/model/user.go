package model

type User struct {
	ID        uint64 `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name" validate:"required,alpha,min=3,max=50"`
	LastName  string `json:"last_name" db:"last_name" validate:"required,alpha,min=3,max=50"`
}
