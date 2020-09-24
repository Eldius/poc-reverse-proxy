package config

import (
	"regexp"
	"net/http"
)

// Headers header filter
type Headers struct {
	In  []string
	Out []string
}

// Route request route definition
type Route struct {
	Methods  []string
	Headers  Headers
	Backends []string
}

// RoutesConfig request route configuration
type RoutesConfig struct {
	Routes   map[string]Route
	patterns map[*regexp.Regexp]Route
}

// GetRoute finds the right Route to this request
func (c *RoutesConfig) GetRoute(req *http.Request) *Route {
	return match(req.URL.RequestURI(), c.patterns)
}

func match(path string, routes map[*regexp.Regexp]Route) *Route {
	for p, r := range routes {
		if p.MatchString(path) {
			return &r
		}
	}
	return nil
}
