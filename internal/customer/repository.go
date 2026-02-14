package customer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	database "github.com/Amir-Golmoradi/Customer-Management-System/internal/database/generated"
)

var ErrCustomerNotFound = errors.New("customer not found")

// Repository is the concrete repository for customer-related database operations
type Repository struct {
	queries *database.Queries
}

// NewCustomerRepository is the constructor for CustomerRepository
func NewCustomerRepository(q *database.Queries) *Repository {
	return &Repository{queries: q}
}

// FindAllCustomers returns all customers
func (r *Repository) FindAllCustomers(ctx context.Context) ([]database.Customer, error) {
	return r.queries.ListCustomers(ctx)
}

// FindCustomerByID returns a customer by ID
func (r *Repository) FindCustomerByID(ctx context.Context, id int32) (*database.Customer, error) {
	customer, err := r.queries.GetCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("get customer by id: %w", err)
	}
	return &customer, nil
}

// FindCustomerByEmail returns a customer by email
func (r *Repository) FindCustomerByEmail(ctx context.Context, email string) (*database.Customer, error) {
	customer, err := r.queries.GetCustomerByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("get customer by email: %w", err)
	}
	return &customer, nil
}

// CreateNewCustomer creates a new customer
func (r *Repository) CreateNewCustomer(ctx context.Context, name, email, password string) (*database.Customer, error) {
	params := database.CreateCustomerParams{
		Name:     name,
		Email:    email,
		Password: password,
	}
	customer, err := r.queries.CreateCustomer(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}
	return &customer, nil
}

// UpdateExistingCustomer updates an existing customer
func (r *Repository) UpdateExistingCustomer(ctx context.Context, id int32, name, email, password string) (*database.Customer, error) {
	params := database.UpdateCustomerParams{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
	updatedCustomer, err := r.queries.UpdateCustomer(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("update customer: %w", err)
	}
	return &updatedCustomer, nil
}

// DeleteCustomerByEmail deletes a customer by email
func (r *Repository) DeleteCustomerByEmail(ctx context.Context, email string) error {
	rows, err := r.queries.DeleteCustomerByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("delete customer: %w", err)
	}
	if rows == 0 {
		return ErrCustomerNotFound
	}
	return nil
}
