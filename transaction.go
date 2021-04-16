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

type CSVHunk struct {
	df           dataframe.DataFrame
	columnMapper map[string]int
	rowSubsetter []int
}

func NewCSVHunk(r io.Reader, columnMapper map[string]int, rowSubsetter []int) *CSVHunk {
	return &CSVHunk{
		dataframe.ReadCSV(r),
		columnMapper,
		rowSubsetter,
	}
}

func (t *CSVHunk) Transactions() ([]Transaction, error) {
	records, err := t.records()
	if err != nil {
		return []Transaction{}, err
	}
	transactions := make([]Transaction, len(records))
	for i, dfr := range records {
		transactions[i] = *transaction(dfr, t.columnMapper)
	}
	return transactions, nil
}

func (t *CSVHunk) records() ([][]string, error) {
	if t.df.Err != nil {
		log.Printf("error while reading CSV: %s", t.df.Err)
		return [][]string{}, t.df.Err
	}
	if err := t.checkDims(); err != nil {
		return [][]string{}, err
	}
	if len(t.rowSubsetter) > 0 {
		t.df = t.df.Subset(t.rowSubsetter)
	}

	return t.df.Records()[1:], t.df.Err
}

func (t *CSVHunk) checkDims() error {
	rows, cols := t.df.Dims()
	if err := validateColumnMapper(t, cols); err != nil {
		return err
	}
	if err := validateRowSubsetter(t, rows); err != nil {
		return err
	}
	return nil
}

func validateRowSubsetter(t *CSVHunk, rows int) error {
	min, max := minMax(t.rowSubsetter)
	if min < 0 || max >= rows {
		err := fmt.Sprintf("invalid row subsetter. indices out of range: rows=%d, min=%d, max=%d", rows, min, max)
		log.Print(err)
		return errors.New(err)
	}
	return nil
}

func validateColumnMapper(t *CSVHunk, cols int) error {
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

func minMax(s []int) (int, int) {
	if len(s) == 0 {
		return 0, 0
	}
	currentMax := 0
	currentMin := s[0]
	for _, e := range s {
		if e > currentMax {
			currentMax = e
		}
		if e < currentMin {
			currentMin = e
		}
	}
	return currentMin, currentMax
}

func transaction(record []string, m map[string]int) *Transaction {
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
	parsed, err := strconv.ParseFloat(num, 64)
	if err != nil {
		log.Printf("error when parsing float: %s", err)
	}
	return parsed
}
