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

type ErrorT1 struct {
	ID        int64
	Val       string
	CreatedAt time.Time
}

func (t *ErrorT1) Create(ctx context.Context, tx *sql.Tx) error {
	if err := tx.QueryRowContext(
		ctx, `insert into error_t1 (val ,created_at) values ($1, $2) returning id`,
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

func CreateTwoT1s(ctx context.Context, db *sql.DB, t1, t2 *T1) error {
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

	tx2, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin tx2")
	}
	defer tx2.Rollback()
	if err := t2.Create(ctx, tx2); err != nil {
		tx2.Rollback()
		tx3, err := db.Begin()
		if err != nil {
			return errors.Wrap(err, "failed to begin tx3")
		}
		defer tx3.Rollback()
		et2 := ErrorT1{
			Val:       t2.Val,
			CreatedAt: t2.CreatedAt,
		}
		if err := et2.Create(ctx, tx3); err != nil {
			return errors.Wrap(err, "failed to create ErrorT1")
		}
		if err := tx3.Commit(); err != nil {
			return errors.Wrap(err, "failed to commit ErrorT1")
		}
		return nil
	}
	if err := tx2.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit the first t1")
	}
	return nil
}
