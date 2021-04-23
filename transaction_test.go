// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"strings"
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
)

var mapper map[string]int = map[string]int{
	"date":        0,
	"description": 1,
	"credit":      4,
	"refund":      5,
	"balance":     6,
}

func TestEmptyTransactionFromEmptyCSV(t *testing.T) {
	r := strings.NewReader("")
	hunk := NewCardTransactions(r, mapper)
	transactions, _ := hunk.Transactions()
	if l := len(transactions); l > 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestSkipsBadDateRecord(t *testing.T) {
	r := helpers.OpenFixture(t, "bad-date.csv")
	defer r.Close()
	hunk := NewCardTransactions(r, mapper)
	transactions, err := hunk.Transactions()
	helpers.Check(t, err)
	if l := len(transactions); l != 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestSkipsBadRecords(t *testing.T) {
	r := helpers.OpenFixture(t, "bad-multi.csv")
	defer r.Close()
	hunk := NewCardTransactions(r, mapper)
	transactions, err := hunk.Transactions()
	helpers.Check(t, err)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := &Transaction{
		Date:        helpers.UTCDate(t, 2021, 03, 18),
		Description: "pizza1",
		Credit:      5.0,
		Balance:     150.0,
	}
	helpers.CheckEquals(t, &transactions[0], expected)
}

func TestSingleTransactionFromSingleRowCSV(t *testing.T) {
	r := helpers.OpenFixture(t, "single-row.csv")
	defer r.Close()
	hunk := NewCardTransactions(r, mapper)
	transactions, err := hunk.Transactions()
	helpers.Check(t, err)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := &Transaction{
		Date:        helpers.UTCDate(t, 2021, 03, 18),
		Description: "pizza",
		Credit:      5.0,
		Refund:      0.0,
		Balance:     150.0,
	}
	helpers.CheckEquals(t, &transactions[0], expected)
}

func TestMapsColumnsByGivenIndices(t *testing.T) {
	r := helpers.OpenFixture(t, "single-row.csv")
	defer r.Close()
	customMapper := map[string]int{
		"date":        7,
		"description": 1,
		"credit":      4,
		"refund":      5,
		"balance":     6,
	}
	hunk := NewCardTransactions(r, customMapper)
	transactions, err := hunk.Transactions()
	helpers.Check(t, err)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := helpers.UTCDate(t, 2021, 04, 21)
	if transactions[0].Date != expected {
		t.Errorf("ColumnMapper was not respected, got: %s expected: %s", transactions[0].Date, expected)
	}
}

func TestIgnoresBadMapper(t *testing.T) {
	r := helpers.OpenFixture(t, "single-row.csv")
	defer r.Close()
	customMapper := map[string]int{"date": 9, "not_exist": 23, "credit": 2}
	hunk := NewCardTransactions(r, customMapper)
	transactions, err := hunk.Transactions()
	helpers.ExpectError(t, err)
	if l := len(transactions); l != 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestTransactionsFromCSV(t *testing.T) {
	r := helpers.OpenFixture(t, "multiple-rows.csv")
	defer r.Close()
	hunk := NewCardTransactions(r, mapper)
	transactions, err := hunk.Transactions()
	helpers.Check(t, err)
	if l := len(transactions); l != 4 {
		t.Errorf("expected 4 transactions got %d", l)
	}
	for i := 0; i < len(transactions); i++ {
		description := fmt.Sprintf("pizza%d", i)
		expected := &Transaction{
			Date:        helpers.UTCDate(t, 2021, 03, 18),
			Description: description,
			Credit:      5.0,
			Refund:      0.0,
			Balance:     150.0,
		}
		helpers.CheckEquals(t, &transactions[i], expected)
	}
}
