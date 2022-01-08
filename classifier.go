// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"regexp"
)

type Classifier struct {
	descriptionToClass map[string]string
}

func NewClassifier(classesToDescriptions map[string][]string) *Classifier {
	descriptionToClass := make(map[string]string)
	for class, descriptions := range classesToDescriptions {
		for _, desc := range descriptions {
			descriptionToClass[desc] = class
		}
	}
	return &Classifier{descriptionToClass}
}

func (c *Classifier) Class(description string) string {
	if class, ok := c.descriptionToClass[description]; ok {
		return class
	}
	for regex, class := range c.descriptionToClass {
		if matched, _ := regexp.MatchString(regex, description); matched {
			return class
		}
	}
	return description
}
