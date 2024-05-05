package main

import (
	"context"
	"testing"
	"time"
)

func TestCurrentTime(t *testing.T) {
	db, cleanup := TestingNewDB(t)
	defer cleanup()
	ctx := context.Background()
	tm, err := CurrentTime(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tm)
}

func TestT1Create(t *testing.T) {
	tx, cleanup := TestingNewTx(t)
	defer cleanup()
	ctx := context.Background()
	t1 := T1{
		Val:       "val",
		CreatedAt: time.Now(),
	}
	if err := t1.Create(ctx, tx); err != nil {
		t.Fatal(err)
	}
	t.Log(t1)
	t2 := T1{
		Val:       "newval",
		CreatedAt: time.Now(),
	}
	if err := t2.Create(ctx, tx); err != nil {
		t.Fatal(err)
	}
	t.Log(t2)
}

func TestT1CreateAndCommit(t *testing.T) {
	db, cleanup := TestingNewDB(t)
	defer cleanup()
	ctx := context.Background()

	t1 := &T1{
		Val:       "val",
		CreatedAt: time.Now(),
	}
	if err := CreateT1AndCommit(ctx, db, t1); err != nil {
		t.Fatal(err)
	}
	t2 := &T1{
		Val:       "val2",
		CreatedAt: time.Now(),
	}
	if err := CreateT1AndCommit(ctx, db, t2); err != nil {
		t.Fatal(err)
	}
	rows, err := db.Query(`select id, val, created_at from t1`)
	if err != nil {
		t.Fatal(err)
	}

	for rows.Next() {
		var ta T1
		err := rows.Scan(
			&ta.ID,
			&ta.Val,
			&ta.CreatedAt,
		)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ta)
	}
}

func TestT1CreateAndRollback(t *testing.T) {
	db, cleanup := TestingNewDB(t)
	defer cleanup()
	ctx := context.Background()

	t1 := &T1{
		Val:       "val",
		CreatedAt: time.Now(),
	}
	t2 := &T1{
		Val:       "val",
		CreatedAt: time.Now(),
	}
	if err := CreateTwoT1s(ctx, db, t1, t2); err != nil {
		t.Fatal(err)
	}

	rows, err := db.Query(`select id, val, created_at from t1`)
	if err != nil {
		t.Fatal(err)
	}

	for rows.Next() {
		var ta T1
		err := rows.Scan(
			&ta.ID,
			&ta.Val,
			&ta.CreatedAt,
		)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("t1=%v", ta)
	}

	rows, err = db.Query(`select id, val, created_at from error_t1`)
	if err != nil {
		t.Fatal(err)
	}

	for rows.Next() {
		var ta T1
		err := rows.Scan(
			&ta.ID,
			&ta.Val,
			&ta.CreatedAt,
		)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("error_t1=%v", ta)
	}
}
