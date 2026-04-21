package main

import (
	"github.com/IPGeolocation/steampipe-plugin-ipgeolocation/ipgeolocation"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: ipgeolocation.Plugin})
}
