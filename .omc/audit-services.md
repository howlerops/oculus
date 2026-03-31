# TS → Go Services Audit

**Date**: 2026-03-31
**Scope**: `old-src/services/` (130 TS files) vs Go packages under `pkg/services/` and `pkg/api/`
**Auditor**: Executor agent

---

## Summary Table

| TS Directory | TS Files | Go Package | Rating |
|---|---|---|---|
| services/api/ | 20 files | pkg/api/ | PARTIAL |
| services/analytics/ | 8 files | pkg/services/analytics/ | PARTIAL |
| services/compact/ | 10 files | pkg/services/compact/ | PARTIAL |
| services/mcp/ | 20 files | pkg/services/mcp/ | PARTIAL |
| services/oauth/ | 5 files | pkg/auth/ | PARTIAL |
| services/plugins/ | 3 files | pkg/services/plugins/ | PARTIAL |
| services/policyLimits/ | 2 files | (none) | MISSING |
| services/remoteManagedSettings/ | 4 files | (none) | MISSING |
| services/AgentSummary/ | 1 file | (none) | MISSING |
| services/toolUseSummary/ | 1 file | (none) | MISSING |
| services/autoDream/ | 4 files | (none) | MISSING |
| services/SessionMemory/ | 3 files | (none) | MISSING |
| services/extractMemories/ | 2 files | (none) | MISSING |
| services/MagicDocs/ | 2 files | (none) | MISSING |
| services/lsp/ | 7 files | (covered in pkg/tools/lsp/) | STUB |
| services/tools/ | 4 files | (covered in pkg/tools/ + pkg/task/) | PARTIAL |
| services/tips/ | 3 files | (none) | MISSING |
| services/settingsSync/ | 2 files | (none) | MISSING |
| services/teamMemorySync/ | 5 files | (none) | MISSING |
| Root-level TS files | 12 files | scattered | PARTIAL/MISSING |

---

## Detailed Analysis

---

### services/api/ → pkg/api/

**TS files (20)**:
1. `adminRequests.ts` — Admin request API (limit_increase, seat_upgrade) via OAuth
2. `bootstrap.ts` — Bootstrap API call to fetch client_data, model options
3. `claude.ts` — Core LLM query function, streaming, multi-provider (Bedrock, Vertex, direct)
4. `client.ts` — Anthropic SDK client factory (multi-provider: API key, Bedrock, Vertex, OAuth)
5. `dumpPrompts.ts` — Debug prompt caching/dumping utility
6. `emptyUsage.ts` — Zero-initialized usage object constant
7. `errors.ts` — Comprehensive API error classification, human-readable messages
8. `errorUtils.ts` — SSL/TLS error code classification, connection error details
9. `filesApi.ts` — Files API (upload/download file attachments)
10. `firstTokenDate.ts` — Fetches/caches user's first Claude Code token date
11. `grove.ts` — Grove (sharing) feature API — enabled status, account settings
12. `logging.ts` — API response logging, telemetry spans, analytics events
13. `metricsOptOut.ts` — Metrics opt-out check with disk + memory TTL caching
14. `overageCreditGrant.ts` — Overage credit grant eligibility API
15. `promptCacheBreakDetection.ts` — Detects prompt cache misses, diffs, logs
16. `referral.ts` — Referral eligibility, redemptions, campaign API
17. `sessionIngress.ts` — Session transcript logging to remote API (sequential, retry)
18. `ultrareviewQuota.ts` — Ultra-review quota fetch
19. `usage.ts` — Rate limit utilization fetch (five-hour, seven-day windows)
20. `withRetry.ts` — Retry logic with backoff, fast-mode handling, OAuth 401 refresh

**Go files (4)**:
- `pkg/api/client.go` — HTTP client, streaming messages, `NewClient`, `SendMessage`, `StreamMessage` (303 lines)
- `pkg/api/errors.go` — `APIError` struct, `IsRateLimited`, `IsOverloaded`, `IsPromptTooLong`, `parseAPIError`
- `pkg/api/retry.go` — Exponential backoff with jitter, capped at 60s
- `pkg/api/client_test.go` — Tests

