// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Reporter struct {
}

func (r *Reporter) Report(expenses *Expenses) string {
	total := totalAmount(expenses)
	return makeReportTable(expenses, total)
}

func totalAmount(expenses *Expenses) float64 {
	total := 0.0
	for _, e := range *expenses {
		total += e.Amount
	}
	return total
}

func makeReportTable(expenses *Expenses, total float64) string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "Date", "Amount", "Class", "Tags"})
	appendTableBody(expenses, t)
	t.AppendSeparator()
	t.AppendFooter(table.Row{"", "Total", total})
	return t.Render()
}

func appendTableBody(expenses *Expenses, t table.Writer) {
	for i, e := range *expenses {
		dateWithoutTime := strings.Split(e.Date.String(), " ")[0]
		t.AppendRows([]table.Row{
			{i, dateWithoutTime, e.Amount, e.Class, e.Tags},
		})
	}
}
