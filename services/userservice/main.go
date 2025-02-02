package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zain2323/cronium/services/userservice/config"
	"github.com/zain2323/cronium/services/userservice/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")
	if port == "" {
		port = "8081"
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

	// cors setup
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(middleware.Logger)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok\n"))
	})
	router.Post("/users", handler.CreateUser)

	router.Mount("/v1/api", router)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}
	fmt.Printf("Starting user service on port %s .....\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
