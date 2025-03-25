package android_mcp_client

import (
	"context"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewAndroidMcpClient(ctx context.Context, cmd string, env []string, args ...string) (*client.StdioMCPClient, error) {
	cli, err := client.NewStdioMCPClient(
		cmd,
		env,
		args...,
	)

	if err != nil {
		return nil, err
	}
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "android_mcp_client",
		Version: "0.0.1",
	}

	_, err = cli.Initialize(ctx, initReq)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
