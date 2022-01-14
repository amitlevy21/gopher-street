package main

import (
	"path/filepath"
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
)

func TestFactoryUnsupportedReader(t *testing.T) {
	_, err := ReaderFactory("non-supported")
	helpers.ExpectError(t, err)
}

func TestCSVReaderEmpty(t *testing.T) {
	r, err := ReaderFactory(".csv")
	helpers.FailTestIfErr(t, err)
	_, err = r.Read(filepath.Join(CSVTransactionsPath, "empty.csv"))
	helpers.ExpectError(t, err)
}

func TestCSVReaderNotExist(t *testing.T) {
	r, err := ReaderFactory(".csv")
	helpers.FailTestIfErr(t, err)
	_, err = r.Read(filepath.Join(CSVTransactionsPath, "not-exist.csv"))
	helpers.ExpectError(t, err)
}

func TestCSVReaderSingleLine(t *testing.T) {
	r, err := ReaderFactory(".csv")
	helpers.FailTestIfErr(t, err)
	_, err = r.Read(filepath.Join(CSVTransactionsPath, "bad-not-maching-num-commas.csv"))
	helpers.ExpectError(t, err)
}

func TestCSVReader(t *testing.T) {
	r, err := ReaderFactory(".csv")
	helpers.FailTestIfErr(t, err)
	_, err = r.Read(filepath.Join(CSVTransactionsPath, "data.csv"))
	helpers.FailTestIfErr(t, err)
}

func TestXLSXReader(t *testing.T) {
	r, err := ReaderFactory(".xlsx")
	helpers.FailTestIfErr(t, err)
	_, err = r.Read(filepath.Join(XLSXTransactionsPath, "empty_2_sheets.xlsx"))
	helpers.FailTestIfErr(t, err)
}
