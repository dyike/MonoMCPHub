package android_mcp_client

import (
	"context"
	"testing"

	eino_mcp "github.com/cloudwego/eino-ext/components/tool/mcp"
)

func TestNewAndroidMcpClient(t *testing.T) {
	ctx := context.Background()
	cmd := "/Users/bytedance/Code/go/bin/android_mcp_server"
	env := []string{
		"ANDROID_DEVICE_ID=100.80.162.9:5555",
		"ANDROID_WORK_DIR=/Users/bytedance/.mcp_tmp",
	}
	args := []string{}
	cli, err := NewAndroidMcpClient(ctx, cmd, env, args...)
	if err != nil {
		t.Fatal(err)
	}

	tools, _ := eino_mcp.GetTools(ctx, &eino_mcp.Config{Cli: cli})
	for _, tool := range tools {
		info, err := tool.Info(ctx)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(info)
	}
	defer cli.Close()
}
