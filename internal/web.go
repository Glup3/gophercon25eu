package internal

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/glup3/gophercon25eu/public"
)

func NewServer(logger *slog.Logger, config *Config) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(logger, handlePublicAssets(config)))
	return mux
}

func handlePublicAssets(config *Config) http.Handler {
	if config.UseEmbeddedFS {
		return http.FileServer(http.FS(public.StaticFiles))
	}

	return http.FileServer(http.FS(os.DirFS(config.StaticFilesPath)))
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
	Hostname        string
	Port            string
	UseEmbeddedFS   bool
	StaticFilesPath string
}

func NewConfig(getenv func(string) string) *Config {
	hostname := getenv("HOSTNAME")
	if hostname == "" {
		hostname = "127.0.0.1"
	}
	port := getenv("PORT")
	if port == "" {
		port = "8080"
	}
	useEmbeddedFs := getenv("USE_EMBEDDED_FS") == "true"
	staticPath := getenv("STATIC_FILES_PATH")
	if staticPath == "" {
		staticPath = "public"
	}

	return &Config{
		Hostname:        hostname,
		Port:            port,
		UseEmbeddedFS:   useEmbeddedFs,
		StaticFilesPath: staticPath,
	}
}
