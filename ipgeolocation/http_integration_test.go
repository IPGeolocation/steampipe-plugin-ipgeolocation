package ipgeolocation

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func TestHTTPSIntegration(t *testing.T) {
	if os.Getenv("STEAMPIPE_HTTP_TEST") != "1" {
		t.Skip("run make test-http with Postman CLI installed")
	}
	if os.Getenv("HTTPS_PROXY") == "" {
		t.Fatal("the local proxy is required")
	}
	certificate, err := os.ReadFile(os.Getenv("SSL_CERT_FILE"))
	if err != nil {
		t.Fatal(err)
	}
	roots, err := x509.SystemCertPool()
	if err != nil {
		t.Fatal(err)
	}
	if !roots.AppendCertsFromPEM(certificate) {
		t.Fatal("invalid local proxy certificate")
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
	original := http.DefaultTransport
	http.DefaultTransport = transport
	defer func() { transport.CloseIdleConnections(); http.DefaultTransport = original }()

	ctx := testContext()
	for _, tc := range []struct {
		name, table, ip, asn string
		want                 map[string]interface{}
	}{
		{"IPv4 geolocation", "ipgeolocation_ip", "198.51.100.42", "", map[string]interface{}{"ip": "198.51.100.42", "city": "New Delhi", "timezone_offset": 5.5, "timezone_offset_with_dst": 5.5}},
		{"IPv6 geolocation", "ipgeolocation_ip", "2001:db8::1", "", map[string]interface{}{"ip": "2001:db8::1", "timezone_offset": 5.5}},
		{"security", "ipgeolocation_security", "198.51.100.42", "", map[string]interface{}{"ip": "198.51.100.42", "is_proxy": false}},
		{"abuse contacts", "ipgeolocation_abuse", "198.51.100.42", "", map[string]interface{}{"ip": "198.51.100.42", "emails": []interface{}{"abuse@example.test"}}},
		{"ASN by IP", "ipgeolocation_asn", "198.51.100.42", "", map[string]interface{}{"as_number": "AS15169"}},
		{"ASN prefix", "ipgeolocation_asn", "", " as15169 ", map[string]interface{}{"as_number": "AS15169"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			key := "TEST_KEY"
			rows := []interface{}{}
			d := &plugin.QueryData{
				Connection:     &plugin.Connection{Config: IPGeolocationConfig{APIKey: &key}},
				EqualsQuals:    plugin.KeyColumnEqualsQualMap{},
				StreamListItem: func(_ context.Context, values ...interface{}) { rows = append(rows, values...) },
			}
			for name, value := range map[string]string{"ip": tc.ip, "asn": tc.asn} {
				if value != "" {
					d.EqualsQuals[name] = &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: value}}
				}
			}
			table := Plugin(ctx).TableMap[tc.table]
			if _, err := table.List.Hydrate(ctx, d, nil); err != nil {
				t.Fatal(err)
			}
			if len(rows) != 1 {
				t.Fatalf("expected one result row, got %d", len(rows))
			}
			for name, want := range tc.want {
				found := false
				for _, column := range table.Columns {
					if column.Name != name {
						continue
					}
					found = true
					got, err := column.Transform.Execute(ctx, &transform.TransformData{HydrateItem: rows[0]})
					if err != nil || !reflect.DeepEqual(got, want) {
						t.Fatalf("%s=%v, want %v, error=%v", name, got, want, err)
					}
				}
				if !found {
					t.Fatalf("missing column %s", name)
				}
			}
		})
	}
	t.Run("keyless caller IP", func(t *testing.T) {
		c := &Client{httpClient: &http.Client{Timeout: time.Second}}
		got, err := c.GetIPGeo(ctx, "")
		if err != nil || got["ip"] != "198.51.100.42" {
			t.Fatalf("response=%v error=%v", got, err)
		}
	})
	for _, tc := range []struct{ key, message string }{
		{"BAD_KEY", "API error 401"}, {"FORBIDDEN", "API error 403"},
		{"RATE_LIMIT", "API error 429"}, {"SERVER_ERROR", "API error 503"},
		{"BAD_JSON", "decoding response"},
	} {
		t.Run(tc.key, func(t *testing.T) {
			key, rows := tc.key, 0
			d := &plugin.QueryData{
				Connection:     &plugin.Connection{Config: IPGeolocationConfig{APIKey: &key}},
				EqualsQuals:    plugin.KeyColumnEqualsQualMap{"ip": &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: "198.51.100.42"}}},
				StreamListItem: func(_ context.Context, values ...interface{}) { rows += len(values) },
			}
			_, err := listIPSecurity(ctx, d, nil)
			if rows != 0 || err == nil || !strings.Contains(err.Error(), tc.message) {
				t.Fatalf("rows=%d error=%v", rows, err)
			}
		})
	}
	t.Run("timeout", func(t *testing.T) {
		c := &Client{apiKey: "SLOW", httpClient: &http.Client{Timeout: 25 * time.Millisecond}}
		_, err := c.GetIPGeo(ctx, "198.51.100.42")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected timeout, got %v", err)
		}
	})
}
