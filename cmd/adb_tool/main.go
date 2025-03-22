package main

import (
	"log/slog"
	"os"

	"github.com/dyike/MonoMCPHub/internal/adb/tools"
	"github.com/dyike/MonoMCPHub/repo/adb_repo"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"android adb mcp server",
		"0.0.1",
	)

	deviceID := os.Getenv("ANDROID_DEVICE_ID")
	if deviceID == "" {
		slog.Error("ANDROID_DEVICE_ID is not set")
		os.Exit(1)
	}
	adbRepo := adb_repo.NewAdbRepo(deviceID)

	getPackageTools := tools.NewGetPackagesTool()
	s.AddTool(getPackageTools, tools.HandleGetPackages(adbRepo))

	// getDevicesTool := tools.NewGetDevicesTool()
	// s.AddTool(getDevicesTool, tools.HandleGetDevices())

	// getAppLogTool := tools.NewAdbLogcatTool()
	// s.AddTool(getAppLogTool, tools.HandleAdbLogcat())

	if err := server.ServeStdio(s); err != nil {
		slog.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}
