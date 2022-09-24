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
		expectError     bool
	}{
		{
			desc:        "EmptyClassWhenEmptyDescription",
			expectError: true,
		},
		{
			desc:            "DescriptionAsClassIfNoMatch",
			transactionDesc: "description",
			expectedClass:   "description",
			expectError:     true,
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
			class, error := tC.classifier.Class(tC.transactionDesc)
			if class != tC.expectedClass {
				t.Errorf("expected class %s received %s", tC.expectedClass, class)
			}
			if error != nil && !tC.expectError {
				helpers.FailTestIfErr(t, error)
			}
		})
	}
}

func TestNewClassifier(t *testing.T) {
	classesToDescriptions := map[string][]string{
		"Eating outside": {"^pizza"},
		"Living":         {"for rent"},
	}
	descriptionToClasses := map[string]string{
		"^pizza":   "Eating outside",
		"for rent": "Living",
	}
	cl := NewClassifier(classesToDescriptions)
	helpers.ExpectEquals(t, cl.descriptionToClass, descriptionToClasses)
}
