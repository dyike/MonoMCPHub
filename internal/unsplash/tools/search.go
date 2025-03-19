package tools

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/dyike/MonoMCPHub/internal/unsplash/config"
	"github.com/dyike/MonoMCPHub/repo/api/unsplash"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewSearchPhotosTool() mcp.Tool {
	return mcp.NewTool("search_photos",
		mcp.WithDescription("Search for Unsplash photos"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search keyword"),
		),
		mcp.WithNumber("page",
			mcp.Required(),
			mcp.DefaultNumber(1),
			mcp.Description("Page number (1-based)"),
		),
		mcp.WithNumber("per_page",
			mcp.Required(),
			mcp.DefaultNumber(5),
			mcp.Description("Results per page (1-30)"),
		),
		mcp.WithString("order_by",
			mcp.Required(),
			mcp.DefaultString("relevant"),
			mcp.Description("Sort method (relevant or latest)"),
		),
		mcp.WithString("color",
			mcp.Description("Color filter (black_and_white, black, white, yellow, orange, red, purple, magenta, green, teal, blue)"),
		),
		mcp.WithString("orientation",
			mcp.Description("Orientation filter (landscape, portrait, squarish)"),
		),
	)
}

func HandleSearchPhotos(cfg *config.Config) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		conf := &unsplash.UnsplashConfig{
			AccessKey: cfg.UnsplashAPIKey,
			Timeout:   cfg.Timeout,
		}
		client := unsplash.NewUnsplashClient(conf)

		params := url.Values{}

		query, ok := req.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("Search query is required"), nil
		}
		params.Add("query", query)

		page := req.Params.Arguments["page"].(int)
		if page < 1 {
			page = 1
		}
		params.Add("page", strconv.Itoa(page))

		perPage := req.Params.Arguments["per_page"].(int)
		if perPage < 1 || perPage > 30 {
			perPage = 5
		}
		params.Add("per_page", strconv.Itoa(perPage))

		orderBy := req.Params.Arguments["order_by"].(string)
		params.Add("order_by", orderBy)

		color, ok := req.Params.Arguments["color"].(string)
		if ok && color != "" {
			params.Add("color", color)
		}

		orientation, ok := req.Params.Arguments["orientation"].(string)
		if ok && orientation != "" {
			params.Add("orientation", orientation)
		}

		photos, err := client.SearchPhotos(params)
		if err != nil {
			return mcp.NewToolResultError("Failed to search photos"), nil
		}

		if len(photos.Results) == 0 {
			return mcp.NewToolResultError("No photos found"), nil
		}

		payload, _ := json.Marshal(photos.Results)

		return mcp.NewToolResultText(string(payload)), nil
	}
}
