package main

import (
	"log/slog"
	"os"

	"github.com/dyike/MonoMCPHub/internal/unsplash/config"
	"github.com/dyike/MonoMCPHub/internal/unsplash/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	s := server.NewMCPServer(
		"unsplash mcp server",
		"0.0.1",
	)

	searchTool := tools.NewSearchPhotosTool()
	s.AddTool(searchTool, tools.HandleSearchPhotos(cfg))

	if err := server.ServeStdio(s); err != nil {
		slog.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}
