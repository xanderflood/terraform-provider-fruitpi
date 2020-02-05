package fruitpi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var resourceDevice = &schema.Resource{
	Create: resourceDeviceCreate,
	Read:   resourceDeviceRead,
	Update: resourceDeviceUpdate,
	Delete: resourceDeviceDelete,
	Schema: resourceDeviceSchema,
}

var resourceDeviceSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"config": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "{}",
	},
}

//////////////////
// Main callbacks
//////////////////

func resourceDeviceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*FruitPiMeta).API
	name := d.Get("name").(string)
	config := d.Get("config").(string)

	device, err := client.RegisterDevice(context.Background(), name, config)
	if err != nil {
		return err
	}

	d.SetId(device.DeviceUUID)
	return resourceDeviceRead(d, m)
}

func resourceDeviceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*FruitPiMeta).API
	name := d.Get("name").(string)
	config := d.Get("config").(string)

	_, err := client.ConfigureDevice(context.Background(), d.Id(), name, config)
	if err != nil {
		return err
	}

	return resourceDeviceRead(d, m)
}

func resourceDeviceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*FruitPiMeta).API

	device, err := client.GetDeviceConfigFor(context.Background(), d.Id())
	if err != nil {
		return err
	}

	setData(d, "name", safelyDereferenceString(device.Name))
	setData(d, "config", string(safelyDereferenceRawJSON(device.Config)))
	return nil
}

func resourceDeviceDelete(d *schema.ResourceData, m interface{}) error {
	// noop
	return nil
}
