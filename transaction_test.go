// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var fixtures string = filepath.Join("test", "fixtures")
var mapper ColumnMapper = ColumnMapper{
	"date":        0,
	"description": 1,
	"credit":      4,
	"refund":      5,
	"balance":     6,
}

func TestEmptyTransactionFromEmptyCSV(t *testing.T) {
	r := strings.NewReader("")
	transactions := TransactionsFromCSV(r, mapper)
	if l := len(transactions); l > 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestBadDate(t *testing.T) {
	r := openFixture(t, "bad-date.csv")
	defer r.Close()
	transactions := TransactionsFromCSV(r, mapper)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := &Transaction{
		Date:        UTCDate(t, 2021, 03, 18),
		Description: "pizza",
		Credit:      5.0,
		Refund:      0.0,
		Balance:     150.0,
	}
	expected.Date, _ = time.Parse("02.01.2006", "123")
	checkEquals(t, &transactions[0], expected)
}

func TestSingleTransactionFromSingleRowCSV(t *testing.T) {
	r := openFixture(t, "single-row.csv")
	defer r.Close()
	transactions := TransactionsFromCSV(r, mapper)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := &Transaction{
		Date:        UTCDate(t, 2021, 03, 18),
		Description: "pizza",
		Credit:      5.0,
		Refund:      0.0,
		Balance:     150.0,
	}
	checkEquals(t, &transactions[0], expected)
}

func TestMapsColumnsByGivenIndices(t *testing.T) {
	r := openFixture(t, "single-row.csv")
	defer r.Close()
	customMapper := ColumnMapper{"date": 7}
	transactions := TransactionsFromCSV(r, customMapper)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := UTCDate(t, 2021, 04, 21)
	if transactions[0].Date != expected {
		t.Errorf("ColumnMapper was not respected, got: %s expected: %s", transactions[0].Date, expected)
	}
}

func TestTransactionsFromCSV(t *testing.T) {
	r := openFixture(t, "multiple-rows.csv")
	defer r.Close()
	transactions := TransactionsFromCSV(r, mapper)
	if l := len(transactions); l != 4 {
		t.Errorf("expected 4 transactions got %d", l)
	}
	for i := 0; i < len(transactions); i++ {
		description := fmt.Sprintf("pizza%d", i)
		expected := &Transaction{
			Date:        UTCDate(t, 2021, 03, 18),
			Description: description,
			Credit:      5.0,
			Refund:      0.0,
			Balance:     150.0,
		}
		checkEquals(t, &transactions[i], expected)
	}

}

func openFixture(t *testing.T, fixtureFileName string) *os.File {
	r, err := os.Open(filepath.Join(fixtures, fixtureFileName))
	if err != nil {
		t.Fatalf("err while opening fixture file: %s", err)
	}
	return r
}

func UTCDate(t *testing.T, year int, month time.Month, day int) time.Time {
	timeZone, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatalf("err while loading location: %s", err)
	}
	return time.Date(year, month, day, 0, 0, 0, 0, timeZone)
}

func checkEquals(t *testing.T, actual *Transaction, expected *Transaction) {
	if *actual != *expected {
		t.Errorf("expected %v, received %v", expected, actual)
	}
}
