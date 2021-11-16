package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	// Savage global mutable variable for setting the documentation syntax.
	// This will be used in document generation and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}
