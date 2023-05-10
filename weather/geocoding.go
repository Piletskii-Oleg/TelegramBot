package weather

import (
	"encoding/json"
	"net/url"
	"path"
	"strconv"
)

const (
	host           = "api.openweathermap.org"
	geoApiAddress  = "geo"
	geoVersion     = "1.0"
	geoCallAddress = "direct"
)

// MakeGeoRequest returns array of GeocodingResponse objects that contain
// latitudes and longitudes of the locations that fit the name criterion
func (c *Client) MakeGeoRequest(location string) ([]GeocodingResponse, error) {
	values := url.Values{}
	values.Add("q", location)
	values.Add("appid", c.token)
	values.Add("limit", strconv.Itoa(4))

	apiUrl := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path.Join(geoApiAddress, geoVersion, geoCallAddress),
	}

	content, err := c.doRequest(values, apiUrl.String())
	if err != nil {
		return nil, err
	}

	var response []GeocodingResponse
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, err
	}

	return response, nil
}
