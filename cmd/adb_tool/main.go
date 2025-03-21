package main

import (
	"log/slog"
	"os"

	"github.com/dyike/MonoMCPHub/internal/adb/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"android adb mcp server",
		"0.0.1",
	)

	getDevicesTool := tools.NewGetDevicesTool()
	s.AddTool(getDevicesTool, tools.HandleGetDevices())

	getAppLogTool := tools.NewAdbLogcatTool()
	s.AddTool(getAppLogTool, tools.HandleAdbLogcat())

	if err := server.ServeStdio(s); err != nil {
		slog.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}
