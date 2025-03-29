package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	// AI provider configuration
	AIProvider string `json:"aiProvider"` // e.g., "openai", "localai"
	AIModel    string `json:"aiModel"`    // e.g., "gpt-3.5-turbo"
	AIEndpoint string `json:"aiEndpoint"` // Custom endpoint if using a local model

	// Kubernetes configuration
	KubeConfig     string `json:"kubeConfig"`     // Path to kubeconfig file
	CurrentContext string `json:"currentContext"` // Current Kubernetes context

	// Application settings
	AnonymizeByDefault bool   `json:"anonymizeByDefault"` // Whether to anonymize sensitive data by default
	OutputFormat       string `json:"outputFormat"`       // e.g., "text", "json", "yaml"
}

// DefaultConfig returns a config with default values
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = ""
	}

	// Default kubeconfig path
	kubeConfigPath := os.Getenv("KUBECONFIG")
	if kubeConfigPath == "" {
		kubeConfigPath = filepath.Join(homeDir, ".kube", "config")
	}

	return &Config{
		AIProvider:         "openai",
		AIModel:            "gpt-3.5-turbo",
		AIEndpoint:         "",
		KubeConfig:         kubeConfigPath,
		CurrentContext:     "",
		AnonymizeByDefault: false,
		OutputFormat:       "text",
	}
}

// LoadConfig loads the configuration from environment variables or defaults
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Override with environment variables if present
	if provider := os.Getenv("K8SMED_AI_PROVIDER"); provider != "" {
		config.AIProvider = provider
	}

	if model := os.Getenv("K8SMED_AI_MODEL"); model != "" {
		config.AIModel = model
	}

	if endpoint := os.Getenv("K8SMED_AI_ENDPOINT"); endpoint != "" {
		config.AIEndpoint = endpoint
	}

	if anonymize := os.Getenv("K8SMED_ANONYMIZE_DEFAULT"); anonymize == "true" {
		config.AnonymizeByDefault = true
	}

	if format := os.Getenv("K8SMED_OUTPUT_FORMAT"); format != "" {
		config.OutputFormat = format
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// validateConfig ensures the configuration is valid
func validateConfig(config *Config) error {
	// Validate AI provider
	switch config.AIProvider {
	case "openai":
		if os.Getenv("OPENAI_API_KEY") == "" {
			return fmt.Errorf("OPENAI_API_KEY environment variable is required when using OpenAI")
		}
	case "localai":
		if config.AIEndpoint == "" {
			return fmt.Errorf("AIEndpoint is required when using LocalAI")
		}
	default:
		return fmt.Errorf("unsupported AI provider: %s", config.AIProvider)
	}

	return nil
}
