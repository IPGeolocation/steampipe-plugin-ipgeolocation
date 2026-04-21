package ipgeolocation

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const baseURL = "https://api.ipgeolocation.io/v3"

// Client wraps an http.Client with the API key.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// newClient builds a Client from the plugin connection.
// Resolution order for the API key:
//  1. api_key in the .spc connection config
//  2. IPGEOLOCATION_API_KEY environment variable
func newClient(ctx context.Context, d *plugin.QueryData) (*Client, error) {
	cfg := GetConfig(d.Connection)

	apiKey := ""
	if cfg.APIKey != nil && *cfg.APIKey != "" {
		apiKey = *cfg.APIKey
	} else if envKey := os.Getenv("IPGEOLOCATION_API_KEY"); envKey != "" {
		apiKey = envKey
		plugin.Logger(ctx).Debug("ipgeolocation: using API key from IPGEOLOCATION_API_KEY env var")
	}

	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// GetSecurity calls GET /v3/security and returns raw JSON decoded into a map.
// Costs 2 credits per request. Requires a paid plan API key.
func (c *Client) GetSecurity(ctx context.Context, ip string) (map[string]interface{}, error) {
	params := url.Values{}
	if c.apiKey != "" {
		params.Set("apiKey", c.apiKey)
	}
	if ip != "" {
		params.Set("ip", ip)
	}
	endpoint := fmt.Sprintf("%s/security?%s", baseURL, params.Encode())
	return c.doGet(ctx, endpoint)
}

// GetAbuse calls GET /v3/abuse and returns raw JSON decoded into a map.
// Costs 1 credit per request. Requires a paid plan API key.
func (c *Client) GetAbuse(ctx context.Context, ip string) (map[string]interface{}, error) {
	params := url.Values{}
	if c.apiKey != "" {
		params.Set("apiKey", c.apiKey)
	}
	if ip != "" {
		params.Set("ip", ip)
	}
	endpoint := fmt.Sprintf("%s/abuse?%s", baseURL, params.Encode())
	return c.doGet(ctx, endpoint)
}

// GetASN calls GET /v3/asn and returns raw JSON decoded into a map.
// Accepts either an IP address (ip param) or an ASN number (asn param).
// Costs 1 credit per request.
func (c *Client) GetASN(ctx context.Context, ip, asn string) (map[string]interface{}, error) {
	params := url.Values{}
	if c.apiKey != "" {
		params.Set("apiKey", c.apiKey)
	}
	if ip != "" {
		params.Set("ip", ip)
	}
	if asn != "" {
		params.Set("asn", asn)
	}
	// Always include all optional relationship fields
	params.Set("include", "peers,upstreams,downstreams,routes,whois_response")
	endpoint := fmt.Sprintf("%s/asn?%s", baseURL, params.Encode())
	return c.doGet(ctx, endpoint)
}

// doGet is a shared helper that executes a GET request and decodes JSON.
func (c *Client) doGet(ctx context.Context, endpoint string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return result, nil
}

// GetIPGeo calls GET /v3/ipgeo and returns raw JSON decoded into a map.
// ip may be empty to geolocate the caller's IP.
func (c *Client) GetIPGeo(ctx context.Context, ip string) (map[string]interface{}, error) {
	params := url.Values{}
	if c.apiKey != "" {
		params.Set("apiKey", c.apiKey)
	}
	if ip != "" {
		params.Set("ip", ip)
	}
	params.Set("include", "*")
	endpoint := fmt.Sprintf("%s/ipgeo?%s", baseURL, params.Encode())
	return c.doGet(ctx, endpoint)
}
