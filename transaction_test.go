// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"path/filepath"
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
var emptySubsetter []int = []int{}

func TestEmptyTransactionFromEmptyCSV(t *testing.T) {
	r := strings.NewReader("")
	c := NewCardTransactions(r, mapper, emptySubsetter)
	transactions, _ := c.Transactions()
	if l := len(transactions); l > 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestSkipsBadDateRecord(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "bad-date.csv"))
	defer r.Close()
	c := NewCardTransactions(r, mapper, emptySubsetter)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestSkipsBadRecords(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "bad-multi.csv"))
	defer r.Close()
	c := NewCardTransactions(r, mapper, emptySubsetter)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := NewTestTransaction(t, "pizza1")
	helpers.CheckEquals(t, &transactions[0], expected)
}

func TestSingleTransactionFromSingleRowCSV(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "single-row.csv"))
	defer r.Close()
	c := NewCardTransactions(r, mapper, emptySubsetter)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := NewTestTransaction(t, "pizza")
	helpers.CheckEquals(t, &transactions[0], expected)
}

func TestMapsColumnsByGivenIndices(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "single-row.csv"))
	defer r.Close()
	customMapper := map[string]int{
		"date":        7,
		"description": 1,
		"credit":      4,
		"refund":      5,
		"balance":     6,
	}
	c := NewCardTransactions(r, customMapper, emptySubsetter)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := helpers.UTCDate(t, 2021, 04, 21)
	if transactions[0].Date != expected {
		t.Errorf("ColumnMapper was not respected, got: %s expected: %s", transactions[0].Date, expected)
	}
}

func TestBadMapper(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "single-row.csv"))
	defer r.Close()
	customMapper := map[string]int{"date": 9, "not_exist": 23, "credit": 2}
	c := NewCardTransactions(r, customMapper, emptySubsetter)
	_, err := c.Transactions()
	helpers.ExpectError(t, err)
}

func TestEmptySubsetterShouldReadAll(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "multiple-rows.csv"))
	defer r.Close()
	c := NewCardTransactions(r, mapper, emptySubsetter)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 4 {
		t.Errorf("expected 4 transactions got %d", l)
	}
}

func TestOutOfUpperRangeSubsetter(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "multiple-rows.csv"))
	defer r.Close()
	subsetter := []int{1, 4}
	c := NewCardTransactions(r, mapper, subsetter)
	_, err := c.Transactions()
	helpers.ExpectError(t, err)
}

func TestOutOfLowerRangeSubsetter(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "single-row.csv"))
	defer r.Close()
	subsetter := []int{-1, 2}
	c := NewCardTransactions(r, mapper, subsetter)
	_, err := c.Transactions()
	helpers.ExpectError(t, err)
}

func TestUnorderedSubsetter(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "multiple-rows.csv"))
	defer r.Close()
	subsetter := []int{3, 1}
	c := NewCardTransactions(r, mapper, subsetter)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	for i, rowIndex := range subsetter {
		description := fmt.Sprintf("pizza%d", rowIndex)
		expected := NewTestTransaction(t, description)
		helpers.CheckEquals(t, &transactions[i], expected)
	}
}

func TestSubsetsRowsByGivenIndices(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "multiple-rows.csv"))
	defer r.Close()
	subsetter := []int{1, 3}
	c := NewCardTransactions(r, mapper, subsetter)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if len(transactions) != len(subsetter) {
		t.Errorf("expected %d but got %d", len(subsetter), len(transactions))
	}
	for i, rowIndex := range subsetter {
		description := fmt.Sprintf("pizza%d", rowIndex)
		expected := NewTestTransaction(t, description)
		helpers.CheckEquals(t, &transactions[i], expected)
	}
}

func TestTransactionsFromCSV(t *testing.T) {
	r := helpers.OpenFixture(t, filepath.Join("transactions", "multiple-rows.csv"))
	defer r.Close()
	c := NewCardTransactions(r, mapper, emptySubsetter)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 4 {
		t.Errorf("expected 4 transactions got %d", l)
	}
	for i := 0; i < len(transactions); i++ {
		description := fmt.Sprintf("pizza%d", i)
		expected := NewTestTransaction(t, description)
		helpers.CheckEquals(t, &transactions[i], expected)
	}
}
