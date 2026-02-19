package customer

import (
	"context"
	"fmt"

	database "github.com/Amir-Golmoradi/Customer-Management-System/internal/database/generated"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetCustomers(ctx context.Context) ([]database.Customer, error) {
	c, err := s.repository.FindAllCustomers(ctx)
	if err != nil {
		return nil, fmt.Errorf("customer not found %w", err)
	}
	return c, nil
}

func (s *Service) GetCustomerByID(ctx context.Context, id int32) (*database.Customer, error) {
	c, err := s.repository.FindCustomerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("no customer with this id has found %w", err)
	}
	return c, nil
}

func (s *Service) GetCustomerByEmail(ctx context.Context, email string) (*database.Customer, error) {
	c, err := s.repository.FindCustomerByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("no customer with this email has found %w", err)
	}
	return c, nil
}

func (s *Service) CreateCustomer(ctx context.Context, name, email, password string) (*database.Customer, error) {
	c, err := s.repository.CreateNewCustomer(ctx, name, email, password)
	if err != nil {
		return nil, fmt.Errorf("no customer created %w", err)
	}
	return c, nil
}

func (s *Service) UpdateCustomer(ctx context.Context, id int32, name, email, password string) (*database.Customer, error) {
	c, err := s.repository.UpdateExistingCustomer(ctx, id, name, email, password)
	if err != nil {
		return nil, fmt.Errorf("no information has changed %w", err)
	}
	return c, nil
}

func (s *Service) DeleteCustomerByEmail(ctx context.Context, email string) error {
	return s.repository.DeleteCustomerByEmail(ctx, email)
}
