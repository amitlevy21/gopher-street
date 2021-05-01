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
			classifier:      Classifier{map[string]string{"description": "class"}},
			transactionDesc: "description",
			expectedClass:   "class",
		},
		{
			desc:            "ClassAccordingToRegexInDict",
			classifier:      Classifier{map[string]string{"^d.*n$": "class"}},
			transactionDesc: "description",
			expectedClass:   "class",
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
	c := NewTestClassifier(t, "classifier.yml")
	expected := map[string]string{"description1": "class1", "^d.*1$": "class1", "description2": "class2"}
	helpers.CheckEquals(t, c.descriptionToClass, expected)
}