**Rating: PARTIAL**

**Implemented in Go**:
- Basic API client (direct API key only)
- Error types: rate limit, overload, prompt-too-long detection
- Retry with exponential backoff + jitter

**Missing in Go**:
- `adminRequests.ts` — Admin request types and API (no equivalent)
- `bootstrap.ts` — Bootstrap API call for client_data/model options
- `claude.ts` — Multi-provider support (Bedrock, Vertex, Foundry); streaming tool use; `queryHaiku`/`queryModelWithoutStreaming`; context window calculation; max-output-tokens per model
- `client.ts` — Multi-provider client factory (Bedrock/Vertex/GCP/AWS credentials, OAuth bearer, proxy settings)
- `dumpPrompts.ts` — Debug prompt dump (internal tool; may not be needed)
- `emptyUsage.ts` — Full `NonNullableUsage` struct with cache fields, server_tool_use, inference_geo
- `errors.ts` — Rich error classification: PDF limits, image resize, rate-limit messages per subscription type, `createAssistantAPIErrorMessage`
- `errorUtils.ts` — SSL/TLS error code detection
- `filesApi.ts` — Files API (upload/download)
- `firstTokenDate.ts` — First-token-date fetch and caching
- `grove.ts` — Grove sharing feature
- `logging.ts` — API response analytics logging, OTel spans, stop-reason tracking
- `metricsOptOut.ts` — Metrics opt-out with disk cache
- `overageCreditGrant.ts` — Overage credit grant
- `promptCacheBreakDetection.ts` — Cache break detection
- `referral.ts` — Referral API
- `sessionIngress.ts` — Session transcript ingress (sequential per-session)
- `ultrareviewQuota.ts` — Ultra-review quota
- `usage.ts` — Utilization fetch (rate limit windows)
- `withRetry.ts` — OAuth 401 retry, fast-mode handling, AWS/GCP credential refresh on error

---

### services/analytics/ → pkg/services/analytics/

**TS files (8)**:
1. `config.ts` — `isAnalyticsDisabled()`, `isFeedbackSurveyDisabled()`
2. `datadog.ts` — Datadog log batching, allowed event allowlist, flush
3. `firstPartyEventLogger.ts` — OpenTelemetry-based 1P event logger, sampling config
4. `firstPartyEventLoggingExporter.ts` — OTel exporter to Anthropic endpoint, file-based retry queue
5. `growthbook.ts` — GrowthBook feature flags, dynamic config, Statsig gate checks
6. `index.ts` — `logEvent()`, `attachAnalyticsSink()`, `stripProtoFields()` — main public API
7. `metadata.ts` — Event metadata enrichment (session ID, model, platform, org, git repo hash)
8. `sink.ts` — Routes events to Datadog and 1P; `initializeAnalyticsSink()`
9. `sinkKillswitch.ts` — Per-sink killswitch via GrowthBook config

**Go file (1)**:
- `pkg/services/analytics/analytics.go` — `Event` struct, `Sink` interface, `FileSink` (writes JSONL locally) (106 lines)

**Rating: PARTIAL**

**Implemented in Go**:
- Basic event struct and sink interface
- File-based JSONL event sink

**Missing in Go**:
- `config.ts` — Analytics disabled checks (env vars: Bedrock, Vertex, Foundry, privacy level)
- `datadog.ts` — Datadog integration, event allowlist, batching
- `firstPartyEventLogger.ts` — OTel-based 1P logging with sampling
- `firstPartyEventLoggingExporter.ts` — 1P exporter with retry queue, proto field hoisting
- `growthbook.ts` — Feature flag system (GrowthBook); all `getFeatureValue_CACHED_MAY_BE_STALE` / `getDynamicConfig_BLOCKS_ON_INIT` / `checkStatsigFeatureGate_CACHED_MAY_BE_STALE` calls
- `index.ts` — `logEvent()` public API with queue-before-sink semantics, `stripProtoFields()`
- `metadata.ts` — Full event metadata enrichment (platform, model, org UUID, git remote hash, WSL version, agent context, subscription type)
- `sink.ts` — Multi-sink dispatch (Datadog + 1P)
- `sinkKillswitch.ts` — Runtime killswitch per sink

