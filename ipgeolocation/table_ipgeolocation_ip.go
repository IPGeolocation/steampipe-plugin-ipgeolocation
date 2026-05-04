package ipgeolocation

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableIpgeolocationIp(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ipgeolocation_ip",
		Description: "IP geolocation lookup using the IPGeolocation.io v3 API. Query a single IP address or domain (paid only)  to retrieve location, network, security, timezone, company, abuse, and hostname data.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "ip",
					Require: plugin.Required,
				},
			},
			Hydrate: listIPGeolocation,
		},
		Columns: ipGeolocationColumns(),
	}
}

// ipGeolocationColumns declares all columns exposed by the table.
func ipGeolocationColumns() []*plugin.Column {
	return []*plugin.Column{
		// ── Core ────────────────────────────────────────────────────────────
		{
			Name:        "ip",
			Type:        proto.ColumnType_STRING,
			Description: "The IP address that was looked up.",
			Transform:   transform.FromField("ip"),
		},
		{
			Name:        "hostname",
			Type:        proto.ColumnType_STRING,
			Description: "Resolved hostname for the IP address (paid plan).",
			Transform:   transform.FromField("hostname"),
		},

		// ── Location ────────────────────────────────────────────────────────
		{
			Name:        "continent_code",
			Type:        proto.ColumnType_STRING,
			Description: "Two-letter continent code (e.g. NA, EU).",
			Transform:   transform.FromField("location.continent_code"),
		},
		{
			Name:        "continent_name",
			Type:        proto.ColumnType_STRING,
			Description: "Full continent name.",
			Transform:   transform.FromField("location.continent_name"),
		},
		{
			Name:        "country_code2",
			Type:        proto.ColumnType_STRING,
			Description: "ISO 3166-1 alpha-2 country code.",
			Transform:   transform.FromField("location.country_code2"),
		},
		{
			Name:        "country_code3",
			Type:        proto.ColumnType_STRING,
			Description: "ISO 3166-1 alpha-3 country code.",
			Transform:   transform.FromField("location.country_code3"),
		},
		{
			Name:        "country_name",
			Type:        proto.ColumnType_STRING,
			Description: "Country name in English.",
			Transform:   transform.FromField("location.country_name"),
		},
		{
			Name:        "country_name_official",
			Type:        proto.ColumnType_STRING,
			Description: "Official country name.",
			Transform:   transform.FromField("location.country_name_official"),
		},
		{
			Name:        "country_capital",
			Type:        proto.ColumnType_STRING,
			Description: "Capital city of the country.",
			Transform:   transform.FromField("location.country_capital"),
		},
		{
			Name:        "state_prov",
			Type:        proto.ColumnType_STRING,
			Description: "State or province name.",
			Transform:   transform.FromField("location.state_prov"),
		},
		{
			Name:        "state_code",
			Type:        proto.ColumnType_STRING,
			Description: "ISO 3166-2 state/province code.",
			Transform:   transform.FromField("location.state_code"),
		},
		{
			Name:        "district",
			Type:        proto.ColumnType_STRING,
			Description: "District or county.",
			Transform:   transform.FromField("location.district"),
		},
		{
			Name:        "city",
			Type:        proto.ColumnType_STRING,
			Description: "City name.",
			Transform:   transform.FromField("location.city"),
		},
		{
			Name:        "locality",
			Type:        proto.ColumnType_STRING,
			Description: "Locality name.",
			Transform:   transform.FromField("location.locality"),
		},
		{
			Name:        "zipcode",
			Type:        proto.ColumnType_STRING,
			Description: "ZIP / postal code.",
			Transform:   transform.FromField("location.zipcode"),
		},
		{
			Name:        "latitude",
			Type:        proto.ColumnType_STRING,
			Description: "Latitude coordinate.",
			Transform:   transform.FromField("location.latitude"),
		},
		{
			Name:        "longitude",
			Type:        proto.ColumnType_STRING,
			Description: "Longitude coordinate.",
			Transform:   transform.FromField("location.longitude"),
		},
		{
			Name:        "accuracy_radius",
			Type:        proto.ColumnType_STRING,
			Description: "Accuracy radius in kilometres.",
			Transform:   transform.FromField("location.accuracy_radius"),
		},
		{

			Name:        "confidence",
			Type:        proto.ColumnType_STRING,
			Description: "Confidence level (low, medium, high).",
			Transform:   transform.FromField("location.confidence"),
		},
		{
			Name:        "is_eu",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is in a European Union member state.",
			Transform:   transform.FromField("location.is_eu"),
		},
		{
			Name:        "country_flag",
			Type:        proto.ColumnType_STRING,
			Description: "URL of the country flag image.",
			Transform:   transform.FromField("location.country_flag"),
		},
		{
			Name:        "country_emoji",
			Type:        proto.ColumnType_STRING,
			Description: "Emoji flag for the country.",
			Transform:   transform.FromField("location.country_emoji"),
		},
		{
			Name:        "geoname_id",
			Type:        proto.ColumnType_STRING,
			Description: "GeoNames database ID for the city.",
			Transform:   transform.FromField("location.geoname_id"),
		},
		{
			Name:        "dma_code",
			Type:        proto.ColumnType_STRING,
			Description: "Designated Market Area (DMA) code (US only).",
			Transform:   transform.FromField("location.dma_code"),
		},

		// ── Network / ASN ───────────────────────────────────────────────────
		{
			Name:        "asn",
			Type:        proto.ColumnType_STRING,
			Description: "Autonomous System Number.",
			Transform:   transform.FromField("asn.as_number"),
		},
		{
			Name:        "organization",
			Type:        proto.ColumnType_STRING,
			Description: "Organization registered to the IP.",
			Transform:   transform.FromField("asn.organization"),
		},
		{
			Name:        "connection_type",
			Type:        proto.ColumnType_STRING,
			Description: "Connection type (e.g. Corporate, ISP, Hosting).",
			Transform:   transform.FromField("network.connection_type"),
		},

		// ── Timezone ────────────────────────────────────────────────────────
		{
			Name:        "timezone_name",
			Type:        proto.ColumnType_STRING,
			Description: "IANA timezone name (e.g. America/New_York).",
			Transform:   transform.FromField("time_zone.name"),
		},
		{
			Name:        "timezone_offset",
			Type:        proto.ColumnType_INT,
			Description: "UTC offset in seconds.",
			Transform:   transform.FromField("time_zone.offset"),
		},
		{
			Name:        "timezone_offset_with_dst",
			Type:        proto.ColumnType_INT,
			Description: "UTC offset including DST adjustment in seconds.",
			Transform:   transform.FromField("time_zone.offset_with_dst"),
		},
		{
			Name:        "timezone_current_time",
			Type:        proto.ColumnType_STRING,
			Description: "Current local time in the timezone.",
			Transform:   transform.FromField("time_zone.current_time"),
		},
		{
			Name:        "timezone_is_dst",
			Type:        proto.ColumnType_BOOL,
			Description: "True if daylight saving time is currently in effect.",
			Transform:   transform.FromField("time_zone.is_dst"),
		},

		// ── Security ────────────────────────────────────────────────────────
		{
			Name:        "threat_score",
			Type:        proto.ColumnType_INT,
			Description: "Threat score from 0 (clean) to 100 (high risk) (paid plan).",
			Transform:   transform.FromField("security.threat_score"),
		},
		{
			Name:        "is_vpn",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is a known VPN exit node (paid plan).",
			Transform:   transform.FromField("security.is_vpn"),
		},
		{
			Name:        "is_proxy",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is a known proxy (paid plan).",
			Transform:   transform.FromField("security.is_proxy"),
		},
		{
			Name:        "is_tor",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is a Tor exit node (paid plan).",
			Transform:   transform.FromField("security.is_tor"),
		},
		{
			Name:        "is_bot",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is associated with bot activity (paid plan).",
			Transform:   transform.FromField("security.is_bot"),
		},
		{
			Name:        "is_spam",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP appears on spam block lists (paid plan).",
			Transform:   transform.FromField("security.is_spam"),
		},
		{
			Name:        "is_cloud_provider",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP belongs to a cloud/hosting provider (paid plan).",
			Transform:   transform.FromField("security.is_cloud_provider"),
		},
		{
			Name:        "is_known_attacker",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is flagged as a known attacker (paid plan).",
			Transform:   transform.FromField("security.is_known_attacker"),
		},
		{
			Name:        "is_anonymous",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is anonymous (paid plan).",
			Transform:   transform.FromField("security.is_anonymous"),
		},
		{
			Name:        "is_residential_proxy",
			Type:        proto.ColumnType_BOOL,
			Description: "True if the IP is a residential proxy (paid plan).",
			Transform:   transform.FromField("security.is_residential_proxy"),
		},

		// ── Company ─────────────────────────────────────────────────────────
		{
			Name:        "company_name",
			Type:        proto.ColumnType_STRING,
			Description: "Company name associated with the IP (paid plan).",
			Transform:   transform.FromField("company.name"),
		},
		{
			Name:        "company_domain",
			Type:        proto.ColumnType_STRING,
			Description: "Company domain associated with the IP (paid plan).",
			Transform:   transform.FromField("company.domain"),
		},
		{
			Name:        "company_type",
			Type:        proto.ColumnType_STRING,
			Description: "Company type, e.g. ISP, Hosting, Business (paid plan).",
			Transform:   transform.FromField("company.type"),
		},

		// ── Abuse Contact ───────────────────────────────────────────────────
		{
			Name:        "abuse_name",
			Type:        proto.ColumnType_STRING,
			Description: "Abuse contact name (paid plan).",
			Transform:   transform.FromField("abuse.name"),
		},
		{
			Name:        "abuse_email",
			Type:        proto.ColumnType_STRING,
			Description: "Abuse contact email address (paid plan).",
			Transform:   transform.FromField("abuse.email"),
		},
		{
			Name:        "abuse_address",
			Type:        proto.ColumnType_STRING,
			Description: "Abuse contact postal address (paid plan).",
			Transform:   transform.FromField("abuse.address"),
		},
		{
			Name:        "abuse_phone",
			Type:        proto.ColumnType_STRING,
			Description: "Abuse contact phone number (paid plan).",
			Transform:   transform.FromField("abuse.phone"),
		},
		{
			Name:        "abuse_network",
			Type:        proto.ColumnType_STRING,
			Description: "CIDR block of the network in the abuse record (paid plan).",
			Transform:   transform.FromField("abuse.route"),
		},
	}
}

// listIPGeolocation is the hydrate function for the List config.
func listIPGeolocation(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := newClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// Pull the ip qual if provided; empty string means "my IP".
	ip := ""
	if d.EqualsQuals["ip"] != nil {
		ip = d.EqualsQuals["ip"].GetStringValue()
	}

	data, err := client.GetIPGeo(ctx, ip)
	if err != nil {
		plugin.Logger(ctx).Error("ipgeolocation_ip.listIPGeolocation", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, flattenResponse(data))
	return nil, nil
}

// flattenResponse hoists nested maps to top-level dotted-key access understood
// by Steampipe's transform.FromField("location.city") paths. It returns the
// original map unchanged — Steampipe's transform layer handles nested lookups
// automatically via dot notation.
func flattenResponse(raw map[string]interface{}) map[string]interface{} {
	return raw
}
