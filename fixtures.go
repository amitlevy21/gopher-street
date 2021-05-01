// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"io/ioutil"
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
	var mapper map[string]int = map[string]int{
		"date":        0,
		"description": 1,
		"credit":      4,
		"refund":      5,
		"balance":     6,
	}
	return NewCardTransactions(r, mapper, []int{})
}

func NewTestClassifier(t *testing.T, fileName string) *Classifier {
	yaml, err := ioutil.ReadFile(filepath.Join(fixtures, "classifiers", fileName))
	helpers.FailTestIfErr(t, err)
	c, err := NewClassifier(yaml)
	helpers.FailTestIfErr(t, err)
	return c
}

func NewTestTagger(t *testing.T, fileName string) *Tagger {
	yaml, err := ioutil.ReadFile(filepath.Join(fixtures, "taggers", fileName))
	helpers.FailTestIfErr(t, err)
	tagger, err := NewTagger(yaml)
	helpers.FailTestIfErr(t, err)
	return tagger
}
