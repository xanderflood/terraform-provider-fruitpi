package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/xanderflood/terraform-provider-fruitpi/fruitpi"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: fruitpi.Provider,
	})
}
