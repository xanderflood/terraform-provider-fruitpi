package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

//Doer impelments the http Do method
//go:generate counterfeiter . Doer
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

//RoundTripper impelments the http RoundTrip method
//go:generate counterfeiter . RoundTripper
type RoundTripper interface {
	RoundTrip(req *http.Request) (*http.Response, error)
}

//RoundTripHandler is a helper type for building RoundTrippers
type RoundTripHandler func(req *http.Request) (*http.Response, error)

//RoundTrip executes the request
func (h RoundTripHandler) RoundTrip(req *http.Request) (*http.Response, error) {
	return h(req)
}

//Client implements API over HTTPs
type Client struct {
	host string
	http Doer
}

//NewClient creates a client that directly uses the provided transport
func NewClient(
	host string,
	http Doer,
) Client {
	return Client{
		host: host,
		http: http,
	}
}

//NewDefaultClient creates a client with default middlewares
func NewDefaultClient(
	host string,
	transport http.RoundTripper,
	token string,
) Client {
	transport = DefaultRetryer(transport)
	transport = DefaultAuthorizer(transport, token)

	client := *http.DefaultClient
	client.Transport = transport
	return NewClient(host, &client)
}

//API is the client interface
//go:generate counterfeiter . API
type API interface {
	//device endpoints - device UUID is inferred from token
	GetDeviceConfig(ctx context.Context) (Device, error)
	InsertReading(ctx context.Context, tCelcius float64, rh float64) (Reading, error)

	//admin endpoints
	RegisterDevice(ctx context.Context, name string, config string) (Device, error)
	ConfigureDevice(ctx context.Context, uuid string, name string, config string) (Device, error)
	GetDeviceTokenFor(ctx context.Context, uuid string) (Device, error)
	GetDeviceConfigFor(ctx context.Context, uuid string) (device Device, err error)
}

////////////////////
// Device endpoints
////////////////////

//GetDeviceConfig gets the current config text for current the device
func (c Client) GetDeviceConfig(ctx context.Context) (device Device, err error) {
	return c.getDeviceConfigHelper(ctx, "/api/v1/get-device-config")
}

//InsertReading records a new sensor reading for current the device
func (c Client) InsertReading(ctx context.Context, tCelcius float64, rh float64) (reading Reading, err error) {
	req, err := c.buildJSONRequest(ctx,
		http.MethodPost, "/api/v1/insert-reading",
		InsertReadingRequest{
			TemperatureCelcius: toJSONNumber(tCelcius),
			RelativeHumidity:   toJSONNumber(rh),
		},
	)
	if err != nil {
		return
	}

	resp, err := c.http.Do(req)
	if err != nil {
		err = fmt.Errorf("failed making insert-reading request: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = responseError(resp, true)
		return
	}

	err = decodeJSONBody(resp, &reading)
	return
}

///////////////////
// Admin endpoints
///////////////////

//GetDeviceConfigFor gets the current config text for the specified device
func (c Client) GetDeviceConfigFor(ctx context.Context, uuid string) (device Device, err error) {
	return c.getDeviceConfigHelper(ctx, "/api/v1/get-device-config/"+uuid)
}

//RegisterDevice registers a new device
func (c Client) RegisterDevice(ctx context.Context, name string, config string) (device Device, err error) {
	req, err := c.buildJSONRequest(
		ctx, http.MethodPost,
		"/api/v1/register-device",
		RegistrationRequest{
			Name:   name,
			Config: json.RawMessage(config),
		},
	)
	if err != nil {
		return
	}

	resp, err := c.http.Do(req)
	if err != nil {
		err = fmt.Errorf("failed making register-device request: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = responseError(resp, true)
		return
	}

	err = decodeJSONBody(resp, &device)
	return
}

//GetDeviceTokenFor generates a token for the device
func (c Client) GetDeviceTokenFor(ctx context.Context, uuid string) (device Device, err error) {
	req, err := c.buildJSONRequest(
		ctx, http.MethodPost,
		"/api/v1/get-device-token",
		GetDeviceTokenRequest{
			DeviceUUID: uuid,
		},
	)
	if err != nil {
		return
	}

	resp, err := c.http.Do(req)
	if err != nil {
		err = fmt.Errorf("failed making device-token-for request: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = responseError(resp, true)
		return
	}

	err = decodeJSONBody(resp, &device)
	return
}

//ConfigureDevice configures the specified device
func (c Client) ConfigureDevice(ctx context.Context, uuid string, name string, config string) (device Device, err error) {
	req, err := c.buildJSONRequest(
		ctx, http.MethodPost,
		"/api/v1/configure-device",
		ConfigureDeviceRequest{
			DeviceUUID: uuid,
			Name:       name,
			Config:     json.RawMessage(config),
		},
	)
	if err != nil {
		return
	}

	resp, err := c.http.Do(req)
	if err != nil {
		err = fmt.Errorf("failed making configure-device request: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = responseError(resp, true)
		return
	}

	err = decodeJSONBody(resp, &device)
	return
}

///////////
// Helpers
///////////

func (c Client) buildJSONRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	buf := &bytes.Buffer{}
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("could not marshal request body into JSON: %w", err)
		}
		buf = bytes.NewBuffer(b)
	}

	//bytes.Buffer never returns an error, so neither will this
	uri := c.host + path
	req, err := http.NewRequestWithContext(ctx, method, uri, buf)
	if err != nil {
		return nil, err
	}

	//Add JSON request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "terminus/1.0")
	req.Header.Set("Accept", "*/*")
	return req, nil
}

func decodeJSONBody(resp *http.Response, body interface{}) error {
	if err := json.NewDecoder(resp.Body).Decode(body); err != nil {
		err = fmt.Errorf("received malformed JSON response from TTD: %w", err)
		return err
	}

	return nil
}

func responseError(resp *http.Response, body bool) error {
	if body {
		var bodyText string
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			bodyText = "<error reading body>"
		} else {
			bodyText = string(bs)
		}

		return fmt.Errorf("HTTP error %d %s with body `%s`", resp.StatusCode, resp.Status, bodyText)
	}

	return fmt.Errorf("HTTP error %d %s", resp.StatusCode, resp.Status)
}

func toJSONNumber(f float64) json.Number {
	return json.Number(strconv.FormatFloat(f, 'f', -1, 64))
}

func (c Client) getDeviceConfigHelper(ctx context.Context, path string) (device Device, err error) {
	req, err := c.buildJSONRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return
	}

	resp, err := c.http.Do(req)
	if err != nil {
		err = fmt.Errorf("failed making get-device-config request: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = responseError(resp, true)
		return
	}

	err = decodeJSONBody(resp, &device)
	return
}
