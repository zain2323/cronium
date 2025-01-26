package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zain2323/cronium/config"
	"github.com/zain2323/cronium/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		return
	}
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	if dbUrl == "" {
		log.Fatal("Please configure DB_URL in the environment")
	}

	apiCfg, err := config.New(dbUrl)
	if err != nil {
		log.Fatal("Failed to initialize API config:", err)
	}

	handler := &handlers.UserHandler{Config: apiCfg}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok\n"))
	})
	router.Post("/users", handler.CreateUser)

	router.Mount("/v1/api", router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Printf("Listening on port %s .....\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
