package tools

import (
	"context"
	"fmt"

	"github.com/dyike/MonoMCPHub/repo/adb_repo"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewGetUILayoutTool() mcp.Tool {
	return mcp.NewTool("get_uilayout",
		mcp.WithDescription("Get the UI layout of the current device"),
	)
}

func HandleGetUILayout(adbRepo adb_repo.AdbRepo) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		layout, err := adbRepo.GetUILayout()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get UI layout: %v", err)
			return mcp.NewToolResultError(errMsg), nil
		}
		return mcp.NewToolResultText(layout), nil
	}
}
