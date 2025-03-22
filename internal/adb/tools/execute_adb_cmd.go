package tools

import (
	"context"
	"fmt"

	"github.com/dyike/MonoMCPHub/repo/adb_repo"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewExecuteAdbCmdTool() mcp.Tool {
	return mcp.NewTool("execute_adb_cmd",
		mcp.WithDescription("Execute an adb command on the current device"),
		mcp.WithString("command",
			mcp.Required(),
			mcp.Description("The adb command to execute"),
		),
	)
}

func HandleExecuteAdbCmd(adbRepo adb_repo.AdbRepo) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cmdArgs := req.Params.Arguments["command"].(string)
		if cmdArgs == "" {
			return mcp.NewToolResultError("Command is required"), nil
		}

		output, err := adbRepo.ExecuteAdbCommand(cmdArgs)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to execute adb command: %v", err)
			return mcp.NewToolResultError(errMsg), nil
		}
		return mcp.NewToolResultText(output), nil
	}
}
