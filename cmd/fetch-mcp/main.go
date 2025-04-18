package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/dyike/MonoMCPHub/internal/fetch"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	logLevel := slog.LevelDebug
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})))

	slog.Info("Starting browser mcp server")
	s := server.NewMCPServer(
		"fetch mcp server",
		"0.0.1",
	)

	ctx := context.Background()
	fs := fetch.NewFetchService(ctx)
	if fs == nil {
		slog.Error("Failed to create fetch service")
		os.Exit(1)
	}

	s.AddTools(fs.Tools()...)

	if err := server.ServeStdio(s); err != nil {
		slog.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}
