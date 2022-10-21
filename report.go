package main

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Reporter struct {
}

func (r *Reporter) Report(expenses *Expenses) string {
	return makeReportTable(expenses)
}

func totalAmount(expenses []*Expense) float64 {
	total := 0.0
	for _, e := range expenses {
		total += e.Amount
	}
	return total
}

func makeReportTable(expenses *Expenses) string {
	report := strings.Builder{}
	report.WriteString("Total report\n")
	report.WriteString(makeMainTable(expenses.ToSlice()))
	if len(expenses.Unclassified) > 0 {
		report.WriteString("There were unclassified expenses, consider adding their classes to the classifier\n")
		unclassifiedReport := makeMainTable(expenses.Unclassified)
		report.WriteString(unclassifiedReport)
	}

	return report.String()
}

func makeMainTable(expenses []*Expense) string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "Date", "Amount", "Class", "Tags"})
	appendTableBody(expenses, t)
	t.AppendSeparator()
	total := totalAmount(expenses)
	t.AppendFooter(table.Row{"", "Total", total})
	return t.Render() + "\n"
}

func appendTableBody(expenses []*Expense, t table.Writer) {
	for i, e := range expenses {
		dateWithoutTime := strings.Split(e.Date.String(), " ")[0]
		t.AppendRows([]table.Row{
			{i, dateWithoutTime, e.Amount, e.Class, e.Tags},
		})
	}
}
