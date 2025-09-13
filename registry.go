package groq

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"sync"
)

// ToolRegistry manages registered tools
type ToolRegistry struct {
	tools map[string]Tool
	mutex sync.RWMutex
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry
func (r *ToolRegistry) Register(tool Tool) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	name := tool.Name()
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s already registered", name)
	}

	r.tools[name] = tool
	return nil
}

// Unregister removes a tool from the registry
func (r *ToolRegistry) Unregister(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tools[name]; !exists {
		return fmt.Errorf("tool %s not found", name)
	}

	delete(r.tools, name)
	return nil
}

// GetTool retrieves a tool by name
func (r *ToolRegistry) GetTool(name string) (Tool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}

	return tool, nil
}

// ListTools returns all registered tool names
func (r *ToolRegistry) ListTools() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}

	return names
}

// GetToolDefinitions returns Groq-compatible tool definitions
func (r *ToolRegistry) GetToolDefinitions() []*ToolDefinition {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	definitions := make([]*ToolDefinition, 0, len(r.tools))
	for _, tool := range r.tools {
		definitions = append(definitions, &ToolDefinition{
			Type: "function",
			Function: &Function{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters:  tool.Parameters(),
			},
		})
	}

	return definitions
}

// Execute runs a tool with the given parameters
func (r *ToolRegistry) Execute(ctx context.Context, toolCall *ToolCall) (*ToolResult, error) {
	if toolCall.Type != "function" {
		return nil, fmt.Errorf("unsupported tool type: %s", toolCall.Type)
	}

	tool, err := r.GetTool(toolCall.Function.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get tool: %w", err)
	}

	result, err := tool.Execute(ctx, toolCall.Function.Arguments)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	return &ToolResult{
		ToolCallID: toolCall.ID,
		Role:       "tool",
		Content:    result,
	}, nil
}

// ToLangChainFormat converts a Parameter to LangChain's expected format
func (p *Parameter) ToLangChainFormat() map[string]any {
	if p == nil {
		return nil
	}

	result := make(map[string]any)

	if p.Type != "" {
		result["type"] = p.Type
	}
	if p.Description != "" {
		result["description"] = p.Description
	}
	if len(p.Enum) > 0 {
		result["enum"] = p.Enum
	}
	if len(p.Required) > 0 {
		result["required"] = p.Required
	}

	if len(p.Properties) > 0 {
		properties := make(map[string]any)
		for key, prop := range p.Properties {
			properties[key] = prop.ToLangChainFormat()
		}
		result["properties"] = properties
	}

	if p.Items != nil {
		result["items"] = p.Items.ToLangChainFormat()
	}

	return result
}

// GetLangChainTools returns tools in LangChain format
func (r *ToolRegistry) GetLangChainTools() []llms.Tool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tools := make([]llms.Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		langchainTool := llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters:  tool.Parameters().ToLangChainFormat(),
			},
		}
		tools = append(tools, langchainTool)
	}

	return tools
}
