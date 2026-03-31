# Audit: Remaining TS Modules → Go Port
Generated: 2026-03-31

---

## Methodology
For each TS file: read first 30-40 lines, identify key exports/logic, compare to Go target.

**Rating scale:**
- **FULL** – All key exports/behaviour present in Go
- **PARTIAL** – Some structure exists but major logic is missing or simplified
- **MISSING** – No corresponding Go code found, or Go file is a stub with different semantics

---

## 1. `old-src/bridge/` (31 files) → `pkg/services/bridge/`

The Go target has exactly **one file**: `bridge.go` (~100 lines). It is a generic in-process message-channel abstraction (inbox/outbox channels, status enum). The entire TS bridge subsystem implements a production-grade Remote Control (CCR / Environments API / Session Ingress / CCR v2) protocol stack. Nothing of the real protocol is ported.

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `bridgeApi.ts` | `createBridgeApiClient`, `BridgeFatalError`, `validateBridgeId`, HTTP wrappers for Environments API (pollForWork, registerBridgeEnvironment, heartbeatWork, reconnectSession, archiveSession) | No Go equivalent | **MISSING** |
| `bridgeConfig.ts` | `getBridgeAccessToken`, `getBridgeBaseUrl`, `getBridgeTokenOverride`, `getBridgeBaseUrlOverride` — auth/URL resolution layer | No Go equivalent | **MISSING** |
| `bridgeDebug.ts` | `BridgeDebugHandle`, fault-injection for test recovery paths, `/bridge-kick` slash command hook | No Go equivalent | **MISSING** |
| `bridgeEnabled.ts` | `isBridgeEnabled`, `isBridgeEnabledBlocking`, `isEnvLessBridgeEnabled`, `isClaudeAISubscriber` — feature-gate checks via GrowthBook | No Go equivalent | **MISSING** |
| `bridgeMain.ts` | Multisession bridge orchestrator: poll loop, session spawner, capacity management, JWT refresh, `createBridgeLogger` wiring | No Go equivalent | **MISSING** |
| `bridgeMessaging.ts` | `isSDKMessage`, `handleIngressMessage`, `handleServerControlRequest`, `makeResultMessage`, `BoundedUUIDSet`, `extractTitleText` — transport-layer message parsing | No Go equivalent | **MISSING** |
| `bridgePermissionCallbacks.ts` | `BridgePermissionCallbacks` type + factory, `isBridgePermissionResponse` predicate — permission request/response lifecycle | No Go equivalent | **MISSING** |
| `bridgePointer.ts` | `writeBridgePointer`, `clearBridgePointer`, `findActiveBridgePointer` — crash-recovery file pointer (mtime-based TTL, worktree scan) | No Go equivalent | **MISSING** |
| `bridgeStatusUtil.ts` | `StatusState` type, `buildBridgeConnectUrl`, `buildBridgeSessionUrl`, `buildActiveFooterText`, `buildIdleFooterText`, `formatDuration`, `timestamp`, `abbreviateActivity` | No Go equivalent | **MISSING** |
| `bridgeUI.ts` | `createBridgeLogger` — terminal output (QR code, shimmer spinner, OSC8 links, status state machine rendering) | No Go equivalent | **MISSING** |
| `capacityWake.ts` | `createCapacityWake`, `CapacityWake`, `CapacitySignal` — AbortSignal merge for at-capacity sleep | No Go equivalent | **MISSING** |
| `codeSessionApi.ts` | `createCodeSession`, `fetchRemoteCredentials` — CCR v2 `/v1/code/sessions` and `/bridge` HTTP wrappers | No Go equivalent | **MISSING** |
| `createSession.ts` | `createBridgeSession` — POST /v1/sessions with event pre-population | No Go equivalent | **MISSING** |
| `debugUtils.ts` | `redactSecrets`, `debugTruncate`, `debugBody`, `extractErrorDetail`, `describeAxiosError` | No Go equivalent | **MISSING** |
| `envLessBridgeConfig.ts` | `EnvLessBridgeConfig` schema + `getEnvLessBridgeConfig` — GrowthBook-tunable config for CCR v2 (retries, heartbeat, timeouts, version floor) | No Go equivalent | **MISSING** |
| `flushGate.ts` | `FlushGate<T>` class — state machine for queuing messages during initial session flush | No Go equivalent | **MISSING** |
| `inboundAttachments.ts` | `resolveInboundAttachments` — download file_uuid attachments from OAuth `/api/oauth/files/{uuid}/content`, write to `~/.claude/uploads/` | No Go equivalent | **MISSING** |
| `inboundMessages.ts` | `extractInboundMessageFields` — parse inbound user messages, normalize camelCase `mediaType` | No Go equivalent | **MISSING** |
| `initReplBridge.ts` | `initReplBridge` — REPL wrapper around initBridgeCore; reads session ID, cwd, OAuth tokens, git context, GrowthBook flags | No Go equivalent | **MISSING** |
| `jwtUtils.ts` | `decodeJwtPayload`, `decodeJwtExpiry`, `createTokenRefreshScheduler` — JWT decode + proactive refresh scheduling | No Go equivalent | **MISSING** |
| `pollConfig.ts` | `getPollIntervalConfig` — GrowthBook-tunable poll interval config with Zod validation | No Go equivalent | **MISSING** |
| `pollConfigDefaults.ts` | `DEFAULT_POLL_CONFIG`, `PollIntervalConfig` — static defaults for poll intervals | No Go equivalent | **MISSING** |
| `remoteBridgeCore.ts` | `initEnvLessBridgeCore` — CCR v2 env-less bridge (direct OAuth→worker_jwt exchange, no Environments API) | No Go equivalent | **MISSING** |
| `replBridge.ts` | `initBridgeCore` — main env-based bridge core (~2400 lines): poll loop, WebSocket/SSE transport, session reconnect, permission callbacks | No Go equivalent | **MISSING** |
| `replBridgeHandle.ts` | `setReplBridgeHandle`, `getReplBridgeHandle`, `getSelfBridgeCompatId` — global singleton bridge handle | No Go equivalent | **MISSING** |
| `replBridgeTransport.ts` | `ReplBridgeTransport` interface, `createV1ReplTransport`, `createV2ReplTransport` — v1 (HybridTransport) / v2 (SSETransport + CCRClient) factory | No Go equivalent | **MISSING** |
| `sessionIdCompat.ts` | `toCompatSessionId`, `toInfraSessionId`, `setCseShimGate` — `cse_*` ↔ `session_*` tag translation | No Go equivalent | **MISSING** |
| `sessionRunner.ts` | `createSessionSpawner`, `safeFilenameId`, `PermissionRequest` — spawn child CLI processes for multisession bridge | No Go equivalent | **MISSING** |
| `trustedDevice.ts` | `getTrustedDeviceToken`, `enrollTrustedDevice`, `isTrustedDeviceEnrollmentEligible` — X-Trusted-Device-Token keychain storage + CCR elevated auth | No Go equivalent | **MISSING** |
| `types.ts` | `WorkData`, `WorkResponse`, `WorkSecret`, `BridgeConfig`, `BridgeApiClient`, `SessionActivity`, `SpawnMode`, `BridgeLogger`, constants | No Go equivalent for protocol types | **MISSING** |
| `workSecret.ts` | `decodeWorkSecret`, `buildSdkUrl`, `buildCCRv2SdkUrl`, `registerWorker`, `sameSessionId` — work secret decode + URL building + CCR worker registration | No Go equivalent | **MISSING** |