---

### services/compact/ → pkg/services/compact/

**TS files (10)**:
1. `apiMicrocompact.ts` — API-side context management strategies (`clear_tool_uses_20250919`)
2. `autoCompact.ts` — Auto-compact trigger: context window check, calls `compactConversation`
3. `compact.ts` — Core compaction: summarize old messages via forked LLM agent, boundary messages
4. `compactWarningHook.ts` — React hook for compact warning suppression state
5. `compactWarningState.ts` — Store for compact warning suppression
6. `grouping.ts` — Groups messages by API round for compaction boundaries
7. `microCompact.ts` — Client-side micro-compaction: clear old tool results in-place
8. `postCompactCleanup.ts` — Post-compact cache/state cleanup
9. `prompt.ts` — Compact summary prompt template (detailed analysis + summary format)
10. `sessionMemoryCompact.ts` — Session memory compaction variant
11. `timeBasedMCConfig.ts` — GrowthBook config for time-based micro-compact

**Go file (1)**:
- `pkg/services/compact/compact.go` — `CompactMessages()`: splits messages, builds text summary, creates boundary message (83 lines, ~simple heuristic)

**Rating: PARTIAL**

**Implemented in Go**:
- Basic message list compaction (simple text summary, boundary message)

**Missing in Go**:
- `apiMicrocompact.ts` — API-side `clear_tool_uses_20250919` context management strategy
- `autoCompact.ts` — Auto-compact trigger (context window monitoring, threshold checks)
- `compact.ts` — LLM-based summarization (forked agent, streaming, proper summary prompt, boundary annotation, tool-use pairing)
- `compactWarningHook.ts` / `compactWarningState.ts` — Warning suppression state
- `grouping.ts` — API-round message grouping
- `microCompact.ts` — Client-side tool-result clearing (time-based, image size limits)
- `postCompactCleanup.ts` — Post-compact cleanup of caches, classifier approvals, beta tracing
- `prompt.ts` — Rich compact summary prompt with `<analysis>` / `<summary>` structure
- `sessionMemoryCompact.ts` — Session memory compaction integration
- `timeBasedMCConfig.ts` — Time-based micro-compact GrowthBook config

---

### services/mcp/ → pkg/services/mcp/

**TS files (20)**:
1. `auth.ts` — Full MCP OAuth 2.0 + PKCE flow (OIDC discovery, token exchange, keychain storage)
2. `channelAllowlist.ts` — Channel plugin allowlist via GrowthBook
3. `channelNotification.ts` — Inbound channel notifications (`notifications/claude/channel`)
4. `channelPermissions.ts` — Permission approval via channel notifications
5. `claudeai.ts` — Fetches claude.ai managed MCP servers via API
6. `client.ts` — Full MCP client: stdio, SSE, HTTP, WebSocket, SDK transports; tool/command/resource fetching; reconnect
7. `config.ts` — MCP config load/save (global + project, scoped, enterprise, plugin-sourced)
8. `elicitationHandler.ts` — MCP elicitation protocol handling
9. `envExpansion.ts` — `${VAR:-default}` env var expansion in server configs
10. `headersHelper.ts` — Dynamic headers from `headersHelper` script
11. `InProcessTransport.ts` — In-process linked transport pair (no subprocess)
12. `mcpStringUtils.ts` — `mcpInfoFromString()` — parse `mcp__server__tool` strings
13. `normalization.ts` — `normalizeNameForMCP()` — sanitize server names for API
14. `oauthPort.ts` — OAuth redirect port helpers (`buildRedirectUri`, `findAvailablePort`)
15. `officialRegistry.ts` — Official MCP registry URL fetch/cache
16. `SdkControlTransport.ts` — SDK↔CLI transport bridge via control messages
17. `types.ts` — Full config schemas: `McpStdioServerConfig`, `McpSSEServerConfig`, `McpHTTPServerConfig`, `McpWebSocketServerConfig`, scopes, transports, `MCPServerConnection`
18. `useManageMCPConnections.ts` — React hook for managing MCP connection lifecycle
19. `utils.ts` — MCP utility functions (permission checks, config lookup, agent info)
20. `vscodeSdkMcp.ts` — VS Code / IDE SDK MCP integration
21. `xaa.ts` — Cross-App Access (XAA): RFC 8693 token exchange, RFC 7523 JWT Bearer
22. `xaaIdpLogin.ts` — XAA IdP OIDC login flow (browser + PKCE)

