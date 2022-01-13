// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"path/filepath"
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
)

var mapper = &ColMapper{
	Date:        0,
	Description: 1,
	Credit:      4,
	Refund:      5,
	Balance:     6,
}
var emptySubsetter = &RowSubsetter{}

var layout = "02.01.2006"

func TestEmptyTransactionFromEmptyCSV(t *testing.T) {
	c := NewCardTransactions([][]string{}, mapper, emptySubsetter, layout)
	transactions, _ := c.Transactions()
	if l := len(transactions); l > 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestSkipsBadDateRecord(t *testing.T) {
	badCSV := helpers.ReadCSVFixture(t, filepath.Join("transactions", "bad-date.csv"))
	c := NewCardTransactions(badCSV, mapper, emptySubsetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestSkipsBadRecords(t *testing.T) {
	badCSV := helpers.ReadCSVFixture(t, filepath.Join("transactions", "bad-multi.csv"))
	c := NewCardTransactions(badCSV, mapper, emptySubsetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 2 {
		t.Fatalf("expected 2 transactions got %d", l)
	}
	expected := NewTestTransaction(t, "pizza1")
	expected2 := NewTestTransaction(t, "pizza3")
	expected2.Balance = 0
	helpers.CheckEquals(t, &transactions[0], expected)
	helpers.CheckEquals(t, &transactions[1], expected2)
}

func TestSingleTransactionFromSingleRowCSV(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join("transactions", "single-row.csv"))
	c := NewCardTransactions(data, mapper, emptySubsetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := NewTestTransaction(t, "pizza")
	helpers.CheckEquals(t, &transactions[0], expected)
}

func TestMapsColumnsByGivenIndices(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join("transactions", "single-row.csv"))
	customMapper := &ColMapper{
		Date:        7,
		Description: 1,
		Credit:      4,
		Refund:      5,
		Balance:     6,
	}
	c := NewCardTransactions(data, customMapper, emptySubsetter, layout)
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

func TestMapperOutOfRange(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join("transactions", "single-row.csv"))
	customMapper := &ColMapper{Date: 23, Credit: 2}
	c := NewCardTransactions(data, customMapper, emptySubsetter, layout)
	_, err := c.Transactions()
	helpers.ExpectError(t, err)
}

func TestRefundTransaction(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join("transactions", "with-refund.csv"))
	c := NewCardTransactions(data, mapper, emptySubsetter, layout)
	trans, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if trans[0].Refund != 5.0 {
		t.Errorf("Refund transaction not created: %v", trans[0])
	}
}

func TestMapperMissingCreditAndRefund(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join("transactions", "single-row.csv"))
	customMapper := &ColMapper{Date: 0}
	c := NewCardTransactions(data, customMapper, emptySubsetter, layout)
	trans, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(trans); l != 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestEmptySubsetterShouldReadAll(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join("transactions", "multiple-rows.csv"))
	c := NewCardTransactions(data, mapper, emptySubsetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 4 {
		t.Errorf("expected 4 transactions got %d", l)
	}
}

func TestOutOfUpperBoundRangeSubsetter(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join("transactions", "multiple-rows.csv"))
	subsetter := &RowSubsetter{1, 5}
	c := NewCardTransactions(data, mapper, subsetter, layout)
	_, err := c.Transactions()
	helpers.ExpectError(t, err)
}

func TestSubsetsRowsByGivenIndices(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join("transactions", "multiple-rows.csv"))
	subsetter := &RowSubsetter{1, 3}
	c := NewCardTransactions(data, mapper, subsetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	expected := int(subsetter.End) - int(subsetter.Start)
	if len(transactions) != expected {
		t.Errorf("expected %d but got %d", expected, len(transactions))
	}
	for i, j := subsetter.Start, 0; i < subsetter.End; i, j = i+1, j+1 {
		description := fmt.Sprintf("pizza%d", i)
		expected := NewTestTransaction(t, description)
		helpers.CheckEquals(t, &transactions[j], expected)
	}
}

func TestTransactionsFromCSV(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join("transactions", "multiple-rows.csv"))
	c := NewCardTransactions(data, mapper, emptySubsetter, layout)
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