**Section summary:** The Go `bridge.go` is a placeholder in-memory channel with no relation to the actual Remote Control protocol. **All 31 TS files are MISSING from Go.**

---

## 2. `old-src/cli/` (19 files) → `cmd/claude-go/`, `pkg/entrypoints/`

### Top-level CLI files

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `cli/exit.ts` | `cliError`, `cliOk` — exit helpers with `: never` return type | `pkg/commands/` has no equivalent; `main.go` uses `os.Exit` inline | **MISSING** |
| `cli/ndjsonSafeStringify.ts` | `ndjsonSafeStringify` — JSON stringify escaping U+2028/U+2029 for NDJSON streams | No Go equivalent | **MISSING** |
| `cli/print.ts` | `runPrintMode` — non-interactive `-p` mode: full pipeline (tools, system prompt, structured IO, streaming, bridge, settings sync) | `main.go` has minimal `-p` flag with simple RunOnce — extremely simplified | **PARTIAL** |
| `cli/remoteIO.ts` | `RemoteIO` extends `StructuredIO` — bidirectional SDK streaming with CCR WebSocket/SSE transport, session state callbacks, flush gate | No Go equivalent | **MISSING** |
| `cli/structuredIO.ts` | `StructuredIO` — full SDK stdio protocol handler: elicitation, permission prompts, hook callbacks, control request/response, MCP tool bridging | No Go equivalent | **MISSING** |
| `cli/update.ts` | `update` — auto-update logic: fetch latest version, install via npm/native installer, regenerate completion cache | No Go equivalent | **MISSING** |

