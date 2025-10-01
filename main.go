package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	version = getenv("VERSION", "0.0.0")
	gitSHA  = getenv("GIT_SHA", "dev")
	builtAt = getenv("BUILD_TIME", time.Now().UTC().Format(time.RFC3339))
	channel = getenv("CHANNEL", "dev") // dev | prod
)

func getenv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Release-Channel", channel)
		fmt.Fprintf(w, "[%s] hello %s (sha:%s)\n", channel, version, gitSHA[:7])
	})
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Release-Channel", channel)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"version":  version,
			"git_sha":  gitSHA,
			"built_at": builtAt,
			"channel":  channel,
		})
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	_ = http.ListenAndServe(":8080", nil)
}
