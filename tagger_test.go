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
		desc            string
		tagger          Tagger
		transactionDesc string
		expectedTag     Tag
	}{
		{
			desc: "EmptyTagWhenEmptyDescription",
		},
		{
			desc:            "NoTagIfNoMatch",
			transactionDesc: "description",
			expectedTag:     "",
		},
		{
			desc:            "TagAccordingToDict",
			tagger:          Tagger{map[string]Tag{"hello": "world"}},
			transactionDesc: "hello",
			expectedTag:     "world",
		},
		{
			desc:            "TagAccordingToRegexInDict",
			tagger:          Tagger{map[string]Tag{"^h.*o$": "world"}},
			transactionDesc: "hello",
			expectedTag:     "world",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tag := tC.tagger.Tag(tC.transactionDesc); tag != tC.expectedTag {
				t.Errorf("expected tag %s received %s", tC.expectedTag, tag)
			}
		})
	}
}

func TestNewTaggerBadFile(t *testing.T) {
	badYAML := []byte("invalid YAML")
	_, err := NewTagger(badYAML)
	helpers.ExpectError(t, err)
}

func TestNewTagger(t *testing.T) {
	tagger := NewTestTagger(t)
	expected := map[string]Tag{"hello": "world", "^h.*o$": "world", "hi": "test"}
	helpers.CheckEquals(t, tagger.tags, expected)
}
