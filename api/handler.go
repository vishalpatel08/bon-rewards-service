package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vishalpatel08/bon-rewards-service/models"
)

type Service interface {
	PayBill(ctx context.Context, billID int64) (*models.Bill, string, error)
	CreateUser(ctx context.Context, name string) (*models.User, error)
	CreateBill(ctx context.Context, userID int64, amount int64, dueDate time.Time) (*models.Bill, error)
}

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) PayBill(w http.ResponseWriter, r *http.Request) {
	billIDStr := chi.URLParam(r, "billID")
	billID, err := strconv.ParseInt(billIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid bill ID"))
		return
	}

	updatedBill, rewardMsg, err := h.service.PayBill(r.Context(), billID)
	if err != nil {
		log.Printf("ERROR paying bill %d: %v", billID, err)
		writeError(w, http.StatusInternalServerError, errors.New("could not process payment"))
		return
	}

	response := map[string]interface{}{
		"status":         "success",
		"bill":           updatedBill,
		"reward_message": rewardMsg,
	}

	writeJSON(w, http.StatusOK, response)
}

type createUserRequest struct {
	Name string `json:"name"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	user, err := h.service.CreateUser(r.Context(), req.Name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

type createBillRequest struct {
	UserID  int64     `json:"user_id"`
	Amount  int64     `json:"amount"`
	DueDate time.Time `json:"due_date"`
}

func (h *Handler) CreateBill(w http.ResponseWriter, r *http.Request) {
	var req createBillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	bill, err := h.service.CreateBill(r.Context(), req.UserID, req.Amount, req.DueDate)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, bill)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ERROR encoding JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
