package main

import (
	"cto-github.cisco.com/jswamina/kvs_infra/src/infra/log"
	"github.com/hashicorp/terraform/plugin"
	"github.com/jswamina/terraform/icf"
)

func main() {
	log.StartLogger("icf-terraform", true)
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: icf.Provider,
	})
}
