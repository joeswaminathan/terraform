package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/joeswaminathan/terraform/icf"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: icf.Provider,
	})
}
