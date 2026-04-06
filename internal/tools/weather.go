package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WeatherTool struct {
	apiKey string
}

func NewWeatherTool(apiKey string) *WeatherTool {
	return &WeatherTool{apiKey: apiKey}
}

func (w *WeatherTool) Name() string {
	return "get_weather"
}

func (w *WeatherTool) Description() string {
	return "Get current weather by city. Input: { city: string }"
}

func (w *WeatherTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"city": map[string]interface{}{
				"type":        "string",
				"description": "The city name to get weather for",
			},
		},
		"required": []string{"city"},
	}
}

func (w *WeatherTool) Execute(input map[string]interface{}) (string, error) {
	city, ok := input["city"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid 'city' field")
	}

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric",
		city, w.apiKey, // ✅ injected, not hardcoded
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

	main, ok := data["main"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	weather, ok := data["weather"].([]any)
	if !ok || len(weather) == 0 {
		return "", fmt.Errorf("no weather data in response")
	}

	temp := main["temp"]
	desc := weather[0].(map[string]any)["description"]

	return fmt.Sprintf("%.1f°C, %v", temp, desc), nil
}
