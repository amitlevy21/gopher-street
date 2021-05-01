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
	expense := NewExpenses([]Transaction{}, &Classifier{}, &Tagger{})
	helpers.CheckEquals(t, expense, &Expenses{})
}

func TestNewFromEmptyTransaction(t *testing.T) {
	cl := NewTestClassifier(t, "classifier.yml")
	expense := NewExpenses([]Transaction{}, cl, &Tagger{})
	helpers.CheckEquals(t, expense, &Expenses{})
}

func TestNewFromUnMatchingClasses(t *testing.T) {
	cl := NewTestClassifier(t, "classifier.yml")
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
	cl := NewTestClassifier(t, "classifier.yml")
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
	cl := NewTestClassifier(t, "classifier.yml")
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
	cl := NewTestClassifier(t, "classifier.yml")
	tr := NewTestTransaction(t, "description1")
	tagger := NewTestTagger(t, "tagger.yml")
	expense := NewExpenses([]Transaction{*tr}, cl, tagger)
	helpers.CheckEquals(t, expense, &Expenses{{
		Date:   tr.Date.Month(),
		Amount: tr.Credit,
		Class:  "class1",
		Tags:   []Tag{"tag1", "tag2"},
	}})
}

