package constants

const (
	OAuthClientID              = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
	OAuthAuthorizeEndpoint     = "https://platform.claude.com/oauth/authorize"
	OAuthClaudeAIAuthorizeURL  = "https://claude.com/cai/oauth/authorize"
	OAuthTokenEndpoint         = "https://platform.claude.com/v1/oauth/token"
	OAuthAPIKeyEndpoint        = "https://api.anthropic.com/api/oauth/claude_cli/create_api_key"
	OAuthRolesEndpoint         = "https://api.anthropic.com/api/oauth/claude_cli/roles"
	OAuthSuccessURL            = "https://platform.claude.com/oauth/code/success?app=claude-code"
	OAuthManualRedirectURL     = "https://platform.claude.com/oauth/code/callback"
	OAuthInferenceScope        = "user:inference"
	OAuthProfileScope          = "user:profile"
	OAuthSessionsScope         = "user:sessions:claude_code"
	OAuthMCPScope              = "user:mcp_servers"
	OAuthFileUploadScope       = "user:file_upload"
	OAuthConsoleScope          = "org:create_api_key"
	OAuthDefaultScopes         = "org:create_api_key user:profile user:inference user:sessions:claude_code user:mcp_servers user:file_upload"
	OAuthBetaHeader            = "oauth-2025-04-20"
)
