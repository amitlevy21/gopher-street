package main

import (
	"fmt"
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

func (c *Classifier) Class(description string) (string, error) {
	if class, ok := c.descriptionToClass[description]; ok {
		return class, nil
	}
	for regex, class := range c.descriptionToClass {
		if matched, _ := regexp.MatchString(regex, description); matched {
			return class, nil
		}
	}
	return description, fmt.Errorf("no class found for %s", description)
}
