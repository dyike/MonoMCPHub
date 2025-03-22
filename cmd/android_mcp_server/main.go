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
	workDir := os.Getenv("ANDROID_WORK_DIR")
	if workDir == "" {
		slog.Error("ANDROID_WORK_DIR is not set")
		os.Exit(1)
	}
	adbRepo := adb_repo.NewAdbRepo(deviceID, workDir)

	getPackageTool := tools.NewGetPackagesTool()
	s.AddTool(getPackageTool, tools.HandleGetPackages(adbRepo))

	getScreenshotTool := tools.NewGetScreenshotTool()
	s.AddTool(getScreenshotTool, tools.HandleGetScreenshot(adbRepo))

	getUILayoutTool := tools.NewGetUILayoutTool()
	s.AddTool(getUILayoutTool, tools.HandleGetUILayout(adbRepo))

	executeAdbCmdTool := tools.NewExecuteAdbCmdTool()
	s.AddTool(executeAdbCmdTool, tools.HandleExecuteAdbCmd(adbRepo))

	// getDevicesTool := tools.NewGetDevicesTool()
	// s.AddTool(getDevicesTool, tools.HandleGetDevices())

	// getAppLogTool := tools.NewAdbLogcatTool()
	// s.AddTool(getAppLogTool, tools.HandleAdbLogcat())

	if err := server.ServeStdio(s); err != nil {
		slog.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}
