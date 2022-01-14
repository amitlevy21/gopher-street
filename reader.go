package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/tealeg/xlsx"
)

type TransactionReader interface {
	Read(string) ([][]string, error)
}

type XLSXReader struct{}
type CSVReader struct{}

func ReaderFactory(fileExtension string) (TransactionReader, error) {
	switch fileExtension {
	case ".xlsx":
		return &XLSXReader{}, nil
	case ".csv":
		return &CSVReader{}, nil
	default:
		return nil, fmt.Errorf("unsupported file extension %s", fileExtension)
	}
}

func (r *XLSXReader) Read(filePath string) ([][]string, error) {
	s, err := xlsx.FileToSlice(filePath)
	return s[0], err
}

func (r *CSVReader) Read(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	// skip first line
	if _, err := reader.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
