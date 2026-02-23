package service

import (
	"context"
	"strings"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/domain"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/repository"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/shared/errorsx"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/shared/validation"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCustomer(ctx context.Context, in domain.CreateCustomerInput) (*domain.Customer, error) {
	if !validation.NonEmpty(in.Name) {
		return nil, errorsx.New(errorsx.KindValidation, "INVALID_NAME", "name is required", nil)
	}
	if len(strings.TrimSpace(in.Name)) > 120 {
		return nil, errorsx.New(errorsx.KindValidation, "INVALID_NAME", "name exceeds max length", nil)
	}
	if !validation.ValidEmail(in.Email) {
		return nil, errorsx.New(errorsx.KindValidation, "INVALID_EMAIL", "email is invalid", nil)
	}
	if len(strings.TrimSpace(in.Password)) < 8 {
		return nil, errorsx.New(errorsx.KindValidation, "INVALID_PASSWORD", "password must be at least 8 characters", nil)
	}

	customer, err := s.repo.Create(ctx, in)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *Service) GetCustomerByID(ctx context.Context, id int32) (*domain.Customer, error) {
	if id <= 0 {
		return nil, errorsx.New(errorsx.KindValidation, "INVALID_ID", "id must be positive", nil)
	}
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *Service) DeleteCustomerByID(ctx context.Context, id int32) (*domain.Customer, error) {
	if id <= 0 {
		return nil, errorsx.New(errorsx.KindValidation, "INVALID_ID", "id must be positive", nil)
	}
	customer, err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customer, nil

}

func (s *Service) ListCustomers(ctx context.Context, limit, offset int32) ([]domain.Customer, error) {
	return s.repo.List(ctx, limit, offset)
}
