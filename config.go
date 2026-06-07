package leap0

import (
	"net/http"
	"os"
	"time"

	"github.com/leap0dev/leap0-go/internal/transport"
)

// Option configures the Client.
type Option func(*config)

type config struct {
	apiKey        string
	baseURL       string
	sandboxDomain string
	timeout       time.Duration
	httpClient    *http.Client
}

func WithAPIKey(key string) Option          { return func(c *config) { c.apiKey = key } }
func WithBaseURL(url string) Option         { return func(c *config) { c.baseURL = url } }
func WithSandboxDomain(d string) Option     { return func(c *config) { c.sandboxDomain = d } }
func WithTimeout(d time.Duration) Option    { return func(c *config) { c.timeout = d } }
func WithHTTPClient(hc *http.Client) Option { return func(c *config) { c.httpClient = hc } }

func resolveConfig(opts []Option) (*config, *transport.Config, error) {
	cfg := &config{
		baseURL:       DefaultBaseURL,
		sandboxDomain: DefaultSandboxDomain,
		timeout:       5 * time.Minute,
	}
	if v := os.Getenv("LEAP0_API_KEY"); v != "" {
		cfg.apiKey = v
	}
	if v := os.Getenv("LEAP0_BASE_URL"); v != "" {
		cfg.baseURL = v
	}
	if v := os.Getenv("LEAP0_SANDBOX_DOMAIN"); v != "" {
		cfg.sandboxDomain = v
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.apiKey == "" {
		return nil, nil, &APIError{Body: "API key required (set LEAP0_API_KEY or use WithAPIKey)", sentinel: ErrForbidden}
	}
	if cfg.httpClient == nil {
		cfg.httpClient = &http.Client{Timeout: cfg.timeout}
	}
	tc := &transport.Config{
		BaseURL:       cfg.baseURL,
		SandboxDomain: cfg.sandboxDomain,
		APIKey:        cfg.apiKey,
		AuthHeader:    "Authorization",
		Bearer:        true,
		Source:        sdkSource,
		Version:       Version,
		HTTPClient:    cfg.httpClient,
	}
	return cfg, tc, nil
}
