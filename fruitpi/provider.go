package fruitpi

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanderflood/fruit-pi-server/lib/api"
)

//Provider is the entrypoint to the fruitpi provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FRUIT_PI_SERVER_URL", nil),
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FRUIT_PI_ADMIN_JWT", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"fruitpi_device": resourceDevice,
		},
		DataSourcesMap: map[string]*schema.Resource{
			"fruitpi_device_token": dataSourceDeviceToken,
		},
		ConfigureFunc: configureProvider,
	}
}

//FruitPiMeta is the internal meta configuration object for the fruit-pi provider
type FruitPiMeta struct {
	api.API
}

func configureProvider(r *schema.ResourceData) (interface{}, error) {
	return &FruitPiMeta{
		API: api.NewDefaultClient(
			stripTrailingSlash(r.Get("server_url").(string)),
			http.DefaultTransport,
			r.Get("token").(string),
		),
	}, nil
}

func stripTrailingSlash(apiURL string) string {
	//Strip the trailing slash, since all of our path constants start with a
	//slash already
	if apiURL[len(apiURL)-1] == '/' {
		apiURL = apiURL[:len(apiURL)-1]
	}

	return apiURL
}
