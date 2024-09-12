package openai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"github.com/vx416/smart-git/ai"
)

const (
	// GPT3Dot5Turbo0125 is the GPT-3.5-turbo-0125 model which is a chat model and the most cheap one.
	GPT3Dot5Turbo0125 = openai.GPT3Dot5Turbo0125
	// GPT4oMini is the GPT-4o-mini model which is a chat model and has balanced performance and cost.
	GPT4oMini = openai.GPT4oMini
)

const (
	SystemAI = "system"
)

// NewAdapter creates a new OpenAI adapter.
func NewAdapter(apiKey string, defaultAIModel string, defaultAIRole string) Adapter {
	client := openai.NewClient(apiKey)
	return Adapter{
		Client:         client,
		defaultAIModel: defaultAIModel,
		defaultAIRole:  defaultAIRole,
	}
}

var (
	_ ai.Adapter = (*Adapter)(nil)
)

// Adapter is an AI adapter that uses the OpenAI API.
type Adapter struct {
	*openai.Client
	defaultAIModel string
	defaultAIRole  string
}

// SimplePrompt prompts the AI model with a prompt message and returns the response without choices.
func (a Adapter) SimplePrompt(ctx context.Context, aiModel string, prompt string) (string, error) {
	if aiModel == "" {
		aiModel = a.defaultAIModel
	}
	aiRole := a.defaultAIRole

	resp, err := a.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0125,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    aiRole,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}
	return resp.Choices[0].Message.Content, nil
}

// GetProvider returns the provider kind of the AI adapter.
func (a Adapter) GetProvider() ai.ProviderKind {
	return ai.OpenAIProvider
}
