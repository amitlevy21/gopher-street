// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"testing"
)

func TestTagString(t *testing.T) {
	tags := []Tag{Recurring, Crucial}
	tagStrings := [...]string{"Recurring", "Crucial"}
	for i, tag := range tags {
		if fmt.Sprint(tag) != tagStrings[i] {
			t.Errorf("expected %s received %s", tagStrings[i], fmt.Sprint(Recurring))
		}
	}
}
