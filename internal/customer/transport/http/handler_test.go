package httptransport

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/domain"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/repository"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/service"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/shared/idempotency"
)

type fakeRepo struct{}

func (f *fakeRepo) Create(_ context.Context, _ domain.CreateCustomerInput) (*domain.Customer, error) {
	return &domain.Customer{ID: 1, Name: "User", Email: "user@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (f *fakeRepo) GetByID(_ context.Context, id int32) (*domain.Customer, error) {
	return &domain.Customer{ID: id, Name: "User", Email: "user@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (f *fakeRepo) List(_ context.Context, _, _ int32) ([]domain.Customer, error) {
	return []domain.Customer{}, nil
}

var _ repository.Repository = (*fakeRepo)(nil)

func TestCreateCustomerMethodValidation(t *testing.T) {
	svc := service.New(&fakeRepo{})
	h := NewHandler(svc, idempotency.NewStore())

	req := httptest.NewRequest(http.MethodPost, "/v1/customers", bytes.NewBufferString("bad-json"))
	rr := httptest.NewRecorder()
	h.CreateCustomer(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}
