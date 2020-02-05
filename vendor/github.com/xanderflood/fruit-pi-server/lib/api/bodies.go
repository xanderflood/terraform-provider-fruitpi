package api

import "encoding/json"

//////////////////
// Device bodies
//////////////////

//RegistrationRequest encodes a single request for user registration
type RegistrationRequest struct {
	Name   string          `json:"name" binding:"required"`
	Config json.RawMessage `json:"config" binding:"required"`
}

//GetDeviceTokenRequest encodes a single request for user registration
type GetDeviceTokenRequest struct {
	DeviceUUID string `json:"device_uuid" binding:"required"`
}

//GetDeviceConfigRequest encodes a single request for user registration
type GetDeviceConfigRequest struct {
	DeviceUUID *string `uri:"uuid"`
}

//ConfigureDeviceRequest encodes a single request for user registration
type ConfigureDeviceRequest struct {
	DeviceUUID string          `json:"device_uuid" binding:"required"`
	Name       string          `json:"name"`
	Config     json.RawMessage `json:"config" binding:"required"`
}

//Device represents the response to a device-related endpoint
type Device struct {
	DeviceUUID string           `json:"device_uuid,omitempty"`
	Name       *string          `json:"name,omitempty"`
	Token      *string          `json:"token,omitempty"`
	Config     *json.RawMessage `json:"config,omitempty"`
}

//////////////////
// Reading bodies
//////////////////

//InsertReadingRequest encodes a single request for user registration
type InsertReadingRequest struct {
	TemperatureCelcius json.Number `json:"temperature_celcius" binding:"required"`
	RelativeHumidity   json.Number `json:"relative_humidity" binding:"required"`
}

//Reading represents the response to a reading-related endpoint
type Reading struct {
	DeviceUUID         string       `json:"device_uuid,omitempty"`
	ReadingUUID        string       `json:"reading_uuid,omitempty"`
	TemperatureCelcius *json.Number `json:"temperature_celcius,omitempty"`
	RelativeHumidity   *json.Number `json:"relative_humidity,omitempty"`
}
