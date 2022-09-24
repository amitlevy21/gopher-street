// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

type Exporter struct{}

type ExpenseToExport struct {
	Date   time.Time
	Amount float64
	Class  string
	Tags   string
}

func (exporter *Exporter) ToXLSX(expenses []*Expense, filename string) error {
	expensesToReport := make([]*ExpenseToExport, len(expenses), len(expenses))

	for i, exp := range expenses {
		expensesToReport[i] = &ExpenseToExport{
			Date:   exp.Date,
			Amount: exp.Amount,
			Class:  exp.Class,
			Tags:   strings.Join(exp.Tags, ","),
		}
	}
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Expenses")
	sheet.AddRow().WriteSlice(&[]string{"Date", "Amount", "Class", "Tags"}, -1)
	for _, exp := range expensesToReport {
		sheet.AddRow().WriteStruct(exp, -1)
	}

	return file.Save(filename)
}

// func (*Exporter) colorCell(sheet *xlsx.Sheet) {
// 	cell := sheet.AddRow().AddCell()
// 	cell.SetFormula("SUM(B2:B6)")
// 	st := cell.GetStyle()
// 	st.Fill.PatternType = "solid"
// 	st.Font.Color = xlsx.RGB_Dark_Red
// 	st.Fill.FgColor = xlsx.RGB_Dark_Green
// 	st.ApplyAlignment = true
// 	st.ApplyFill = true
// }
