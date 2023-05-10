package sunny_day

import (
	"fmt"
	"telegramBot/weather"
	"time"
)

func locationInfo(response *weather.Response) string {
	if response == nil {
		return ""
	}

	info := fmt.Sprintf("The weather in %s:\n%s%s%s%s",
		response.Name,
		mainInfo(&response.Main),
		wind(&response.Wind),
		weatherType(&response.Weather[0]),
		measurementTime(response),
	)
	return info
}

func mainInfo(info *weather.Info) string {
	return fmt.Sprintf("Temperature: %.2f degrees Celcius\n"+
		"Feels like: %.2f degrees Celcius\n"+
		"Humidity: %d%%\n", info.Temperature, info.FeelsLike, info.Humidity)
}

func measurementTime(response *weather.Response) string {
	return fmt.Sprintf("Time of measurement: %s\n", time.Unix(response.Time, 0).Format(time.RFC1123))
}

func wind(windInfo *weather.Wind) string {
	return fmt.Sprintf("Wind speed: %.1f m/s \n", windInfo.Speed)
}

func weatherType(weatherInfo *weather.Weather) string {
	return fmt.Sprintf("Weather type is %s, more precisely: %s\n", weatherInfo.Main, weatherInfo.Description)
}
