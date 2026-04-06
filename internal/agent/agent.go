package agent

import (
	"ai-agent/internal/domain"
	"fmt"
)

type Agent struct {
	llm   domain.LLM
	tools map[string]domain.Tool
}

func NewAgent(llm domain.LLM, tools []domain.Tool) *Agent {
	toolMap := make(map[string]domain.Tool)
	for _, t := range tools {
		toolMap[t.Name()] = t
	}

	return &Agent{
		llm:   llm,
		tools: toolMap,
	}
}

func (a *Agent) Run(input string) (string, error) {
	toolDefs := ToToolDefinitions(a.tools)

	llmRes, err := a.llm.Generate(input, toolDefs)
	if err != nil {
		return "", err
	}

	if llmRes.ToolCall != nil {
		tool, ok := a.tools[llmRes.ToolCall.Name]
		if !ok {
			return "", fmt.Errorf("unknown tool: %s", llmRes.ToolCall.Name)
		}

		result, err := tool.Execute(llmRes.ToolCall.Input)
		if err != nil {
			return "", err
		}

		// ✅ include original question for context
		final, err := a.llm.Generate(
			fmt.Sprintf("User asked: %s\nTool result: %s\nNow answer the user.", input, result),
			toolDefs,
		)
		if err != nil {
			return "", err
		}

		return final.Content, nil
	}

	return llmRes.Content, nil
}
