// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
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

func getExpensesFromFile(conf *ConfigData, transactionFilePath string) (*Expenses, error) {
	expenses := Expenses{}
	tagger := &Tagger{conf.Tags}
	cl := NewClassifier(conf.Classes)
	base := path.Base(transactionFilePath)
	noExt := strings.TrimSuffix(base, filepath.Ext(base))
	reader, err := ReaderFactory(filepath.Ext(base))
	if err != nil {
		return &Expenses{}, err
	}
	data, err := reader.Read(transactionFilePath)
	if err != nil {
		return &Expenses{}, err
	}
	transactionFileConfig, ok := conf.Files[strings.ToLower(noExt)]
	if !ok {
		return &Expenses{}, fmt.Errorf("no configuration found for transaction file: %s", noExt)
	}
	for _, card := range transactionFileConfig.Cards {
		cardTrans := NewCardTransactions(data, card.ColMapper, card.RowSubsetter, card.DateLayout)
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
