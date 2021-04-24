// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"regexp"

	"gopkg.in/yaml.v2"
)

type Classifier struct {
	classes map[string]string
}

func NewClassifier(yamlParseable []byte) (*Classifier, error) {
	classesToDesc := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(yamlParseable, &classesToDesc); err != nil {
		return &Classifier{}, err
	}

	classes := make(map[string]string)
	for _, v := range classesToDesc {
		for class, descriptions := range v.(map[interface{}]interface{}) {
			className := class.(string)
			for _, desc := range descriptions.([]interface{}) {
				description := desc.(string)
				classes[description] = className
			}
		}
	}
	return &Classifier{classes}, nil
}

func (c *Classifier) Class(description string) string {
	if class, ok := c.classes[description]; ok {
		return class
	}
	for regex, class := range c.classes {
		if matched, _ := regexp.MatchString(string(regex), string(description)); matched {
			return class
		}
	}
	return description
}
