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

type BrowserService struct {
	sv.ServiceManager
	config *config.BrowserConfig
	name   string
	ctx    context.Context
	cancel context.CancelFunc
}

func NewBrowserService(ctx context.Context, args []string) (sv.Service, error) {
	bconf := config.NewBrowserConfig()
	bs := &BrowserService{
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

	bs.AddTool(mcp.NewTool(
		"browser_click",
		mcp.WithDescription("Click on an element on the page"),
		mcp.WithString("selector",
			mcp.Description("The CSS selector of the element to click on"),
			mcp.Required(),
		),
	), bs.handleClick)

	bs.AddTool(mcp.NewTool(
		"browser_fill",
		mcp.WithDescription("Fill an input with a value"),
		mcp.WithString("selector",
			mcp.Description("The CSS selector of the input to fill"),
			mcp.Required(),
		),
		mcp.WithString("value",
			mcp.Description("The value to fill the input with"),
			mcp.Required(),
		),
	), bs.handleFill)

	bs.AddTool(mcp.NewTool(
		"browser_select",
		mcp.WithDescription("Select an element on the page with selector tag"),
		mcp.WithString("selector",
			mcp.Description("The CSS selector for element to select"),
			mcp.Required(),
		),
		mcp.WithString("value",
			mcp.Description("The value to select"),
			mcp.Required(),
		),
	), bs.handleSelect)

	bs.AddTool(mcp.NewTool(
		"browser_hover",
		mcp.WithDescription("Hover over an element on the page"),
		mcp.WithString("selector",
			mcp.Description("The CSS selector for element to hover over"),
			mcp.Required(),
		),
	), bs.handleHover)

	bs.AddTool(mcp.NewTool(
		"browser_evaluate",
		mcp.WithDescription("Execute a JavaScript in the browser console"),
		mcp.WithString("script",
			mcp.Description("The JavaScript code to execute"),
			mcp.Required(),
		),
	), bs.handleEvaluate)

	return bs, nil
}

func (bs *BrowserService) handleNavigate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (bs *BrowserService) handleScreenshot(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (bs *BrowserService) handleClick(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector, ok := request.Params.Arguments["selector"].(string)
	result := &mcp.CallToolResult{
		IsError: false,
	}
	if !ok {
		result.IsError = true
		result.Content = []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "selector must be a string",
			},
		}
		return result, nil
	}
	err := chromedp.Run(bs.ctx, chromedp.Click(selector, chromedp.NodeVisible, chromedp.ByQuery))
	if err != nil {
		result.IsError = true
		result.Content = []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Failed to click on %s: %v", selector, err),
			},
		}
		return result, nil
	}
	result.Content = []mcp.Content{
		mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("Clicked on %s", selector),
		},
	}
	return result, nil
}

func (bs *BrowserService) handleFill(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector, ok := request.Params.Arguments["selector"].(string)
	if !ok {
		return nil, fmt.Errorf("selector must be a string")
	}
	value, ok := request.Params.Arguments["value"].(string)
	if !ok {
		return nil, fmt.Errorf("value must be a string")
	}
	err := chromedp.Run(bs.ctx, chromedp.SendKeys(selector, value, chromedp.NodeVisible, chromedp.ByQuery))
	if err != nil {
		return nil, fmt.Errorf("failed to fill %s with %s: %v", selector, value, err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Filled input %s with value %s", selector, value),
			},
		},
	}, nil
}

func (bs *BrowserService) handleSelect(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector, ok := request.Params.Arguments["selector"].(string)
	if !ok {
		return nil, fmt.Errorf("selector must be a string")
	}
	value, ok := request.Params.Arguments["value"].(string)
	if !ok {
		return nil, fmt.Errorf("value must be a string")
	}
	err := chromedp.Run(bs.ctx, chromedp.SetValue(selector, value, chromedp.NodeVisible, chromedp.ByQuery))
	if err != nil {
		return nil, fmt.Errorf("failed to select %s with value %s: %v", selector, value, err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Selected %s with value %s", selector, value),
			},
		},
	}, nil
}

func (bs *BrowserService) handleHover(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector, ok := request.Params.Arguments["selector"].(string)
	if !ok {
		return nil, fmt.Errorf("selector must be a string")
	}
	var res bool
	err := chromedp.Run(bs.ctx, chromedp.Evaluate(`document.querySelector('`+selector+`').dispatchEvent(new Event('mouseover'))`, &res))
	if err != nil {
		return nil, fmt.Errorf("failed to hover over %s: %v", selector, err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Hovered over %s, result: %t", selector, res),
			},
		},
	}, nil
}

func (bs *BrowserService) handleEvaluate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	script, ok := request.Params.Arguments["script"].(string)
	if !ok {
		return nil, fmt.Errorf("script must be a string")
	}
	var result interface{}
	err := chromedp.Run(bs.ctx, chromedp.Evaluate(script, &result))
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate %s: %v", script, err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Evaluated %s, result: %v", script, result),
			},
		},
	}, nil
}

func (bs *BrowserService) Close() error {
	bs.cancel()
	return nil
}

func (bs *BrowserService) Config() string {
	return ""
}

func (bs *BrowserService) Name() string {
	return bs.name
}

func (bs *BrowserService) initBrowser(userDataDir string) error {
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
