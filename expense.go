// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "time"

type Expense struct {
	Date   time.Month
	Amount float64
	Class  string
	Tags   []Tag
}

func NewExpense(transactions []*Transaction, classifier *Classifier) []Expense {
	expenses := make([]Expense, 0)
	for _, tr := range transactions {
		expense := Expense{
			Date:   tr.Date.Month(),
			Amount: tr.Credit,
			Class:  classifier.Class(tr.Description),
			Tags:   []Tag{},
		}
		expenses = append(expenses, expense)
	}
	return expenses
}
