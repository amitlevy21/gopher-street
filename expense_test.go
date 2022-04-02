// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"path/filepath"
	"testing"
	"time"

	helpers "github.com/amitlevy21/gopher-street/test"
)

func TestExpenseCreation(t *testing.T) {
	testCases := []struct {
		desc         string
		transactions []Transaction
		classifier   *Classifier
		tagger       *Tagger
		expected     *Expenses
	}{
		{
			desc:         "From empty data -> empty expenses",
			transactions: []Transaction{},
			classifier:   &Classifier{},
			tagger:       &Tagger{},
			expected:     &Expenses{},
		},
		{
			desc:         "From empty transactions -> empty expenses",
			transactions: []Transaction{},
			classifier:   NewTestClassifier(),
			tagger:       &Tagger{},
			expected:     &Expenses{},
		},
		{
			desc:         "From non-matching classes -> all expenses. Class same as description",
			transactions: []Transaction{*NewTestTransaction(t, "pizza")},
			classifier:   NewTestClassifier(),
			tagger:       &Tagger{},
			expected: &Expenses{{
				Date:   helpers.UTCDate(t, 2021, time.March, 18),
				Amount: 5.0,
				Class:  "pizza",
				Tags:   []Tag{},
			}},
		},
		{
			desc:         "From matching classes -> all expenses. Class set by classifier",
			transactions: []Transaction{*NewTestTransaction(t, "description1")},
			classifier:   NewTestClassifier(),
			tagger:       &Tagger{},
			expected: &Expenses{{
				Date:   helpers.UTCDate(t, 2021, time.March, 18),
				Amount: 5.0,
				Class:  "class1",
				Tags:   []Tag{},
			}},
		},
		{
			desc:         "From non-matching tags -> all expenses. Class set by classifier",
			transactions: []Transaction{*NewTestTransaction(t, "description1")},
			classifier:   NewTestClassifier(),
			tagger:       &Tagger{classesToTags: map[string][]Tag{"nonExistClass": {"someTag"}}},
			expected: &Expenses{{
				Date:   helpers.UTCDate(t, 2021, time.March, 18),
				Amount: 5.0,
				Class:  "class1",
				Tags:   []Tag{},
			}},
		},
		{
			desc:         "From matching tags and classses -> all expenses. Class set by classifier, if tags match they should show",
			transactions: []Transaction{*NewTestTransaction(t, "description1")},
			classifier:   NewTestClassifier(),
			tagger:       NewTestTagger(),
			expected: &Expenses{{
				Date:   helpers.UTCDate(t, 2021, time.March, 18),
				Amount: 5.0,
				Class:  "class1",
				Tags:   []Tag{"tag1", "tag2"},
			}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			expense := NewExpenses(tC.transactions, tC.classifier, tC.tagger)
			helpers.CheckEquals(t, expense, tC.expected)
		})
	}
}

func TestGetExpensesBadFile(t *testing.T) {
	_, err := getExpensesFromFile(&ConfigData{}, "unsupported.json")
	helpers.ExpectError(t, err)
}

func TestGetExpensesNonExist(t *testing.T) {
	_, err := getExpensesFromFile(&ConfigData{}, "non-exist.csv")
	helpers.ExpectError(t, err)
}

func TestGetExpensesFileDontAppearInConfig(t *testing.T) {
	file := filepath.Join(CSVTransactionsPath, "with-refund.csv")
	_, err := getExpensesFromFile(NewTestConfig(), file)
	helpers.ExpectError(t, err)
}

func TestGetExpensesFromFile(t *testing.T) {
	file := filepath.Join(CSVTransactionsPath, "multiple-rows.csv")
	exps, err := getExpensesFromFile(NewTestConfig(), file)
	helpers.FailTestIfErr(t, err)
	helpers.CheckEquals(t, exps, &Expenses{
		{
			Date:   helpers.UTCDate(t, 2021, time.March, 18),
			Amount: 5.0,
			Class:  "pizza1",
			Tags:   []string{},
		},
		{
			Date:   helpers.UTCDate(t, 2021, time.March, 18),
			Amount: 5.0,
			Class:  "pizza2",
			Tags:   []string{},
		},
		{
			Date:   helpers.UTCDate(t, 2021, time.March, 18),
			Amount: 5.0,
			Class:  "pizza3",
			Tags:   []string{},
		},
	})
}

func TestGroupByDate(t *testing.T) {
	cl := NewTestClassifierWithData()
	ct := NewTestCardTransactions(t, "data.csv")
	tagger := NewTestTaggerWithData()
	trs, err := ct.Transactions()
	helpers.FailTestIfErr(t, err)
	expense := NewExpenses(trs, cl, tagger)
	byMonth := expense.GroupByMonth()
	helpers.CheckEquals(t, byMonth, map[time.Month]Expenses{
		time.March: {
			{
				Date:   helpers.UTCDate(t, 2021, time.March, 18),
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
			{
				Date:   helpers.UTCDate(t, 2021, time.March, 22),
				Class:  "shirt",
				Amount: 20,
				Tags:   []Tag{},
			},
		},
		time.April: {
			{
				Date:   helpers.UTCDate(t, 2021, time.April, 24),
				Class:  "Living",
				Amount: 3500,
				Tags:   []Tag{"Crucial"},
			},
		},
		time.May: {
			{
				Date:   helpers.UTCDate(t, 2021, time.May, 5),
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
		},
	})
}

func TestGroupByClass(t *testing.T) {
	cl := NewTestClassifierWithData()
	ct := NewTestCardTransactions(t, "data.csv")
	tagger := NewTestTaggerWithData()
	trs, err := ct.Transactions()
	helpers.FailTestIfErr(t, err)
	expense := NewExpenses(trs, cl, tagger)
	byMonth := expense.GroupByClass()
	helpers.CheckEquals(t, byMonth, map[string]Expenses{
		"Eating outside": {
			{
				Date:   helpers.UTCDate(t, 2021, time.March, 18),
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
			{
				Date:   helpers.UTCDate(t, 2021, time.May, 5),
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
		},
		"shirt": {
			{
				Date:   helpers.UTCDate(t, 2021, time.March, 22),
				Class:  "shirt",
				Amount: 20,
				Tags:   []Tag{},
			},
		},
		"Living": {
			{
				Date:   helpers.UTCDate(t, 2021, time.April, 24),
				Class:  "Living",
				Amount: 3500,
				Tags:   []Tag{"Crucial"},
			},
		},
	})
}

func TestGroupByTag(t *testing.T) {
	cl := NewTestClassifierWithData()
	ct := NewTestCardTransactions(t, "data.csv")
	tagger := NewTestTaggerWithData()
	trs, err := ct.Transactions()
	helpers.FailTestIfErr(t, err)
	expense := NewExpenses(trs, cl, tagger)
	byMonth := expense.GroupByTag()
	helpers.CheckEquals(t, byMonth, map[Tag]Expenses{
		"Crucial": {
			{
				Date:   helpers.UTCDate(t, 2021, time.April, 24),
				Class:  "Living",
				Amount: 3500,
				Tags:   []Tag{"Crucial"},
			},
		},
		"None": {
			{
				Date:   helpers.UTCDate(t, 2021, time.March, 18),
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
			{
				Date:   helpers.UTCDate(t, 2021, time.March, 22),
				Class:  "shirt",
				Amount: 20,
				Tags:   []Tag{},
			},
			{
				Date:   helpers.UTCDate(t, 2021, time.May, 5),
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
		},
	})
}