### `cli/handlers/` (6 files)

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `handlers/agents.ts` | `agentsHandler` — `claude agents` subcommand: lists configured agents with model/source info | No Go equivalent | **MISSING** |
| `handlers/auth.ts` | `loginHandler`, `logoutHandler` — full OAuth flow, API key creation, role fetch, profile fetch, SSL hints | `pkg/auth/auth.go` has basic `GetAuthToken` / `Logout` — no OAuth web flow, no role fetch | **PARTIAL** |
| `handlers/autoMode.ts` | `autoModeDefaultsHandler`, `autoModeCritiqueHandler` — `claude auto-mode` subcommand | No Go equivalent | **MISSING** |
| `handlers/mcp.tsx` | `mcpAddHandler`, `mcpRemoveHandler`, `mcpListHandler` etc. — `claude mcp *` subcommands | `pkg/entrypoints/mcp_entry.go` is a minimal JSON-RPC stub; no `mcp` subcommand management | **MISSING** |
| `handlers/plugins.ts` | Plugin/marketplace subcommand handlers: install, uninstall, enable, disable, list, update | `pkg/services/plugins/plugins.go` has manager but no CLI handlers | **MISSING** |
| `handlers/util.tsx` | `renderAndRunTUI`, dialog/TUI helpers for CLI handlers | No Go equivalent | **MISSING** |

### `cli/transports/` (7 files)

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `transports/ccrClient.ts` | `CCRClient`, `CCRInitError` — CCR v2 worker client (SSE + PUT /worker/state + SerialBatchEventUploader + WorkerStateUploader + heartbeat) | No Go equivalent | **MISSING** |
| `transports/HybridTransport.ts` | `HybridTransport` — WebSocket reads + HTTP POST writes, 100ms batch flush, 3s close grace | No Go equivalent | **MISSING** |
| `transports/SerialBatchEventUploader.ts` | `SerialBatchEventUploader`, `RetryableError` — serial ordered uploader: batch drain, exponential backoff, backpressure | No Go equivalent | **MISSING** |
| `transports/SSETransport.ts` | `SSETransport`, `StreamClientEvent` — SSE reads + POST writes, reconnect (10min budget), liveness timeout, permanent HTTP codes | No Go equivalent | **MISSING** |
| `transports/transportUtils.ts` | `getTransportForUrl` — selects WebSocket / Hybrid / SSE transport based on env vars | No Go equivalent | **MISSING** |
| `transports/WebSocketTransport.ts` | `WebSocketTransport` — WebSocket with circular buffer, keep-alive frames, 10min reconnect budget, mTLS, proxy | No Go equivalent | **MISSING** |
| `transports/WorkerStateUploader.ts` | `WorkerStateUploader` — coalescing PUT /worker uploader: RFC 7396 merge, 2-slot bounded, exponential backoff | No Go equivalent | **MISSING** |

**Section summary:** 17/19 files MISSING; 2 PARTIAL (main `-p` mode, basic login/logout).

---

## 3. `old-src/skills/` (20 files) → `pkg/skills/`

The TS skills module has 3 entry files + 17 bundled skill definitions.

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `bundledSkills.ts` | `BundledSkillDefinition` type, `getBundledSkill`, `loadBundledSkillContent`, `extractSkillFiles` — bundled skills registry with file extraction to disk | `pkg/skills/loader.go` has `Skill` struct + `LoadSkills` but no bundled skill concept, no `BundledSkillDefinition`, no file extraction | **PARTIAL** |
| `loadSkillsDir.ts` | `loadSkillsDir`, `createSkillCommand`, `parseSkillFrontmatterFields` — full skill loader: frontmatter parsing, argument substitution, effort levels, `ignore` patterns, token estimation, MCP registration | `loader.go` reads `.md` files but has no frontmatter parsing, no argument substitution, no effort levels, no token estimation | **PARTIAL** |
| `mcpSkillBuilders.ts` | `MCPSkillBuilders` type, `registerMCPSkillBuilders`, `getMCPSkillBuilders` — write-once registry bridging MCP skill discovery to loadSkillsDir | No Go equivalent | **MISSING** |
| `bundled/*.ts` (17 files) | Individual bundled skill definitions: batch, claudeApi, claudeApiContent, claudeInChrome, debug, index, keybindings, loop, loremIpsum, remember, scheduleRemoteAgents, simplify, skillify, stuck, updateConfig, verify, verifyContent | No Go equivalent for any bundled skills | **MISSING** |

