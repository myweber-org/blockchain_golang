package data

import (
	"strings"
)

type Cleaner struct {
	normalizeFunc func(string) string
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		normalizeFunc: strings.ToLower,
	}
}

func (c *Cleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	
	for _, item := range items {
		normalized := c.normalizeFunc(item)
		if _, exists := seen[normalized]; !exists {
			seen[normalized] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (c *Cleaner) SetNormalizeFunc(fn func(string) string) {
	c.normalizeFunc = fn
}