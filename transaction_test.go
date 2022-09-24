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
var emptySubSetter = &RowSubSetter{}

var layout = "02.01.2006"

func TestEmptyTransactionFromEmptyCSV(t *testing.T) {
	c := NewCardTransactions([][]string{}, mapper, emptySubSetter, layout)
	transactions, _ := c.Transactions()
	if l := len(transactions); l > 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestSkipsBadDateRecord(t *testing.T) {
	badCSV := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "bad-date.csv"))
	c := NewCardTransactions(badCSV, mapper, emptySubSetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestSkipsBadRecords(t *testing.T) {
	badCSV := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "bad-multi.csv"))
	c := NewCardTransactions(badCSV, mapper, emptySubSetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 2 {
		t.Fatalf("expected 2 transactions got %d", l)
	}
	expected := NewTestTransaction(t, "pizza1")
	expected2 := NewTestTransaction(t, "pizza3")
	expected2.Balance = 0
	helpers.ExpectEquals(t, &transactions[0], expected)
	helpers.ExpectEquals(t, &transactions[1], expected2)
}

func TestSingleTransactionFromSingleRowCSV(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "single-row.csv"))
	c := NewCardTransactions(data, mapper, emptySubSetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 1 {
		t.Fatalf("expected 1 transactions got %d", l)
	}
	expected := NewTestTransaction(t, "pizza")
	helpers.ExpectEquals(t, &transactions[0], expected)
}

func TestMapsColumnsByGivenIndices(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "single-row.csv"))
	customMapper := &ColMapper{
		Date:        7,
		Description: 1,
		Credit:      4,
		Refund:      5,
		Balance:     6,
	}
	c := NewCardTransactions(data, customMapper, emptySubSetter, layout)
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
	data := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "single-row.csv"))
	customMapper := &ColMapper{Date: 23, Credit: 2}
	c := NewCardTransactions(data, customMapper, emptySubSetter, layout)
	_, err := c.Transactions()
	helpers.ExpectError(t, err)
}

func TestRefundTransaction(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "with-refund.csv"))
	c := NewCardTransactions(data, mapper, emptySubSetter, layout)
	trans, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if trans[0].Refund != 5.0 {
		t.Errorf("Refund transaction not created: %v", trans[0])
	}
}

func TestMapperMissingCreditAndRefund(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "single-row.csv"))
	customMapper := &ColMapper{Date: 0}
	c := NewCardTransactions(data, customMapper, emptySubSetter, layout)
	trans, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(trans); l != 0 {
		t.Errorf("expected 0 transactions got %d", l)
	}
}

func TestEmptySubSetterShouldReadAll(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "multiple-rows.csv"))
	c := NewCardTransactions(data, mapper, emptySubSetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 4 {
		t.Errorf("expected 4 transactions got %d", l)
	}
}

func TestOutOfUpperBoundRangeSubSetter(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "multiple-rows.csv"))
	subSetter := &RowSubSetter{1, 5}
	c := NewCardTransactions(data, mapper, subSetter, layout)
	_, err := c.Transactions()
	helpers.ExpectError(t, err)
}

func TestSubsetsRowsByGivenIndices(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "multiple-rows.csv"))
	subSetter := &RowSubSetter{1, 3}
	c := NewCardTransactions(data, mapper, subSetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	expected := int(subSetter.End) - int(subSetter.Start)
	if len(transactions) != expected {
		t.Errorf("expected %d but got %d", expected, len(transactions))
	}
	for i, j := subSetter.Start, 0; i < subSetter.End; i, j = i+1, j+1 {
		description := fmt.Sprintf("pizza%d", i)
		expected := NewTestTransaction(t, description)
		helpers.ExpectEquals(t, &transactions[j], expected)
	}
}

func TestTransactionsFromCSV(t *testing.T) {
	data := helpers.ReadCSVFixture(t, filepath.Join(CSVTransactionsPath, "multiple-rows.csv"))
	c := NewCardTransactions(data, mapper, emptySubSetter, layout)
	transactions, err := c.Transactions()
	helpers.FailTestIfErr(t, err)
	if l := len(transactions); l != 4 {
		t.Errorf("expected 4 transactions got %d", l)
	}
	for i := 0; i < len(transactions); i++ {
		description := fmt.Sprintf("pizza%d", i)
		expected := NewTestTransaction(t, description)
		helpers.ExpectEquals(t, &transactions[i], expected)
	}
}
