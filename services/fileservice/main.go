package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/zain2323/cronium/services/fileservice/files"
	"github.com/zain2323/cronium/services/fileservice/handlers"
	"log"
	"net/http"
	"os"
)

var basePath = "./imagestore"

func main() {
	logger := log.New(os.Stdout, "files-api", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	//storage, storageErr := files.NewLocal(basePath, 1024*1000*5) // max file size is kept 5MB
	storage, storageErr := files.NewS3Storage(1024 * 1000 * 5)
	if storageErr != nil {
		logger.Fatal(storageErr)
		return
	}

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Error loading .env file")
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		logger.Println("Defaulting to port to: ", 8082)
		port = "8082"
	}

	fileHandler := handlers.NewFileHandler(storage, logger)

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

	router.Post("/upload/images/{resourceId}/{filename}", fileHandler.Upload)
	router.Get("/upload/images/{resourceId}/{filename}", fileHandler.Download)

	router.Mount("/api/v1", router)

	server := &http.Server{
		Addr:     fmt.Sprintf(":%s", port),
		Handler:  router,
		ErrorLog: logger,
	}
	fmt.Printf("Starting file service on port %s .....\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
