package main

import (
	"ai-agent/internal/agent"
	"ai-agent/internal/config"
	delivery "ai-agent/internal/delivery/http"
	"ai-agent/internal/domain"
	"ai-agent/internal/llm"
	"ai-agent/internal/tools"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	llmClient := llm.NewGeminiClient(config.App.GeminiKey)

	agentService := agent.NewAgent(
		llmClient,
		[]domain.Tool{
			tools.NewWeatherTool(config.App.OpenWeatherKey),
			tools.NewForecastTool(config.App.OpenWeatherKey),
			tools.NewAirQualityTool(config.App.OpenWeatherKey),
		},
	)

	handler := delivery.NewHandler(agentService)

	r.POST("/chat", handler.Chat)

	r.Run(":8080")
}
