package main

import (
	"testing"
)

func TestGetEnvFallbackValue(t *testing.T) {

	result := getEnv("shouldnotExists", "fallbackValue")

	want := "fallbackValue"
	if result != want {
		t.Errorf("getEnv returned %+v, want %+v", result, want)
	}
}
