package ipgeolocation

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// tableASN returns the Steampipe table definition for ipgeolocation_asn.
func tableASN(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "ipgeolocation_asn",
		Description: "ASN lookup using the IPGeolocation.io /v3/asn endpoint. " +
			"Filter by ip (any IPv4/IPv6) or asn (e.g. '15169' or 'AS15169'). " +
			"With no filter the caller's own public IP ASN is returned. " +
			"Each query costs 1 API credit.",
		List: &plugin.ListConfig{
			// 'ip' and 'asn' are pure qual-only input columns — they are NOT
			// mapped to any response field, which is what allows Steampipe to
			// pass them through to the hydrate function cleanly.
			KeyColumns: plugin.OptionalColumns([]string{"ip", "asn"}),
			Hydrate:    listASN,
		},
		Columns: asnColumns(),
	}
}
func asnColumns() []*plugin.Column {
	return []*plugin.Column{
		// ── Qual-only input columns ──────────────────────────────────────────
		// These are the WHERE clause handles. They carry no transform because
		// they are inputs, not outputs. Steampipe will set them to null in rows
		// where the API did not echo them back (e.g. plain ASN lookup has no
		// top-level "ip" field).
		{
			Name:        "ip",
			Type:        proto.ColumnType_STRING,
			Description: "IP address to look up (qual only — use in WHERE clause). The API returns the ASN that owns this IP.",
			Transform:   transform.FromField("ip"),
		},
		{
			Name:        "asn",
			Type:        proto.ColumnType_STRING,
			Description: "ASN to look up (qual only — use in WHERE clause). Accepts '15169', 'AS15169', or 'as15169'.",
			// Not present in the API response; populated from the qual value
			// by the hydrate function so the column is non-null when queried.
			Transform: transform.FromQual("asn"),
		},

		// ── Core ASN response fields ─────────────────────────────────────────
		{
			Name:        "as_number",
			Type:        proto.ColumnType_STRING,
			Description: "Autonomous System Number as returned by the API (e.g. \"AS15169\").",
			Transform:   transform.FromField("asn.as_number"),
		},
		{
			Name:        "asn_name",
			Type:        proto.ColumnType_STRING,
			Description: "Short registered name for the ASN (e.g. \"GOOGLE\").",
			Transform:   transform.FromField("asn.asn_name"),
		},
		{
			Name:        "organization",
			Type:        proto.ColumnType_STRING,
			Description: "Organisation that owns the ASN.",
			Transform:   transform.FromField("asn.organization"),
		},
		{
			Name:        "country",
			Type:        proto.ColumnType_STRING,
			Description: "ISO alpha-2 country code where the ASN is registered.",
			Transform:   transform.FromField("asn.country"),
		},
		{
			Name:        "type",
			Type:        proto.ColumnType_STRING,
			Description: "ASN classification: ISP, HOSTING, EDUCATION, GOVERNMENT, or BUSINESS.",
			Transform:   transform.FromField("asn.type"),
		},
		{
			Name:        "domain",
			Type:        proto.ColumnType_STRING,
			Description: "Primary domain associated with the ASN organisation.",
			Transform:   transform.FromField("asn.domain"),
		},
		{
			Name:        "date_allocated",
			Type:        proto.ColumnType_STRING,
			Description: "Date the ASN was allocated by the RIR (YYYY-MM-DD).",
			Transform:   transform.FromField("asn.date_allocated"),
		},
		{
			Name:        "allocation_status",
			Type:        proto.ColumnType_STRING,
			Description: "Allocation status of the ASN (e.g. \"assigned\").",
			Transform:   transform.FromField("asn.allocation_status"),
		},
		{
			Name:        "rir",
			Type:        proto.ColumnType_STRING,
			Description: "Regional Internet Registry (ARIN, RIPE, APNIC, LACNIC, AFRINIC).",
			Transform:   transform.FromField("asn.rir"),
		},

		// ── Route counts ─────────────────────────────────────────────────────
		{
			Name:        "num_of_ipv4_routes",
			Type:        proto.ColumnType_INT,
			Description: "Number of IPv4 prefixes announced by this ASN.",
			Transform:   transform.FromField("asn.num_of_ipv4_routes"),
		},
		{
			Name:        "num_of_ipv6_routes",
			Type:        proto.ColumnType_INT,
			Description: "Number of IPv6 prefixes announced by this ASN.",
			Transform:   transform.FromField("asn.num_of_ipv6_routes"),
		},

		// ── Relationship arrays ───────────────────────────────────────────────
		{
			Name:        "routes",
			Type:        proto.ColumnType_JSON,
			Description: "Array of CIDR prefixes announced by this ASN.",
			Transform:   transform.FromField("asn.routes"),
		},
		{
			Name:        "peers",
			Type:        proto.ColumnType_JSON,
			Description: "Array of peer ASNs ({as_number, description, country}).",
			Transform:   transform.FromField("asn.peers"),
		},
		{
			Name:        "upstreams",
			Type:        proto.ColumnType_JSON,
			Description: "Array of upstream/transit provider ASNs ({as_number, description, country}).",
			Transform:   transform.FromField("asn.upstreams"),
		},
		{
			Name:        "downstreams",
			Type:        proto.ColumnType_JSON,
			Description: "Array of downstream customer ASNs ({as_number, description, country}).",
			Transform:   transform.FromField("asn.downstreams"),
		},

		// ── WHOIS ────────────────────────────────────────────────────────────
		{
			Name:        "whois_response",
			Type:        proto.ColumnType_STRING,
			Description: "Raw WHOIS text for the ASN.",
			Transform:   transform.FromField("asn.whois_response"),
		},

		// ── Raw JSON ─────────────────────────────────────────────────────────
		{
			Name:        "raw",
			Type:        proto.ColumnType_JSON,
			Description: "Full raw JSON response from the /v3/asn API endpoint.",
			Transform:   transform.FromValue(),
		},
	}
}

// listASN is the hydrate function for ipgeolocation_asn.
//
// Supported query patterns:
//
//	select * from ipgeolocation_asn where ip  = '8.8.8.8';
//	select * from ipgeolocation_asn where asn = '15169';
//	select * from ipgeolocation_asn where asn = 'AS15169';
//	select * from ipgeolocation_asn;   -- returns caller's IP ASN
func listASN(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := newClient(ctx, d)
	if err != nil {
		return nil, err
	}

	ip := ""
	asnVal := ""

	if q := d.EqualsQuals["ip"]; q != nil {
		ip = q.GetStringValue()
	}
	if q := d.EqualsQuals["asn"]; q != nil {
		// normalise: "15169" / "AS15169" / "as15169" → "15169"
		asnVal = normaliseASN(q.GetStringValue())
	}

	data, err := client.GetASN(ctx, ip, asnVal)
	if err != nil {
		plugin.Logger(ctx).Error("ipgeolocation_asn.listASN", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, data)
	return nil, nil
}

// normaliseASN strips the optional "AS" prefix so both "15169" and "AS15169"
// are accepted by the /v3/asn endpoint.
func normaliseASN(s string) string {
	upper := strings.ToUpper(strings.TrimSpace(s))
	if strings.HasPrefix(upper, "AS") {
		return upper[2:]
	}
	return strings.TrimSpace(s)
}
