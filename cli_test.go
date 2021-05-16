// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bytes"
	"path/filepath"
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
	"github.com/spf13/cobra"
)

func NewCmdBuffer(cmd *cobra.Command) *bytes.Buffer {
	b := new(bytes.Buffer)
	cmdLoad.SetOut(b)
	cmdLoad.SetErr(b)
	return b
}

func TestRootCmd(t *testing.T) {
	err := rootCmd.Execute()
	helpers.FailTestIfErr(t, err)
}

func TestLoadCmdMissingFilePath(t *testing.T) {
	err := cmdLoad.Execute()
	helpers.ExpectContains(t, err.Error(), "path missing")
}

func TestLoadCmdFilePathNotExist(t *testing.T) {
	cmdLoad.SetArgs([]string{"not_exist"})
	err := cmdLoad.Execute()
	helpers.ExpectContains(t, err.Error(), "No such file")
}

func TestLoadCmdPathNotFile(t *testing.T) {
	cmdLoad.SetArgs([]string{"test"})
	err := cmdLoad.Execute()
	helpers.ExpectContains(t, err.Error(), "is not file")
}

func TestLoadOutOfRangeMapper(t *testing.T) {
	err := CLIInit(filepath.Join(fixtures, "configs", "out-of-range.yml"))
	helpers.FailTestIfErr(t, err)
	b := NewCmdBuffer(rootCmd)
	rootCmd.SetArgs([]string{"load", filepath.Join(fixtures, "transactions", "data.csv")})
	err = rootCmd.Execute()
	helpers.FailTestIfErr(t, err)
	helpers.ExpectContains(t, b.String(), "out of range")
}

func TestLoadCmd(t *testing.T) {
	err := CLIInit(filepath.Join(fixtures, "configs", "config.yml"))
	helpers.FailTestIfErr(t, err)
	b := NewCmdBuffer(rootCmd)
	rootCmd.SetArgs([]string{"load", filepath.Join(fixtures, "transactions", "data.csv")})
	err = rootCmd.Execute()
	helpers.FailTestIfErr(t, err)
	helpers.ExpectContains(t, b.String(), "{April 3500 living [Crucial]}{March 20 shirt []}{May 5 eating outside []}")
}

func TestCLIInit(t *testing.T) {
	err := CLIInit(filepath.Join(fixtures, "configs", "config.yml"))
	helpers.FailTestIfErr(t, err)
}

func TestGetExpensesBadFile(t *testing.T) {
	_, err := getExpenses(&ConfigData{}, "not_exist.json")
	helpers.ExpectError(t, err)
}
