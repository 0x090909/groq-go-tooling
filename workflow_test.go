package groq

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
)

func TestCompleteWorkflow(t *testing.T) {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		t.Logf("Warning: Error loading .env file: %v", err)
	}

	// Check if API key is available
	if os.Getenv("GROQ_API_KEY") == "" {
		t.Skip("GROQ_API_KEY not set. Please set it in .env file to test the complete workflow.")
	}

	// Create tool registry and register example tools
	registry := NewToolRegistry()

	tools := []Tool{
		&CalculatorTool{},
		&WeatherTool{},
		&TextProcessorTool{},
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			t.Fatalf("Failed to register tool %s: %v", tool.Name(), err)
		}
		t.Logf("✓ Registered tool: %s", tool.Name())
	}

	// Create generator
	generator := NewGroqGenerator()

	// Test scenarios
	testScenarios := []struct {
		name         string
		systemPrompt string
		userMessage  string
	}{
		{
			name:         "Calculator Test",
			systemPrompt: "You are a helpful math assistant. Use tools when needed to provide accurate calculations.",
			userMessage:  "What is 45 plus 37?",
		},
		{
			name:         "Weather Test",
			systemPrompt: "You are a helpful weather assistant. Use tools to get current weather information.",
			userMessage:  "What's the weather like in London?",
		},
		{
			name:         "Text Processing Test",
			systemPrompt: "You are a helpful text processing assistant. Use tools to manipulate text as requested.",
			userMessage:  "Can you reverse the text 'Hello World' for me?",
		},
		{
			name:         "Multiple Operations",
			systemPrompt: "You are a helpful assistant that can use multiple tools to complete complex tasks.",
			userMessage:  "Calculate 25 * 4, then tell me how many characters are in the result when written as text",
		},
	}

	// Execute test scenarios
	for i, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			t.Logf("User: %s", scenario.userMessage)

			result, err := generator.GenerateWithRegistryExecution(
				scenario.systemPrompt,
				scenario.userMessage,
				registry,
				llms.WithTemperature(0.1),
				llms.WithMaxTokens(500),
			)

			if err != nil {
				t.Errorf("Error in scenario %d (%s): %v", i+1, scenario.name, err)
				return
			}

			if result == "" {
				t.Errorf("Expected non-empty result for scenario %s", scenario.name)
				return
			}

			t.Logf("Assistant: %s", result)
		})
	}
}
