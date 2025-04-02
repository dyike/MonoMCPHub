package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	sv "github.com/dyike/MonoMCPHub/pkg/service"
	"github.com/kkdai/youtube/v2"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36,gzip(gfe)"
)

var (
	reXMLTranscript = regexp.MustCompile(`<text start="([^"]*)" dur="([^"]*)">([^<]*)</text>`)
)

type TranscriptLine struct {
	Text     string  `json:"text"`
	Duration float64 `json:"duration"`
	Offset   float64 `json:"offset"`
	Lang     string  `json:"lang"`
}

type TranscriptResponse struct {
	Title string           `json:"title"`
	Lines []TranscriptLine `json:"lines"`
}

type CaptionTrack struct {
	BaseURL      string `json:"baseUrl"`
	LanguageCode string `json:"languageCode"`
}

type CaptionsData struct {
	CaptionTracks []CaptionTrack `json:"captionTracks"`
}

type FetchService struct {
	sv.ServiceManager
	client        *http.Client
	youtubeClient *youtube.Client
}

func NewFetchService(ctx context.Context) *FetchService {
	fs := &FetchService{
		client:        &http.Client{},
		youtubeClient: &youtube.Client{},
	}
	fs.ServiceManager = *sv.NewServiceManager(ctx)

	fs.AddTool(mcp.NewTool("fetch_url",
		mcp.WithDescription("Fetch the content of a URL, can return HTML or Markdown (default)"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL to fetch"),
		),
		mcp.WithBoolean("as_html",
			mcp.Description("Return the content as HTML"),
			mcp.DefaultBool(false),
		),
	), fs.handleFetchURL)

	return fs
}

func (fs *FetchService) handleFetchURL(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url must be a string")
	}
	asHTML, ok := request.Params.Arguments["as_html"].(bool)
	if !ok {
		asHTML = false
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := fs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if asHTML {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "html",
					Text: string(body),
				},
			},
		}, nil
	}

	// Convert HTML to Markdown
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	// Basic HTML to Markdown conversion
	var markdown strings.Builder
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		s.Find("h1,h2,h3,h4,h5,h6").Each(func(i int, h *goquery.Selection) {
			markdown.WriteString("#" + strings.Repeat("#", i))
			markdown.WriteString(" " + h.Text() + "\n\n")
		})
		s.Find("p").Each(func(i int, p *goquery.Selection) {
			markdown.WriteString(p.Text() + "\n\n")
		})
		s.Find("a").Each(func(i int, a *goquery.Selection) {
			href, exists := a.Attr("href")
			if exists {
				markdown.WriteString("[" + a.Text() + "](" + href + ")\n")
			}
		})
	})

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: markdown.String(),
			},
		},
	}, nil
}
