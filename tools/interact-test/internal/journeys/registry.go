package journeys

import (
	"fmt"
	"sort"
	"strings"
)

// Journey is a single test journey
type Journey struct {
	Name   string
	Suite  string
	Desc   string
	Run    func(*Context) error
}

// Context carries test context
type Context struct {
	BaseURL string
	Page    interface{} // playwright.Page
	Gates   interface{} // uxgates.Recorder
	Logf    func(format string, args ...interface{})
}

var registry []Journey

// Register adds a journey to the registry
func Register(j Journey) {
	// Check for duplicate names
	for _, existing := range registry {
		if existing.Name == j.Name {
			panic(fmt.Sprintf("duplicate journey name: %s", j.Name))
		}
	}
	registry = append(registry, j)
}

// All returns all journeys sorted by name
func All() []Journey {
	sort.Slice(registry, func(i, j int) bool {
		return registry[i].Name < registry[j].Name
	})
	return registry
}

// BySuite returns journeys filtered by suite
func BySuite(suite string) []Journey {
	var result []Journey
	for _, j := range registry {
		if j.Suite == suite || suite == "" {
			result = append(result, j)
		}
	}
	return result
}

// ByGrep returns journeys matching a filter
func ByGrep(filter string) []Journey {
	var result []Journey
	for _, j := range registry {
		if strings.Contains(strings.ToLower(j.Name), strings.ToLower(filter)) ||
			strings.Contains(strings.ToLower(j.Desc), strings.ToLower(filter)) {
			result = append(result, j)
		}
	}
	return result
}
