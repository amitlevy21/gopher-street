// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"path/filepath"
	"testing"

	helpers "github.com/amitlevy21/gopher-street/test"
)

func TestConfigNotFound(t *testing.T) {
	err := initConfig("non_exist_config.yml")
	helpers.ExpectError(t, err)
}

func TestCLIInitNoOptionalVars(t *testing.T) {
	err := initConfig(filepath.Join(fixtures, "configs", "without-refund-balance.yml"))
	helpers.FailTestIfErr(t, err)
}

func TestCLIInitBadConfig(t *testing.T) {
	err := initConfig(filepath.Join(fixtures, "configs", "bad.yml"))
	helpers.ExpectError(t, err)
}

func TestConfig(t *testing.T) {
	err := initConfig("config.yml")
	helpers.FailTestIfErr(t, err)
	_, err = GetConfigData()
	helpers.FailTestIfErr(t, err)
}
