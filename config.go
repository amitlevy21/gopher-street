package main

import (
	"github.com/spf13/viper"
)

type ConfigData struct {
	Files    Files
	Classes  Classes
	Tags     Tags
	Database Database
	Report   ReportConf
}

type Files map[string]File

type File struct {
	Cards FileCards
}

type FileCards map[string]Card

type Card struct {
	RowSubSetter *RowSubSetter
	ColMapper    *ColMapper
	DateLayout   string
}

type Classes map[string][]string
type Tags map[string][]string

type Database struct {
	URI string
}

type ReportConf struct {
	RightToLeftLanguage bool
	Headers             *struct {
		Date   string
		Amount string
		Class  string
		Tags   string
	}
}

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
