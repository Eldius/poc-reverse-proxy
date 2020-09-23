package config

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/spf13/viper"
)

type Headers struct {
	In  []string
	Out []string
}

type Route struct {
	Methods  []string
	Headers  Headers
	Backends []string
}

type RoutesConfig struct {
	Routes   map[string]Route
	patterns map[*regexp.Regexp]Route
}

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

func LoadRoutes() (cfg RoutesConfig, err error) {
	routes := make(map[string]Route)
	patterns := make(map[*regexp.Regexp]Route)
	for k := range viper.GetStringMap("routes") {
		var r Route
		_ = viper.UnmarshalKey(fmt.Sprintf("routes.%s", k), &r)
		patterns[regexp.MustCompile(k)] = r
		routes[k] = r
	}
	cfg = RoutesConfig{
		Routes:   routes,
		patterns: patterns,
	}
	return
}
