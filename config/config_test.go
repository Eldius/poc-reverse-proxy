package config

import (
	"log"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("samples/config.yaml")
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func TestLoadRoutes(t *testing.T) {
	if cfg, err := LoadRoutes(); err != nil {
		t.Error("Failed to load configuration from file")
	} else {
		if len(cfg.Routes) != 4 {
			t.Errorf("Should have 2 routes, but has %d", len(cfg.Routes))
		}
	}
}

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

func executeGetRoute(cfg RoutesConfig, url string) *Route {
	req := httptest.NewRequest("GET", url, nil)
	return cfg.GetRoute(req)
}

var (
	validPrefixes   = []string{"/app01", "/app02", "/app03/xpto", "/app04/abc"}
	invalidPrefixes = []string{"/app", "/xpto", "/app0/xpto", "/abc"}
)

func generateRandomPath(prefix string) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return prefix + "/" + string(b)
}

func randValidPrefix() string {
	return validPrefixes[rand.Intn(len(validPrefixes))]
}

func randInvalidPrefix() string {
	return invalidPrefixes[rand.Intn(len(invalidPrefixes))]
}

func BenchmarkMatch(b *testing.B) {
	cfg, err := LoadRoutes()
	if err != nil {
		b.Error("Failed to load config")
	}

	qtd := b.N
	start := time.Now()
	for n := 0; n < qtd; n++ {
		if n%2 == 0 {
			prefix := randValidPrefix()
			path := generateRandomPath(prefix)
			r := match(path, cfg.patterns)
			if r == nil {
				b.Errorf("Result Route should not be nil for path '%s', but it was", path)
			}
		} else {
			prefix := randInvalidPrefix()
			path := generateRandomPath(prefix)
			r := match(path, cfg.patterns)
			if r != nil {
				b.Errorf("Result Route should be nil for path '%s', but it was not", path)
			}
		}

	}
	b.Logf("iterations: %d (%d ms)", qtd, time.Since(start).Milliseconds())

}

func BenchmarkGetRoute(b *testing.B) {
	cfg, err := LoadRoutes()
	if err != nil {
		b.Error("Failed to load config")
	}

	qtd := b.N
	start := time.Now()
	for n := 0; n < qtd; n++ {
		if n%2 == 0 {
			prefix := randValidPrefix()
			path := generateRandomPath(prefix)
			r := executeGetRoute(cfg, "http://test.com"+path)
			if r == nil {
				b.Errorf("Result Route should not be nil for path '%s', but it was", path)
			}
		} else {
			prefix := randInvalidPrefix()
			path := generateRandomPath(prefix)
			r := executeGetRoute(cfg, "http://test.com"+path)
			if r != nil {
				b.Errorf("Result Route should be nil for path '%s', but it was not", path)
			}
		}

	}
	b.Logf("iterations: %d (%d ms)", qtd, time.Since(start).Milliseconds())
}
