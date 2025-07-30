package dup

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func Dup() {
	db, err := Open(DefaultPostgresConfig())
	if err != nil {
		return
	}

	uid := "51780aa2-0263-4579-b2ac-53f63dadd7ae"
	row := db.QueryRow(`
	INSERT INTO sessions
	VALUES($1, $2, $3) RETURNING id`, uid, uid, uid)

	if err := row.Scan(); err != nil {
		log.Fatalf("sessions service create: %v", err)
	}
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "devuser",
		Password: "devpass",
		Database: "weblensgo",
		SSLMode:  "disable",
	}
}

func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}
	return db, nil
}
