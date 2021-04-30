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
)

type Tagger struct {
	classesToTags map[string][]Tag
}

func NewTagger(yamlParseable []byte) (*Tagger, error) {
	tagsToClasses := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(yamlParseable, &tagsToClasses); err != nil {
		return &Tagger{}, err
	}

	classesToTags := make(map[string][]Tag)
	for _, v := range tagsToClasses {
		for class, tgs := range v.(map[interface{}]interface{}) {
			className := class.(string)
			ts := make([]Tag, len(tgs.([]interface{})))
			for i, tg := range tgs.([]interface{}) {
				ts[i] = Tag(tg.(string))
			}
			classesToTags[className] = ts
		}
	}
	return &Tagger{classesToTags}, nil
}

func (t *Tagger) Tags(class string) []Tag {
	if tags, ok := t.classesToTags[class]; ok {
		return tags
	}
	for regex, tag := range t.classesToTags {
		if matched, _ := regexp.MatchString(string(regex), string(class)); matched {
			return tag
		}
	}
	return []Tag{}
}
