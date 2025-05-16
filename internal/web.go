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

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *customResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(logger *slog.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := &customResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default status code before WriteHeader is called
		}
		h.ServeHTTP(wr, r)
		logger.Info(
			"web request",
			slog.Int("status", wr.statusCode),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.String("referer", r.Referer()),
			slog.String("userAgent", r.UserAgent()),
		)
	})
}

type Config struct {
	Hostname string
	Port     string
}

func NewConfig(getenv func(string) string) *Config {
	config := &Config{Hostname: "127.0.0.1", Port: "8080"}
	hostname := getenv("HOSTNAME")
	if hostname != "" {
		config.Hostname = hostname
	}
	port := getenv("PORT")
	if port != "" {
		config.Port = port
	}
	return config
}
