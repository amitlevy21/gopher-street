// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"testing"
	"time"

	helpers "github.com/amitlevy21/gopher-street/test"
)

func TestTagString(t *testing.T) {
	tags := []Tag{Recurring, Crucial, None}
	tagStrings := [...]string{"Recurring", "Crucial", "None"}
	for i, tag := range tags {
		if fmt.Sprint(tag) != tagStrings[i] {
			t.Errorf("expected %s received %s", tagStrings[i], fmt.Sprint(Recurring))
		}
	}
}

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
		Date:   tr.Date.Month(),
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
		Date:   tr.Date.Month(),
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
		Date:   tr.Date.Month(),
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
		Date:   tr.Date.Month(),
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
				Date:   time.March,
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
			{
				Date:   time.March,
				Class:  "shirt",
				Amount: 20,
				Tags:   []Tag{},
			},
		},
		time.April: {
			{
				Date:   time.April,
				Class:  "Living",
				Amount: 3500,
				Tags:   []Tag{"Crucial"},
			},
		},
		time.May: {
			{
				Date:   time.May,
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
				Date:   time.March,
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
			{
				Date:   time.May,
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
		},
		"shirt": {
			{
				Date:   time.March,
				Class:  "shirt",
				Amount: 20,
				Tags:   []Tag{},
			},
		},
		"Living": {
			{
				Date:   time.April,
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
				Date:   time.April,
				Class:  "Living",
				Amount: 3500,
				Tags:   []Tag{"Crucial"},
			},
		},
		"None": {
			{
				Date:   time.March,
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
			{
				Date:   time.March,
				Class:  "shirt",
				Amount: 20,
				Tags:   []Tag{},
			},
			{
				Date:   time.May,
				Class:  "Eating outside",
				Amount: 5,
				Tags:   []Tag{},
			},
		},
	})
}
