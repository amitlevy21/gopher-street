// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package test_helpers

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var fixtures string = filepath.Join("test", "fixtures")

func OpenFixture(t *testing.T, fixtureFileName string) *os.File {
	r, err := os.Open(filepath.Join(fixtures, fixtureFileName))
	if err != nil {
		t.Fatalf("err while opening fixture file: %s", err)
	}
	return r
}

func UTCDate(t *testing.T, year int, month time.Month, day int) time.Time {
	timeZone, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatalf("err while loading location: %s", err)
	}
	return time.Date(year, month, day, 0, 0, 0, 0, timeZone)
}

func CheckEquals(t *testing.T, actual interface{}, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, received %v", expected, actual)
	}
}

func ExpectError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("expected to throw error, received nil")
	}
}

func FailTestIfErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error %s while running test %s", err, t.Name())
	}
}