**Section summary:** Core loader PARTIAL; MCP bridge and all 17 bundled skills MISSING.

---

## 4. `old-src/memdir/` (8 files) → `pkg/memdir/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `findRelevantMemories.ts` | `findRelevantMemories` — scan memory headers + LLM side-query (Sonnet) to select top-5 relevant memories; `RelevantMemory` type | No Go equivalent | **MISSING** |
| `memdir.ts` | `buildMemoryPrompt`, `buildCombinedMemoryPrompt`, `ENTRYPOINT_NAME`, `MAX_ENTRYPOINT_LINES`, `getAutoMemPath` re-export, team mem feature flag — main memory system prompt builder | `memdir.go` has `ListMemoryFiles`, `SaveMemoryFile`, `DeleteMemoryFile` — file ops only, no prompt building | **PARTIAL** |
| `memoryAge.ts` | `memoryAgeDays`, `memoryAge`, `memoryFreshnessText` — human-readable age + staleness caveat | No Go equivalent | **MISSING** |
| `memoryScan.ts` | `scanMemoryFiles`, `formatMemoryManifest`, `MemoryHeader` — scan `.md` files, parse frontmatter, return sorted header list (cap 200) | No Go equivalent (Go only reads content, not frontmatter/headers) | **MISSING** |
| `memoryTypes.ts` | `MEMORY_TYPES`, `MemoryType`, `parseMemoryType`, `TYPES_SECTION_INDIVIDUAL`, `TYPES_SECTION_COMBINED`, `MEMORY_FRONTMATTER_EXAMPLE`, `TRUSTING_RECALL_SECTION`, `WHAT_NOT_TO_SAVE_SECTION`, `WHEN_TO_ACCESS_SECTION` — taxonomy + prompt copy | No Go equivalent | **MISSING** |
| `paths.ts` | `isAutoMemoryEnabled`, `getAutoMemPath`, `getProjectMemPath` — multi-source memory path resolution (env, bare mode, CCR, settings, git root) | `memdir.go` has `IsAutoMemoryEnabled` (env only) + simplified `GetMemoryDir` — lacks settings/CCR/git resolution | **PARTIAL** |
| `teamMemPaths.ts` | `PathTraversalError`, `getTeamMemPath`, `validateTeamMemPath`, `sanitizePathKey` — team memory path with traversal protection | No Go equivalent | **MISSING** |
| `teamMemPrompts.ts` | `buildCombinedMemoryPrompt` — combined private+team memory system prompt with XML-scoped type guidance | No Go equivalent | **MISSING** |

**Section summary:** File-level CRUD exists; all semantic logic (relevance search, prompt building, frontmatter, age, team mem) MISSING.

---

## 5. `old-src/migrations/` (11 files) → `pkg/migrations/`

The TS migrations are each single-purpose one-shot config key migrations. The Go implementation has a generic versioned migration runner.

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `migrateAutoUpdatesToSettings.ts` | Move `autoUpdates` config key to `settings.json` | Generic framework exists; this specific migration not registered | **MISSING** |
| `migrateBypassPermissionsAcceptedToSettings.ts` | Move `bypassPermissionsModeAccepted` → `skipDangerousModePermissionPrompt` in settings | Not registered | **MISSING** |
| `migrateEnableAllProjectMcpServersToSettings.ts` | Move MCP server approval fields from project config to local settings | Not registered | **MISSING** |
| `migrateFennecToOpus.ts` | Rename `fennec-latest` / `fennec-fast-latest` / `opus-4-5-fast` aliases to Opus 4.6 equivalents | Not registered | **MISSING** |
| `migrateLegacyOpusToCurrent.ts` | Remap explicit Opus 4.0/4.1 model strings to `opus` alias | Not registered | **MISSING** |
| `migrateOpusToOpus1m.ts` | Migrate `opus` → `opus[1m]` for eligible Max/Team Premium users | Not registered | **MISSING** |
| `migrateReplBridgeEnabledToRemoteControlAtStartup.ts` | Rename config key `replBridgeEnabled` → `remoteControlAtStartup` | Not registered | **MISSING** |
| `migrateSonnet1mToSonnet45.ts` | Pin `sonnet[1m]` users to `sonnet-4-5-20250929[1m]` | Not registered | **MISSING** |
| `migrateSonnet45ToSonnet46.ts` | Migrate Pro/Max/Team Premium users from Sonnet 4.5 → `sonnet` alias (4.6) | Not registered | **MISSING** |
| `resetAutoModeOptInForDefaultOffer.ts` | Clear `skipAutoPermissionPrompt` for users who saw old 2-option dialog | Not registered | **MISSING** |
| `resetProToOpusDefault.ts` | Reset Pro users to Opus default | Not registered | **MISSING** |
| *(framework)* | `RunMigrations`, `Migration`, `MigrationState` — versioned runner with JSON state file | `migrations.go` has `RunMigrations`, `BuiltInMigrations` (1 init migration) | **PARTIAL** |

