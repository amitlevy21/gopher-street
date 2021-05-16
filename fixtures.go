// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"path/filepath"
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
)

var fixtures string = filepath.Join("test", "fixtures")

func NewTestTransaction(t *testing.T, description string) *Transaction {
	return &Transaction{
		Date:        helpers.UTCDate(t, 2021, 03, 18),
		Description: description,
		Credit:      5.0,
		Refund:      0.0,
		Balance:     150.0,
	}
}

func NewTestCardTransactions(t *testing.T, fileName string) *CardTransactions {
	r := helpers.OpenFixture(t, filepath.Join("transactions", fileName))
	mapper := map[string]int{
		"date":        0,
		"description": 1,
		"credit":      4,
		"refund":      5,
		"balance":     6,
	}
	layout := "02.01.2006"
	return NewCardTransactions(r, mapper, []int{}, layout)
}

func NewTestClassifier() *Classifier {
	return &Classifier{map[string]string{
		"description1": "class1",
		"^d.*1$":       "class1",
		"description2": "class2",
	}}
}

func NewTestClassifierWithData() *Classifier {
	return &Classifier{map[string]string{
		"^pizza":   "Eating outside",
		"for rent": "Living",
	}}
}

func NewTestTagger() *Tagger {
	return &Tagger{map[string][]Tag{
		"class1": {"tag1", "tag2"},
		"class2": {"tag3"},
		"^c.*3$": {"tag4"},
	}}
}

func NewTestTaggerWithData() *Tagger {
	return &Tagger{map[string][]Tag{
		"Living": {Crucial},
	}}
}
