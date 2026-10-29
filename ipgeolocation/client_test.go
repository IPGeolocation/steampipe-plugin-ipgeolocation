package ipgeolocation

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/context_key"
)

func testContext() context.Context {
	return context.WithValue(context.Background(), context_key.Logger, hclog.NewNullLogger())
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestClientRequests(t *testing.T) {
	for _, tc := range []struct {
		name, path string
		params     url.Values
		call       func(*Client) (map[string]interface{}, error)
	}{
		{"geolocation", "/v3/ipgeo", url.Values{"apiKey": {"test & key"}, "ip": {"2001:4860:4860::8888"}, "include": {"*"}}, func(c *Client) (map[string]interface{}, error) {
			return c.GetIPGeo(context.Background(), "2001:4860:4860::8888")
		}},
		{"caller IP", "/v3/ipgeo", url.Values{"include": {"*"}}, func(c *Client) (map[string]interface{}, error) {
			c.apiKey = ""
			return c.GetIPGeo(context.Background(), "")
		}},
		{"security", "/v3/security", url.Values{"apiKey": {"test & key"}, "ip": {"8.8.8.8"}}, func(c *Client) (map[string]interface{}, error) { return c.GetSecurity(context.Background(), "8.8.8.8") }},
		{"abuse", "/v3/abuse", url.Values{"apiKey": {"test & key"}, "ip": {"8.8.8.8"}}, func(c *Client) (map[string]interface{}, error) { return c.GetAbuse(context.Background(), "8.8.8.8") }},
		{"ASN by IP", "/v3/asn", url.Values{"apiKey": {"test & key"}, "ip": {"8.8.8.8"}, "include": {"peers,upstreams,downstreams,routes,whois_response"}}, func(c *Client) (map[string]interface{}, error) { return c.GetASN(context.Background(), "8.8.8.8", "") }},
		{"ASN number", "/v3/asn", url.Values{"apiKey": {"test & key"}, "asn": {"15169"}, "include": {"peers,upstreams,downstreams,routes,whois_response"}}, func(c *Client) (map[string]interface{}, error) { return c.GetASN(context.Background(), "", "15169") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			c := &Client{apiKey: "test & key", httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.Method != "GET" || r.URL.Scheme != "https" || r.URL.Host != "api.ipgeolocation.io" || r.URL.Path != tc.path || r.URL.Query().Encode() != tc.params.Encode() {
					t.Fatalf("unexpected request: %s %s", r.Method, r.URL)
				}
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"ip":"8.8.8.8","security":{"is_proxy":false}}`)), Header: make(http.Header)}, nil
			})}}
			got, err := tc.call(c)
			if err != nil || got["ip"] != "8.8.8.8" || calls != 1 {
				t.Fatalf("result=%v calls=%d error=%v", got, calls, err)
			}
		})
	}
}

func TestClientResponseHandling(t *testing.T) {
	for _, tc := range []struct {
		name            string
		status          int
		body, wantError string
	}{
		{"success", 200, `{"ip":"8.8.8.8"}`, ""},
		{"authentication", 401, `{"message":"Invalid API key"}`, "API error 401"},
		{"forbidden", 403, `{"message":"Access denied"}`, "API error 403"},
		{"rate limit", 429, `{"message":"Too many requests"}`, "API error 429"},
		{"server error", 503, `{"message":"Unavailable"}`, "API error 503"},
		{"invalid JSON", 200, `{`, "decoding response"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}))
			defer server.Close()
			c := &Client{httpClient: server.Client()}
			got, err := c.doGet(context.Background(), server.URL)
			if tc.wantError == "" {
				if err != nil || got["ip"] != "8.8.8.8" {
					t.Fatalf("result=%v error=%v", got, err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.wantError) || got != nil {
				t.Fatalf("result=%v error=%v", got, err)
			}
		})
	}
	t.Run("truncated response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "100")
			_, _ = io.WriteString(w, `{}`)
		}))
		defer server.Close()
		_, err := (&Client{httpClient: server.Client()}).doGet(context.Background(), server.URL)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("expected truncated-body error, got %v", err)
		}
	})
}

func TestClientTransportErrorRedactsKey(t *testing.T) {
	cause := errors.New("connection failed")
	c := &Client{apiKey: "test-only-secret", httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) { return nil, cause })}}
	_, err := c.GetIPGeo(context.Background(), "8.8.8.8")
	if !errors.Is(err, cause) || strings.Contains(err.Error(), c.apiKey) {
		t.Fatalf("transport error must retain its cause without the API key: %v", err)
	}
}

func TestClientCancellation(t *testing.T) {
	c := &Client{apiKey: "test-only-secret", httpClient: &http.Client{}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.GetIPGeo(ctx, "8.8.8.8")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func TestAPIKeyConfiguration(t *testing.T) {
	t.Setenv("IPGEOLOCATION_API_KEY", "environment-key")
	configured := "configured-key"
	empty := ""
	for _, tc := range []struct {
		name   string
		config interface{}
		want   string
	}{
		{"configured key wins", IPGeolocationConfig{APIKey: &configured}, configured},
		{"empty key uses environment", IPGeolocationConfig{APIKey: &empty}, "environment-key"},
		{"missing config uses environment", nil, "environment-key"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, err := newClient(testContext(), &plugin.QueryData{Connection: &plugin.Connection{Config: tc.config}})
			if err != nil || c.apiKey != tc.want || c.httpClient.Timeout != 30*time.Second {
				t.Fatalf("client=%+v error=%v", c, err)
			}
		})
	}
	t.Setenv("IPGEOLOCATION_API_KEY", "")
	c, err := newClient(context.Background(), &plugin.QueryData{})
	if err != nil || c.apiKey != "" {
		t.Fatalf("caller-IP client=%+v error=%v", c, err)
	}
}

func TestTLSVerification(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, `{"ip":"198.51.100.42"}`) }))
	defer server.Close()
	_, err := (&Client{httpClient: &http.Client{}}).doGet(context.Background(), server.URL)
	var certificateError *tls.CertificateVerificationError
	if !errors.As(err, &certificateError) {
		t.Fatalf("must reject an untrusted certificate: %v", err)
	}
	result, err := (&Client{httpClient: server.Client()}).doGet(context.Background(), server.URL)
	if err != nil || result["ip"] != "198.51.100.42" {
		t.Fatalf("trusted TLS response=%v error=%v", result, err)
	}
}
