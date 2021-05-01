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

	return expenses
}
