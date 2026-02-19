package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/domain"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/shared/errorsx"
)

type mockRepo struct {
	created    *domain.Customer
	createErr  error
	getResult  *domain.Customer
	getErr     error
	listResult []domain.Customer
	listErr    error
}

func (m *mockRepo) Create(_ context.Context, _ domain.CreateCustomerInput) (*domain.Customer, error) {
	return m.created, m.createErr
}

func (m *mockRepo) GetByID(_ context.Context, _ int32) (*domain.Customer, error) {
	return m.getResult, m.getErr
}

func (m *mockRepo) List(_ context.Context, _, _ int32) ([]domain.Customer, error) {
	return m.listResult, m.listErr
}

func TestCreateCustomerValidation(t *testing.T) {
	svc := New(&mockRepo{})
	_, err := svc.CreateCustomer(context.Background(), domain.CreateCustomerInput{Name: "", Email: "bad", Password: "123"})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !errorsx.IsKind(err, errorsx.KindValidation) {
		t.Fatalf("expected validation kind, got %v", err)
	}
}

func TestCreateCustomerPassesRepositoryError(t *testing.T) {
	expected := errors.New("db error")
	svc := New(&mockRepo{createErr: expected})

	_, err := svc.CreateCustomer(context.Background(), domain.CreateCustomerInput{Name: "Test", Email: "test@example.com", Password: "password1"})
	if !errors.Is(err, expected) {
		t.Fatalf("expected repository error to pass through, got %v", err)
	}
}
