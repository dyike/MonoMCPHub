package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
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

	fs.AddTool(mcp.NewTool("fetch_youtube_transcript",
		mcp.WithDescription("Fetch the transcript of a YouTube video"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL of the YouTube video"),
		),
		mcp.WithString("lang",
			mcp.Description("Language code for the transcript (e.g., 'en', 'zh')"),
		),
	), fs.handleFetchYouTubeTranscript)

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

func (fs *FetchService) handleFetchYouTubeTranscript(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url must be a string")
	}

	lang, _ := request.Params.Arguments["lang"].(string)

	videoID, err := extractVideoID(url)
	if err != nil {
		return nil, fmt.Errorf("failed to extract video ID: %v", err)
	}

	// 创建请求获取视频页面
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	if lang != "" {
		req.Header.Set("Accept-Language", lang)
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

	pageContent := string(body)

	// 检查是否需要验证码
	if strings.Contains(pageContent, "class=\"g-recaptcha\"") {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "YouTube 需要验证码，请稍后再试",
				},
			},
			IsError: true,
		}, nil
	}

	// 检查视频是否可用
	if !strings.Contains(pageContent, "\"playabilityStatus\":") {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("视频不可用 (%s)", videoID),
				},
			},
			IsError: true,
		}, nil
	}

	// 提取标题
	titleMatch := regexp.MustCompile(`<title>(.*)</title>`).FindStringSubmatch(pageContent)
	title := "unknown"
	if len(titleMatch) > 1 {
		title = titleMatch[1]
	}

	// 提取字幕数据
	parts := strings.Split(pageContent, "\"captions\":")
	if len(parts) <= 1 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("该视频未启用字幕 (%s)", videoID),
				},
			},
			IsError: true,
		}, nil
	}

	captionsJSON := strings.Split(parts[1], ",\"videoDetails")[0]
	var captionsData struct {
		PlayerCaptionsTracklistRenderer struct {
			CaptionTracks []CaptionTrack `json:"captionTracks"`
		} `json:"playerCaptionsTracklistRenderer"`
	}

	if err := json.Unmarshal([]byte(captionsJSON), &captionsData); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "无法解析字幕数据",
				},
			},
			IsError: true,
		}, nil
	}

	tracks := captionsData.PlayerCaptionsTracklistRenderer.CaptionTracks
	if len(tracks) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("该视频没有可用的字幕 (%s)", videoID),
				},
			},
			IsError: true,
		}, nil
	}

	// 选择字幕轨道
	var selectedTrack *CaptionTrack
	if lang != "" {
		for _, track := range tracks {
			if track.LanguageCode == lang {
				selectedTrack = &track
				break
			}
		}
		if selectedTrack == nil {
			availableLangs := make([]string, len(tracks))
			for i, track := range tracks {
				availableLangs[i] = track.LanguageCode
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("该视频没有 %s 语言的字幕。可用语言: %s", lang, strings.Join(availableLangs, ", ")),
					},
				},
				IsError: true,
			}, nil
		}
	} else {
		selectedTrack = &tracks[0]
	}

	// 获取字幕内容
	req, err = http.NewRequestWithContext(ctx, "GET", selectedTrack.BaseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err = fs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	transcriptBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析字幕XML
	matches := reXMLTranscript.FindAllStringSubmatch(string(transcriptBody), -1)
	if len(matches) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "无法解析字幕内容",
				},
			},
			IsError: true,
		}, nil
	}

	// 构建字幕行
	var lines []TranscriptLine
	for _, match := range matches {
		offset, _ := strconv.ParseFloat(match[1], 64)
		duration, _ := strconv.ParseFloat(match[2], 64)
		lines = append(lines, TranscriptLine{
			Text:     match[3],
			Duration: duration,
			Offset:   offset,
			Lang:     selectedTrack.LanguageCode,
		})
	}

	// 将字幕内容格式化为易读的文本
	var formattedText strings.Builder
	formattedText.WriteString(fmt.Sprintf("标题: %s\n\n", title))
	for _, line := range lines {
		formattedText.WriteString(fmt.Sprintf("[%s] %s\n", formatTime(line.Offset), line.Text))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: formattedText.String(),
			},
		},
	}, nil
}

func formatTime(seconds float64) string {
	minutes := int(seconds) / 60
	remainingSeconds := int(seconds) % 60
	return fmt.Sprintf("%02d:%02d", minutes, remainingSeconds)
}

func extractVideoID(url string) (string, error) {
	patterns := []string{
		`(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&\n?]+)`,
		`youtube\.com\/embed\/([^&\n?]+)`,
		`youtube\.com\/v\/([^&\n?]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}
	return "", fmt.Errorf("invalid YouTube URL")
}
