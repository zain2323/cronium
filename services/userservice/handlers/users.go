package handlers

import (
	"github.com/google/uuid"
	"github.com/zain2323/cronium/services/userservice/config"
	"github.com/zain2323/cronium/services/userservice/data"
	"github.com/zain2323/cronium/services/userservice/internal/database"
	"github.com/zain2323/cronium/services/userservice/utils"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	Config *config.ApiConfig
	Logger *log.Logger
}

// NewUser returns new user handler with the provided logger and Config
func NewUser(config *config.ApiConfig, logger *log.Logger) *UserHandler {
	return &UserHandler{config, logger}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := &data.User{}
	err := data.FromJSON(user, r.Body)
	if err != nil {
		h.Logger.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check user credentials and generate JWT token
	currentUser, err := h.Config.DB.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		h.Logger.Println("User not found: ", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !utils.VerifyPassword(currentUser.Password, user.Password) {
		h.Logger.Println("Invalid email or password")
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	token, err := utils.GenerateToken(currentUser.ID)
	if err != nil {
		h.Logger.Println("JWT generation error: ", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"token": token,
	}
	err = data.ToJSON(&response, w)
	if err != nil {
		h.Logger.Println("JSON marshal error: ", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
	h.Logger.Println("Login success for user: ", user.Email)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &database.User{}
	err := data.FromJSON(user, r.Body)
	if err != nil {
		h.Logger.Println("Failed to parse json: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		h.Logger.Println("Password hash error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.Config.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  hashedPassword,
	})

	if err != nil {
		h.Logger.Println("Error creating user", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = data.ToJSON(user, w)
	if err != nil {
		h.Logger.Println("JSON marshal error: ", err)
		return
	}
	h.Logger.Println("Created user: ", user.Email)
}
