package ai

import (
	"context"
)

// ProviderKind is an enum that defines the different AI providers.
type ProviderKind string

const (
	// OpenAIProvider is the OpenAI provider.
	OpenAIProvider ProviderKind = "openai"
)

// Adapter is an interface that defines the methods that an AI adapter should implement.
type Adapter interface {
	// SimplePrompt prompts the AI model with a  prompt message and returns the response without choices.
	SimplePrompt(ctx context.Context, aiModel string, prompt string) (string, error)
	// GetProvider returns the provider kind of the AI adapter.
	GetProvider() ProviderKind
}
