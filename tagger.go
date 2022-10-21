package main

import (
	"regexp"
)

type Tag = string

const (
	Recurring Tag = "Recurring"
	Crucial   Tag = "Crucial"
)

type Tagger struct {
	classesToTags map[string][]Tag
}

func (t *Tagger) Tags(class string) []Tag {
	if tags, ok := t.classesToTags[class]; ok {
		return tags
	}
	for regex, tag := range t.classesToTags {
		if matched, _ := regexp.MatchString(regex, class); matched {
			return tag
		}
	}
	return []Tag{}
}
