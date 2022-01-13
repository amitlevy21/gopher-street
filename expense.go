// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"encoding/csv"
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
	r, err := readCSV(transactionFileName)
	if err != nil {
		return &Expenses{}, err
	}
	base := path.Base(transactionFileName)
	noExt := strings.TrimSuffix(base, filepath.Ext(base))
	for _, card := range conf.Files[noExt].Cards {
		cardTrans := NewCardTransactions(r, card.ColMapper, card.RowSubsetter, card.DateLayout)
		trans, err := cardTrans.Transactions()
		if err != nil {
			return &Expenses{}, err
		}
		expenses = append(expenses, *NewExpenses(trans, cl, tagger)...)
	}

	return &expenses, nil
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

func readCSV(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	// skip first line
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
