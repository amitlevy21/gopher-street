// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
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

type CardTransactions struct {
	df           dataframe.DataFrame
	columnMapper map[string]int
	rowSubsetter []int
}

func NewCardTransactions(r io.Reader, columnMapper map[string]int, rowSubsetter []int) *CardTransactions {
	return &CardTransactions{
		dataframe.ReadCSV(r),
		columnMapper,
		rowSubsetter,
	}
}

func (t *CardTransactions) Transactions() ([]Transaction, error) {
	records, err := t.records()
	if err != nil {
		return []Transaction{}, err
	}
	transactions := make([]Transaction, 0, len(records))
	for _, record := range records {
		trans, err := transaction(record, t.columnMapper)
		if err == nil {
			transactions = append(transactions, *trans)
		}
	}
	return transactions, nil
}

func (t *CardTransactions) records() ([][]string, error) {
	if t.df.Err != nil {
		log.Printf("error while reading CSV: %s", t.df.Err)
		return [][]string{}, t.df.Err
	}
	if err := t.checkDims(); err != nil {
		return [][]string{}, err
	}
	if len(t.rowSubsetter) == 0 {
		return t.df.Records()[1:], nil
	}
	d := t.df.Subset(t.rowSubsetter)

	return d.Records()[1:], t.df.Err
}

func (t *CardTransactions) checkDims() error {
	rows, cols := t.df.Dims()
	if err := t.validateRowSubsetter(rows); err != nil {
		return err
	}
	if err := t.validateColumnMapper(cols); err != nil {
		return err
	}
	return nil
}

func (t *CardTransactions) validateRowSubsetter(rows int) error {
	min, max := minMax(t.rowSubsetter)
	if min < 0 || max > rows-1 {
		return errors.New("RowSubsetter indices out of range")
	}
	return nil
}

func minMax(s []int) (min int, max int) {
	if len(s) == 0 {
		return 0, 0
	}
	min = s[0]
	max = s[0]
	for _, e := range s[1:] {
		if e > max {
			max = e
		}
		if e < min {
			min = e
		}
	}
	return min, max
}

func (t *CardTransactions) validateColumnMapper(cols int) error {
	invalid := make(map[string]int)
	for k, v := range t.columnMapper {
		if v < 0 || v >= cols {
			invalid[k] = v
		}
	}
	if len(invalid) > 0 {
		err := fmt.Sprintf("invalid column mapper. invalid values: %v", invalid)
		log.Print(err)
		return errors.New(err)
	}
	return nil
}

func transaction(record []string, m map[string]int) (*Transaction, error) {
	time, err := parseTime(record[m["date"]])
	if err != nil {
		return &Transaction{}, err
	}
	balance, err := parseFloat(record[m["balance"]])
	if !validFloat(balance, err) {
		return &Transaction{}, errors.New("must have valid balance")
	}
	credit, err := parseFloat(record[m["credit"]])
	refund, err2 := parseFloat(record[m["refund"]])
	if !validFloat(credit, err) && !validFloat(refund, err2) {
		return &Transaction{}, errors.New("must have valid credit or refund")
	}
	description := record[m["description"]]
	return &Transaction{time, description, credit, refund, balance}, nil
}

func parseTime(timeStr string) (time.Time, error) {
	layout := "02.01.2006"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Printf("couldn't parse time: %s", err)
	}
	return t, err
}

func parseFloat(num string) (float64, error) {
	parsed, err := strconv.ParseFloat(num, 64)
	if err != nil {
		log.Printf("error when parsing float: %s", err)
	}
	return parsed, err
}

func validFloat(n float64, err error) bool {
	return err == nil && !math.IsNaN(n)
}
