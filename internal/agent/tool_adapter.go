package agent

import "ai-agent/internal/domain"

func ToToolDefinitions(tools map[string]domain.Tool) []domain.ToolDefinition {
	defs := []domain.ToolDefinition{}

	for _, t := range tools {
		defs = append(defs, domain.ToolDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			Parameters:  t.Parameters(),
		})
	}

	return defs
}
