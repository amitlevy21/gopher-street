// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
	"github.com/spf13/cobra"
)

func NewCmdBuffer(cmd *cobra.Command) *bytes.Buffer {
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	return b
}

func TestRootCmd(t *testing.T) {
	err := rootCmd.Execute()
	helpers.FailTestIfErr(t, err)
}

func TestLoadCmdMissingFilePath(t *testing.T) {
	err := loadcmd.Execute()
	helpers.ExpectContains(t, err.Error(), "path missing")
}

func TestLoadCmdFilePathNotExist(t *testing.T) {
	loadcmd.SetArgs([]string{"not_exist"})
	err := loadcmd.Execute()
	helpers.ExpectContains(t, err.Error(), "No such file")
}

func TestLoadCmdPathNotFile(t *testing.T) {
	loadcmd.SetArgs([]string{"test"})
	err := loadcmd.Execute()
	helpers.ExpectContains(t, err.Error(), "is not file")
}

func TestLoadOutOfRangeMapper(t *testing.T) {
	err := CLIInit(filepath.Join(fixtures, "configs", "out-of-range.yml"))
	helpers.FailTestIfErr(t, err)
	b := NewCmdBuffer(rootCmd)
	rootCmd.SetArgs([]string{"load", filepath.Join(CSVTransactionsPath, "data.csv")})
	err = rootCmd.Execute()
	helpers.FailTestIfErr(t, err)
	helpers.ExpectContains(t, b.String(), "out of range")
}

func TestLoadCmd(t *testing.T) {
	err := CLIInit(filepath.Join(fixtures, "configs", "config.yml"))
	helpers.FailTestIfErr(t, err)
	b := NewCmdBuffer(rootCmd)
	rootCmd.SetArgs([]string{"load", filepath.Join(CSVTransactionsPath, "data.csv")})
	err = rootCmd.Execute()
	helpers.FailTestIfErr(t, err)
	helpers.ExpectContains(t, b.String(), "eating outside")
}

func TestGetEmptyExpenses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := Instance(ctx)
	helpers.FailTestIfErr(t, db.dropDB(ctx))
	err := CLIInit(filepath.Join(fixtures, "configs", "config.yml"))
	helpers.FailTestIfErr(t, err)
	b := NewCmdBuffer(rootCmd)
	rootCmd.SetArgs([]string{"get"})
	err = rootCmd.Execute()
	helpers.FailTestIfErr(t, err)
	helpers.ExpectContains(t, b.String(), "0 expenses")
}

func TestGetExpenses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := Instance(ctx)
	defer helpers.FailTestIfErr(t, db.dropDB(ctx))
	helpers.FailTestIfErr(t, db.dropDB(ctx))
	err := CLIInit(filepath.Join(fixtures, "configs", "config.yml"))
	helpers.FailTestIfErr(t, err)
	rootCmd.SetArgs([]string{"load", filepath.Join(CSVTransactionsPath, "data.csv")})
	err = rootCmd.Execute()
	helpers.FailTestIfErr(t, err)
	b := NewCmdBuffer(rootCmd)
	rootCmd.SetArgs([]string{"get"})
	err = rootCmd.Execute()
	helpers.FailTestIfErr(t, err)
	helpers.ExpectContains(t, b.String(), "3 expenses")
}

func TestCLIInit(t *testing.T) {
	err := CLIInit(filepath.Join(fixtures, "configs", "config.yml"))
	helpers.FailTestIfErr(t, err)
}

func TestBadCMDWriter(t *testing.T) {
	writeCmd(&helpers.BadWriter{}, "")
}
