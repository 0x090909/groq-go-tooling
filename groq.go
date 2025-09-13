package groq

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"log"
	"os"
)

type GroqGenerator struct {
	client *openai.LLM
}

func NewGroqGenerator() *GroqGenerator {
	llm, err := openai.New(
		openai.WithModel("meta-llama/llama-4-scout-17b-16e-instruct"),
		openai.WithBaseURL("https://api.groq.com/openai/v1"),
		openai.WithToken(os.Getenv("GROQ_API_KEY")),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &GroqGenerator{
		client: llm,
	}
}

func NewGroqGeneratorWithModel(model string) *GroqGenerator {
	llm, err := openai.New(
		openai.WithModel(model),
		openai.WithBaseURL("https://api.groq.com/openai/v1"),
		openai.WithToken(os.Getenv("GROQ_API_KEY")),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &GroqGenerator{
		client: llm,
	}
}

func (g *GroqGenerator) Generate(prompt string) string {
	ctx := context.Background()
	output, err := llms.GenerateFromSinglePrompt(ctx,
		g.client,
		prompt,
		llms.WithTemperature(0.15),
		llms.WithMaxTokens(100),
		llms.WithJSONMode(),
	)
	if err != nil {
		log.Fatalf("Error generating output: %v", err)
	}

	return output
}

func (g *GroqGenerator) GenerateWithOptions(system_prompt string, user_message string, option ...llms.CallOption) string {
	ctx := context.Background()
	output, err := g.client.GenerateContent(ctx,
		[]llms.MessageContent{
			{
				Role: "system",
				Parts: []llms.ContentPart{
					llms.TextContent{
						Text: system_prompt,
					},
				},
			},
			{
				Role: llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{
					llms.TextContent{
						Text: user_message,
					},
				},
			},
		},

		option...,
	)
	if err != nil {
		log.Fatalf("Error generating output: %v", err)
	}

	return output.Choices[0].Content
}

func (g *GroqGenerator) GenerateWithToolsNoExecution(system_prompt string, user_message string, tools []llms.Tool, option ...llms.CallOption) (*llms.ContentResponse, error) {
	ctx := context.Background()

	// Try tool choice auto to let the model decide when to use tools
	toolOptions := []llms.CallOption{
		llms.WithTools(tools),
		llms.WithToolChoice("auto"),
	}
	toolOptions = append(toolOptions, option...)

	output, err := g.client.GenerateContent(ctx,
		[]llms.MessageContent{
			{
				Role: "system",
				Parts: []llms.ContentPart{
					llms.TextContent{
						Text: system_prompt,
					},
				},
			},
			{
				Role: llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{
					llms.TextContent{
						Text: user_message,
					},
				},
			},
		},

		toolOptions...,
	)
	if err != nil {
		return nil, err
	}

	return output, nil
}