**Go files (2)**:
- `pkg/services/mcp/client.go` — `Client` struct managing stdio connections; `Connect`, `ListTools`, `CallTool`, `Disconnect` (265 lines)
- `pkg/services/mcp/types.go` — `ServerConfig`, `ServerConnection`, `ConnectionStatus`, `ServerTool`, `ServerResource`, `ToolCallRequest`

**Rating: PARTIAL**

**Implemented in Go**:
- Stdio transport only (JSON-RPC over subprocess stdin/stdout)
- `Connect`/`ListTools`/`CallTool`/`Disconnect`
- Basic server config and connection status types

**Missing in Go**:
- `auth.ts` — Full MCP OAuth 2.0/PKCE/OIDC flow
- `channelAllowlist.ts` / `channelNotification.ts` / `channelPermissions.ts` — Channel system
- `claudeai.ts` — Claude.ai managed server fetching
- `client.ts` — SSE, HTTP, WebSocket, SDK transports; reconnect logic; command/resource fetching; elicitation; IDE RPC
- `config.ts` — Scoped config (global/project/local/enterprise/plugin/claudeai/managed), config read/write, enterprise file path
- `elicitationHandler.ts` — Elicitation protocol
- `envExpansion.ts` — `${VAR:-default}` expansion
- `headersHelper.ts` — Dynamic headers script
- `InProcessTransport.ts` — In-process transport
- `mcpStringUtils.ts` — `mcp__server__tool` string parsing
- `normalization.ts` — Name normalization (in Go types but not as standalone utility)
- `oauthPort.ts` — Port allocation for OAuth redirect
- `officialRegistry.ts` — Official registry check
- `SdkControlTransport.ts` — SDK↔CLI bridge transport
- `types.ts` — Full config schema variants (SSE, HTTP, WS, SDK transports; all scopes)
- `useManageMCPConnections.ts` — Connection lifecycle management
- `utils.ts` — MCP utilities, permission checks
- `vscodeSdkMcp.ts` — VS Code MCP integration
- `xaa.ts` / `xaaIdpLogin.ts` — Enterprise XAA/IdP flows

---

### services/oauth/ → pkg/auth/

**TS files (5)**:
1. `auth-code-listener.ts` — `AuthCodeListener` class: localhost HTTP server captures OAuth redirect
2. `client.ts` — OAuth client: full token exchange, refresh, profile fetch, `isOAuthTokenExpired`
3. `crypto.ts` — PKCE helpers: `generateCodeVerifier`, `generateCodeChallenge`, `generateState`
4. `getOauthProfile.ts` — Fetch OAuth profile from API key or OAuth token
5. `index.ts` — `OAuthService` class orchestrating the full PKCE authorization code flow

**Go files (4)**:
- `pkg/auth/auth.go` — `GetAuthToken()`: env var → keychain → OAuth flow (58 lines)
- `pkg/auth/oauth.go` — Full PKCE flow: code verifier/challenge, state, local HTTP server, browser open, token exchange (full implementation)
- `pkg/auth/keychain.go` — Secure storage for API keys and tokens
- `pkg/auth/trust.go` — Trust dialog acceptance

**Rating: PARTIAL**

**Implemented in Go**:
- Full PKCE authorization code flow
- Code verifier/challenge/state generation
- Local HTTP server for OAuth redirect capture
- Token storage in keychain
- Token refresh
- `GetAuthToken()` with env var / stored key / OAuth flow priority

