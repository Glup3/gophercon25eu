package internal

import (
	"log/slog"
	"net/http"

	"github.com/glup3/gophercon25eu/public"
)

func NewServer(logger *slog.Logger, config *Config) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(logger, handlePublicAssets()))
	return mux
}

func handlePublicAssets() http.Handler {
	staticFS := http.FS(public.StaticFiles)
	fs := http.FileServer(staticFS)
	return fs
}

func loggingMiddleware(logger *slog.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(
			"web request",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.String("userAgent", r.UserAgent()),
			slog.String("referer", r.Referer()),
		)
		h.ServeHTTP(w, r)
	})
}

type Config struct {
	Port string
}

func NewConfig(getenv func(string) string) *Config {
	config := &Config{Port: "8080"}
	port := getenv("PORT")
	if port != "" {
		config.Port = port
	}
	return config
}
