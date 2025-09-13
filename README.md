# Groq Go Tooling Library

A Go library for creating and managing function calling tools with Groq's API using LangChain Go. This library provides an easy way to define custom tools and execute them automatically with Groq's language models.

## Features

- **Easy Tool Creation**: Simple interface for defining custom tools with JSON schema parameters
- **Tool Registry**: Thread-safe centralized management of registered tools  
- **Groq Integration**: Direct integration with Groq API via LangChain Go
- **Automatic Tool Execution**: Seamless tool calling and result integration
- **Built-in Example Tools**: Calculator, Weather, and Text Processing tools included
- **Type Safety**: Strong typing for tool parameters and execution

## Installation

```bash
go mod init your-project
go get github.com/0x090909/groq-go-tooling
go get github.com/tmc/langchaingo
go get github.com/joho/godotenv
```

## Quick Start

### 1. Set up your environment

Create a `.env` file with your Groq API key:
```
GROQ_API_KEY=your_groq_api_key_here
```

### 2. Create and register tools

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/joho/godotenv"
    "github.com/tmc/langchaingo/llms"
    groq "github.com/0x090909/groq-go-tooling"
)

func main() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    
    // Create tool registry and register tools
    registry := groq.NewToolRegistry()
    
    // Register built-in tools
    registry.Register(&groq.CalculatorTool{})
    registry.Register(&groq.WeatherTool{})
    registry.Register(&groq.TextProcessorTool{})
    
    // Create Groq generator
    generator := groq.NewGroqGenerator()
    
    // Generate with automatic tool execution
    result, err := generator.GenerateWithRegistryExecution(
        "You are a helpful assistant. Use tools when needed.",
        "What is 45 plus 37?",
        registry,
        llms.WithTemperature(0.1),
    )
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(result)
}
```

## Core Components

### Tool Interface

All tools must implement the `Tool` interface:

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() *Parameter
    Execute(ctx context.Context, args json.RawMessage) (string, error)
}
```

### Tool Registry

The `ToolRegistry` manages tool registration and execution with these methods:

- `Register(tool Tool) error` - Register a new tool
- `Unregister(name string) error` - Remove a tool  
- `GetTool(name string) (Tool, error)` - Retrieve a specific tool
- `ListTools() []string` - Get all registered tool names
- `GetToolDefinitions() []*ToolDefinition` - Get Groq-compatible tool definitions
- `GetLangChainTools() []llms.Tool` - Get tools in LangChain format
- `Execute(ctx context.Context, toolCall *ToolCall) (*ToolResult, error)` - Execute a tool call

### GroqGenerator Methods

The `GroqGenerator` provides multiple ways to generate content:

- `NewGroqGenerator() *GroqGenerator` - Create with default model (meta-llama/llama-4-scout-17b-16e-instruct)
- `NewGroqGeneratorWithModel(model string) *GroqGenerator` - Create with custom model
- `Generate(prompt string) string` - Simple generation with JSON mode
- `GenerateWithOptions(systemPrompt, userMessage string, options ...llms.CallOption) string` - Generate with system/user prompts
- `GenerateWithToolsNoExecution(systemPrompt, userMessage string, tools []llms.Tool, options ...llms.CallOption) (*llms.ContentResponse, error)` - Generate with tools but don't execute them
- `GenerateWithRegistryExecution(systemPrompt, userMessage string, registry *ToolRegistry, options ...llms.CallOption) (string, error)` - Generate and automatically execute tools

### Built-in Example Tools

#### CalculatorTool
Performs basic mathematical operations (add, subtract, multiply, divide):
```go
// Parameters: operation (string), a (string), b (string)
// Returns: result as formatted string
```

#### WeatherTool  
Simulates weather information retrieval:
```go
// Parameters: location (string), unit (optional: "celsius"|"fahrenheit")
// Returns: simulated weather data
```

#### TextProcessorTool
Performs text manipulation operations:
```go
// Parameters: text (string), operation ("uppercase"|"lowercase"|"reverse"|"word_count"|"char_count")
// Returns: processed text or count
```

## Creating Custom Tools

Define a struct that implements the `Tool` interface:

```go
type MyCustomTool struct{}

func (t *MyCustomTool) Name() string {
    return "my_custom_tool"
}

func (t *MyCustomTool) Description() string {
    return "Description of what this tool does"
}

func (t *MyCustomTool) Parameters() *Parameter {
    return &Parameter{
        Type: "object",
        Properties: map[string]*Parameter{
            "param_name": {
                Type:        "string",
                Description: "Parameter description",
                Enum:        []interface{}{"option1", "option2"}, // Optional
            },
        },
        Required: []string{"param_name"},
    }
}

func (t *MyCustomTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        ParamName string `json:"param_name"`
    }
    
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }
    
    // Your tool logic here
    return "Tool result", nil
}
```

## Environment Setup

The library requires a `GROQ_API_KEY` environment variable. You can set this in a `.env` file:

```
GROQ_API_KEY=your_actual_api_key_here
```

## Available Models

You can use any Groq-supported model by creating a generator with a custom model:

```go
generator := groq.NewGroqGeneratorWithModel("llama-2-70b-4096")
```

## Error Handling

All operations return detailed errors. Tool execution failures are wrapped with context about which tool failed and why.