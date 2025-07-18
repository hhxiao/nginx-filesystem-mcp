package pkg

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"log"
	"net/url"
	"os"
	"path"
	"resty.dev/v3"
	"strings"
	"time"
)

var contentCache = cache.New(5*time.Minute, 10*time.Minute)

type Content struct {
	Data   string
	Type   string
	Length string
}

type Client struct {
	baseURL            string
	insecureSkipVerify bool
}

type Option func(*Client)

func NewClient(opts ...Option) *Client {
	client := &Client{
		baseURL:            os.Getenv("CONTENT_BASEURL"),
		insecureSkipVerify: os.Getenv("INSECURE_SKIP_VERIFY") == "true",
	}
	for _, opt := range opts {
		opt(client)
	}
	if client.baseURL == "" {
		panic("CONTENT_BASEURL environment variable not set")
	}
	log.Printf("baseURL: %s", client.baseURL)
	return client
}

// WithBaseURL sets the base URL
func WithBaseURL(baseURL string) Option {
	return func(s *Client) {
		s.baseURL = baseURL
	}
}

// WithInsecureSkipVerify sets the insecureSkipVerify flag
func WithInsecureSkipVerify(insecureSkipVerify bool) Option {
	return func(s *Client) {
		s.insecureSkipVerify = insecureSkipVerify
	}
}

func (p *Client) GetContent(ctx context.Context, filepath string) (*Content, error) {
	filepath = strings.TrimPrefix(filepath, "/")
	rawUrl := fmt.Sprintf("%s/%s", p.baseURL, filepath)
	u, _ := url.Parse(rawUrl)
	u.Path = path.Clean(u.Path)
	rawUrl = u.String()
	if !strings.HasPrefix(rawUrl, p.baseURL) {
		rawUrl = p.baseURL
	}

	if content, found := contentCache.Get(rawUrl); found {
		return content.(*Content), nil
	}

	r := resty.New().EnableDebug().SetDebugLogFormatter(debugLogCustomFormatter)
	defer r.Close()

	if val := ctx.Value("Authorization"); val != nil {
		if auth, ok := val.(string); ok {
			r.SetHeader("Authorization", auth)
		}
	}
	if p.insecureSkipVerify {
		r.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := r.R().
		SetContext(ctx).
		Get(rawUrl)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, errors.New(res.Status())
	}
	content := Content{
		Data:   res.String(),
		Type:   res.Header().Get("Content-Type"),
		Length: res.Header().Get("Content-Length"),
	}
	contentCache.Set(rawUrl, &content, cache.DefaultExpiration)
	return &content, nil
}

func debugLogCustomFormatter(dl *resty.DebugLog) string {
	return fmt.Sprintf("%s %s ... %s", dl.Request.Method, dl.Request.URI, dl.Response.Status)
}