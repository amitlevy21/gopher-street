package main

import (
	"fmt"
	"os"
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
)

func TestExporterXLSX(t *testing.T) {
	exporter := &Exporter{}
	testCases := []struct {
		desc     string
		expenses []*Expense
		expected [][]string
	}{
		{
			desc:     "EmptyExpensesToEmptyFile",
			expenses: []*Expense{},
			expected: [][]string{{"Date", "Amount", "Class", "Tags"}},
		},
		{
			desc:     "SingleExpense",
			expenses: []*Expense{NewTestExpense(t)},
			expected: [][]string{
				{"Date", "Amount", "Class", "Tags"}, {"3/18/21 00:00", "5", "class1", "tag1"},
			},
		},
		{
			desc:     "ManyExpenses",
			expenses: NewTestExpenses(t).Classified,
			expected: [][]string{
				{"Date", "Amount", "Class", "Tags"}, {"3/18/21 00:00", "5", "class1", "tag1"},
				{"4/19/21 00:00", "5", "class1", "tag2"},
				{"5/20/21 00:00", "5", "class2", "tag1"},
			},
		},
		{
			desc: "ExpenseWithManyTags",
			expenses: []*Expense{{
				Date:   helpers.UTCDate(t, 2022, 9, 24),
				Amount: 5,
				Class:  "class",
				Tags:   []Tag{"tag1", "tag2", "tag3"},
			}},
			expected: [][]string{
				{"Date", "Amount", "Class", "Tags"},
				{"9/24/22 00:00", "5", "class", "tag1,tag2,tag3"},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			filename := fmt.Sprintf("%s.xlsx", tC.desc)
			err := exporter.ToXLSX(tC.expenses, filename)
			helpers.FailTestIfErr(t, err)
			reader := &XLSXReader{}
			content, err := reader.Read(filename)
			helpers.FailTestIfErr(t, err)
			helpers.ExpectEquals(t, content, tC.expected)
			os.Remove(filename)
		})
	}
}
