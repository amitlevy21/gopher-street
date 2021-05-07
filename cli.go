// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gstreet",
	Short: "",
	Long:  ``,
}
var cmdLoad = &cobra.Command{
	Use:   "load [file path]",
	Short: "Load data from file",
	Long:  ``,
	Args:  validateFilePath,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := cmd.OutOrStdout().Write([]byte("TODO"))
		log.Printf("couldn't write cmd output: %s", err)
	},
}

func CLIInit() {
	rootCmd.AddCommand(cmdLoad)
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
