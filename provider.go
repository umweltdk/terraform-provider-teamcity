package main

import (
		"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["user"],
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["password"],
			},

			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["url"],
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["insecure"],
			},

			"skip_credentials_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_credentials_validation"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"teamcity_project": resourceProject(),
			"teamcity_build_configuration": resourceBuildConfiguration(),
			"teamcity_build_template": resourceBuildTemplate(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"user": "The username used for API operations. This must be \n" +
			"an admin user created on the TeamCity Server.",

		"password": "Password of the TeamCity user. The password for the user \n" +
			"specified in the user option.",

		"url": "URL of the TeamCity server to connect to. If not set, the default profile\n" +
			"created with `aws configure` will be used.",

		"insecure": "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted," +
			"default value is `false`",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		User:                d.Get("user").(string),
		Password:            d.Get("password").(string),
		URL:                 d.Get("url").(string),
		Insecure:            d.Get("insecure").(bool),
		SkipCredsValidation: d.Get("skip_credentials_validation").(bool),
	}

	return config.Client()
}
