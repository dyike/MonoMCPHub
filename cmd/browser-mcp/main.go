package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/dyike/MonoMCPHub/internal/browser/service"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	logLevel := slog.LevelDebug
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})))

	slog.Info("Starting browser mcp server")
	s := server.NewMCPServer(
		"browser mcp server",
		"0.0.1",
	)

	ctx := context.Background()
	bs, err := service.NewBrowserService(ctx, []string{})
	if err != nil {
		slog.Error("Failed to create browser service", "error", err)
		os.Exit(1)
	}

	s.AddTools(bs.Tools()...)

	// tools := bs.Tools()
	// slog.Info("Available tools", "tools count", len(tools))

	// for _, tool := range tools {
	// 	s.AddTool(tool.Tool, tool.Handler)
	// }

	if err := server.ServeStdio(s); err != nil {
		slog.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}
