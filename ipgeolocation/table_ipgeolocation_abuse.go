package ipgeolocation

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// tableIPAbuse returns the Steampipe table definition for ipgeolocation_abuse.
func tableIPAbuse(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "ipgeolocation_abuse",
		Description: "Abuse contact lookup using the IPGeolocation.io /v3/abuse endpoint. " +
			"Returns the responsible organisation, role, emails, phone numbers, " +
			"postal address, CIDR route, and registered country for any IPv4 or IPv6 address. " +
			"Requires a paid plan. Each query costs 1 API credit.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"ip"}),
			Hydrate:    listIPAbuse,
		},
		Columns: ipAbuseColumns(),
	}
}

func ipAbuseColumns() []*plugin.Column {
	return []*plugin.Column{
		// ── Top-level ───────────────────────────────────────────────────────
		{
			Name:        "ip",
			Type:        proto.ColumnType_STRING,
			Description: "The IP address that was looked up.",
			Transform:   transform.FromField("ip"),
		},

		// ── Abuse contact ───────────────────────────────────────────────────
		{
			Name:        "route",
			Type:        proto.ColumnType_STRING,
			Description: "CIDR block that covers the queried IP (e.g. 1.0.0.0/24).",
			Transform:   transform.FromField("abuse.route"),
		},
		{
			Name:        "country",
			Type:        proto.ColumnType_STRING,
			Description: "ISO alpha-2 country code of the abuse registrant.",
			Transform:   transform.FromField("abuse.country"),
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "Name of the abuse contact or IRT (Incident Response Team).",
			Transform:   transform.FromField("abuse.name"),
		},
		{
			Name:        "organization",
			Type:        proto.ColumnType_STRING,
			Description: "Organisation responsible for the IP block.",
			Transform:   transform.FromField("abuse.organization"),
		},
		{
			Name:        "kind",
			Type:        proto.ColumnType_STRING,
			Description: "Kind of contact record (e.g. \"group\", \"individual\").",
			Transform:   transform.FromField("abuse.kind"),
		},
		{
			Name:        "address",
			Type:        proto.ColumnType_STRING,
			Description: "Postal address of the abuse contact.",
			Transform:   transform.FromField("abuse.address"),
		},
		{
			Name:        "emails",
			Type:        proto.ColumnType_JSON,
			Description: "Array of abuse contact email addresses.",
			Transform:   transform.FromField("abuse.emails"),
		},
		{
			Name:        "phone_numbers",
			Type:        proto.ColumnType_JSON,
			Description: "Array of abuse contact phone numbers.",
			Transform:   transform.FromField("abuse.phone_numbers"),
		},

		// ── Raw JSON ────────────────────────────────────────────────────────
		{
			Name:        "raw",
			Type:        proto.ColumnType_JSON,
			Description: "Full raw JSON response from the /v3/abuse API endpoint.",
			Transform:   transform.FromValue(),
		},
	}
}

// listIPAbuse is the hydrate function for ipgeolocation_abuse.
func listIPAbuse(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := newClient(ctx, d)
	if err != nil {
		return nil, err
	}

	ip := ""
	if d.EqualsQuals["ip"] != nil {
		ip = d.EqualsQuals["ip"].GetStringValue()
	}

	data, err := client.GetAbuse(ctx, ip)
	if err != nil {
		plugin.Logger(ctx).Error("ipgeolocation_abuse.listIPAbuse", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, data)
	return nil, nil
}