**Section summary:** Migration framework PARTIAL; all 11 specific migrations MISSING.

---

## 6. `old-src/entrypoints/` (8 files) → `pkg/entrypoints/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `agentSdkTypes.ts` | Re-exports all SDK public types: `SDKControlRequest`, `SDKControlResponse`, `SDKMessage` (via sdk/coreTypes), runtime types, settings types, tool types | `pkg/entrypoints/types.go` has `EntrypointMode`, `InitOptions`, `SDKStatus` — a fraction of the type surface | **PARTIAL** |
| `cli.tsx` | Bootstrap entrypoint: env setup (`COREPACK_ENABLE_AUTO_PIN=0`, `NODE_OPTIONS` for CCR, ablation baseline), dynamic import main for fast `--version` path | `cmd/claude-go/main.go` — basic cobra CLI, no env pre-setup, no fast-path | **PARTIAL** |
| `init.ts` | `init`, `initializeTelemetryAfterTrust` — full app init: telemetry, policy limits, remote managed settings, OAuth prefetch, CA certs, graceful shutdown, LSP shutdown, repository detection, diagnostics | `pkg/entrypoints/` has no init function | **MISSING** |
| `mcp.ts` | `startMCPServer` — full MCP server: `tools/list` returns real tool list with JSON schemas, `tools/call` dispatches real tool execution with permission context | `mcp_entry.go` returns empty `tools/list`; no tool dispatch | **PARTIAL** |
| `sandboxTypes.ts` | `SandboxNetworkConfigSchema`, `SandboxConfig`, network/filesystem sandbox config with Zod validation | No Go equivalent | **MISSING** |
| `sdk/*.ts` (3 files) | `controlSchemas.ts`, `coreSchemas.ts`, `coreTypes.ts` — SDK wire types: `SDKMessage`, `StdoutMessage`, `StdinMessage`, `SDKControlRequest`, elicitation schemas | No Go equivalent | **MISSING** |

**Section summary:** Mostly PARTIAL or MISSING; core SDK type surface and real MCP dispatch absent.

---

## 7. `old-src/bootstrap/` (1 file) → `pkg/auth/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `bootstrap/state.ts` | Global mutable session state: `getSessionId`, `setSessionId`, `getOriginalCwd`, `getProjectRoot`, `getSessionCounter`, `setMeter` (OpenTelemetry), `switchSession`, `setStatsStore`, `getKairosActive`, `getAllowedChannels`, `AttributedCounter`, `ChannelEntry` — ~60+ getters/setters | `pkg/auth/` has OAuth/keychain. No session state module exists anywhere in Go | **MISSING** |

**Section summary:** The bootstrap/state module (global session state) is entirely absent. MISSING.

---

## 8. `old-src/buddy/` (6 files) → `pkg/services/bridge/`

The TS "buddy" is an optional companion sprite system (pet companion in the TUI). It has no relation to the bridge.

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `buddy/companion.ts` | `getCompanion`, `generateCompanion` — seeded PRNG companion generation (species, rarity, stats, hat, eyes) | No Go equivalent | **MISSING** |
| `buddy/CompanionSprite.tsx` | React component rendering the companion sprite in TUI | No Go equivalent | **MISSING** |
| `buddy/prompt.ts` | `companionIntroText`, `getCompanionIntroAttachment` — injects companion into system prompt | No Go equivalent | **MISSING** |
| `buddy/sprites.ts` | ASCII sprite data for all species | No Go equivalent | **MISSING** |
| `buddy/types.ts` | `RARITIES`, `SPECIES`, `HATS`, `EYES`, `STAT_NAMES`, `Companion`, `CompanionBones` types | No Go equivalent | **MISSING** |
| `buddy/useBuddyNotification.tsx` | React hook for companion speech bubble notifications | No Go equivalent | **MISSING** |

