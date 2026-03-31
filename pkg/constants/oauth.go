package constants

const (
	OAuthClientID          = "claude-code"
	OAuthAuthorizeEndpoint = "https://console.anthropic.com/oauth/authorize"
	OAuthTokenEndpoint     = "https://console.anthropic.com/oauth/token"
	OAuthProfileEndpoint   = "https://console.anthropic.com/api/profile"
	OAuthInferenceScope    = "user:inference"
	OAuthProfileScope      = "user:profile"
	OAuthDefaultScopes     = "user:inference user:profile"
)
