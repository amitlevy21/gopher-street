// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"time"
)

type Expense struct {
	Date   time.Month
	Amount float64
	Class  string
	Tags   []Tag
}

type Expenses []Expense

func NewExpenses(transactions []Transaction, classifier *Classifier, tagger *Tagger) *Expenses {
	expenses := Expenses{}
	for _, tr := range transactions {
		class := classifier.Class(tr.Description)
		expense := Expense{
			Date:   tr.Date.Month(),
			Amount: tr.Credit,
			Class:  class,
			Tags:   tagger.Tags(class),
		}
		expenses = append(expenses, expense)
	}
	return &expenses
}

func (exps *Expenses) GroupByMonth() map[time.Month]Expenses {
	expenses := make(map[time.Month]Expenses)
	for _, exp := range *exps {
		expenses[exp.Date] = append(expenses[exp.Date], exp)
	}
	return expenses
}

func (exps *Expenses) GroupByClass() map[string]Expenses {
	expenses := make(map[string]Expenses)
	for _, exp := range *exps {
		expenses[exp.Class] = append(expenses[exp.Class], exp)
	}
	return expenses
}

func (exps *Expenses) GroupByTag() map[Tag]Expenses {
	expenses := make(map[Tag]Expenses)
	for _, exp := range *exps {
		if len(exp.Tags) == 0 {
			expenses[Tag("None")] = append(expenses[Tag("None")], exp)
		}
		for _, tag := range exp.Tags {
			expenses[tag] = append(expenses[tag], exp)
		}
	}
	return expenses
}
