package groq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"math"
	"strings"
)

// GenerateWithRegistryExecution generates content with tools and executes them automatically
func (g *GroqGenerator) GenerateWithRegistryExecution(systemPrompt string, userMessage string, registry *ToolRegistry, options ...llms.CallOption) (string, error) {
	ctx := context.Background()

	// Get tools from registry in langchain format
	tools := registry.GetLangChainTools()

	//manage iterations & Execute all tool calls
	var toolResults []string
	executed_iterations := 0
	MAX_ITERATIONS := 10
	returned_tool_calls := math.MaxInt
	for executed_iterations < MAX_ITERATIONS && returned_tool_calls > 0 {
		// Generate with tools
		response, err := g.GenerateWithToolsNoExecution(systemPrompt, userMessage, tools, options...)
		if err != nil {
			return "", fmt.Errorf("failed to generate with tools: %w", err)
		}

		// Check if there are any tool calls to execute
		if len(response.Choices) == 0 {
			return "", fmt.Errorf("no response choices received")
		}

		choice := response.Choices[0]

		// If no tool calls, return the content directly
		returned_tool_calls = len(choice.ToolCalls)
		/*if len(choice.ToolCalls) == 0 {
			return choice.Content, nil
		}*/

		// Execute each tool call
		for _, toolCall := range choice.ToolCalls {
			// Convert langchain tool call to our format for execution
			groqToolCall := &ToolCall{
				ID:   toolCall.ID,
				Type: "function",
				Function: &FunctionCall{
					Name:      toolCall.FunctionCall.Name,
					Arguments: json.RawMessage(toolCall.FunctionCall.Arguments),
				},
			}

			// Execute the tool
			result, err := registry.Execute(ctx, groqToolCall)
			if err != nil {
				return "", fmt.Errorf("tool execution failed for %s: %w", toolCall.FunctionCall.Name, err)
			}

			toolResults = append(toolResults, fmt.Sprintf("Tool '%s' returned: %s", toolCall.FunctionCall.Name, result.Content))
		}

		executed_iterations++
	}

	// Create a follow-up prompt with tool results
	followUpPrompt := fmt.Sprintf("Based on the following tool results, please provide a comprehensive answer to the user's question:\n\n%s\n\nOriginal question: %s",
		strings.Join(toolResults, "\n"), userMessage)

	// Generate final response
	finalResponse := g.GenerateWithOptions(systemPrompt, followUpPrompt, options...)
	return finalResponse, nil
}
