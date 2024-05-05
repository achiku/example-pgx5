package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

func CurrentTime(ctx context.Context, db *sql.DB) (*time.Time, error) {
	var t time.Time
	if err := db.QueryRowContext(ctx, `select now()`).Scan(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

type T1 struct {
	ID        int64
	Val       string
	CreatedAt time.Time
}

func (t *T1) Create(ctx context.Context, tx *sql.Tx) error {
	if err := tx.QueryRowContext(
		ctx, `insert into t1 (val ,created_at) values ($1, $2) returning id`,
		&t.Val, &t.CreatedAt).Scan(&t.ID); err != nil {
		return err
	}
	return nil
}

type DuplicateT1 struct {
	ID        int64
	Val       string
	CreatedAt time.Time
}

func (t *DuplicateT1) Create(ctx context.Context, tx *sql.Tx) error {
	if err := tx.QueryRowContext(
		ctx, `insert into t1 (val ,created_at) values ($1, $2) returning id`,
		&t.Val, &t.CreatedAt).Scan(&t.ID); err != nil {
		return err
	}
	return nil
}

func CreateT1AndCommit(ctx context.Context, db *sql.DB, t1 *T1) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback()

	if err := t1.Create(ctx, tx); err != nil {
		return errors.Wrap(err, "failed to create the first t1")
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit the first t1")
	}
	return nil
}
