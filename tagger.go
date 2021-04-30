// Copyright (c) 2021 Amit Levy
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"regexp"

	"gopkg.in/yaml.v2"
)

type Tag string

const (
	Recurring Tag = "Recurring"
	Crucial   Tag = "Crucial"
	None      Tag = ""
)

type Tagger struct {
	classesToTags map[string][]Tag
}

func NewTagger(yamlParseable []byte) (*Tagger, error) {
	tagsToClasses := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(yamlParseable, &tagsToClasses); err != nil {
		return &Tagger{}, err
	}

	tags := make(map[string]Tag)
	for _, v := range tagsToDesc {
		for tag, descriptions := range v.(map[interface{}]interface{}) {
			tagName := tag.(string)
			for _, desc := range descriptions.([]interface{}) {
				description := desc.(string)
				tags[description] = Tag(tagName)
			}
		}
	}
	return &Tagger{tags}, nil
}

func (t *Tagger) Tag(description string) Tag {
	if tags, ok := t.classesToTags[class]; ok {
		return tags
	}
	for regex, tag := range t.classesToTags {
		if matched, _ := regexp.MatchString(string(regex), string(class)); matched {
			return tag
		}
	}
	return None
}