**Section summary:** Feature entirely absent in Go. All 6 MISSING (feature parity not required for core functionality, but noted).

---

## 9. `old-src/plugins/` (2 files) → `pkg/services/plugins/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `plugins/builtinPlugins.ts` | `registerBuiltinPlugin`, `getBuiltinPlugins`, `getBuiltinPluginsForLoading`, `BUILTIN_MARKETPLACE_NAME` — built-in plugin registry with enable/disable state | `plugins.go` has `Manager` with `LoadFromDirectory` — external plugins only; no built-in registry concept | **PARTIAL** |
| `plugins/bundled/` | Bundled plugin definitions (directory with plugin implementations) | Not present | **MISSING** |

**Section summary:** Plugin manager framework PARTIAL; built-in plugin registry and bundled plugins MISSING.

---

## 10. `old-src/remote/` (4 files) → `pkg/services/bridge/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `remote/remotePermissionBridge.ts` | `createSyntheticAssistantMessage` — wraps CCR permission requests as synthetic AssistantMessage for the REPL permission UI | No Go equivalent | **MISSING** |
| `remote/RemoteSessionManager.ts` | `RemoteSessionManager` — manages a remote REPL session viewed locally: subscribes to `SessionsWebSocket`, dispatches SDK messages, handles permission requests, forwards control responses | No Go equivalent | **MISSING** |
| `remote/sdkMessageAdapter.ts` | `convertSDKMessageToMessage`, `adaptSDKMessagesToMessages` — converts CCR SDK-format messages to internal REPL `Message` types | No Go equivalent | **MISSING** |
| `remote/SessionsWebSocket.ts` | `SessionsWebSocket` — WebSocket client for `/v1/sessions/{id}/events` stream; OAuth-authed, reconnect with 4001/4003/4004/4010 close code handling, ping keepalive | No Go equivalent | **MISSING** |

**Section summary:** All 4 MISSING.

---

## 11. `old-src/server/` (3 files) → `pkg/services/bridge/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `server/createDirectConnectSession.ts` | `createDirectConnectSession`, `DirectConnectError` — POST `/sessions` to a local `claude serve` server, validate response | No Go equivalent | **MISSING** |
| `server/directConnectManager.ts` | `DirectConnectManager`, `DirectConnectConfig`, `DirectConnectCallbacks` — WebSocket client for local server session; message/permission dispatch | No Go equivalent | **MISSING** |
| `server/types.ts` | `ServerConfig`, `SessionState`, `SessionInfo`, `connectResponseSchema` — local server types | No Go equivalent | **MISSING** |

**Section summary:** All 3 MISSING. (Note: a `claude serve` mode / local server is not implemented in Go at all.)

---

## 12. `old-src/assistant/` (1 file) → `pkg/entrypoints/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `assistant/sessionHistory.ts` | `createHistoryAuthCtx`, `fetchHistoryPage`, `HISTORY_PAGE_SIZE`, `HistoryPage`, `HistoryAuthCtx` — paginated session event history via `/v1/sessions/{id}/events` (OAuth-authed, 100 events/page) | `pkg/services/history/history.go` stores local conversation list (title, turns) — no remote session event fetching | **MISSING** |

---

## 13. `old-src/coordinator/` (1 file) → `pkg/services/bridge/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `coordinator/coordinatorMode.ts` | `isCoordinatorMode`, `getCoordinatorSystemPrompt`, `getCoordinatorUserContext` — multi-agent coordinator mode system prompt and tool allowlist | No Go equivalent | **MISSING** |

---

## 14. `old-src/voice/` (1 file) → `pkg/services/bridge/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `voice/voiceModeEnabled.ts` | `isVoiceGrowthBookEnabled`, `hasVoiceAuth`, `isVoiceModeEnabled` — voice mode feature gate (GrowthBook kill-switch + OAuth check) | No Go equivalent | **MISSING** |

---

## 15. `old-src/upstreamproxy/` (2 files) → `pkg/services/bridge/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `upstreamproxy/relay.ts` | `startNodeRelay`, `startBunRelay` — CONNECT-over-WebSocket relay: TCP listen, HTTP CONNECT from local tools, tunnel bytes over WS to CCR upstreamproxy with protobuf framing | No Go equivalent | **MISSING** |
| `upstreamproxy/upstreamproxy.ts` | `startUpstreamProxy` — container-side wiring: read session token from `/run/ccr/session_token`, `prctl(PR_SET_DUMPABLE)`, download CA cert, start relay, set `HTTPS_PROXY` / `SSL_CERT_FILE` | No Go equivalent | **MISSING** |

