package repository

import (
	"context"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/domain"
)

type Repository interface {
	Create(ctx context.Context, in domain.CreateCustomerInput) (*domain.Customer, error)
	GetByID(ctx context.Context, id int32) (*domain.Customer, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Customer, error)
	DeleteByID(ctx context.Context, id int32) (*domain.Customer, error)
}
