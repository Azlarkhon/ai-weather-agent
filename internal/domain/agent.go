package domain

type AgentRequest struct {
	Message string
}

type AgentResponse struct {
	Message string
}

type ToolCall struct {
	Name  string         `json:"name"`
	Input map[string]any `json:"input"`
}

type LLMResponse struct {
	Content  string    `json:"content"`
	ToolCall *ToolCall `json:"tool_call,omitempty"`
}
