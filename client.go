package weatherkit

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// DefaultUserAgent to send along with requests.
const DefaultUserAgent = "nicknassar-weatherkit"

// Client is a WeatherKit API client.
type Client struct {
	HttpClient  *http.Client
	Credentials *credentials

	// The UserAgent header value to send along with requests.
	UserAgent string
}

func NewClient(keyId, teamId, serviceId, privateKey string) (*Client, error) {
	if creds, err := newCredentials(keyId, teamId, serviceId, privateKey); err != nil {
		return nil, err
	} else {
		return &Client{
			HttpClient:  &http.Client{},
			Credentials: creds,
		}, nil
	}
}

// Weather obtains weather data for the specified location.
func (d *Client) Weather(ctx context.Context, request WeatherRequest) (*WeatherResponse, error) {
	response := WeatherResponse{}
	err := d.get(ctx, request, &response)
	return &response, err
}

// Availability determines the data sets available for the specified location.
// The token parameter is a JWT developer token.
func (d *Client) Availability(ctx context.Context, request AvailabilityRequest) (*AvailabilityResponse, error) {
	response := AvailabilityResponse{}
	err := d.get(ctx, request, &response)
	return &response, err
}

// Alert receives information on an active weather alert.
// The token parameter is a JWT developer token.
func (d *Client) Alert(ctx context.Context, request WeatherAlertRequest) (*WeatherAlertResponse, error) {
	response := WeatherAlertResponse{}
	err := d.get(ctx, request, &response)
	return &response, err
}

func (d *Client) token() (string, error) {
	if d.Credentials == nil {
		return "", nil
	} else if token, err := d.Credentials.Token(); err != nil {
		return "", err
	} else {
		return token, nil
	}
}

func (d *Client) get(ctx context.Context, request urlBuilder, output interface{}) error {
	if d.HttpClient == nil {
		d.HttpClient = &http.Client{}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, request.url(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", d.userAgent())
	req.Header.Add("Accept", "application/json; charset=utf-8")
	req.Header.Add("Accept-Encoding", "gzip")

	// get the auth token and add it to the headers if present
	if token, err := d.token(); err != nil {
		return fmt.Errorf("weatherkit failed to get token: %w", err)
	} else if len(token) > 0 {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	response, err := d.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = validateResponse(response)
	if err != nil {
		return err
	}

	return decode(response, &output)
}

func (d *Client) userAgent() string {
	if len(d.UserAgent) > 0 {
		return d.UserAgent
	}

	return DefaultUserAgent
}

func validateResponse(response *http.Response) error {
	if response.StatusCode == http.StatusOK {
		return nil
	}

	errorResponse := ErrorResponse{}

	err := decode(response, &errorResponse)
	if err != nil {
		return err
	}

	return &RestError{
		Response:      response,
		ErrorResponse: &errorResponse,
	}
}

func decode(response *http.Response, into interface{}) error {
	body, err := decompress(response)
	if err != nil {
		return err
	}

	return unmarshal(body, into)
}

func decompress(response *http.Response) (io.Reader, error) {
	header := response.Header.Get("Content-Encoding")
	if len(header) < 1 {
		return response.Body, nil
	}

	return gzip.NewReader(response.Body)
}

func unmarshal(body io.Reader, into interface{}) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	if bytes == nil || len(bytes) < 1 {
		return nil
	}

	return json.Unmarshal(bytes, &into)
}
