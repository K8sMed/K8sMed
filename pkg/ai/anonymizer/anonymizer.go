package anonymizer

import (
	"regexp"
	"strings"
)

// These are the regular expressions used to detect sensitive information
var (
	// Email pattern
	emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// IP address pattern (both IPv4 and simple IPv6)
	ipRegex = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}`)

	// API key/token patterns (common formats)
	apiKeyRegex = regexp.MustCompile(`(api[-_]?key|token|secret|password|apikey)[\s]*[=:][\s]*["']?[a-zA-Z0-9_\-\.]{16,}["']?`)

	// UUID pattern
	uuidRegex = regexp.MustCompile(`[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}`)

	// Kubernetes namespace pattern
	namespaceRegex = regexp.MustCompile(`namespace\s+['"]?([a-z0-9]([-a-z0-9]*[a-z0-9])?)['"]?`)

	// Base64 pattern (for potential secrets)
	base64Regex = regexp.MustCompile(`(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})`)
)

// Anonymizer provides methods to anonymize sensitive data
type Anonymizer struct {
	// Additional patterns to detect (beyond the default ones)
	additionalPatterns []*regexp.Regexp

	// Custom replacements (pattern -> replacement)
	customReplacements map[string]string
}

// NewAnonymizer creates a new Anonymizer with default patterns
func NewAnonymizer() *Anonymizer {
	return &Anonymizer{
		additionalPatterns: make([]*regexp.Regexp, 0),
		customReplacements: make(map[string]string),
	}
}

// AddPattern adds a custom regex pattern to detect
func (a *Anonymizer) AddPattern(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	a.additionalPatterns = append(a.additionalPatterns, re)
	return nil
}

// AddReplacement adds a custom replacement for a specific string
func (a *Anonymizer) AddReplacement(original, replacement string) {
	a.customReplacements[original] = replacement
}

// Anonymize anonymizes sensitive data in the input text
func (a *Anonymizer) Anonymize(input string) string {
	// Apply custom replacements first
	output := input
	for original, replacement := range a.customReplacements {
		output = strings.ReplaceAll(output, original, replacement)
	}

	// Apply regex patterns
	patterns := []*regexp.Regexp{
		emailRegex,
		ipRegex,
		apiKeyRegex,
		uuidRegex,
		namespaceRegex,
	}

	// Add custom patterns
	patterns = append(patterns, a.additionalPatterns...)

	// Process all patterns
	for _, pattern := range patterns {
		output = pattern.ReplaceAllStringFunc(output, replaceWithType)
	}

	// Special handling for base64 (only replace long strings that look like they could be secrets)
	output = base64Regex.ReplaceAllStringFunc(output, func(match string) string {
		if len(match) >= 20 {
			return "[BASE64_DATA]"
		}
		return match
	})

	return output
}

// replaceWithType replaces a matched string with a type indicator
func replaceWithType(match string) string {
	lower := strings.ToLower(match)

	switch {
	case emailRegex.MatchString(match):
		return "[EMAIL]"
	case ipRegex.MatchString(match):
		return "[IP_ADDRESS]"
	case strings.Contains(lower, "api") || strings.Contains(lower, "key") || strings.Contains(lower, "token") || strings.Contains(lower, "secret"):
		return "[API_KEY]"
	case uuidRegex.MatchString(match):
		return "[UUID]"
	case strings.Contains(lower, "namespace"):
		parts := namespaceRegex.FindStringSubmatch(match)
		if len(parts) > 1 {
			return strings.Replace(match, parts[1], "[NAMESPACE_NAME]", 1)
		}
	}

	return "[REDACTED]"
}
