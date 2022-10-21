package main

import (
	"fmt"
	"os"
)

func main() {
	LoadCLIFromConfig("config.yml")
	Execute()
}

func LoadCLIFromConfig(configPath string) {
	if err := CLIInit(configPath); err != nil {
		panic(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
