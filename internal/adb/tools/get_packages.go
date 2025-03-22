package tools

import (
	"context"
	"fmt"

	"github.com/dyike/MonoMCPHub/repo/adb_repo"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewGetPackagesTool() mcp.Tool {
	return mcp.NewTool("get_packages",
		mcp.WithDescription("Get all packages of your android device"),
	)
}

func HandleGetPackages(adbRepo adb_repo.AdbRepo) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packages, err := adbRepo.GetPackages()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get packages: %v", err)
			return mcp.NewToolResultError(errMsg), nil
		}
		return mcp.NewToolResultText(packages), nil
	}
}
