package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zain2323/cronium/internal/database"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	if dbUrl == "" {
		log.Fatal("Please configure DB_URL in the environment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}
	fmt.Println(apiCfg)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World\n"))
	})
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
