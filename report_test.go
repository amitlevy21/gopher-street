// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"testing"
	"time"

	helpers "github.com/amitlevy21/gopher-street/test"
)

func TestReportTable(t *testing.T) {
	testCases := []struct {
		desc     string
		expenses *Expenses
	}{
		{
			desc:     "Empty table for empty expenses",
			expenses: &Expenses{},
		},
		{
			desc: "Single row for single expense",
			expenses: &Expenses{
				{Date: time.April, Amount: 53.6, Class: "Food Outside", Tags: []Tag{Crucial}},
			},
		},
		{
			desc: "Multiple rows with total for many expense",
			expenses: &Expenses{
				{Date: time.April, Amount: 53.6, Class: "Food Outside", Tags: []Tag{Crucial}},
				{Date: time.April, Amount: 26.4, Class: "Food Outside", Tags: []Tag{Crucial}},
				{Date: time.April, Amount: 30.0, Class: "Food Outside", Tags: []Tag{Crucial}},
			},
		},
	}
	for _, tC := range testCases {
		r := Reporter{}
		t.Run(tC.desc, func(t *testing.T) {
			fixture := "reports"
			actual := r.Report(tC.expenses)
			helpers.CheckUpdateFlag(t, fixture, actual)
			helpers.ExpectEqualsGolden(t, fixture, actual)
		})
	}
}
