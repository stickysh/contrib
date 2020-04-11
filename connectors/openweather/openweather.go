/*
Copyright 2015 Sticky Contrib Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package openweather

import (
	"encoding/json"
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

// UnmarshalJSON reads in a SunTime from its JSON format.
// Used to parse unix epoc to time.time
func (o *SunTime) UnmarshalJSON(data []byte) error {
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil { return err; }
	m := f.(map[string]interface{})
	for k, v := range m {
		switch k {
		case "country": o.Country  = v.(string)
		case "sunrise": o.Sunrise = time.Unix(int64(v.(float64)), 0)
		case "sunset": o.Sunset = time.Unix(int64(v.(float64)), 0)
		}
	}

	return nil
}

type MainWeather struct {
	Temp float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
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

// New create a new openweather api client
func New(apiKey string, unitSys string) *conn{
	if unitSys == "" {
		unitSys = "imperial"
	}

	return &conn{
		fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?units=%v&appid=%v", unitSys, apiKey),
		http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *conn) buildEp(query string) string {
	return fmt.Sprint(c.baseURI, "&", query)
}

func (c *conn) openWeatherQuery(query string) (Response, error) {

	// Create an endpoint
	ep :=  c.buildEp(query)

	// Send the request
	res, err := c.client.Get(ep)
	if err != nil {
		log.Printf("[__sticky]connector=openweather|level=error,msg=couldn't send request,err=%s\n", err)
		return Response{}, err
	}
	defer res.Body.Close()

	// Error on
	if res.StatusCode != 200 {
		log.Printf("[__sticky]connector=openweather|level=error,msg=http status not ok,err=%s\n", err)
		return Response{}, fmt.Errorf("error: http status not ok, errormsg=%v", res)
	}

	// Response from open get_weather
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("[__sticky]connector=openweather|level=error,msg=failed to query,errormsg=%s,payload=%s\n", err, raw)
		return Response{}, err
	}

	// Parse get_weather response
	var owData Response
	err = json.Unmarshal(raw, &owData)
	if err != nil {
		log.Printf("[__sticky]connector=openweather|level=error,msg=failed to unmarshal,errormsg=%s,payload=%s\n", err, raw)
		return Response{}, err
	}

	return owData, nil
}

// QueryByName queries the openweather api by city, state, country name
func (c *conn) QueryByName(city string, state string, country string) (Response, error){
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
