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
