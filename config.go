// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"github.com/spf13/viper"
)

type ConfigData struct {
	Files   Files
	Classes Classes
	Tags    Tags
}

type Files map[string]File

type File struct {
	Cards FileCards
}

type FileCards map[string]Card

type Card struct {
	RowSubsetter RowSubsetter
	ColMapper    ColMapper
	DateLayout   string
}

type RowSubsetter struct {
	Start int
	End   int
}

type ColMapper struct {
	Date        int
	Description int
	Credit      int `mapstructure:",omitempty"`
	Refund      int `mapstructure:",omitempty"`
	Balance     int `mapstructure:",omitempty"`
}

type Classes map[string][]string
type Tags map[string][]string

func initConfig(path string) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	return err
}

func GetConfigData() (*ConfigData, error) {
	configData := &ConfigData{}
	err := viper.Unmarshal(&configData)
	return configData, err
}
