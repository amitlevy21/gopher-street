// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "time"

type Expense struct {
	Date   time.Time
	Amount float64
	Class  string
	Tags   []Tag
}

type Tag int

const (
	Recurring Tag = iota
	Crucial
)

func (t Tag) String() string {
	return [...]string{"Recurring", "Crucial"}[t]
}
