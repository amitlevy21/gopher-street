// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Expense struct {
	Date   time.Time
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
			Date:   tr.Date,
			Amount: tr.Credit,
			Class:  class,
			Tags:   tagger.Tags(class),
		}
		expenses = append(expenses, expense)
	}
	return &expenses
}

func getExpensesFromFile(conf *ConfigData, transactionFileName string) (*Expenses, error) {
	expenses := Expenses{}
	tagger := &Tagger{conf.Tags}
	cl := NewClassifier(conf.Classes)
	r, err := os.Open(transactionFileName)
	if err != nil {
		return &Expenses{}, err
	}
	base := path.Base(transactionFileName)
	noExt := strings.TrimSuffix(base, filepath.Ext(base))
	for _, card := range conf.Files[noExt].Cards {
		cm := makeColMapper(card)
		rs := makeRowSubsetter(card)
		cardTrans := NewCardTransactions(r, cm, rs, card.DateLayout)
		trans, err := cardTrans.Transactions()
		if err != nil {
			return &Expenses{}, err
		}
		expenses = append(expenses, *NewExpenses(trans, cl, tagger)...)
	}

	return &expenses, nil
}

func makeRowSubsetter(card Card) []int {
	low := card.RowSubsetter.Start
	high := card.RowSubsetter.End
	bound := int(math.Abs(float64(low) - float64(high)))
	rs := make([]int, bound)
	for i, j := low, 0; i < high; i, j = i+1, j+1 {
		rs[j] = i
	}
	return rs
}

func makeColMapper(card Card) map[string]int {
	m := map[string]int{
		"date":        card.ColMapper.Date,
		"description": card.ColMapper.Description,
	}
	if card.ColMapper.Credit != 0 {
		m["credit"] = card.ColMapper.Credit
	}
	if card.ColMapper.Refund != 0 {
		m["refund"] = card.ColMapper.Refund
	}
	if card.ColMapper.Balance != 0 {
		m["balance"] = card.ColMapper.Balance
	}
	return m
}

func (exps *Expenses) GroupByMonth() map[time.Month]Expenses {
	expenses := make(map[time.Month]Expenses)
	for _, exp := range *exps {
		expenses[exp.Date.Month()] = append(expenses[exp.Date.Month()], exp)
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
