package config

import (
	"math/rand"
	"testing"
)

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
	b.StopTimer()
	cfg, err := LoadRoutes()
	if err != nil {
		b.Error("Failed to load config")
	}

	for n := 0; n < b.N; n++ {
		if n%2 == 0 {
			prefix := randValidPrefix()
			path := generateRandomPath(prefix)
			b.StartTimer()
			r := match(path, cfg.patterns)
			if r == nil {
				b.Errorf("Result Route should not be nil for path '%s', but it was", path)
			}
		} else {
			prefix := randInvalidPrefix()
			path := generateRandomPath(prefix)
			b.StartTimer()
			r := match(path, cfg.patterns)
			if r != nil {
				b.Errorf("Result Route should be nil for path '%s', but it was not", path)
			}
		}
	}
}

func BenchmarkGetRoute(b *testing.B) {
	b.StopTimer()
	cfg, err := LoadRoutes()
	if err != nil {
		b.Error("Failed to load config")
	}

	for n := 0; n < b.N; n++ {
		if n%2 == 0 {
			prefix := randValidPrefix()
			path := generateRandomPath(prefix)
			b.StartTimer()
			r := executeGetRoute(cfg, "http://test.com"+path)
			if r == nil {
				b.Errorf("Result Route should not be nil for path '%s', but it was", path)
			}
		} else {
			prefix := randInvalidPrefix()
			path := generateRandomPath(prefix)
			b.StartTimer()
			r := executeGetRoute(cfg, "http://test.com"+path)
			if r != nil {
				b.Errorf("Result Route should be nil for path '%s', but it was not", path)
			}
		}
	}
}