**Missing in Go**:
- `client.ts` — `isClaudeAISubscriber()`, `getSubscriptionType()`, `hasProfileScope()`, `getOauthAccountInfo()` — subscription/account info checks used throughout the codebase
- `getOauthProfile.ts` — Profile fetch from API (subscription type, billing type, rate limit tier, referral info)
- `index.ts` — Manual auth code flow (non-browser environments), re-use of existing verifier
- Token expiry check callable from outside auth package
- Scoped token access (inference scope vs profile scope)

---

### services/plugins/ → pkg/services/plugins/

**TS files (3)**:
1. `pluginCliCommands.ts` — CLI wrappers: install, uninstall, enable, disable, update commands with analytics
2. `PluginInstallationManager.ts` — Background auto-install of plugins/marketplaces from trusted sources
3. `pluginOperations.ts` — Core plugin operations: install, uninstall, enable, disable, update (pure library)

**Go file (1)**:
- `pkg/services/plugins/plugins.go` — `PluginManifest`, `LoadedPlugin`, `Manager`, `LoadPlugins`, `GetPlugin`, `ListPlugins` (147 lines)

**Rating: PARTIAL**

**Implemented in Go**:
- Plugin manifest loading from disk
- Plugin enable/disable state
- List/get plugins by name

**Missing in Go**:
- `pluginCliCommands.ts` — CLI install/uninstall/enable/disable commands with analytics, process exit
- `PluginInstallationManager.ts` — Background auto-install, marketplace reconciliation, cache management
- `pluginOperations.ts` — Install from URL/path, version management, reverse dependency checking, orphaned version cleanup, marketplace support

---

### services/policyLimits/ → (no Go equivalent)

**TS files (2)**:
1. `index.ts` — Fetches org-level policy restrictions from API; ETag caching; background polling; retry; fails open
2. `types.ts` — `PolicyLimitsResponseSchema`, `PolicyLimitsFetchResult`

**Rating: MISSING**

No equivalent in Go. Policy limits control which CLI features are disabled per organization. Needed for enterprise deployments.

---

### services/remoteManagedSettings/ → (no Go equivalent)

**TS files (4)**:
1. `index.ts` — Fetches remote-managed settings for enterprise; checksum-based validation; ETag caching; background polling
2. `syncCache.ts` — Eligibility check for remote managed settings (auth-touching layer)
3. `syncCacheState.ts` — Leaf state for sync cache (no auth import; disk cache read)
4. `types.ts` — `RemoteManagedSettingsResponseSchema`, `RemoteManagedSettingsFetchResult`

**Rating: MISSING**

No equivalent in Go. Remote managed settings allow enterprises to push settings to CLI users without a release. Critical for enterprise/team deployments.

---

### services/AgentSummary/ → (no Go equivalent)

**TS files (1)**:
1. `agentSummary.ts` — Periodic (30s) background sub-agent progress summarization for coordinator mode

**Rating: MISSING**

No equivalent in Go. Needed for multi-agent coordinator mode UI.

---

### services/toolUseSummary/ → (no Go equivalent)

**TS files (1)**:
1. `toolUseSummaryGenerator.ts` — Generates human-readable tool-batch summaries via Haiku for SDK clients

**Rating: MISSING**

No equivalent in Go.

---

### services/autoDream/ → (no Go equivalent)

**TS files (4)**:
1. `autoDream.ts` — Background memory consolidation (time-gate + session-count gate + lock)
2. `config.ts` — `isAutoDreamEnabled()` from settings + GrowthBook
3. `consolidationLock.ts` — Lock file (mtime = lastConsolidatedAt, PID-guarded stale check)
4. `consolidationPrompt.ts` — Consolidation prompt builder

**Rating: MISSING**

No equivalent in Go. Related to the memory system (`pkg/memdir/memdir.go` exists but auto-dream background consolidation is absent).

---

### services/SessionMemory/ → (no Go equivalent)

**TS files (3)**:
1. `prompts.ts` — Session memory template, `isSessionMemoryEmpty`, `truncateSessionMemoryForCompact`
2. `sessionMemory.ts` — Background session memory extraction (forked agent, periodic)
3. `sessionMemoryUtils.ts` — `getSessionMemoryContent`, `waitForSessionMemoryExtraction`, `getLastSummarizedMessageId`

