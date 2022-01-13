// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package testhelpers

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type BadWriter struct{}

var Fixtures = filepath.Join("test", "fixtures")
var update = flag.Bool("update", false, "update .golden files")

func OpenFixture(t *testing.T, fixtureFileName string) *os.File {
	r, err := os.Open(filepath.Join(Fixtures, fixtureFileName))
	if err != nil {
		t.Fatalf("err while opening fixture file: %s", err)
	}
	return r
}

func ReadCSVFixture(t *testing.T, fixtureFileName string) [][]string {
	f := OpenFixture(t, fixtureFileName)
	defer f.Close()

	r := csv.NewReader(f)
	// skip first line
	if _, err := r.Read(); err != nil {
		t.Fatalf("failed to read first line of file: %s", err)
	}

	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("failed to read csv file: %s", err)
	}

	return records
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

func CheckUpdateFlag(t *testing.T, fixture string, actual string) {
	goldenPath := getGoldenPath(t, fixture)
	if *update {
		updateGoldenFile(t, goldenPath, actual)
	}
}

func ExpectEqualsGolden(t *testing.T, fixture string, actual string) {
	gp := getGoldenPath(t, fixture)
	g, err := ioutil.ReadFile(gp)
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	t.Log(string([]byte(actual)))
	if !bytes.Equal([]byte(actual), g) {
		t.Errorf("bytes do not match .golden file")
	}
}

func getGoldenPath(t *testing.T, fixture string) string {
	return filepath.Join(Fixtures, fixture, t.Name()+".golden")
}

func updateGoldenFile(t *testing.T, goldenPath string, actual string) {
	t.Log("updating golden file")
	if err := ioutil.WriteFile(goldenPath, []byte(actual), 0644); err != nil {
		t.Fatalf("failed to update golden file: %s", err)
	}
}

func (w *BadWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("i always get it wrong")
}
