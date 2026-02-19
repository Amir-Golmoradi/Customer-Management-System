package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/domain"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/service"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/platform/httpserver"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/shared/errorsx"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/shared/idempotency"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/shared/pagination"
)

type Handler struct {
	service          *service.Service
	idempotencyStore *idempotency.Store
}

func NewHandler(service *service.Service, store *idempotency.Store) *Handler {
	return &Handler{service: service, idempotencyStore: store}
}

type createCustomerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type customerResponse struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req createCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, errorsx.New(errorsx.KindValidation, "INVALID_REQUEST", "invalid request body", err))
		return
	}

	if key := strings.TrimSpace(r.Header.Get("Idempotency-Key")); key != "" {
		if cached, ok := h.idempotencyStore.Get(key); ok {
			if response, ok := cached.(customerResponse); ok {
				writeJSON(w, http.StatusCreated, response)
				return
			}
		}
	}

	customer, err := h.service.CreateCustomer(r.Context(), domain.CreateCustomerInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	response := toResponse(customer)
	if key := strings.TrimSpace(r.Header.Get("Idempotency-Key")); key != "" {
		h.idempotencyStore.Set(key, response)
	}
	writeJSON(w, http.StatusCreated, response)
}

func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	params := pagination.Parse(r.URL.Query().Get("limit"), r.URL.Query().Get("offset"))
	customers, err := h.service.ListCustomers(r.Context(), params.Limit, params.Offset)
	if err != nil {
		writeError(w, r, err)
		return
	}

	response := make([]customerResponse, 0, len(customers))
	for _, customer := range customers {
		c := customer
		response = append(response, toResponse(&c))
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, r, errorsx.New(errorsx.KindValidation, "INVALID_ID", "id must be numeric", err))
		return
	}

	customer, err := h.service.GetCustomerByID(r.Context(), int32(id))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(customer))
}

func toResponse(customer *domain.Customer) customerResponse {
	return customerResponse{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: customer.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	httpErr := errorsx.ToError(err)
	response := errorsx.ErrorResponse{
		Code:      httpErr.Code,
		Message:   httpErr.Message,
		RequestID: httpserver.RequestIDFromContext(r.Context()),
	}
	writeJSON(w, errorsx.HTTPStatus(err), response)
}
