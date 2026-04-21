package ipgeolocation

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

// IPGeolocationConfig holds the connection configuration for the plugin.
type IPGeolocationConfig struct {
	// APIKey is the ipgeolocation.io API key (required for most endpoints).
	APIKey *string `hcl:"api_key"`
}

// ConfigInstance returns a new, empty config instance.
func ConfigInstance() interface{} {
	return &IPGeolocationConfig{}
}

// GetConfig fetches the parsed config from the connection.
func GetConfig(connection *plugin.Connection) IPGeolocationConfig {
	if connection == nil || connection.Config == nil {
		return IPGeolocationConfig{}
	}
	config, _ := connection.Config.(IPGeolocationConfig)
	return config
}

// ConnectionConfigSchema wires HCL attribute definitions to the struct.
var ConfigSchema = map[string]*schema.Attribute{
	"api_key": {
		Type: schema.TypeString,
	},
}
