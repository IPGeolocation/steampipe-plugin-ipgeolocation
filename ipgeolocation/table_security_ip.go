package ipgeolocation

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// tableIPSecurity returns the Steampipe table definition for ipgeolocation_security.
func tableIPSecurity(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "ipgeolocation_security",
		Description: "IP threat intelligence using the IPGeolocation.io /v3/security endpoint. " +
			"Returns VPN, proxy, Tor, relay, bot, spam, attacker, and cloud-provider signals " +
			"with confidence scores, provider names, and last-seen dates. Requires a paid plan. " +
			"Each query costs 2 API credits.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"ip"}),
			Hydrate:    listIPSecurity,
		},
		Columns: ipSecurityColumns(),
	}
}

func ipSecurityColumns() []*plugin.Column {
	return []*plugin.Column{
		// ── Top-level ───────────────────────────────────────────────────────
		{
			Name:        "ip",
			Type:        proto.ColumnType_STRING,
			Description: "The IP address that was looked up.",
			Transform:   transform.FromField("ip"),
		},
		// ── Threat overview ─────────────────────────────────────────────────
		{
			Name:        "threat_score",
			Type:        proto.ColumnType_INT,
			Description: "Overall threat score from 0 (clean) to 100 (high risk). Summarises all security signals.",
			Transform:   transform.FromField("security.threat_score"),
		},
		{
			Name:        "is_anonymous",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is considered anonymous (VPN, proxy, Tor, or relay detected).",
			Transform:   transform.FromField("security.is_anonymous"),
		},

		// ── VPN ─────────────────────────────────────────────────────────────
		{
			Name:        "is_vpn",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is a known VPN exit node.",
			Transform:   transform.FromField("security.is_vpn"),
		},
		{
			Name:        "vpn_provider_names",
			Type:        proto.ColumnType_JSON,
			Description: "Array of VPN provider names identified for this IP (e.g. [\"Nord VPN\"]).",
			Transform:   transform.FromField("security.vpn_provider_names"),
		},
		{
			Name:        "vpn_confidence_score",
			Type:        proto.ColumnType_INT,
			Description: "Confidence score (0–100) for the VPN detection.",
			Transform:   transform.FromField("security.vpn_confidence_score"),
		},
		{
			Name:        "vpn_last_seen",
			Type:        proto.ColumnType_STRING,
			Description: "Date the IP was last observed as a VPN exit node (YYYY-MM-DD).",
			Transform:   transform.FromField("security.vpn_last_seen"),
		},

		// ── Proxy ───────────────────────────────────────────────────────────
		{
			Name:        "is_proxy",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is a known proxy.",
			Transform:   transform.FromField("security.is_proxy"),
		},
		{
			Name:        "proxy_provider_names",
			Type:        proto.ColumnType_JSON,
			Description: "Array of proxy provider names identified for this IP (e.g. [\"Zyte Proxy\"]).",
			Transform:   transform.FromField("security.proxy_provider_names"),
		},
		{
			Name:        "proxy_confidence_score",
			Type:        proto.ColumnType_INT,
			Description: "Confidence score (0–100) for the proxy detection.",
			Transform:   transform.FromField("security.proxy_confidence_score"),
		},
		{
			Name:        "proxy_last_seen",
			Type:        proto.ColumnType_STRING,
			Description: "Date the IP was last observed as a proxy (YYYY-MM-DD).",
			Transform:   transform.FromField("security.proxy_last_seen"),
		},

		// ── Residential proxy ───────────────────────────────────────────────
		{
			Name:        "is_residential_proxy",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is a residential proxy (harder to detect and block than datacenter proxies).",
			Transform:   transform.FromField("security.is_residential_proxy"),
		},

		// ── Tor ─────────────────────────────────────────────────────────────
		{
			Name:        "is_tor",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is a Tor exit node.",
			Transform:   transform.FromField("security.is_tor"),
		},

		// ── Relay ───────────────────────────────────────────────────────────
		{
			Name:        "is_relay",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is part of a relay network (e.g. Apple Private Relay, iCloud Relay).",
			Transform:   transform.FromField("security.is_relay"),
		},
		{
			Name:        "relay_provider_name",
			Type:        proto.ColumnType_STRING,
			Description: "Name of the relay provider if is_relay is true.",
			Transform:   transform.FromField("security.relay_provider_name"),
		},

		// ── Cloud / hosting ─────────────────────────────────────────────────
		{
			Name:        "is_cloud_provider",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP belongs to a cloud or hosting provider.",
			Transform:   transform.FromField("security.is_cloud_provider"),
		},
		{
			Name:        "cloud_provider_name",
			Type:        proto.ColumnType_STRING,
			Description: "Name of the cloud/hosting provider (e.g. \"Amazon AWS\", \"Packethub S.A.\").",
			Transform:   transform.FromField("security.cloud_provider_name"),
		},

		// ── Attacker / abuse signals ─────────────────────────────────────────
		{
			Name:        "is_known_attacker",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP has been flagged for brute force, credential stuffing, or scanning activity.",
			Transform:   transform.FromField("security.is_known_attacker"),
		},
		{
			Name:        "is_bot",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is associated with automated bot activity.",
			Transform:   transform.FromField("security.is_bot"),
		},
		{
			Name:        "is_spam",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP appears on spam block lists.",
			Transform:   transform.FromField("security.is_spam"),
		},

		// ── Raw JSON ────────────────────────────────────────────────────────
		{
			Name:        "raw",
			Type:        proto.ColumnType_JSON,
			Description: "Full raw JSON response from the /v3/security API endpoint.",
			Transform:   transform.FromValue(),
		},
	}
}

// listIPSecurity is the hydrate function for ipgeolocation_security.
func listIPSecurity(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := newClient(ctx, d)
	if err != nil {
		return nil, err
	}

	ip := ""
	if d.EqualsQuals["ip"] != nil {
		ip = d.EqualsQuals["ip"].GetStringValue()
	}

	data, err := client.GetSecurity(ctx, ip)
	if err != nil {
		plugin.Logger(ctx).Error("ipgeolocation_security.listIPSecurity", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, data)
	return nil, nil
}
