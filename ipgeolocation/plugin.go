package ipgeolocation

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns the plugin definition.
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-ipgeolocation",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"ipgeolocation_ip":       tableIPGeolocation(ctx),
			"ipgeolocation_security": tableIPSecurity(ctx),
			"ipgeolocation_abuse":    tableIPAbuse(ctx),
			"ipgeolocation_asn":      tableASN(ctx),
		},
	}
	return p
}
