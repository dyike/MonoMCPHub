package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func NewGetDevicesTool() mcp.Tool {
	return mcp.NewTool("get_devices",
		mcp.WithDescription("Get all devices"),
		mcp.WithBoolean("show_detail",
			mcp.DefaultBool(true),
			mcp.Description("Show device details (-l)"),
		),
	)
}

func HandleGetDevices() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"devices"}
		showDetail := req.Params.Arguments["show_detail"].(bool)

		if showDetail {
			args = append(args, "-l")
		}

		output, err := ExecuteAdbCommand(args)
		if err != nil {
			return mcp.NewToolResultError("Failed to get devices"), nil
		}

		return mcp.NewToolResultText(output), nil
	}
}
