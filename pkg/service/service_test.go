package service

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestServiceAddResource(t *testing.T) {
	ctx := context.Background()
	sm := NewServiceManager(ctx)

	resource := mcp.Resource{Name: "test", URI: "test"}
	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				Text:     "test",
				URI:      "test",
				MIMEType: "text/plain",
			},
		}, nil
	}
	sm.AddResourceHandler(resource, handler)

	if len(sm.Resources()) != 1 {
		t.Errorf("expected 1 resource, got %d", len(sm.Resources()))
	}
	if sm.Resources()[resource] == nil {
		t.Errorf("expected resource to be added")
	}
}

func TestServiceAddResourceTemplate(t *testing.T) {
	ctx := context.Background()
	sm := NewServiceManager(ctx)

	resourceTemplate := mcp.ResourceTemplate{Name: "testPemplate"}
	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				Text:     "test",
				URI:      "test",
				MIMEType: "text/plain",
			},
		}, nil

	}
	sm.AddResourceTemplateHandler(resourceTemplate, handler)

	if len(sm.ResourceTemplates()) != 1 {
		t.Errorf("expected 1 resource template, got %d", len(sm.ResourceTemplates()))
	}
	if sm.ResourceTemplates()[resourceTemplate] == nil {
		t.Errorf("expected resource template to be added")
	}
}

func TestServiceAddNotificationHandler(t *testing.T) {
	ctx := context.Background()
	sm := NewServiceManager(ctx)

	handler := func(ctx context.Context, notification mcp.JSONRPCNotification) {
		t.Logf("notification: %v", notification)
	}
	name := "testNotification"
	sm.AddNotificationHandler(name, handler)

	if len(sm.NotificationHandlers()) != 1 {
		t.Errorf("expected 1 notification handler, got %d", len(sm.NotificationHandlers()))
	}
	if sm.NotificationHandlers()[name] == nil {
		t.Errorf("expected notification handler to be added")
	}
}

func TestServiceAddPrompt(t *testing.T) {
	ctx := context.Background()
	sm := NewServiceManager(ctx)

	prompt := mcp.Prompt{Name: "testPrompt"}
	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		pms := make([]mcp.PromptMessage, 0)
		pms = append(pms, mcp.PromptMessage{
			Role: mcp.RoleUser,
			Content: mcp.TextContent{
				Type: "text",
				Text: "Prompt msg",
			},
		})
		return &mcp.GetPromptResult{
			Description: "Prompt description",
			Messages:    pms,
		}, nil
	}
	sm.AddPrompt(prompt, handler)

	if len(sm.Prompts()) != 1 {
		t.Errorf("expected 1 prompt, got %d", len(sm.Prompts()))
	}

	for _, p := range sm.Prompts() {
		if p.Prompt().Name != prompt.Name {
			t.Errorf("expected prompt to be added")
		}
	}
}

func TestServiceAddTool(t *testing.T) {
	ctx := context.Background()
	sm := NewServiceManager(ctx)

	tool := mcp.Tool{Name: "testTool"}
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Tool msg",
				},
			},
		}, nil
	}
	sm.AddTool(tool, handler)

	if len(sm.Tools()) != 1 {
		t.Errorf("expected 1 tool, got %d", len(sm.Tools()))
	}

	for _, tl := range sm.Tools() {
		if tl.Tool.Name != tool.Name {
			t.Errorf("expected tool to be added")
		}
		if tl.Handler == nil {
			t.Errorf("expected tool handler to be added")
		}
	}
}
