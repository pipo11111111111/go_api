package db

import (
	"database/sql"
	"facturadirecta_api/configs"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func OpenConnection() (*sql.DB, error) {
	conf := configs.GetDB()
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Pass, conf.Database)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return conn, nil
}
