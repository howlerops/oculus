package api

// ProviderConfig holds provider-specific settings for routing API requests
type ProviderConfig struct {
	Type      string // "firstParty", "bedrock", "vertex"
	BaseURL   string
	APIKey    string
	Region    string // AWS region for Bedrock
	ProjectID string // GCP project for Vertex
}

// GetProviderBaseURL returns the API base URL for the named provider.
// An empty or unrecognised provider name falls back to the default Anthropic URL.
func GetProviderBaseURL(provider string) string {
	switch provider {
	case "bedrock":
		return "https://bedrock-runtime.us-east-1.amazonaws.com"
	case "vertex":
		return "https://us-central1-aiplatform.googleapis.com"
	default:
		return DefaultBaseURL
	}
}

// IsFirstPartyProvider reports whether the given base URL targets the default
// Anthropic API (as opposed to Bedrock or Vertex).
func IsFirstPartyProvider(baseURL string) bool {
	return baseURL == "" || baseURL == DefaultBaseURL
}
