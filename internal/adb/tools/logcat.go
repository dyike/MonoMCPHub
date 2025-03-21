package tools

import (
	"context"
	"errors"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

func NewAdbLogcatTool() mcp.Tool {
	return mcp.NewTool("adb_logcat",
		mcp.WithDescription("Get logcat of android device"),
		mcp.WithString("device_id",
			mcp.Required(),
			mcp.Description("Device ID"),
		),
		mcp.WithNumber("line",
			mcp.DefaultNumber(100),
			mcp.Description("Number of lines to get"),
		),
		// mcp.WithString("log_level",
		// 	mcp.Required(),
		// 	mcp.DefaultString("verbose"),
		// 	mcp.Description("Log level"),
		// ),
		mcp.WithString("keyword",
			mcp.Required(),
			mcp.Description("Service name to filter"),
		),
	)
}

func HandleAdbLogcat() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		deviceID := req.Params.Arguments["device_id"].(string)
		if deviceID == "" {
			return nil, errors.New("device_id is required")
		}
		keyword := req.Params.Arguments["keyword"].(string)
		if keyword == "" {
			return mcp.NewToolResultError("keyword is required"), nil
		}
		lines := req.Params.Arguments["line"].(float64)

		args := []string{"-s", deviceID, "logcat", "-d", "-v", "time", "-s", keyword, "-n", strconv.Itoa(int(lines))}
		output, err := ExecuteAdbCommand(args)
		if err != nil {
			return mcp.NewToolResultError("Failed to get logcat"), nil
		}

		return mcp.NewToolResultText(output), nil
	}
}
