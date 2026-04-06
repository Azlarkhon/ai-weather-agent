package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ForecastTool struct {
	apiKey string
}

func NewForecastTool(apiKey string) *ForecastTool {
	return &ForecastTool{apiKey: apiKey}
}

func (f *ForecastTool) Name() string {
	return "get_forecast"
}

func (f *ForecastTool) Description() string {
	return "Get 5-day weather forecast for a city. Input: { city: string }"
}

func (f *ForecastTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"city": map[string]interface{}{
				"type":        "string",
				"description": "The city name to get forecast for",
			},
		},
		"required": []string{"city"},
	}
}

func (f *ForecastTool) Execute(input map[string]interface{}) (string, error) {
	city, ok := input["city"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid 'city' field")
	}

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric&cnt=5",
		city, f.apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	list, ok := data["list"].([]interface{})
	if !ok || len(list) == 0 {
		return "", fmt.Errorf("no forecast data in response")
	}

	result := fmt.Sprintf("5-day forecast for %s:\n", city)
	for _, item := range list {
		entry := item.(map[string]interface{})
		main := entry["main"].(map[string]interface{})
		weather := entry["weather"].([]interface{})[0].(map[string]interface{})
		dt := entry["dt_txt"].(string)
		temp := main["temp"]
		desc := weather["description"]
		result += fmt.Sprintf("- %s: %.1f°C, %v\n", dt, temp, desc)
	}

	return result, nil
}
