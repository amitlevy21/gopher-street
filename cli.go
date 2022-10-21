package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gst",
	Short: "",
	Long:  ``,
}
var loadcmd = &cobra.Command{
	Use:   "load [file path]",
	Short: "Load data from file",
	Long:  ``,
	Args:  validateFilePath,
	Run:   loadFile,
}
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get expenses from DB",
	Long:  ``,
	Run:   getExpensesFromDB,
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

func CLIInit(configPath string) error {
	err := initConfig(configPath)
	rootCmd.AddCommand(loadcmd)
	rootCmd.AddCommand(getCmd)
	return err
}

func loadFile(cmd *cobra.Command, args []string) {
	conf, err := GetConfigData()
	writeErrCmd(cmd, err)
	expenses := getExpenses(conf, args, cmd)
	if len(expenses.ToSlice()) > 0 {
		writeExpensesToDB(conf, expenses, cmd)
	}
	writeCmd(cmd.OutOrStdout(), "\n")
	writeReport(cmd, expenses)
	writeCmd(cmd.OutOrStdout(), "\n\nDone!\n")
}

func getExpenses(conf *ConfigData, args []string, cmd *cobra.Command) *Expenses {
	expenses, err := getExpensesFromFile(conf, args[0])
	writeErrCmd(cmd, err)
	return expenses
}

func writeExpensesToDB(conf *ConfigData, expenses *Expenses, cmd *cobra.Command) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := Instance(ctx, conf.Database.URI)
	defer db.closeDB(ctx)
	writeErrCmd(cmd, db.WriteExpenses(ctx, expenses))
}

func writeReport(cmd *cobra.Command, expenses *Expenses) {
	r := Reporter{}
	writeCmd(cmd.OutOrStdout(), r.Report(expenses))
}

func getExpensesFromDB(cmd *cobra.Command, args []string) {
	conf, err := GetConfigData()
	writeErrCmd(cmd, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := Instance(ctx, conf.Database.URI)
	defer db.closeDB(ctx)
	expenses, err := db.GetExpenses(ctx)
	writeErrCmd(cmd, err)
	writeReport(cmd, expenses)
	writeCmd(cmd.OutOrStdout(), fmt.Sprintf("\n%d expenses found\n", len(expenses.ToSlice())))
}

func writeErrCmd(cmd *cobra.Command, err error) {
	if err != nil {
		writeCmd(cmd.ErrOrStderr(), err.Error())
	}
}

func writeCmd(writer io.Writer, msg string) {
	if _, err := writer.Write([]byte(msg)); err != nil {
		log.Printf("error while writing to cmd, %s\n", err)
	}
}
