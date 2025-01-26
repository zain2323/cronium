package handlers

import (
	"net/http"
)

func cr(w http.ResponseWriter, r *http.Request) {
	createUser()
}
