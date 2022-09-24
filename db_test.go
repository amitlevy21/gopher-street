package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	helpers "github.com/amitlevy21/gopher-street/test"
)

type BadDecoderCursor struct{}

type EndDecodeErrCursor struct {
	BadDecoderCursor
}

func (f *BadDecoderCursor) Next(ctx context.Context) bool {
	return true
}

func (f *BadDecoderCursor) Decode(val interface{}) error {
	return fmt.Errorf("always fail")
}

func (f *BadDecoderCursor) Err() error {
	return nil
}

func (f *BadDecoderCursor) Close(ctx context.Context) error {
	return fmt.Errorf("always fail")
}

func (f *EndDecodeErrCursor) Next(ctx context.Context) bool {
	return false
}

func (f *EndDecodeErrCursor) Decode(val interface{}) error {
	return nil
}

func (f *EndDecodeErrCursor) Err() error {
	return fmt.Errorf("always fail")
}

func TestBadURI(t *testing.T) {
	defer func() { _ = recover() }()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	openDB(ctx, "")
	t.Errorf("did not panic")
}

func TestBadPing(t *testing.T) {
	defer func() { _ = recover() }()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()
	openDB(ctx, "mongodb://hangup")
	t.Errorf("did not panic")
}

func TestClosedTooEarly(t *testing.T) {
	defer func() { _ = recover() }()
	ctx, cancel := context.WithCancel(context.Background())
	db := Instance(ctx, helpers.MongoURI)
	cancel()
	db.closeDB(ctx)
	db.closeDB(ctx)
	t.Errorf("did not panic")
}

func TestBadDecodeCursor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := Instance(ctx, helpers.MongoURI)
	defer db.closeDB(ctx)
	c := &BadDecoderCursor{}
	_, err := db.getExpensesFromCur(ctx, c)
	helpers.ExpectError(t, err)
}

func TestBadEndDecodeCursor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := Instance(ctx, helpers.MongoURI)
	defer db.closeDB(ctx)
	c := &EndDecodeErrCursor{}
	_, err := db.getExpensesFromCur(ctx, c)
	helpers.ExpectError(t, err)
}

func TestCloseBadCursor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := Instance(ctx, helpers.MongoURI)
	defer db.closeDB(ctx)
	c := &BadDecoderCursor{}
	db.closeCursor(ctx, c)
}

func TestErrWhileGetting(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := Instance(ctx, helpers.MongoURI)
	defer db.closeDB(ctx)

	ctx2, cancel := context.WithCancel(context.Background())
	cancel()
	exp, err := db.GetExpenses(ctx2)
	helpers.ExpectError(t, err)
	helpers.ExpectEquals(t, &Expenses{}, exp)
}

func TestDB(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := Instance(ctx, helpers.MongoURI)
	defer db.closeDB(ctx)
	t.Run("TestGetEmptyExpenses", func(t *testing.T) {
		helpers.FailTestIfErr(t, db.dropDB(ctx))
		expenses, err := db.GetExpenses(ctx)
		helpers.FailTestIfErr(t, err)
		helpers.ExpectEquals(t, expenses, &Expenses{})
	})
	t.Run("TestWriteEmptyExpense", func(t *testing.T) {
		helpers.FailTestIfErr(t, db.dropDB(ctx))
		err := db.WriteExpenses(ctx, &Expenses{})
		helpers.ExpectError(t, err)
	})
	t.Run("TestWriteNonEmptyExpense", func(t *testing.T) {
		helpers.FailTestIfErr(t, db.dropDB(ctx))
		err := db.WriteExpenses(ctx, &Expenses{Classified: []*Expense{NewTestExpense(t)}})
		helpers.FailTestIfErr(t, err)
	})
	t.Run("TestGetNonEmptyExpense", func(t *testing.T) {
		helpers.FailTestIfErr(t, db.dropDB(ctx))
		e := NewTestExpense(t)
		err := db.WriteExpenses(ctx, &Expenses{Classified: []*Expense{e}})
		helpers.FailTestIfErr(t, err)
		exp, err := db.GetExpenses(ctx)
		helpers.FailTestIfErr(t, err)
		helpers.ExpectEquals(t, exp, &Expenses{Classified: []*Expense{e}})
	})
	t.Run("TestWriteExpenses", func(t *testing.T) {
		helpers.FailTestIfErr(t, db.dropDB(ctx))
		err := db.WriteExpenses(ctx, NewTestExpenses(t))
		helpers.FailTestIfErr(t, err)
	})
	t.Run("TestGetExpenses", func(t *testing.T) {
		helpers.FailTestIfErr(t, db.dropDB(ctx))
		exps := NewTestExpenses(t)
		err := db.WriteExpenses(ctx, exps)
		helpers.FailTestIfErr(t, err)
		res, err := db.GetExpenses(ctx)
		helpers.FailTestIfErr(t, err)
		helpers.ExpectEquals(t, res, exps)
	})
}
