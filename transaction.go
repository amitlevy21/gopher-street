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

type CSVHunk struct {
	df           dataframe.DataFrame
	columnMapper map[string]int
}

func NewCSVHunk(r io.Reader, columnMapper map[string]int) *CSVHunk {
	return &CSVHunk{
		dataframe.ReadCSV(r),
		columnMapper,
	}
}

func (t *CSVHunk) Transactions() ([]Transaction, error) {
	records, err := t.records()
	if err != nil {
		return []Transaction{}, err
	}
	transactions := make([]Transaction, 0, len(records))
	for _, dfr := range records {
		trans, err := transaction(dfr, t.columnMapper)
		if err == nil {
			transactions = append(transactions, *trans)
		}
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

	return t.df.Records()[1:], t.df.Err
}

func (t *CSVHunk) checkDims() error {
	_, cols := t.df.Dims()
	if err := validateColumnMapper(t, cols); err != nil {
		return err
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