**Rating: MISSING**

No equivalent in Go. Session memory is used by compact, auto-compact, and away-summary.

---

### services/extractMemories/ → (no Go equivalent)

**TS files (2)**:
1. `extractMemories.ts` — End-of-query background memory extraction (forked agent)
2. `prompts.ts` — Memory extraction prompt templates

**Rating: MISSING**

No equivalent in Go. Part of the memory system.

---

### services/MagicDocs/ → (no Go equivalent)

**TS files (2)**:
1. `magicDocs.ts` — Auto-maintains markdown files with `# MAGIC DOC:` headers
2. `prompts.ts` — Magic Docs update prompt template

**Rating: MISSING**

No equivalent in Go.

---

### services/lsp/ → pkg/tools/lsp/

**TS files (7)**:
1. `config.ts` — Get all configured LSP servers from plugins
2. `LSPClient.ts` — Full LSP client (JSON-RPC over stdio, initialize, textDocument/diagnostics)
3. `LSPDiagnosticRegistry.ts` — LRU cache for pending LSP diagnostics
4. `LSPServerInstance.ts` — Per-server instance: initialize, shutdown, request routing, retry on content-modified
5. `LSPServerManager.ts` — Multi-server manager: route by file extension
6. `manager.ts` — Global singleton manager initialization
7. `passiveFeedback.ts` — Maps LSP severity, registers `publishDiagnostics` notifications

**Go file (1)**:
- `pkg/tools/lsp/lsp.go` — Tool definition for LSP diagnostics (tool interface, not a service)

**Rating: STUB**

The Go `pkg/tools/lsp/lsp.go` is a tool definition (Claude tool), not an LSP service client. The full LSP client infrastructure (server process management, JSON-RPC protocol, diagnostic registry, multi-server routing) is entirely absent.

---

### services/tools/ → pkg/tools/ + pkg/task/

**TS files (4)**:
1. `StreamingToolExecutor.ts` — Streaming tool execution with concurrency, queuing, yielding
2. `toolExecution.ts` — `runToolUse()`: execute a single tool call, analytics logging
3. `toolHooks.ts` — Pre/post tool hooks, permission prompts, analytics events
4. `toolOrchestration.ts` — `runTools()`: parallel tool execution up to concurrency limit

**Go packages**: `pkg/tools/*/` (individual tool implementations), `pkg/task/task.go`

**Rating: PARTIAL**

**Implemented in Go**:
- Individual tool implementations exist (bash, fileread, filewrite, fileedit, glob, grep, etc.)
- `pkg/task/task.go` has task execution

**Missing in Go**:
- `StreamingToolExecutor.ts` — Streaming interleaving of tool results with assistant chunks
- `toolExecution.ts` — Unified `runToolUse()` with analytics, telemetry, MCP tool routing
- `toolHooks.ts` — Pre/post tool hook execution pipeline
- `toolOrchestration.ts` — Concurrent tool execution (`CLAUDE_CODE_MAX_TOOL_USE_CONCURRENCY`)

---

### services/tips/ → (no Go equivalent)

**TS files (3)**:
1. `tipHistory.ts` — Records tip shown history in global config by startup count
2. `tipRegistry.ts` — Registry of available tips with eligibility conditions
3. `tipScheduler.ts` — Selects tip with longest time since shown

**Rating: MISSING**

No equivalent in Go.

---

### services/settingsSync/ → (no Go equivalent)

**TS files (2)**:
1. `index.ts` — Syncs user settings and memory files across environments (upload local → remote)
2. `types.ts` — `UserSyncContentSchema`, `UserSyncDataSchema`

**Rating: MISSING**

No equivalent in Go.

---

### services/teamMemorySync/ → (no Go equivalent)

**TS files (5)**:
1. `index.ts` — Syncs team memory files per-repo (GET/PUT API, delta upload)
2. `secretScanner.ts` — Client-side gitleaks-based secret scanner before upload
3. `teamMemSecretGuard.ts` — Guard called from FileWriteTool/FileEditTool
4. `types.ts` — `TeamMemoryContentSchema`, `TeamMemoryDataSchema`
5. `watcher.ts` — fs.watch for team memory dir, debounced push on change

