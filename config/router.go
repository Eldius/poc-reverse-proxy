package config

import (
	"io/ioutil"
	"net/http"
	"regexp"
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
	client   *http.Client
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

func (r *Route) Redirect(w http.ResponseWriter, req *http.Request) {
	req1, err := http.NewRequest(req.Method, r.Backends[0]+req.URL.Path, req.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	res, err := r.client.Do(req1)
	if err != nil {
		w.WriteHeader(502)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer res.Body.Close()
	w.WriteHeader(res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(502)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write(body)
}
