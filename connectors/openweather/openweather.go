package openweather

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Coordinate struct {
	Lon float64 `json:"lon"`
	lat float64 `json:"lon"`
}

type WeatherCond struct {
	Main string `json:"main"`
	Description string `json:"description"`
	Icon string `json:"icon"`
}

type SunTime struct {
	Country string `json:"country"`
	Sunrise time.Time `json:"sunrise"`
	Sunset  time.Time `json:"sunset"`
}

type MainWeather struct {
	Temp float64 `json:"temp"`
	FeelsLike string `json:"feels_like"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
	Pressure float64 `json:"pressure"`
	Humidity float64 `json:"humidity"`
}

type Response struct {
	Name    string        `json:"name"`
	Coor    Coordinate    `json:"coord"`
	Weather []WeatherCond `json:"weather"`
	SunTime SunTime       `json:"sys"`
	Main    MainWeather   `json:"main"`
}

type conn struct {
	baseURI string
	client http.Client
}

func New(apiKey string, unitSys string) *conn{
	return &conn{
		fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?units=%v&appid=%v", unitSys, apiKey),
		http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *conn) buildEndpoint(query string) string {
	return fmt.Sprint(c.baseURI, "&", query)
}

func (c *conn) openWeatherQuery(query string) ([]byte, error) {

	// Create and endpoint
	ep :=  c.buildEndpoint(query)

	log.Println(ep)
	// Send the request
	res, err := c.client.Get(ep)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Error on
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error: http status not ok, %v", res)
	}

	// Response from open get_weather
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (c *conn) QueryByName(city string, state string, country string) ([]byte, error){
	temp := []string{city, state, country}
	q := strings.Builder{}
	q.WriteString("q=")
	count := 0
	for _, v := range temp {
		if  v != "" {
			count = count + 1

		}

		if count > 1 {
			q.WriteString(",")
		}

		q.WriteString(v)

	}

	return c.openWeatherQuery(q.String())
}
