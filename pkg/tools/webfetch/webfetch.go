package webfetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type WebFetchTool struct {
	tool.BaseTool
	httpClient *http.Client
}

func NewWebFetchTool() *WebFetchTool {
	return &WebFetchTool{
		BaseTool: tool.BaseTool{
			ToolName:          "WebFetch",
			ToolSearchHint:    "fetch download url web page content",
			ToolMaxResultSize: 100000,
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(_ *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

func (t *WebFetchTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"url":    map[string]interface{}{"type": "string", "description": "URL to fetch"},
			"prompt": map[string]interface{}{"type": "string", "description": "What to extract from the page"},
		},
		Required: []string{"url"},
	}
}

func (t *WebFetchTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Fetches a URL and returns its content as markdown.", nil
}

func (t *WebFetchTool) Prompt(_ context.Context) (string, error) {
	return "Fetches content from a URL, converts HTML to markdown, and processes it.\n\nUsage:\n- URL must be fully-formed and valid\n- Prompt describes what to extract from the page\n- Read-only, does not modify files\n- For GitHub URLs, prefer using gh CLI via Bash", nil
}

func (t *WebFetchTool) IsConcurrencySafe(_ map[string]interface{}) bool { return true }
func (t *WebFetchTool) IsReadOnly(_ map[string]interface{}) bool        { return true }

func (t *WebFetchTool) Call(ctx context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	url, _ := input["url"].(string)
	if url == "" {
		return &tool.Result{Data: "Error: url is required"}, nil
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return &tool.Result{Data: "Error: url must start with http:// or https://"}, nil
	}

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error creating request: %v", err)}, nil
	}
	req.Header.Set("User-Agent", "Claude-Code-Go/0.1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,text/plain;q=0.8,*/*;q=0.7")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error fetching URL: %v", err)}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &tool.Result{Data: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)}, nil
	}

	// Read body with size limit (5MB)
	body, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024))
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error reading response: %v", err)}, nil
	}

	content := string(body)
	contentType := resp.Header.Get("Content-Type")

	// Convert HTML to simplified text/markdown
	if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/xhtml") {
		content = htmlToMarkdown(content)
	}

	// Truncate if too long
	const maxChars = 80000
	if len(content) > maxChars {
		content = content[:maxChars] + "\n\n... (truncated)"
	}

	duration := time.Since(start)

	result := fmt.Sprintf("URL: %s\nStatus: %d\nContent-Type: %s\nFetched in: %v\nContent-Length: %d bytes\n\n%s",
		url, resp.StatusCode, contentType, duration.Round(time.Millisecond), len(body), content)

	return &tool.Result{Data: result}, nil
}

// htmlToMarkdown does a basic HTML to text conversion
func htmlToMarkdown(html string) string {
	// Remove script and style tags
	reScript := regexp.MustCompile(`(?is)<script.*?</script>`)
	html = reScript.ReplaceAllString(html, "")
	reStyle := regexp.MustCompile(`(?is)<style.*?</style>`)
	html = reStyle.ReplaceAllString(html, "")

	// Convert common tags
	html = regexp.MustCompile(`(?i)<br\s*/?>|<br>`).ReplaceAllString(html, "\n")
	html = regexp.MustCompile(`(?i)</?p>`).ReplaceAllString(html, "\n\n")
	html = regexp.MustCompile(`(?i)<h[1-6][^>]*>`).ReplaceAllString(html, "\n## ")
	html = regexp.MustCompile(`(?i)</h[1-6]>`).ReplaceAllString(html, "\n")
	html = regexp.MustCompile(`(?i)<li[^>]*>`).ReplaceAllString(html, "\n- ")
	html = regexp.MustCompile(`(?i)<a[^>]*href="([^"]*)"[^>]*>`).ReplaceAllString(html, "[")
	html = regexp.MustCompile(`(?i)</a>`).ReplaceAllString(html, "]")
	html = regexp.MustCompile(`(?i)<code[^>]*>`).ReplaceAllString(html, "`")
	html = regexp.MustCompile(`(?i)</code>`).ReplaceAllString(html, "`")
	html = regexp.MustCompile(`(?i)<pre[^>]*>`).ReplaceAllString(html, "\n```\n")
	html = regexp.MustCompile(`(?i)</pre>`).ReplaceAllString(html, "\n```\n")

	// Remove remaining HTML tags
	html = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(html, "")

	// Decode HTML entities
	html = strings.ReplaceAll(html, "&amp;", "&")
	html = strings.ReplaceAll(html, "&lt;", "<")
	html = strings.ReplaceAll(html, "&gt;", ">")
	html = strings.ReplaceAll(html, "&quot;", "\"")
	html = strings.ReplaceAll(html, "&#39;", "'")
	html = strings.ReplaceAll(html, "&nbsp;", " ")

	// Clean up whitespace
	html = regexp.MustCompile(`\n{3,}`).ReplaceAllString(html, "\n\n")
	return strings.TrimSpace(html)
}
