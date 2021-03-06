package main

import (
    "github.com/hashicorp/terraform/terraform"
    "github.com/hashicorp/terraform/plugin"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: func() terraform.ResourceProvider {
            return Provider()
        },
    })
}