// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"testing"
	"time"

	helpers "github.com/amitlevy21/gopher-street/test"
)

func TestNewFromEmptyTransactionAndEmptyClasses(t *testing.T) {
	expense := NewExpenses([]Transaction{}, &Classifier{}, &Tagger{})
	helpers.CheckEquals(t, expense, &Expenses{})
}

func TestNewFromEmptyTransaction(t *testing.T) {
	cl := NewTestClassifier()
	expense := NewExpenses([]Transaction{}, cl, &Tagger{})
	helpers.CheckEquals(t, expense, &Expenses{})
}

func TestNewFromUnMatchingClasses(t *testing.T) {
	cl := NewTestClassifier()
	tr := NewTestTransaction(t, "pizza")
	expense := NewExpenses([]Transaction{*tr}, cl, &Tagger{})
	helpers.CheckEquals(t, expense, &Expenses{{
		Date:   tr.Date,
		Amount: tr.Credit,
		Class:  tr.Description,
		Tags:   []Tag{},
	}})
}

func TestNewFromMatchingClasses(t *testing.T) {
	cl := NewTestClassifier()
	tr := NewTestTransaction(t, "description1")
	expense := NewExpenses([]Transaction{*tr}, cl, &Tagger{})
	helpers.CheckEquals(t, expense, &Expenses{{
		Date:   tr.Date,
		Amount: tr.Credit,
		Class:  "class1",
		Tags:   []Tag{},
	}})
}

func TestNewFromUnMatchingTaggerAndClassifier(t *testing.T) {
	cl := NewTestClassifier()
	tr := NewTestTransaction(t, "description1")
	tagger := &Tagger{classesToTags: map[string][]Tag{"nonExistClass": {"someTag"}}}
	expense := NewExpenses([]Transaction{*tr}, cl, tagger)
	helpers.CheckEquals(t, expense, &Expenses{{
		Date:   tr.Date,
		Amount: tr.Credit,
		Class:  "class1",
		Tags:   []Tag{},
	}})
}

func TestNewFromMatchingTaggerAndClassifier(t *testing.T) {
	cl := NewTestClassifier()
	tr := NewTestTransaction(t, "description1")
	tagger := NewTestTagger()
	expense := NewExpenses([]Transaction{*tr}, cl, tagger)
	helpers.CheckEquals(t, expense, &Expenses{{
		Date:   tr.Date,
		Amount: tr.Credit,
		Class:  "class1",
		Tags:   []Tag{"tag1", "tag2"},
	}})
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
