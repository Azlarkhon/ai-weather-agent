package domain

type ToolDefinition struct {
	Name        string
	Description string
	Parameters  map[string]any
}

type LLM interface {
	Generate(input string, tools []ToolDefinition) (LLMResponse, error)
}
