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

type CardTransactions struct {
	df           dataframe.DataFrame
	columnMapper map[string]int
	rowSubsetter []int
	dateLayout   string
}

func NewCardTransactions(r io.Reader, columnMapper map[string]int, rowSubsetter []int, dateLayout string) *CardTransactions {
	return &CardTransactions{
		dataframe.ReadCSV(r),
		columnMapper,
		rowSubsetter,
		dateLayout,
	}
}

func (t *CardTransactions) Transactions() ([]Transaction, error) {
	records, err := t.records()
	if err != nil {
		return []Transaction{}, err
	}
	transactions := make([]Transaction, 0, len(records))
	for _, record := range records {
		trans, err := t.transaction(record)
		if err != nil {
			log.Printf("Failed to create transaction: %s", err)
			continue
		}
		transactions = append(transactions, *trans)
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

func (t *CardTransactions) transaction(record []string) (*Transaction, error) {
	time, err := parseTime(record[t.columnMapper["date"]], t.dateLayout)
	if err != nil {
		return &Transaction{}, err
	}
	credit, refund, err := parseCreditAndRefund(t, record)
	if err != nil {
		return &Transaction{}, err
	}
	description := parseDescription(t, record)
	balance := t.parseBalance(record)

	return &Transaction{time, description, credit, refund, balance}, nil
}

func parseTime(timeStr string, layout string) (time.Time, error) {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Printf("couldn't parse time: %s", err)
	}
	return t, err
}

func parseCreditAndRefund(t *CardTransactions, record []string) (float64, float64, error) {
	credit := 0.0
	refund := 0.0
	creditIndex, creditOk := t.columnMapper["credit"]
	refundIndex, refundOk := t.columnMapper["refund"]
	if !creditOk && !refundOk {
		return 0, 0, errors.New("ColumnMapper missing both credit and refund")
	}
	creditStr := ""
	if creditOk {
		creditStr = record[creditIndex]
	}
	refundStr := ""
	if refundOk {
		refundStr = record[refundIndex]
	}
	hasCredit := isValidField(creditStr)
	hasRefund := isValidField(refundStr)
	if hasCredit && hasRefund || !hasCredit && !hasRefund {
		return 0, 0, errors.New("must define credit or refund but no both")
	}
	if hasCredit {
		credit, _ = strconv.ParseFloat(creditStr, 64)
	}
	if hasRefund {
		refund, _ = strconv.ParseFloat(refundStr, 64)
	}
	return credit, refund, nil
}

func parseDescription(t *CardTransactions, record []string) string {
	descriptionIndex, ok := t.columnMapper["description"]
	description := ""
	if ok {
		description = record[descriptionIndex]
	}
	return description
}

func isValidField(field string) bool {
	return field != "" && field != "NaN"
}

func (t *CardTransactions) parseBalance(record []string) float64 {
	balance := 0.0
	balanceIndex, ok := t.columnMapper["balance"]
	if ok {
		balanceStr := record[balanceIndex]
		hasBalance := isValidField(balanceStr)
		if hasBalance {
			balance, _ = strconv.ParseFloat(record[balanceIndex], 64)
		}
	}
	return balance
}
