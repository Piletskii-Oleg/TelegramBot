package weather

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"telegramBot/bot/errors"
)

const (
	apiAddress  = "data"
	version     = "2.5"
	callAddress = "weather"
)

// Client is OpenWeatherMap client.
type Client struct {
	token  string
	client http.Client
}

// NewClient makes a new OpenWeatherMap client using the token provided.
func NewClient(token string) *Client {
	return &Client{token: token, client: http.Client{}}
}

// MakeRequest calls OpenWeatherMap API to get weather data for the specified location
// and returns the data.
func (c *Client) MakeRequest(location string) (*Response, error) {
	geo, err := c.MakeGeoRequest(location)
	if err != nil {
		return nil, err
	}
	if len(geo) == 0 {
		return nil, err
	}

	values := url.Values{}
	values.Add("lat", strconv.FormatFloat(geo[0].Latitude, 'f', -1, 64))
	values.Add("lon", strconv.FormatFloat(geo[0].Longitude, 'f', -1, 64))
	values.Add("appid", c.token)
	values.Add("units", "metric")

	apiUrl := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path.Join(apiAddress, version, callAddress),
	}

	content, err := c.doRequest(values, apiUrl.String())
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) doRequest(query url.Values, url string) (content []byte, err error) {
	defer func() {
		err = errors.WrapIfError("unable to do request: %w", err)
	}()

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap("unable to do request: %w", err)
	}
	request.URL.RawQuery = query.Encode()

	response, err := c.client.Do(request)
	if err != nil {
		return nil, errors.Wrap("unable to do request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	content, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap("unable to read content: %w", err)
	}

	return content, err
}
