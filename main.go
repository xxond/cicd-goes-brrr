package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type VersionInfo struct {
	Version string `json:"version"`
	GitSHA  string `json:"git_sha"`
	BuiltAt string `json:"built_at"`
	Channel string `json:"channel"`
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func shortSHA(sha string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(sha) < n {
		return sha
	}
	return sha[:n]
}

func main() {
	version := getenv("VERSION", "0.0.0")
	gitSHA := getenv("GIT_SHA", "dev")
	builtAt := getenv("BUILD_TIME", time.Now().UTC().Format(time.RFC3339))
	channel := getenv("CHANNEL", "unknown")

	vi := VersionInfo{
		Version: version,
		GitSHA:  gitSHA,
		BuiltAt: builtAt,
		Channel: channel,
	}

	mux := http.NewServeMux()

	// Root: human-friendly banner
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Release-Channel", vi.Channel)
		w.Header().Set("Content-T
