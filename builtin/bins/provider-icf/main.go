package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/joeswaminathan/terraform/builtin/providers/icf"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: icf.Provider,
	})
}
