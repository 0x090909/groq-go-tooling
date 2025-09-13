package groq

import (
	"context"
	"encoding/json"
)

// Parameter represents a parameter in a tool's JSON schema
type Parameter struct {
	Type        string                `json:"type"`
	Description string                `json:"description,omitempty"`
	Properties  map[string]*Parameter `json:"properties,omitempty"`
	Required    []string              `json:"required,omitempty"`
	Items       *Parameter            `json:"items,omitempty"`
	Enum        []interface{}         `json:"enum,omitempty"`
}

// ToolDefinition represents a tool definition for Groq API
type ToolDefinition struct {
	Type     string    `json:"type"`
	Function *Function `json:"function"`
}

// Function represents the function part of a tool definition
type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  *Parameter `json:"parameters"`
}

// ToolCall represents a tool call from the model
type ToolCall struct {
	ID       string        `json:"id"`
	Type     string        `json:"type"`
	Function *FunctionCall `json:"function"`
}

// FunctionCall represents the function call details
type FunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ToolResult represents the result of executing a tool
type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Role       string `json:"role"`
	Content    string `json:"content"`
}

// Tool interface that all tools must implement
type Tool interface {
	Name() string
	Description() string
	Parameters() *Parameter
	Execute(ctx context.Context, args json.RawMessage) (string, error)
}

// ToolExecutor handles tool execution
type ToolExecutor interface {
	Execute(ctx context.Context, toolCall *ToolCall) (*ToolResult, error)
}
