package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/domain"
	generated "github.com/Amir-Golmoradi/Customer-Management-System/internal/database/generated"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/shared/errorsx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryObserver interface {
	ObserveDBQuery(name string, duration time.Duration, failed bool)
}

type PostgresRepository struct {
	queries  *generated.Queries
	observer QueryObserver
}

func NewPostgresRepository(queries *generated.Queries, observer QueryObserver) *PostgresRepository {
	return &PostgresRepository{queries: queries, observer: observer}
}

func (r *PostgresRepository) Create(ctx context.Context, in domain.CreateCustomerInput) (*domain.Customer, error) {
	start := time.Now()
	created, err := r.queries.CreateCustomer(ctx, generated.CreateCustomerParams{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	})
	r.observe("create_customer", start, err)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errorsx.New(errorsx.KindConflict, "CUSTOMER_ALREADY_EXISTS", "customer already exists", err)
		}
		return nil, fmt.Errorf("create customer: %w", err)
	}

	customer := toDomain(created)
	return &customer, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int32) (*domain.Customer, error) {
	start := time.Now()
	customer, err := r.queries.GetCustomerByID(ctx, id)
	r.observe("get_customer_by_id", start, err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsx.New(errorsx.KindNotFound, "CUSTOMER_NOT_FOUND", "customer not found", err)
		}
		return nil, fmt.Errorf("get customer by id: %w", err)
	}

	result := toDomain(customer)
	return &result, nil
}

func (r *PostgresRepository) List(ctx context.Context, limit, offset int32) ([]domain.Customer, error) {
	start := time.Now()
	customers, err := r.queries.ListCustomersPaginated(ctx, generated.ListCustomersPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	r.observe("list_customers_paginated", start, err)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}

	result := make([]domain.Customer, 0, len(customers))
	for _, customer := range customers {
		result = append(result, toDomain(customer))
	}
	return result, nil
}

func (r *PostgresRepository) DeleteByID(ctx context.Context, id int32) (*domain.Customer, error){
	
}

func (r *PostgresRepository) observe(name string, start time.Time, err error) {
	if r.observer != nil {
		r.observer.ObserveDBQuery(name, time.Since(start), err != nil)
	}
}

func toDomain(customer generated.Customer) domain.Customer {
	return domain.Customer{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt.Time,
		UpdatedAt: customer.UpdatedAt.Time,
	}
}
