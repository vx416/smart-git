package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/vx416/smart-git/ai"
	"github.com/vx416/smart-git/ai/openai"
)

// NewAIAdapter creates a new AI adapter based on the provider kind and AI model.
func NewAIAdapter(providerKind ai.ProviderKind, aiModel string) (ai.Adapter, error) {
	if providerKind == "" {
		providerKind = ai.ProviderKind(strings.ToLower(os.Getenv("AI_PROVIDER")))
	}
	if providerKind == "" {
		return nil, fmt.Errorf("ai provider is required, please set AI_PROVIDER environment variable (e.g export AI_PROVIDER=openai)")
	}
	apiKey := os.Getenv("AI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ai api key is required, please set AI_API_KEY environment variable (e.g export AI_API_KEY=your_api_key)")
	}
	switch providerKind {
	case ai.OpenAIProvider:
		adapter := openai.NewAdapter(apiKey, openai.GPT3Dot5Turbo0125, openai.SystemAI)
		return adapter, nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", providerKind)
	}
}
