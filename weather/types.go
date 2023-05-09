package weather

type GeocodingRequest struct {
	Name  string `json:"q"`
	Token string `json:"appid"`
	Limit int    `json:"limit"` // up to 5
}

type GeocodingResponse struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type LocationRequest struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Token     string  `json:"appid"`
	Units     string  `json:"units"`
}

type Response struct {
	Weather []Weather `json:"weather"`
	Main    Info      `json:"main"`
	Time    int64     `json:"dt"`
	Wind    Wind      `json:"wind"`
	Name    string    `json:"name"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}

type Info struct {
	Temperature float64 `json:"temp"`
	FeelsLike   float64 `json:"feels_like"`
	Humidity    int     `json:"humidity"`
}

type Weather struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}
