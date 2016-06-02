package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/joeswaminathan/terraform/builtin/providers/icf"
	"cto-github.cisco.com/jswamina/kvs_infra/src/infra/log"
)

func main() {
	log.StartLogger("icf-terraform.log", true)
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: icf.Provider,
	})

}
