package main

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Reporter struct {
	reportConf ReportConf
}

func (r *Reporter) Report(expenses *Expenses) string {
	return r.makeReportTable(expenses)
}

func totalAmount(expenses []*Expense) float64 {
	total := 0.0
	for _, e := range expenses {
		total += e.Amount
	}
	return total
}

func (r *Reporter) makeReportTable(expenses *Expenses) string {
	report := strings.Builder{}
	report.WriteString("Total report\n")
	report.WriteString(r.makeMainTable(expenses.Classified))
	if len(expenses.Unclassified) > 0 {
		report.WriteString("There were unclassified expenses, consider adding their classes to the classifier\n")
		unclassifiedReport := r.makeMainTable(expenses.Unclassified)
		report.WriteString(unclassifiedReport)
	}

	return report.String()
}

func (r *Reporter) makeMainTable(expenses []*Expense) string {
	t := table.NewWriter()
	if r.reportConf.RightToLeftLanguage {
		t.Style().Format.Direction = text.LeftToRight
	}
	header := r.getHeader()
	t.AppendHeader(header)
	appendTableBody(expenses, t)
	t.AppendSeparator()
	total := totalAmount(expenses)
	t.AppendFooter(table.Row{"", "Total", total})
	return t.Render() + "\n"
}

func (r *Reporter) getHeader() table.Row {
	customHeaders := r.reportConf.Headers
	if customHeaders != nil {
		return table.Row{"#", customHeaders.Date, customHeaders.Amount, customHeaders.Class, customHeaders.Tags}
	}
	return table.Row{"#", "Date", "Amount", "Class", "Tags"}
}

func appendTableBody(expenses []*Expense, t table.Writer) {
	for i, e := range expenses {
		dateWithoutTime := strings.Split(e.Date.String(), " ")[0]
		t.AppendRows([]table.Row{
			{i, dateWithoutTime, e.Amount, e.Class, e.Tags},
		})
	}
}
