package domain

import "time"

type Customer struct {
	ID        int32
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateCustomerInput struct {
	Name     string
	Email    string
	Password string
}
