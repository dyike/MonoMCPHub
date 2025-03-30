package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dyike/MonoMCPHub/internal/browser/config"
	sv "github.com/dyike/MonoMCPHub/pkg/service"
	"github.com/mark3labs/mcp-go/mcp"
)

type BrowserSevice struct {
	sv.ServiceManager
	config *config.BrowserConfig
	name   string
	ctx    context.Context
	cancel context.CancelFunc
}

func NewBrowserService(ctx context.Context, args []string) (sv.Service, error) {
	bconf := config.NewBrowserConfig()
	bs := &BrowserSevice{
		ctx:    ctx,
		config: bconf,
		name:   "browser_mcp",
	}
	bs.ServiceManager = *sv.NewServiceManager(ctx)
	err := bs.initBrowser(bconf.DataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init browser: %v", err)
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(bconf.UserAgent),
		chromedp.Flag("headless", bconf.Headless),
		chromedp.Flag("lang", bconf.DefaultLanguage),
		chromedp.WindowSize(1312, 848),
	)

	bs.ctx, bs.cancel = chromedp.NewExecAllocator(ctx, opts...)
	bs.ctx, bs.cancel = chromedp.NewContext(bs.ctx)

	bs.AddTool(mcp.NewTool(
		"browser_navigate",
		mcp.WithDescription("Navigate to a URL"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL to navigate to"),
		),
	), bs.handleNavigate)

	bs.AddTool(mcp.NewTool(
		"browser_screenshot",
		mcp.WithDescription("Take a screenshot of the current page"),
		mcp.WithString("name",
			mcp.Description("The name of the screenshot"),
			mcp.Required(),
		),
		mcp.WithString("selector",
			mcp.Description("The CSS selector of the element to screenshot"),
		),
		mcp.WithNumber("width",
			mcp.Description("The width of the screenshot in pixels, default: 1600"),
		),
		mcp.WithNumber("height",
			mcp.Description("The height of the screenshot in pixels, default: 1000"),
		),
	), bs.handleScreenshot)

	return bs, nil
}

func (bs *BrowserSevice) handleNavigate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url must be a string")
	}
	err := chromedp.Run(bs.ctx, chromedp.Navigate(url))
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Failed to navigate to %s: %v", url, err),
				},
			},
			IsError: true,
		}, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Navigated to %s", url),
			},
		},
	}, nil
}

func (bs *BrowserSevice) handleScreenshot(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name must be a string")
	}
	selector, ok := request.Params.Arguments["selector"].(string)
	width, _ := request.Params.Arguments["width"].(float64)
	height, _ := request.Params.Arguments["height"].(float64)

	if width == 0 {
		width = 1600
	}
	if height == 0 {
		height = 1000
	}

	var buf []byte
	var err error
	if selector == "" {
		err = chromedp.Run(bs.ctx, chromedp.FullScreenshot(&buf, 90))
	} else {
		// TODO: add width and height
		err = chromedp.Run(bs.ctx, chromedp.Screenshot(selector, &buf, chromedp.NodeVisible, chromedp.ByQuery))
	}
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Failed to take screenshot: %v", err),
				},
			},
			IsError: true,
		}, nil
	}
	newName := filepath.Join(bs.config.DataPath, fmt.Sprintf("%s_%d.png", name, time.Now().Unix()))
	err = os.WriteFile(newName, buf, 0644)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Failed to save screenshot: %v", err),
				},
			},
			IsError: true,
		}, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Screenshot saved to %s", newName),
			},
		},
	}, nil
}

func (bs *BrowserSevice) Close() error {
	bs.cancel()
	return nil
}

func (bs *BrowserSevice) Config() string {
	return ""
}

func (bs *BrowserSevice) Name() string {
	return bs.name
}

func (bs *BrowserSevice) initBrowser(userDataDir string) error {
	_, err := os.Stat(userDataDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check user data directory: %v", err)
	}
	if err == nil {
		err = os.RemoveAll(userDataDir)
		if err != nil {
			return fmt.Errorf("failed to remove user data directory: %v", err)
		}
	}
	err = os.MkdirAll(userDataDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create user data directory: %v", err)
	}

	return nil
}