**Rating: MISSING**

No equivalent in Go. The `pkg/memdir/memdir.go` exists but team memory sync (API upload/download, secret scanning, file watching) is absent.

---

### Root-level TS service files

| File | Key Exports | Go Equivalent | Rating |
|---|---|---|---|
| `claudeAiLimits.ts` | `currentLimits()`, `processRateLimitHeaders()`, rate limit quota tracking, `getRateLimitWarning()` | None | MISSING |
| `claudeAiLimitsHook.ts` | React hook for limits state | None (UI concern) | N/A |
| `awaySummary.ts` | `generateAwaySummary()` — 1-sentence recap of session for "away" screen | None | MISSING |
| `diagnosticTracking.ts` | `DiagnosticTrackingService` — IDE diagnostic fetch via MCP/LSP | None | MISSING |
| `internalLogging.ts` | `logInternalToolEvent()` — Kubernetes namespace detection, internal usage logging | None | MISSING |
| `mockRateLimits.ts` | Mock rate limit headers for testing/demo | None | MISSING |
| `notifier.ts` | `sendNotification()` — system notification + hooks | None | MISSING |
| `preventSleep.ts` | macOS `caffeinate` wrapper to prevent sleep during work | None | MISSING |
| `rateLimitMessages.ts` | `getRateLimitErrorMessage()`, `getRateLimitWarning()`, `getUsingOverageText()` | None | MISSING |
| `rateLimitMocking.ts` | Facade: `processRateLimitHeaders()` with mock injection | None | MISSING |
| `tokenEstimation.ts` | `roughTokenCountEstimation()`, `countTokens()` — multi-provider token counting (API, Bedrock, Vertex) | `pkg/services/tokencount/counter.go` — character-based estimation only | STUB |
| `vcr.ts` | VCR-style API request recording/replay for testing | None | MISSING |
| `voice.ts` | Push-to-talk audio recording (cpal native, sox, arecord fallback) | None | MISSING |
| `voiceKeyterms.ts` | Key term extraction for voice input | None | MISSING |
| `voiceStreamSTT.ts` | Streaming speech-to-text | None | MISSING |

---

## Cross-cutting Gaps

### 1. Feature Flag System (GrowthBook)
The TS codebase uses `getFeatureValue_CACHED_MAY_BE_STALE()`, `getDynamicConfig_BLOCKS_ON_INIT()`, `checkStatsigFeatureGate_CACHED_MAY_BE_STALE()` pervasively across all services. **No Go equivalent exists.** This means any behavior controlled by GrowthBook flags is hardcoded or absent in Go.

### 2. OAuth Account/Subscription Awareness
Many TS services gate behavior on `isClaudeAISubscriber()`, `getSubscriptionType()`, `hasProfileScope()`, `getOauthAccountInfo()`. The Go `pkg/auth/` has token acquisition but no subscription/profile awareness.

### 3. Multi-provider API
The Go `pkg/api/` only supports direct API key auth to `api.anthropic.com`. The TS codebase supports Bedrock (AWS SigV4), Vertex AI (GCP), Foundry, and claude.ai OAuth as drop-in providers.

### 4. Memory System
The TS codebase has a rich memory system: `SessionMemory`, `extractMemories`, `autoDream`, `teamMemorySync`, `settingsSync`, `MagicDocs`. Go has `pkg/memdir/memdir.go` (directory paths) only.

### 5. Remote Configuration
`remoteManagedSettings` and `policyLimits` are both absent in Go. These are critical for enterprise deployments.

### 6. Analytics Depth
Go has a local file sink only. The TS system has: Datadog batching, 1P OTel exporter with retry queue, proto-field PII tagging, per-event sampling, per-sink killswitch, GrowthBook gate.

---

## Ratings Legend

- **FULL**: All key functionality ported with equivalent behavior
- **PARTIAL**: Core structure exists but significant features missing (listed above)
- **STUB**: File/package exists but is a minimal placeholder or unrelated implementation
- **MISSING**: No Go equivalent exists at all
