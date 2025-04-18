package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/dyike/MonoMCPHub/internal/browser/service"
	"github.com/mark3labs/mcp-go/server"
)

var (
	transport string
	port      string
)

func init() {
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&port, "p", "8080", "Port to listen on")
	flag.StringVar(&port, "port", "8080", "Port to listen on")
}

func main() {
	flag.Parse()

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

	if transport == "sse" {
		sseServer := server.NewSSEServer(s)
		if err := sseServer.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}

	} else {
		if err := server.ServeStdio(s); err != nil {
			slog.Error("Failed to serve", "error", err)
			os.Exit(1)
		}
	}
}
