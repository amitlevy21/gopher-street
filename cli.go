// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gst",
	Short: "",
	Long:  ``,
}
var cmdLoad = &cobra.Command{
	Use:   "load [file path]",
	Short: "Load data from file",
	Long:  ``,
	Args:  validateFilePath,
	Run:   loadFile,
}

func CLIInit(configPath string) error {
	err := initConfig(configPath)
	rootCmd.AddCommand(cmdLoad)
	return err
}

func loadFile(cmd *cobra.Command, args []string) {
	r, _ := os.Open(args[0])
	conf, _ := GetConfigData()
	for _, file := range conf.Files {
		cards := file.Cards
		for _, card := range cards {
			cm := makeColMapper(card)
			rs := makeRowSubsetter(card)
			cardTrans := NewCardTransactions(r, cm, rs, card.DateLayout)
			_, err := cardTrans.Transactions()
			if err != nil {
				cmd.OutOrStderr().Write([]byte(err.Error()))
			}
		}
	}
	cmd.OutOrStdout().Write([]byte("Done!"))
}

func validateFilePath(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("File path missing")
	}
	stat, err := os.Stat(args[0])
	if os.IsNotExist(err) {
		return fmt.Errorf("No such file: %s", args[0])
	}
	if !stat.Mode().IsRegular() {
		return fmt.Errorf("Path is not file: %s", args[0])
	}
	return nil
}

func makeRowSubsetter(card Card) []int {
	low := card.RowSubsetter.Start
	high := card.RowSubsetter.End
	bound := int(math.Abs(float64(low) - float64(high)))
	rs := make([]int, bound)
	for i, j := low, 0; i < high; i, j = i+1, j+1 {
		rs[j] = i
	}
	return rs
}

func makeColMapper(card Card) map[string]int {
	m := map[string]int{
		"date":        card.ColMapper.Date,
		"description": card.ColMapper.Description,
	}
	if card.ColMapper.Credit != 0 {
		m["credit"] = card.ColMapper.Credit
	}
	if card.ColMapper.Refund != 0 {
		m["refund"] = card.ColMapper.Refund
	}
	if card.ColMapper.Balance != 0 {
		m["balance"] = card.ColMapper.Balance
	}
	return m
}
