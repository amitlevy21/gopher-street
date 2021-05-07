// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bytes"
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

func TestLoadCmd(t *testing.T) {
	b := NewCmdBuffer(cmdLoad)
	cmdLoad.SetArgs([]string{"cli_test.go"})
	err := cmdLoad.Execute()
	helpers.FailTestIfErr(t, err)
	helpers.ExpectContains(t, b.String(), "TODO")
}

func TestCLIInit(t *testing.T) {
	CLIInit()
	b := NewCmdBuffer(rootCmd)
	rootCmd.SetArgs([]string{"load", "cli_test.go"})
	err := rootCmd.Execute()
	helpers.FailTestIfErr(t, err)
	helpers.ExpectContains(t, b.String(), "TODO")
}
