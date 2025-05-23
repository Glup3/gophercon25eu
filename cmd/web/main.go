package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	internal "github.com/glup3/gophercon25eu/internal"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	config := internal.NewConfig(getenv)
	srv := internal.NewServer(logger, config)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Hostname, config.Port),
		Handler: srv,
	}

	go func() {
		logger.Info("http server started", slog.String("addr", httpServer.Addr), slog.Bool("useEmbeddedFS", config.UseEmbeddedFS))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error listening and serving", slog.Any("error", err))
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		logger.Info("shutting down http server")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("error shutting down http server", slog.Any("error", err))
		}
	}()
	wg.Wait()
	return nil
}
