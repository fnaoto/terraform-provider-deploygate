package main

import (
	"terraform-provider-deploygate/deploygate"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: deploygate.Provider,
	})
}
