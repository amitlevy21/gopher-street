// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
)

type Transaction struct {
	Date        time.Time
	Description string
	Credit      float64
	Refund      float64
	Balance     float64
}

type CardTransactions struct {
	data         [][]string
	columnMapper *ColMapper
	rowSubsetter *RowSubsetter
	dateLayout   string
}

type ColMapper struct {
	Date        uint32
	Description uint32
	Credit      uint32 `mapstructure:",omitempty"`
	Refund      uint32 `mapstructure:",omitempty"`
	Balance     uint32 `mapstructure:",omitempty"`
}

type RowSubsetter struct {
	Start uint32
	End   uint32
}

func NewCardTransactions(data [][]string, columnMapper *ColMapper, rowSubsetter *RowSubsetter, dateLayout string) *CardTransactions {
	return &CardTransactions{
		data,
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
	if err := t.checkDims(); err != nil {
		return [][]string{}, err
	}
	if t.rowSubsetter.Start == t.rowSubsetter.End {
		return t.data, nil
	}

	return t.data[t.rowSubsetter.Start:t.rowSubsetter.End], nil
}

func (t *CardTransactions) checkDims() error {
	rows := len(t.data)
	if err := t.validateRowSubsetter(uint32(rows)); err != nil {
		return err
	}
	if rows > 0 {
		cols := len(t.data[t.rowSubsetter.Start])
		if err := t.validateColumnMapper(uint32(cols)); err != nil {
			return err
		}
	}
	return nil
}

func (t *CardTransactions) validateRowSubsetter(rows uint32) error {
	if t.rowSubsetter.Start == t.rowSubsetter.End {
		return nil
	}

	if t.rowSubsetter.End > rows {
		return errors.New("RowSubsetter indices out of range")
	}
	return nil
}

func (t *CardTransactions) validateColumnMapper(cols uint32) error {
	ref := reflect.ValueOf(*t.columnMapper)
	invalid := make(map[string]uint32)
	for i := 0; i < ref.NumField(); i++ {
		value := ref.Field(i).Interface().(uint32)
		if value > cols {
			invalid[ref.Type().Field(i).Name] = value
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
	time, err := parseTime(record[t.columnMapper.Date], t.dateLayout)
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
	creditIndex := t.columnMapper.Credit
	refundIndex := t.columnMapper.Refund
	if creditIndex == refundIndex {
		return 0, 0, errors.New("ColumnMapper has overlapping indexes for credit and refund")
	}
	creditStr := record[creditIndex]
	refundStr := record[refundIndex]
	hasCredit := isValidField(creditStr, creditIndex)
	hasRefund := isValidField(refundStr, refundIndex)
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
	descriptionIndex := t.columnMapper.Description
	return record[descriptionIndex]
}

func isValidField(field string, index uint32) bool {
	return field != "" && field != "NaN" && index > 0
}

func (t *CardTransactions) parseBalance(record []string) float64 {
	balance := 0.0
	balanceIndex := t.columnMapper.Balance
	balanceStr := record[balanceIndex]
	hasBalance := isValidField(balanceStr, balanceIndex)
	if hasBalance {
		balance, _ = strconv.ParseFloat(record[balanceIndex], 64)
	}

	return balance
}
