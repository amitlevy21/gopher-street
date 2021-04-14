// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"io"
	"log"
	"strconv"
	"time"

	"github.com/go-gota/gota/dataframe"
)

type Transaction struct {
	Date        time.Time
	Description string
	Credit      float64
	Refund      float64
	Balance     float64
}

type ColumnMapper map[string]int

func TransactionsFromCSV(r io.Reader, m ColumnMapper) []Transaction {
	df := dataframe.ReadCSV(r)
	records := df.Records()[1:]
	transactions := make([]Transaction, len(records))
	for i, dfr := range records {
		transactions[i] = *transactionFromRecord(dfr, m)
	}
	return transactions
}

func transactionFromRecord(record []string, m ColumnMapper) *Transaction {
	time := parseTime(record[m["date"]])
	description := record[m["description"]]
	credit := parseFloat(record[m["credit"]])
	refund := parseFloat(record[m["refund"]])
	balance := parseFloat(record[m["balance"]])
	return &Transaction{time, description, credit, refund, balance}
}

func parseTime(timeStr string) time.Time {
	layout := "02.01.2006"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Printf("couldn't parse time: %s", err)
	}
	return t
}

func parseFloat(num string) float64 {
	credit, err := strconv.ParseFloat(num, 64)
	if err != nil {
		log.Printf("error when parsing float: %s", err)
	}
	return credit
}
