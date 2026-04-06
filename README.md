# 🌤️ AI Weather Agent

A smart AI-powered weather agent built with **Go**, **Gin**, and **Gemini 2.5 Flash**. Ask it anything weather-related in natural language and it will use the right tool to answer you.

---

## 🧠 How It Works

```
User Message → Gin HTTP Server → AI Agent → Gemini LLM
                                               ↓
                                     Decides which tool to call
                                               ↓
                                     Executes tool (OpenWeatherMap)
                                               ↓
                                     Sends result back to Gemini
                                               ↓
                                     Returns natural language answer
```

---

## 🚀 Features

- 💬 Natural language interface via REST API
- 🌡️ Current weather by city
- 📅 5-day weather forecast
- 🌫️ Air quality index (AQI)
- 🤖 Powered by Gemini 2.5 Flash (free tier)

---

## 📁 Project Structure

```
ai-agent/
├── cmd/
│   └── main.go                  # Entry point
├── internal/
│   ├── agent/
│   │   ├── agent.go             # Core agent loop
│   │   └── tools.go             # Tool definition builder
│   ├── config/
│   │   └── config.go            # Env loader (godotenv)
│   ├── delivery/
│   │   └── http/
│   │       └── handler.go       # Gin HTTP handler
│   ├── domain/
│   │   ├── domain.go            # LLM, ToolDefinition types
│   │   ├── models.go            # Request/Response/ToolCall types
│   │   └── tool.go              # Tool interface
│   ├── llm/
│   │   └── gemini.go            # Gemini LLM client
│   └── tools/
│       ├── weather.go           # Current weather tool
│       ├── forecast.go          # 5-day forecast tool
│       └── airquality.go        # Air quality tool
├── .env                         # API keys (not committed)
├── .gitignore
├── go.mod
└── go.sum
```

---

## ⚙️ Setup

### 1. Clone the repo

```bash
git clone https://github.com/yourusername/ai-agent.git
cd ai-agent
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Get API keys

| Service | URL | Cost |
|---|---|---|
| Gemini | [aistudio.google.com](https://aistudio.google.com) | Free |
| OpenWeatherMap | [openweathermap.org](https://openweathermap.org) | Free |

> ⚠️ For Gemini, make sure to get the key from **AI Studio**, not Google Cloud Console.

### 4. Create `.env` file

```dotenv
GEMINI_API_KEY=your_gemini_key_here
OPENWEATHER_API_KEY=your_openweather_key_here
```

### 5. Build and run

```bash
go build -o ai-agent.exe ./cmd/main.go
.\ai-agent.exe
```

Server starts at `http://localhost:8080`.

---

## 📡 API

### `POST /chat`

**Request:**
```json
{
    "message": "What is the weather in Tokyo?"
}
```

**Response:**
```json
{
    "response": "The current weather in Tokyo is 18.5°C, partly cloudy."
}
```

---

## 💡 Example Queries

```
"What is the weather in London?"
"Give me the 5-day forecast for New York"
"What is the air quality in Beijing?"
"Is it going to rain in Paris this week?"
```

---

## 🛠️ Adding a New Tool

1. Create a new file in `internal/tools/`
2. Implement the `domain.Tool` interface:

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() map[string]interface{}
    Execute(input map[string]interface{}) (string, error)
}
```

3. Register it in `main.go`:

```go
agentService := agent.NewAgent(
    llmClient,
    []domain.Tool{
        tools.NewWeatherTool(config.App.OpenWeatherKey),
        tools.NewForecastTool(config.App.OpenWeatherKey),
        tools.NewAirQualityTool(config.App.OpenWeatherKey),
        tools.NewYourTool(...), // 👈 add here
    },
)
```

---

## 📦 Dependencies

- [gin-gonic/gin](https://github.com/gin-gonic/gin) — HTTP framework
- [google/generative-ai-go](https://github.com/google/generative-ai-go) — Gemini SDK
- [joho/godotenv](https://github.com/joho/godotenv) — .env loader

---

## 📄 License

MIT
