package config

import (
	"net/http/httptest"
	"testing"
)

func TestMatch(t *testing.T) {
	if cfg, err := LoadRoutes(); err != nil {
		t.Error("Failed to load configuration from file")
	} else {

		r := match("/app01", cfg.patterns)
		if r == nil {
			t.Error("Should return a non nil route")
		} else if r.Backends[0] != "http://localhost:1111" {
			t.Errorf("Backend should be 'http://localhost:1111', but was '%s'", r.Backends[0])
		}
	}
}

func TestNotMatch(t *testing.T) {
	if cfg, err := LoadRoutes(); err != nil {
		t.Error("Failed to load configuration from file")
	} else {

		r := match("/app", cfg.patterns)
		if r != nil {
			t.Error("Should return a nil route")
		}
	}
}

func TestMatchPrefix01(t *testing.T) {
	if cfg, err := LoadRoutes(); err != nil {
		t.Error("Failed to load configuration from file")
	} else {
		r := match("/app01", cfg.patterns)
		if r.Backends[0] != "http://localhost:1111" {
			t.Errorf("Backend should be 'http://localhost:1111', but was '%s'", r.Backends[0])
		}
		r1 := match("/app01/xpto", cfg.patterns)
		if r1.Backends[0] != "http://localhost:1111" {
			t.Errorf("Backend should be 'http://localhost:1111', but was '%s'", r1.Backends[0])
		}
	}
}

func TestMatchPrefix02(t *testing.T) {
	if cfg, err := LoadRoutes(); err != nil {
		t.Error("Failed to load configuration from file")
	} else {
		r := match("/app02", cfg.patterns)
		if r.Backends[0] != "http://localhost:2222" {
			t.Errorf("Backend should be 'http://localhost:2222', but was '%s'", r.Backends[0])
		}
		r1 := match("/app02/xpto", cfg.patterns)
		if r1.Backends[0] != "http://localhost:2222" {
			t.Errorf("Backend should be 'http://localhost:2222', but was '%s'", r1.Backends[0])
		}
	}
}

func TestMatchRequest03(t *testing.T) {
	cfg, err := LoadRoutes()
	if err != nil {
		t.Error("Failed to load routes", err.Error())
	}
	r := executeGetRoute(cfg, "http://example.com/app03/abcd")
	if r == nil {
		t.Errorf("Should return a non nil route")
	} else if r.Backends[0] != "http://localhost:3333" {
		t.Errorf("Should return 'http://localhost:3333' as backend, but was '%s'", r.Backends[0])
	}
}

func TestMatchRequest01(t *testing.T) {
	cfg, err := LoadRoutes()
	if err != nil {
		t.Error("Failed to load routes", err.Error())
	}
	r := executeGetRoute(cfg, "http://example.com/app01/abcd")
	if r == nil {
		t.Errorf("Should return a non nil route")
	} else if r.Backends[0] != "http://localhost:1111" {
		t.Errorf("Should return 'http://localhost:1111' as backend, but was '%s'", r.Backends[0])
	}
}

func TestMatchRequestNotMatch(t *testing.T) {
	cfg, err := LoadRoutes()
	if err != nil {
		t.Error("Failed to load routes", err.Error())
	}

	r := executeGetRoute(cfg, "http://example.com/xpto/abcd")
	if r != nil {
		t.Errorf("Should return a nil route")
	}
}

func TestRedirect(t *testing.T) {
	gock.
}

func executeGetRoute(cfg RoutesConfig, url string) *Route {
	req := httptest.NewRequest("GET", url, nil)
	return cfg.GetRoute(req)
}
