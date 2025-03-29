package llm

import (
	"context"
	"errors"
)

// Common errors
var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrAuthFailure    = errors.New("authentication failure")
	ErrAPILimit       = errors.New("API rate limit exceeded")
	ErrModelNotFound  = errors.New("specified model not found")
)

// Message represents a message in a conversation
type Message struct {
	Role    string `json:"role"`    // e.g., "system", "user", "assistant"
	Content string `json:"content"` // The message content
}

// CompletionRequest represents a request to an LLM for completion
type CompletionRequest struct {
	Messages    []Message `json:"messages"`              // The conversation history
	Model       string    `json:"model"`                 // The model to use
	MaxTokens   int       `json:"maxTokens,omitempty"`   // Maximum tokens to generate
	Temperature float64   `json:"temperature,omitempty"` // Sampling temperature (0.0-2.0)
}

// CompletionResponse represents the response from an LLM completion request
type CompletionResponse struct {
	Content      string `json:"content"`      // The generated content
	Model        string `json:"model"`        // The model used
	TokensUsed   int    `json:"tokensUsed"`   // Total tokens used (prompt + completion)
	FinishReason string `json:"finishReason"` // Why the model stopped generating
}

// Client is an interface for interacting with LLM providers
type Client interface {
	// Complete sends a completion request to the LLM and returns the response
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
}

// ClientOptions contains options for configuring an LLM client
type ClientOptions struct {
	APIKey   string // API key for authentication
	Endpoint string // Optional custom endpoint URL
	Timeout  int    // Timeout in seconds
}

// NewClient creates a new LLM client based on the provider name
func NewClient(provider string, options ClientOptions) (Client, error) {
	switch provider {
	case "openai":
		return NewOpenAIClient(options)
	case "localai":
		return NewLocalAIClient(options)
	default:
		return nil, errors.New("unsupported LLM provider: " + provider)
	}
}
