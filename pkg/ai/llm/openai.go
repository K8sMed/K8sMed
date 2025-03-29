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

// OpenAIClient implements the Client interface for OpenAI
type OpenAIClient struct {
	apiKey     string
	endpoint   string
	httpClient *http.Client
}

// openAIRequest represents the request structure for OpenAI API
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

// openAIMessage represents a message in the OpenAI API format
type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIResponse represents the response structure from OpenAI API
type openAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(options ClientOptions) (*OpenAIClient, error) {
	if options.APIKey == "" {
		return nil, ErrAuthFailure
	}

	endpoint := "https://api.openai.com/v1/chat/completions"
	if options.Endpoint != "" {
		endpoint = options.Endpoint
	}

	timeout := 30
	if options.Timeout > 0 {
		timeout = options.Timeout
	}

	return &OpenAIClient{
		apiKey:   options.APIKey,
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, nil
}

// Complete implements the Client interface for OpenAI
func (c *OpenAIClient) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	// Convert to OpenAI format
	openAIReq := openAIRequest{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Messages:    make([]openAIMessage, len(req.Messages)),
	}

	// Copy messages
	for i, msg := range req.Messages {
		openAIReq.Messages[i] = openAIMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Marshal request to JSON
	reqBody, err := json.Marshal(openAIReq)
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
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

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
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error.Message != "" {
			return nil, fmt.Errorf("API error: %s", errorResp.Error.Message)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var openAIResp openAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if we have choices
	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no completions returned")
	}

	// Convert to our response format
	response := &CompletionResponse{
		Content:      openAIResp.Choices[0].Message.Content,
		Model:        openAIResp.Model,
		TokensUsed:   openAIResp.Usage.TotalTokens,
		FinishReason: openAIResp.Choices[0].FinishReason,
	}

	return response, nil
}
