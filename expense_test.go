// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
)

func TestTagString(t *testing.T) {
	tags := []Tag{Recurring, Crucial}
	tagStrings := [...]string{"Recurring", "Crucial"}
	for i, tag := range tags {
		if fmt.Sprint(tag) != tagStrings[i] {
			t.Errorf("expected %s received %s", tagStrings[i], fmt.Sprint(Recurring))
		}
	}
}

func TestNewFromEmptyTransactionAndEmptyClasses(t *testing.T) {
	expense := NewExpense([]*Transaction{}, &Classifier{})
	helpers.CheckEquals(t, expense, []Expense{})
}

func TestNewFromEmptyTransaction(t *testing.T) {
	cl := NewTestClassifier(t)
	expense := NewExpense([]*Transaction{}, cl)
	helpers.CheckEquals(t, expense, []Expense{})
}

func TestNewFromUnmatchingClasses(t *testing.T) {
	cl := NewTestClassifier(t)
	tr := NewTestTransaction(t, "pizza")
	expense := NewExpense([]*Transaction{tr}, cl)
	helpers.CheckEquals(t, expense, []Expense{{
		Date:   tr.Date.Month(),
		Amount: tr.Credit,
		Class:  tr.Description,
		Tags:   []Tag{},
	}})
}
