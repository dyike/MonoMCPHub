package tools

import (
	"context"
	"fmt"

	"github.com/dyike/MonoMCPHub/repo/adb_repo"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewGetScreenshotTool() mcp.Tool {
	return mcp.NewTool("get_screenshot",
		mcp.WithDescription("Get a screenshot of the current device"),
	)
}

func HandleGetScreenshot(adbRepo adb_repo.AdbRepo) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		err := adbRepo.TakeScreenshot()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to take screenshot: %v", err)
			return mcp.NewToolResultError(errMsg), nil
		}
		return mcp.NewToolResultText("Screenshot take successfully"), nil
	}
}
