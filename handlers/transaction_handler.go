package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"kasir-api/utils"
	"net/http"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(services *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: services}
}

func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.ResponseSuccess(w, transaction, http.StatusCreated, "Transaction created successfully")
}
