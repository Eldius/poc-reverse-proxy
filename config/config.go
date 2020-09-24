package config

import (
	"fmt"
	"regexp"

	"github.com/spf13/viper"
)

// LoadRoutes loads configuration from file
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
