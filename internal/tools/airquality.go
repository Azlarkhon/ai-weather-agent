package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AirQualityTool struct {
	apiKey string
}

func NewAirQualityTool(apiKey string) *AirQualityTool {
	return &AirQualityTool{apiKey: apiKey}
}

func (a *AirQualityTool) Name() string {
	return "get_air_quality"
}

func (a *AirQualityTool) Description() string {
	return "Get air quality index (AQI) for a city. Input: { city: string }"
}

func (a *AirQualityTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"city": map[string]interface{}{
				"type":        "string",
				"description": "The city name to get air quality for",
			},
		},
		"required": []string{"city"},
	}
}

func (a *AirQualityTool) Execute(input map[string]interface{}) (string, error) {
	city, ok := input["city"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid 'city' field")
	}

	// Step 1: get coordinates for the city
	geoURL := fmt.Sprintf(
		"http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s",
		city, a.apiKey,
	)

	geoResp, err := http.Get(geoURL)
	if err != nil {
		return "", err
	}
	defer geoResp.Body.Close()

	var geoData []map[string]interface{}
	if err := json.NewDecoder(geoResp.Body).Decode(&geoData); err != nil {
		return "", err
	}
	if len(geoData) == 0 {
		return "", fmt.Errorf("city not found: %s", city)
	}

	lat := geoData[0]["lat"].(float64)
	lon := geoData[0]["lon"].(float64)

	// Step 2: get air quality using coordinates
	aqURL := fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/air_pollution?lat=%f&lon=%f&appid=%s",
		lat, lon, a.apiKey,
	)

	aqResp, err := http.Get(aqURL)
	if err != nil {
		return "", err
	}
	defer aqResp.Body.Close()

	var aqData map[string]interface{}
	if err := json.NewDecoder(aqResp.Body).Decode(&aqData); err != nil {
		return "", err
	}

	list, ok := aqData["list"].([]interface{})
	if !ok || len(list) == 0 {
		return "", fmt.Errorf("no air quality data in response")
	}

	entry := list[0].(map[string]interface{})
	aqi := entry["main"].(map[string]interface{})["aqi"].(float64)

	aqiLabels := map[float64]string{
		1: "Good",
		2: "Fair",
		3: "Moderate",
		4: "Poor",
		5: "Very Poor",
	}

	return fmt.Sprintf("Air quality in %s: AQI %v (%s)", city, aqi, aqiLabels[aqi]), nil
}