---

## 16. `old-src/native-ts/` (4 files) → `pkg/utils/platform/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `native-ts/color-diff/index.ts` | Color diff computation (native binding) | `pkg/utils/platform/platform.go` has `GetPlatform`, `IsMacOS`, etc. — no color diff | **MISSING** |
| `native-ts/file-index/index.ts` | Native file indexer (binding) | No Go equivalent | **MISSING** |
| `native-ts/yoga-layout/index.ts` | Yoga layout engine bindings | No Go equivalent | **MISSING** |
| `native-ts/yoga-layout/enums.ts` | Yoga layout enums | No Go equivalent | **MISSING** |

**Section summary:** The Go `platform.go` covers OS detection only. Native bindings not ported (expected — different Go equivalents would be needed).

---

## 17. `old-src/moreright/` (1 file) → `pkg/utils/permissions/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `moreright/useMoreRight.tsx` | `useMoreRight` — stub for internal-only ant hook (external build stub: `onBeforeQuery` returns true, `onTurnComplete` no-op, `render` returns null) | `pkg/utils/permissions/checker.go` has `CheckToolPermission`, `IsPathAllowed` — no `useMoreRight` equivalent needed externally | **N/A** (stub only; intentionally no-op in external builds) |

---

## 18. `old-src/outputStyles/` (1 file) → `pkg/tui/`

| TS File | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|
| `outputStyles/loadOutputStylesDir.ts` | `getOutputStyleDirStyles` — memoized loader: scans `.claude/output-styles/*.md` (project + user), parses frontmatter, returns `OutputStyleConfig[]` | No Go equivalent | **MISSING** |

---

## 19. Top-Level Files

| TS File | Go Target | Key Exports / Purpose | Go Status | Rating |
|---|---|---|---|---|
| `main.tsx` | `cmd/claude-go/main.go` | Full CLI entry: commander setup, 40+ subcommands, GrowthBook init, keychain prefetch, MDM read, OAuth, all flag handling | `main.go` has cobra CLI, 6 flags, basic subcommands (print/interactive), no subcommand routing | **PARTIAL** |
| `commands.ts` | `pkg/commands/` | `getCommands` — aggregates 50+ slash commands (clear, compact, config, cost, doctor, init, login, logout, mcp, memory, review, etc.) | `command.go` has `Registry` + `Command` struct — no actual commands registered | **PARTIAL** |
| `tools.ts` | `cmd/claude-go/` | `getTools`, `assembleToolPool`, `filterToolsByDenyRules` — dynamic tool pool assembly with feature flags | `main.go` hardcodes 6 tools (bash, fileread, filewrite, fileedit, glob, grep) — no dynamic assembly | **PARTIAL** |
| `cost-tracker.ts` | `pkg/services/cost/` | `formatTotalCost`, `saveCurrentSessionCosts`, `addToTotalCostState`, `resetCostState` — session cost tracking + formatting; reads from bootstrap state | `tracker.go` has `Tracker` with token add/cost calc — no session persistence, no state integration | **PARTIAL** |
| `history.ts` | `pkg/services/history/` | `addToHistory`, `getHistory`, `getSessionHistory` — JSONL append-only history with paste store, reverse-read, locking | `history.go` has `AddEntry`/`ListEntries` as JSON per-file — different format, no paste store, no locking | **PARTIAL** |
| `setup.ts` | `pkg/auth/` | `setup` — full session setup: cwd, project root, git root, hooks watchers, memory init, settings validation, env vars, deep link, MCP server approvals | `auth.go` has `CheckTrustDialog`/`CheckOnboarding` only — tiny fraction | **PARTIAL** |
| `ink.ts` | `pkg/tui/` | `render`, `createRoot`, re-exports: `Box`, `Text`, `ThemeProvider`, `color`, `BoxProps`, `TextProps` — Ink TUI framework wrappers with ThemeProvider | `pkg/tui/app.go` uses Bubble Tea — different TUI framework; no Ink compatibility | **PARTIAL** |
| `replLauncher.tsx` | `pkg/tui/` | `launchRepl` — dynamically imports `App` + `REPL` components and renders them with `renderAndRun` | No Go equivalent (Go TUI is `tui.NewModel`/`tea.NewProgram`) | **MISSING** |
| `costHook.ts` | `pkg/hooks/` | `useCostSummary` — React hook: on process exit, print cost summary + save session costs | No Go equivalent hook; cost printing done inline in main | **MISSING** |
| `dialogLaunchers.tsx` | `pkg/tui/` | `launchSnapshotUpdateDialog`, `launchResumeConversationDialog`, `launchOnboardingDialog`, `launchTrustDialog` etc. — dialog launcher wrappers | No Go equivalent | **MISSING** |
| `interactiveHelpers.tsx` | `pkg/tui/` | `showDialog`, `renderAndRun`, `showSetupDialog`, `completeOnboarding`, `showTrustDialog`, `showOnboardingDialog`, `handleMcpApprovals` — core TUI session management | No Go equivalent | **MISSING** |
| `projectOnboardingState.ts` | `pkg/auth/` | `getSteps`, `Step` — onboarding step state (workspace empty check, CLAUDE.md exists check) | No Go equivalent | **MISSING** |

