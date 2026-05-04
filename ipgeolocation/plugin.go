package ipgeolocation

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns the plugin definition.
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "ipgeolocation",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		DefaultTransform: transform.FromCamel(),
		TableMap: map[string]*plugin.Table{
			"ipgeolocation_ip":       tableIpgeolocationIp(ctx),
			"ipgeolocation_security": tableIpgeolocationSecurity(ctx),
			"ipgeolocation_abuse":    tableIpgeolocationAbuse(ctx),
			"ipgeolocation_asn":      tableIpgeolocationAsn(ctx),
		},
	}
	return p
}
