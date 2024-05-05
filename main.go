package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"
)

type dbConfig struct {
	DBName       string
	UserName     string
	RootUserName string
	SSLMode      string
	Host         string
	Port         string
}

func (c dbConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s@%s:%s/%s?sslmode=%s",
		c.UserName, c.Host, c.Port, c.DBName, c.SSLMode)
}

func (c dbConfig) RootURL() string {
	return fmt.Sprintf(
		"postgres://%s@%s:%s/%s?sslmode=%s",
		c.RootUserName, c.Host, c.Port, c.DBName, c.SSLMode)
}

func main() {
	cfg := dbConfig{
		DBName:       "pgx",
		RootUserName: "pgx_root",
		UserName:     "pgx_api",
		Host:         "localhost",
		Port:         "5432",
		SSLMode:      "disable",
	}
	db, err := sql.Open("pgx", cfg.URL())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	t2, err := CurrentTime(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(t2)
}
