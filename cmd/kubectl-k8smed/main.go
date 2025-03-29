package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/k8smed/k8smed/pkg/ai/anonymizer"
	"github.com/k8smed/k8smed/pkg/ai/llm"
	"github.com/k8smed/k8smed/pkg/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config
var version = "0.1.0-alpha" // This can be set during build with -ldflags

var rootCmd = &cobra.Command{
	Use:   "kubectl-k8smed",
	Short: "K8sMed - AI-Powered Kubernetes First Responder",
	Long: `K8sMed is an AI-powered troubleshooting assistant designed to act as a first responder for Kubernetes clusters.
It leverages LLMs to diagnose issues, provide natural language explanations, and generate actionable remediation commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				fmt.Println("Error displaying help:", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	},
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze [query]",
	Short: "Analyze Kubernetes resources and provide troubleshooting insights",
	Long: `Analyze Kubernetes resources based on your query and provide AI-powered insights.
The analysis will include diagnostic information and potential remediation steps.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Extract flags
		explain, _ := cmd.Flags().GetBool("explain")
		anonymizeFlag, _ := cmd.Flags().GetBool("anonymize")

		// Override default anonymization based on flag
		anonymize := cfg.AnonymizeByDefault || anonymizeFlag

		// Join all arguments as the query
		query := args[0]
		for i := 1; i < len(args); i++ {
			query += " " + args[i]
		}

		fmt.Printf("Analyzing query: %s\n", query)
		fmt.Printf("Explain: %t, Anonymize: %t\n", explain, anonymize)

		// Create LLM client based on configuration
		llmClient, err := createLLMClient(cfg)
		if err != nil {
			fmt.Printf("Error creating LLM client: %v\n", err)
			os.Exit(1)
		}

		// Create a context with a reasonable timeout
		ctx := context.Background()

		// If anonymize is enabled, anonymize the query
		if anonymize {
			anon := anonymizer.NewAnonymizer()
			query = anon.Anonymize(query)
			fmt.Printf("Anonymized query: %s\n", query)
		}

		// Create a simple completion request
		req := llm.CompletionRequest{
			Model: cfg.AIModel,
			Messages: []llm.Message{
				{
					Role: "system",
					Content: `You are K8sMed, an AI-powered Kubernetes troubleshooting assistant.
Help diagnose issues in Kubernetes clusters based on the query provided.
Provide clear explanations and suggest remediation steps or commands when possible.`,
				},
				{
					Role:    "user",
					Content: query,
				},
			},
			MaxTokens:   500,
			Temperature: 0.7,
		}

		// Send the request to the LLM
		resp, err := llmClient.Complete(ctx, req)
		if err != nil {
			fmt.Printf("Error getting LLM response: %v\n", err)
			os.Exit(1)
		}

		// Print the response
		fmt.Println("\n" + resp.Content)

		// If explain flag is set, print additional details
		if explain {
			fmt.Printf("\nModel: %s\n", resp.Model)
			fmt.Printf("Tokens used: %d\n", resp.TokensUsed)
		}
	},
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start an interactive troubleshooting session",
	Long: `Start an interactive, conversational troubleshooting session with the AI assistant.
The assistant will maintain context throughout the session to provide more targeted help.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting interactive troubleshooting session...")
		fmt.Println("Interactive mode will be implemented in future versions.")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("K8sMed v" + version)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display or modify configuration",
	Run: func(cmd *cobra.Command, args []string) {
		// Pretty print the configuration
		cfgJSON, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(cfgJSON))
	},
}

func init() {
	// Initialize configuration
	cobra.OnInitialize(initConfig)

	// Add flags to analyze command
	analyzeCmd.Flags().BoolP("explain", "e", false, "Provide detailed explanations for the analysis")
	analyzeCmd.Flags().BoolP("anonymize", "a", false, "Anonymize sensitive information in queries")

	// Add commands to root command
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	// Skip config validation for version command which doesn't require API keys
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--version" || os.Args[1] == "-v") {
		// Use default config without validation for version command
		cfg = config.DefaultConfig()
		return
	}

	// Load and validate configuration for other commands
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}
}

// createLLMClient creates an LLM client based on the configuration
func createLLMClient(cfg *config.Config) (llm.Client, error) {
	var options llm.ClientOptions

	switch strings.ToLower(cfg.AIProvider) {
	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required for OpenAI provider")
		}
		options = llm.ClientOptions{
			APIKey:   apiKey,
			Endpoint: cfg.AIEndpoint, // May be empty for default
			Timeout:  30,
		}
	case "localai":
		if cfg.AIEndpoint == "" {
			return nil, fmt.Errorf("AIEndpoint must be set for LocalAI provider")
		}
		options = llm.ClientOptions{
			APIKey:   "", // LocalAI may not need an API key
			Endpoint: cfg.AIEndpoint,
			Timeout:  60, // Local models may be slower
		}
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.AIProvider)
	}

	return llm.NewClient(cfg.AIProvider, options)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
