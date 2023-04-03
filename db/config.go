package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

// Config options to connect to DB
type Config struct {
	User         string
	Password     string
	Host         string
	DatabaseName string
	Port         string
}

// TestDB creates a db when testing
func TestDB() (*pgx.Conn, error) {
	port := "5433"
	if os.Getenv("CI") != "" {
		port = "5432"
	}
	return FromConfig(Config{
		User:         "postgres",
		Password:     "supersecret",
		Host:         "127.0.0.1",
		Port:         port,
		DatabaseName: "breakdown_test",
	})
}

// FromConfig generates a sql.DB based on a Config
func FromConfig(cfg Config) (*pgx.Conn, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	)

	return pgx.Connect(context.Background(), url)
}
