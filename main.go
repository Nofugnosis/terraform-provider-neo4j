package main

import (
	"flag"
	"terraform-provider-neo4j/neo4j"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return neo4j.Provider()
		},
		ProviderAddr: "registry.terraform.io/headyj/neo4j",
		Debug:        debug,
	})
}
