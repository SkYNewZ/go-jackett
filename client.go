package jackett

//go:generate ifacemaker --file=client.go --struct=client --iface=Client --pkg=jackett -y "Client interface describes wrapped Jackett client." --doc=true --output=generated.go
//go:generate jackett-sync --output=trackers/trackers.go --package trackers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type client struct {
	apiURL     string
	apiKey     string
	httpClient *http.Client
}

// Config describes available config for Jackett.
type Config struct {
	APIUrl     string
	APIKey     string
	HTTPClient *http.Client
}

// New creates a new Jackett client.
func New(config *Config) Client {
	if config == nil {
		config = new(Config)
	}

	if config.APIUrl == "" {
		config.APIUrl = os.Getenv("JACKETT_API_URL")
	}

	if config.APIKey == "" {
		config.APIKey = os.Getenv("JACKETT_API_KEY")
	}

	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}

	return &client{
		apiURL:     config.APIUrl,
		apiKey:     config.APIKey,
		httpClient: config.HTTPClient,
	}
}

// Fetch make a query to Jackett and returns found results.
// categoryID must match an existent category ID from Jackett. You can check them on your Jackett instance's indexers configuration.
// You can specify multiples categories.
// Custom indexer can be specified. Default will be "all".
func (c *client) Fetch(ctx context.Context, options ...Option) (*Response, error) {
	// /api/v2.0/indexers/all/results?apikey=XXX&Query=foo&Category[]=102183&Category[]=102184&Tracker[]=yggtorrent&_=1656077812668
	req, err := c.makeURL(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("error making query: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making Jackett request: %w", err)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("cannot read Jackett response: %w", err)
	}

	return &result, nil
}

func (c *client) makeURL(ctx context.Context, options []Option) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL+"/api/v2.0/indexers/all/results", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("cannot make request: %w", err)
	}

	q := req.URL.Query()
	q.Set("apikey", c.apiKey)
	for _, opt := range options {
		opt(&q)
	}

	req.URL.RawQuery = q.Encode()
	return req, nil
}
