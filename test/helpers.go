// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package test_helpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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
	if !cmp.Equal(actual, expected) {
		t.Errorf("expected %v, received %v\nDiff: %s", expected, actual, cmp.Diff(actual, expected))
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

func ExpectContains(t *testing.T, s, subString string) {
	if !strings.Contains(s, subString) {
		t.Errorf("expected %s to contain %s", s, subString)
	}
}
