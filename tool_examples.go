package groq

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// CalculatorTool performs basic mathematical operations
type CalculatorTool struct{}

func (c *CalculatorTool) Name() string {
	return "calculator"
}

func (c *CalculatorTool) Description() string {
	return "Perform basic mathematical operations (add, subtract, multiply, divide)"
}

func (c *CalculatorTool) Parameters() *Parameter {
	return &Parameter{
		Type: "object",
		Properties: map[string]*Parameter{
			"operation": {
				Type:        "string",
				Description: "The mathematical operation to perform",
				Enum:        []interface{}{"add", "subtract", "multiply", "divide"},
			},
			"a": {
				Type:        "string",
				Description: "First number",
			},
			"b": {
				Type:        "string",
				Description: "Second number",
			},
		},
		Required: []string{"operation", "a", "b"},
	}
}

type CalculatorArgs struct {
	Operation string `json:"operation"`
	A         string `json:"a"`
	B         string `json:"b"`
}

func (c *CalculatorTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var calcArgs CalculatorArgs
	if err := json.Unmarshal(args, &calcArgs); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	// Parse string values to float64
	a, err := strconv.ParseFloat(calcArgs.A, 64)
	if err != nil {
		return "", fmt.Errorf("invalid number for a: %s", calcArgs.A)
	}

	b, err := strconv.ParseFloat(calcArgs.B, 64)
	if err != nil {
		return "", fmt.Errorf("invalid number for b: %s", calcArgs.B)
	}

	var result float64
	switch calcArgs.Operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return "", fmt.Errorf("unsupported operation: %s", calcArgs.Operation)
	}

	return fmt.Sprintf("%.2f", result), nil
}

// WeatherTool simulates weather information retrieval
type WeatherTool struct{}

func (w *WeatherTool) Name() string {
	return "get_weather"
}

func (w *WeatherTool) Description() string {
	return "Get current weather information for a specified location"
}

func (w *WeatherTool) Parameters() *Parameter {
	return &Parameter{
		Type: "object",
		Properties: map[string]*Parameter{
			"location": {
				Type:        "string",
				Description: "The city and state/country for weather lookup",
			},
			"unit": {
				Type:        "string",
				Description: "Temperature unit (celsius or fahrenheit)",
				Enum:        []interface{}{"celsius", "fahrenheit"},
			},
		},
		Required: []string{"location"},
	}
}

type WeatherArgs struct {
	Location string `json:"location"`
	Unit     string `json:"unit,omitempty"`
}

func (w *WeatherTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var weatherArgs WeatherArgs
	if err := json.Unmarshal(args, &weatherArgs); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	if weatherArgs.Unit == "" {
		weatherArgs.Unit = "celsius"
	}

	// Simulate weather data
	temp := 20.0 + math.Sin(float64(time.Now().Unix()%86400)/86400*2*math.Pi)*10
	if weatherArgs.Unit == "fahrenheit" {
		temp = temp*9/5 + 32
	}

	conditions := []string{"sunny", "cloudy", "rainy", "partly cloudy"}
	condition := conditions[time.Now().Unix()%int64(len(conditions))]

	unit := "°C"
	if weatherArgs.Unit == "fahrenheit" {
		unit = "°F"
	}

	return fmt.Sprintf("Weather in %s: %.1f%s, %s",
		weatherArgs.Location, temp, unit, condition), nil
}

// TextProcessorTool performs text manipulation operations
type TextProcessorTool struct{}

func (t *TextProcessorTool) Name() string {
	return "text_processor"
}

func (t *TextProcessorTool) Description() string {
	return "Process text with various operations like uppercase, lowercase, reverse, or word count"
}

func (t *TextProcessorTool) Parameters() *Parameter {
	return &Parameter{
		Type: "object",
		Properties: map[string]*Parameter{
			"text": {
				Type:        "string",
				Description: "The text to process",
			},
			"operation": {
				Type:        "string",
				Description: "The operation to perform on the text",
				Enum:        []interface{}{"uppercase", "lowercase", "reverse", "word_count", "char_count"},
			},
		},
		Required: []string{"text", "operation"},
	}
}

type TextProcessorArgs struct {
	Text      string `json:"text"`
	Operation string `json:"operation"`
}

func (t *TextProcessorTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var textArgs TextProcessorArgs
	if err := json.Unmarshal(args, &textArgs); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	switch textArgs.Operation {
	case "uppercase":
		return strings.ToUpper(textArgs.Text), nil
	case "lowercase":
		return strings.ToLower(textArgs.Text), nil
	case "reverse":
		runes := []rune(textArgs.Text)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes), nil
	case "word_count":
		words := strings.Fields(textArgs.Text)
		return fmt.Sprintf("%d", len(words)), nil
	case "char_count":
		return fmt.Sprintf("%d", len(textArgs.Text)), nil
	default:
		return "", fmt.Errorf("unsupported operation: %s", textArgs.Operation)
	}
}
