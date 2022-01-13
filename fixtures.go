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

var fixtures = filepath.Join("test", "fixtures")

func NewTestConfig() *ConfigData {
	return &ConfigData{
		Files: map[string]File{
			"multiple-rows": {
				Cards: map[string]Card{
					"card1": {
						RowSubsetter: &RowSubsetter{
							Start: 1,
							End:   4,
						},
						ColMapper: &ColMapper{
							Date:        0,
							Description: 1,
							Credit:      4,
							// Refund:      5,
							// Balance:     6,
						},
						DateLayout: "02.01.2006",
					},
				},
			},
		},
	}
}

func NewTestExpense(t *testing.T) *Expense {
	return &Expense{
		Date:   helpers.UTCDate(t, 2021, 03, 18),
		Amount: 5.0,
		Class:  "class1",
		Tags:   []Tag{"tag1"},
	}
}

func NewTestExpenses(t *testing.T) *Expenses {
	return &Expenses{
		{
			Date:   helpers.UTCDate(t, 2021, 03, 18),
			Amount: 5.0,
			Class:  "class1",
			Tags:   []Tag{"tag1"},
		},
		{
			Date:   helpers.UTCDate(t, 2021, 04, 19),
			Amount: 5.0,
			Class:  "class1",
			Tags:   []Tag{"tag2"},
		},
		{
			Date:   helpers.UTCDate(t, 2021, 05, 20),
			Amount: 5.0,
			Class:  "class2",
			Tags:   []Tag{"tag1"},
		},
	}
}

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
	data := helpers.ReadCSVFixture(t, filepath.Join("transactions", fileName))
	mapper := &ColMapper{
		Date:        0,
		Description: 1,
		Credit:      4,
		Refund:      5,
		Balance:     6,
	}
	layout := "02.01.2006"
	return NewCardTransactions(data, mapper, &RowSubsetter{}, layout)
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
