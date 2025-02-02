package config

import (
	"database/sql"
	"github.com/zain2323/cronium/services/internal/database"
	"log"
)

type ApiConfig struct {
	DB *database.Queries
}

func New(dbUrl string) (*ApiConfig, error) {
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &ApiConfig{
		DB: database.New(conn),
	}, nil
}
