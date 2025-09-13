package groq

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestGroq(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		t.Skip("Warning: Error loading .env file, skipping test")
	}

	groq := NewGroqGenerator()
	output := groq.GenerateWithOptions("you are an expert that answers in a rude way", "I want to get the best grade what to do.")

	if output == "" {
		t.Error("Expected non-empty output from Groq generator")
	}

	t.Logf("Groq output: %s", output)
}
