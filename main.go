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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "[%s] hello %s (sha:%s)\n", vi.Channel, vi.Version, shortSHA(vi.GitSHA, 7))
	})

	// JSON version endpoint
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Release-Channel", vi.Channel)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(vi)
	})

	// Simple healthcheck
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Optionally show env (handy for debugging; safe to keep or remove)
	mux.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"VERSION":   vi.Version,
			"GIT_SHA":   vi.GitSHA,
			"BUILD_TIME": vi.BuiltAt,
			"CHANNEL":   vi.Channel,
		})
	})

	// Basic server with sane timeouts
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           logging(mux),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("listening on %s (version=%s, sha=%s, channel=%s)", srv.Addr, vi.Version, shortSHA(vi.GitSHA, 7), vi.Channel)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &logWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(lw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lw.status, time.Since(start))
	})
}

type logWriter struct {
	http.ResponseWriter
	status int
}

func (lw *logWriter) WriteHeader(code int) {
	lw.status = code
	lw.ResponseWriter.WriteHeader(code)
}
