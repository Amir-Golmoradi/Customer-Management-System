package handler

import (
	"encoding/json"
	"net/http"
)

type getCustomerRequest struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	// 1. Ensure method in GET
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Map domain to response
	// var request getCustomerRequest
	// customer := &database.Customer{
	// 	ID:    request.ID,
	// 	Name:  request.Name,
	// 	Email: request.Email,
	// }
	customers, err := h.service.GetCustomers(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch customers: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customers)
}
