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

func NewTestTransaction(t *testing.T, description string) *Transaction {
	return &Transaction{
		Date:        helpers.UTCDate(t, 2021, 03, 18),
		Description: description,
		Credit:      5.0,
		Refund:      0.0,
		Balance:     150.0,
	}
}

func NewTestClassifier(t *testing.T) *Classifier {
	yaml, err := ioutil.ReadFile(filepath.Join("test", "fixtures", "classifier.yml"))
	helpers.FailTestIfErr(t, err)
	c, err := NewClassifier(yaml)
	helpers.FailTestIfErr(t, err)
	return c
}

func NewTestTagger(t *testing.T) *Tagger {
	yaml, err := ioutil.ReadFile(filepath.Join("test", "fixtures", "tagger.yml"))
	helpers.FailTestIfErr(t, err)
	tagger, err := NewTagger(yaml)
	helpers.FailTestIfErr(t, err)
	return tagger
}
