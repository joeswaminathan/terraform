package icf

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/joeswaminathan/icf-sdk-go/icf"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// TODO: Move the validation to this, requires conditional schemas
	// TODO: Move the configuration to this, requires validation

	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["username"],
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["password"],
			},

			"icfb": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ICF_SERVER", nil),
				Description: descriptions["icfb"],
			},

			"vdc": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"ICF_VDC",
				}, nil),
				Description: descriptions["vdc"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"icf_instance": resourceIcfInstance(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"icfb": "ICFB IP address or domain name",
		"vdc":  "The OID of the virtual data center on which VMs will be deployed",

		"username": "The username for API operations",

		"password": "The password for API operations",
	}
}

func providerConfigure(d *schema.ResourceData) (metadata interface{}, err error) {
	config := Config{
		Credentials: icf.Credentials{
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
		},
		EndPoint: d.Get("icfb").(string),
		Protocol: "http",
		Root:     "icfb/v1",
	}

	if metadata = config.Client(); metadata == nil {
		err = fmt.Errorf("Unable to create Client")
	}

	return
}
