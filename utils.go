package main

import "os"

func EnvOr(key string, fallback string) string {
	val := os.Getenv(key)

	if val == "" {
		return fallback
	}
	return val
}


