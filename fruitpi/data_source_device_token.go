package fruitpi

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var dataSourceDeviceToken = &schema.Resource{
	Read:   dataSourceDeviceTokenRead,
	Schema: dataSourceDeviceTokenSchema,
}

var dataSourceDeviceTokenSchema = map[string]*schema.Schema{
	"device_uuid": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},

	//computed fields
	"token": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
}

//////////////////
// Main callbacks
//////////////////

func dataSourceDeviceTokenRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*FruitPiMeta).API
	uuid := d.Get("device_uuid").(string)

	device, err := client.GetDeviceTokenFor(context.Background(), uuid)
	if err != nil {
		return err
	}
	if device.Token == nil {
		return errors.New("did not obtain token")
	}

	setData(d, "token", *device.Token)
	d.SetId(uuid)
	return nil
}
