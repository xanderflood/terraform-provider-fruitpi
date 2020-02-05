package fruitpi

import (
	//nolint:gosec

	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//setData is a helper function to conceal an ugly piece of logic. The general
//TF rule is that you can/should assume that things have the appropriate types
//and panicking if they don't is A-OK. However, when setting a complex data
//type on a ResourceData struct, an error could be returned if the code that
//build that complex object built something that doesn't match the schema. It's
//preferable IMHO to panic in that case, and we should assume that the provider
//was written correctly and that this never ever happens, just like we do with
//all sorts of other things in TF.
//
//By putting this unpleasant little panic right here, we avoid needing to get
//test coverage of it all over the dang place.
func setData(d *schema.ResourceData, key string, val interface{}) {
	if err := d.Set(key, val); err != nil {
		panic(err)
	}
}

func safelyDereferenceString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func safelyDereferenceRawJSON(ptr *json.RawMessage) json.RawMessage {
	if ptr == nil {
		return json.RawMessage(nil)
	}
	return *ptr
}
