// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Expense struct {
	Date   time.Time
	Amount float64
	Class  string
	Tags   []Tag
}

type Expenses struct {
	Classified   []*Expense
	Unclassified []*Expense
}

func NewEmptyExpenses() *Expenses {
	return &Expenses{
		Classified:   []*Expense{},
		Unclassified: []*Expense{},
	}
}

func NewExpenses(transactions []Transaction, classifier *Classifier, tagger *Tagger) *Expenses {
	classified := []*Expense{}
	unclassified := []*Expense{}
	for _, tr := range transactions {
		class, error := classifier.Class(tr.Description)
		expense := Expense{
			Date:   tr.Date,
			Amount: tr.Credit,
			Class:  class,
			Tags:   tagger.Tags(class),
		}
		if error != nil {
			unclassified = append(unclassified, &expense)
			continue
		}
		classified = append(classified, &expense)
	}
	return &Expenses{classified, unclassified}
}

func getExpensesFromFile(conf *ConfigData, transactionFilePath string) (*Expenses, error) {
	expenses := NewEmptyExpenses()
	tagger := &Tagger{conf.Tags}
	cl := NewClassifier(conf.Classes)
	base := path.Base(transactionFilePath)
	noExt := strings.TrimSuffix(base, filepath.Ext(base))
	reader, err := ReaderFactory(filepath.Ext(base))
	if err != nil {
		return expenses, err
	}
	data, err := reader.Read(transactionFilePath)
	if err != nil {
		return expenses, err
	}
	transactionFileConfig, ok := conf.Files[strings.ToLower(noExt)]
	if !ok {
		return expenses, fmt.Errorf("no configuration found for transaction file: %s", noExt)
	}
	for _, card := range transactionFileConfig.Cards {
		cardTrans := NewCardTransactions(data, card.ColMapper, card.RowSubSetter, card.DateLayout)
		trans, err := cardTrans.Transactions()
		if err != nil {
			return expenses, err
		}
		cardExpenses := *NewExpenses(trans, cl, tagger)
		expenses.Classified = append(expenses.Classified, cardExpenses.Classified...)
		expenses.Unclassified = append(expenses.Unclassified, cardExpenses.Unclassified...)
	}

	return expenses, nil
}

func (exps *Expenses) ToSlice() []*Expense {
	expenses := []*Expense{}
	expenses = append(expenses, exps.Classified...)
	expenses = append(expenses, exps.Unclassified...)
	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date.Before(expenses[j].Date)
	})
	return expenses
}

func (exps *Expenses) GroupByMonth() map[time.Month][]Expense {
	expenses := make(map[time.Month][]Expense)
	for _, exp := range exps.ToSlice() {
		expenses[exp.Date.Month()] = append(expenses[exp.Date.Month()], *exp)
	}
	return expenses
}

func (exps *Expenses) GroupByClass() map[string][]Expense {
	expenses := make(map[string][]Expense)
	for _, exp := range exps.ToSlice() {
		expenses[exp.Class] = append(expenses[exp.Class], *exp)
	}
	return expenses
}

func (exps *Expenses) GroupByTag() map[Tag][]Expense {
	expenses := make(map[Tag][]Expense)
	for _, exp := range exps.ToSlice() {
		if len(exp.Tags) == 0 {
			expenses[Tag("None")] = append(expenses[Tag("None")], *exp)
		}
		for _, tag := range exp.Tags {
			expenses[tag] = append(expenses[tag], *exp)
		}
	}
	return expenses
}