---

## Summary Table

| Module | TS Files | Go Files | Overall Rating |
|---|---|---|---|
| `bridge/` | 31 | 1 (placeholder) | **MISSING** |
| `cli/` (top-level) | 6 | ~2 (minimal) | **MISSING** |
| `cli/handlers/` | 6 | 0 | **MISSING** |
| `cli/transports/` | 7 | 0 | **MISSING** |
| `skills/` | 3+17 bundled | 1 | **PARTIAL** |
| `memdir/` | 8 | 1 | **PARTIAL** |
| `migrations/` | 11 specific + framework | 1 (framework only) | **PARTIAL** |
| `entrypoints/` | 8 | 5 | **PARTIAL** |
| `bootstrap/state.ts` | 1 | 0 | **MISSING** |
| `buddy/` | 6 | 0 | **MISSING** |
| `plugins/` | 2 | 1 | **PARTIAL** |
| `remote/` | 4 | 0 | **MISSING** |
| `server/` | 3 | 0 | **MISSING** |
| `assistant/` | 1 | 0 | **MISSING** |
| `coordinator/` | 1 | 0 | **MISSING** |
| `voice/` | 1 | 0 | **MISSING** |
| `upstreamproxy/` | 2 | 0 | **MISSING** |
| `native-ts/` | 4 | 1 (OS-only) | **MISSING** |
| `moreright/` | 1 | — | N/A (stub) |
| `outputStyles/` | 1 | 0 | **MISSING** |
| Top-level files | 12 | ~5 | **PARTIAL** |
| **TOTAL** | **~141** | **~18 meaningful** | — |

---

## Critical Gaps (Blocking for Core Feature Parity)

1. **`bootstrap/state.ts`** — Global session state is the dependency hub for nearly everything. No equivalent exists in Go. This needs to be implemented first as all other modules depend on it.

2. **`bridge/` (31 files)** — Remote Control (CCR) is entirely absent. The Go `bridge.go` is a placeholder that cannot replace the protocol stack.

3. **`cli/structuredIO.ts` + `cli/remoteIO.ts`** — SDK stdio protocol (control messages, elicitation, permission prompts, hook callbacks) is absent. Required for non-interactive and SDK modes.

4. **`cli/transports/` (7 files)** — All transport implementations (WebSocket, Hybrid, SSE, CCR) are missing. The current Go code has no WebSocket transport.

5. **`entrypoints/init.ts`** — Application initialization (telemetry, policy limits, remote settings, CA certs, graceful shutdown) is absent.

6. **`memdir/findRelevantMemories.ts`** — LLM-powered memory recall is absent. The Go memdir can only list files.

7. **`migrations/` (11 files)** — All actual config migrations absent; only the versioned runner framework exists.

8. **`outputStyles/loadOutputStylesDir.ts`** — Output styles not ported.

9. **`commands.ts`** — No slash commands are registered in the Go command registry.

---

## Notes on Go Implementation Quality

- The Go code that exists is well-structured and idiomatic Go (proper error handling, `sync.Mutex`, interfaces).
- Framework layers (command registry, hook registry, plugin manager, migration runner, cost tracker, auth flow) are sound starting points.
- The gap is almost entirely in **protocol implementations** (bridge/transport layer), **LLM-side features** (memory relevance, coordinator mode, voice), and **full CLI completeness** (50+ slash commands, all subcommand handlers, structured IO).
- The TUI (Bubble Tea) is a reasonable Go-native replacement for the Ink/React TUI — but the session management, dialog launchers, and screen routing logic from `interactiveHelpers.tsx` and `dialogLaunchers.tsx` have not been ported.
