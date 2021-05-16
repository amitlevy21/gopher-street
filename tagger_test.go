// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
)

func TestTag(t *testing.T) {
	testCases := []struct {
		desc         string
		tagger       Tagger
		class        string
		expectedTags []Tag
	}{
		{
			desc:         "EmptyTagWhenEmptyClass",
			expectedTags: []Tag{},
		},
		{
			desc:         "NoTagIfNoMatch",
			class:        "class",
			expectedTags: []Tag{},
		},
		{
			desc:         "TagAccordingToDict",
			tagger:       Tagger{map[string][]Tag{"class": {"tag"}}},
			class:        "class",
			expectedTags: []Tag{"tag"},
		},
		{
			desc:         "TagAccordingToRegexInDict",
			tagger:       Tagger{map[string][]Tag{"^c.*s$": {"tag"}}},
			class:        "class",
			expectedTags: []Tag{"tag"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tags := tC.tagger.Tags(tC.class)
			helpers.CheckEquals(t, tags, tC.expectedTags)
		})
	}
}
