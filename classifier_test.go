// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
)

func TestClass(t *testing.T) {
	testCases := []struct {
		desc            string
		classifier      Classifier
		transactionDesc string
		expectedClass   string
	}{
		{
			desc: "EmptyClassWhenEmptyDescription",
		},
		{
			desc:            "DescriptionAsClassIfNoMatch",
			transactionDesc: "description",
			expectedClass:   "description",
		},
		{
			desc:            "ClassAccordingToDict",
			classifier:      Classifier{map[string]string{"hello": "world"}},
			transactionDesc: "hello",
			expectedClass:   "world",
		},
		{
			desc:            "ClassAccordingToRegexInDict",
			classifier:      Classifier{map[string]string{"^h.*o$": "world"}},
			transactionDesc: "hello",
			expectedClass:   "world",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if class := tC.classifier.Class(tC.transactionDesc); class != tC.expectedClass {
				t.Errorf("expected class %s received %s", tC.expectedClass, class)
			}
		})
	}
}

func TestNewClassifierBadFile(t *testing.T) {
	badYAML := []byte("invalid YAML")
	_, err := NewClassifier(badYAML)
	helpers.ExpectError(t, err)
}

func TestNewClassifier(t *testing.T) {
	c := NewTestClassifier(t)
	expected := map[string]string{"hello": "world", "^h.*o$": "world", "hi": "test"}
	helpers.CheckEquals(t, c.classes, expected)
}
