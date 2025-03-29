package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LocalAIClient implements the Client interface for LocalAI/Ollama and similar local LLM servers
type LocalAIClient struct {
	endpoint   string
	apiKey     string // Some local servers may support API keys
	httpClient *http.Client
}

// NewLocalAIClient creates a new LocalAI client
func NewLocalAIClient(options ClientOptions) (*LocalAIClient, error) {
	if options.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for LocalAI")
	}

	timeout := 60 // Local models may be slower
	if options.Timeout > 0 {
		timeout = options.Timeout
	}

	return &LocalAIClient{
		endpoint: options.Endpoint,
		apiKey:   options.APIKey, // May be empty for many local servers
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, nil
}

// Complete implements the Client interface for LocalAI
func (c *LocalAIClient) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	// Convert our request to a format compatible with LocalAI/Ollama
	// This uses the same format as OpenAI for compatibility
	localAIReq := struct {
		Model       string          `json:"model"`
		Messages    []openAIMessage `json:"messages"`
		MaxTokens   int             `json:"max_tokens,omitempty"`
		Temperature float64         `json:"temperature,omitempty"`
	}{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Messages:    make([]openAIMessage, len(req.Messages)),
	}

	// Copy messages
	for i, msg := range req.Messages {
		localAIReq.Messages[i] = openAIMessage(msg)
	}

	// Marshal request to JSON
	reqBody, err := json.Marshal(localAIReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.endpoint,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response (assuming OpenAI-compatible format)
	var localAIResp openAIResponse
	if err := json.Unmarshal(body, &localAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if we have choices
	if len(localAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no completions returned")
	}

	// Convert to our response format
	response := &CompletionResponse{
		Content:      localAIResp.Choices[0].Message.Content,
		Model:        localAIResp.Model,
		TokensUsed:   localAIResp.Usage.TotalTokens,
		FinishReason: localAIResp.Choices[0].FinishReason,
	}

	return response, nil
}
