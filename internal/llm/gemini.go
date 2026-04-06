package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"ai-agent/internal/domain"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

func NewGeminiClient(apiKey string) *GeminiClient {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		panic(fmt.Sprintf("failed to create Gemini client: %v", err))
	}

	return &GeminiClient{
		client: client,
		model:  "gemini-2.5-flash",
	}
}

func (g *GeminiClient) Generate(input string, tools []domain.ToolDefinition) (domain.LLMResponse, error) {
	ctx := context.Background()
	model := g.client.GenerativeModel(g.model)

	// Convert tools to Gemini format
	var geminiTools []*genai.Tool
	for _, t := range tools {
		schema := toGeminiSchema(t.Parameters)
		geminiTools = append(geminiTools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        t.Name,
					Description: t.Description,
					Parameters:  schema,
				},
			},
		})
	}
	model.Tools = geminiTools

	resp, err := model.GenerateContent(ctx, genai.Text(input))
	if err != nil {
		return domain.LLMResponse{}, err
	}

	candidate := resp.Candidates[0]

	for _, part := range candidate.Content.Parts {
		// Check for tool call
		if fc, ok := part.(genai.FunctionCall); ok {
			args := map[string]interface{}{}
			for k, v := range fc.Args {
				args[k] = v
			}
			return domain.LLMResponse{
				ToolCall: &domain.ToolCall{
					Name:  fc.Name,
					Input: args,
				},
			}, nil
		}

		// Text response
		if txt, ok := part.(genai.Text); ok {
			return domain.LLMResponse{
				Content: string(txt),
			}, nil
		}
	}

	return domain.LLMResponse{}, fmt.Errorf("empty response from Gemini")
}

func toGeminiSchema(params map[string]interface{}) *genai.Schema {
	schema := &genai.Schema{
		Type:       genai.TypeObject,
		Properties: map[string]*genai.Schema{},
	}

	props, ok := params["properties"].(map[string]interface{})
	if !ok {
		return schema
	}

	for name, val := range props {
		prop := val.(map[string]interface{})
		desc, _ := prop["description"].(string)
		schema.Properties[name] = &genai.Schema{
			Type:        genai.TypeString,
			Description: desc,
		}
	}

	if required, ok := params["required"].([]string); ok {
		schema.Required = required
	}

	// handle required as []interface{} too
	if required, ok := params["required"].([]interface{}); ok {
		for _, r := range required {
			schema.Required = append(schema.Required, r.(string))
		}
	}

	_ = json.Marshal
	return schema
}
