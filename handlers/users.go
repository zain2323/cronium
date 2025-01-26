package handlers

import (
	"github.com/google/uuid"
	"github.com/zain2323/cronium/config"
	"github.com/zain2323/cronium/data"
	"github.com/zain2323/cronium/internal/database"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	Config *config.ApiConfig
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &database.User{}
	err := data.FromJSON(user, r.Body)
	if err != nil {
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
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)

		return
	}
}
