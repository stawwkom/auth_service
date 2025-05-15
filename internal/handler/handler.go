package handler

import (
	"encoding/json"
	"net/http"

	"github.com/stawwkom/auth_service/internal/service"
)

// Handler представляет HTTP-обработчики
type Handler struct {
	services *service.Service
}

// NewHandler создает новый экземпляр Handler
func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

// Register обрабатывает регистрацию нового пользователя
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user service.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := (*h.services).Register(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login обрабатывает вход пользователя
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := (*h.services).Login(r.Context(), credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}
