package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-txdb"
)

func TestingNewDB(t *testing.T) (*sql.DB, func()) {
	cfg := dbConfig{
		DBName:       "pgx",
		RootUserName: "pgx_root",
		UserName:     "pgx_api_test",
		Host:         "localhost",
		Port:         "5432",
		SSLMode:      "disable",
	}
	db := sql.OpenDB(txdb.New("pgx", cfg.URL()))
	return db, func() {
		db.Close()
	}
}

func TestingNewTx(t *testing.T) (*sql.Tx, func()) {
	cfg := dbConfig{
		DBName:       "pgx",
		RootUserName: "pgx_root",
		UserName:     "pgx_api_test",
		Host:         "localhost",
		Port:         "5432",
		SSLMode:      "disable",
	}
	db := sql.OpenDB(txdb.New("pgx", cfg.URL()))
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	return tx, func() {
		tx.Rollback()
		db.Close()
	}
}

func TestingSetupDB() (func(), error) {
	cfg := dbConfig{
		DBName:       "pgx",
		RootUserName: "pgx_root",
		UserName:     "pgx_api_test",
		Host:         "localhost",
		Port:         "5432",
		SSLMode:      "disable",
	}
	// prep db as root
	rootCon, err := sql.Open("pgx", cfg.RootURL())
	if err != nil {
		log.Fatal(err)
	}
	defer rootCon.Close()
	// create schema
	_, err = rootCon.Exec(fmt.Sprintf("CREATE SCHEMA %s AUTHORIZATION %s", cfg.UserName, cfg.UserName))
	if err != nil {
		log.Fatalf("failed to create schema: %s", err)
	}
	// prep db as test
	testCon, err := sql.Open("pgx", cfg.URL())
	if err != nil {
		log.Fatal(err)
	}
	defer testCon.Close()
	// create tables
	f, err := os.Open("./schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	s, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := testCon.Exec(string(s)); err != nil {
		log.Fatalf("failed to create tables: %s", err)
	}

	cleanup := func() {
		// prep db as root
		rootCon, err := sql.Open("pgx", cfg.RootURL())
		if err != nil {
			log.Fatal(err)
		}
		defer rootCon.Close()
		// create schema
		_, err = rootCon.Exec(fmt.Sprintf("DROP SCHEMA %s CASCADE;", cfg.UserName))
		if err != nil {
			log.Fatalf("failed to drop schema: %s", err)
		}
	}
	return cleanup, nil
}

func TestMain(m *testing.M) {
	flag.Parse()

	teardown, err := TestingSetupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer teardown()

	m.Run()
}
