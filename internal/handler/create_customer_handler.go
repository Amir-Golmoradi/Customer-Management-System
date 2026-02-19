package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer"
	database "github.com/Amir-Golmoradi/Customer-Management-System/internal/database/generated"
)

type Handler struct {
	service *customer.Service
}

func NewHandler(service *customer.Service) *Handler {
	return &Handler{service: service}
}

type createCustomerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST
func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	// 1. Ensure method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// 2. Decode the JSON request
	var request createCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Map request to domain entity
	customer := &database.Customer{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
	createdCustomer, err := h.service.CreateCustomer(r.Context(), customer.Name, customer.Email, customer.Password)
	if err != nil {
		http.Error(w, "could not create customer", http.StatusInternalServerError)
		return
	}
	resp := struct {
		ID    int32  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		ID:    createdCustomer.ID,
		Name:  createdCustomer.Name,
		Email: createdCustomer.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
