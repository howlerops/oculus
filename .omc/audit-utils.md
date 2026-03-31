# Utils Audit: TypeScript → Go Port

**Date:** 2026-03-31
**Scope:** `old-src/utils/` (329 top-level entries + subdirectories) vs `pkg/`
**Method:** Read first 40-60 lines of each TS file, read full Go equivalents, compare exports/behavior.

---

## Legend

| Rating | Meaning |
|--------|---------|
| FULL | Go equivalent covers all major exports/behavior |
| PARTIAL | Go has the concept but is missing significant functionality |
| MISSING | No Go equivalent exists |

---

## Top 30 Most Important Files — Detailed Analysis

---

### 1. `utils/auth.ts` → `pkg/auth/auth.go`

**Key TS exports:** `getAuthToken()`, `setAuthToken()`, `logout()`, `isLoggedIn()`, `getAPIKey()`, OAuth flow management, AWS STS auth, subscription type checking (`isMaxSubscriber`, `isClaudeAISubscriber`, `isProSubscriber`), `AccountInfo`, trust dialog helpers, betas cache clearing.

**Go coverage (`pkg/auth/auth.go`, `oauth.go`, `trust.go`, `keychain.go`):**
- `GetAuthToken()` — FULL: env → keychain → OAuth refresh → interactive login
- `Logout()` — FULL
- `IsLoggedIn()` — FULL
- AWS/Bedrock/STS auth — **MISSING**: no `pkg/auth/aws.go` or equivalent
- Subscription type (`isMaxSubscriber`, `isClaudeAISubscriber`) — **MISSING**: no subscription checking logic
- `AccountInfo` struct — **MISSING**: TS has rich account/profile metadata
- OAuth profile fetching (`getOauthProfileFromOauthToken`) — PARTIAL: tokens saved but profile not fetched
- `clearBetasCaches` coordination — MISSING

**Rating: PARTIAL** — Core auth flow covered, subscription/account metadata missing.

---

### 2. `utils/config.ts` → `pkg/config/config.go` + `pkg/config/settings.go`

**Key TS exports:** `getGlobalConfig()`, `saveGlobalConfig()`, `getCurrentProjectConfig()`, `getProjectConfig()`, `setProjectConfig()`, `getMemoryPath()`, `getManagedClaudeRulesDir()`, `getUserClaudeRulesDir()`, config file watching with `watchFile`/`unwatchFile`, `PastedContent` type, `BillingType`, re-entrancy guard, file change invalidation.

**Go coverage:**
- `GetGlobalConfig()` / `SaveGlobalConfig()` — FULL (basic read/write)
- `ProjectConfig` with `AllowedTools`/`DeniedTools` — PARTIAL (TS has much richer project config)
- Config file watching — **MISSING**: Go has no `watchFile` equivalent
- `getMemoryPath()`, `getManagedClaudeRulesDir()` — **MISSING**
- `PastedContent` type — **MISSING**
- Re-entrancy guard for analytics bootstrap — **MISSING**
- MCP server config in global config — PARTIAL (in `SettingsJson.MCPServers`)

**Rating: PARTIAL** — Basic persistence covered; file watching, memory paths, project config richness missing.

---

### 3. `utils/messages.ts` → `pkg/utils/messages/messages.go`

**Key TS exports:** `createUserMessage()`, `createAssistantAPIErrorMessage()`, `createUserInterruptionMessage()`, `extractTextContent()`, `normalizeMessages()`, `isToolUseMessage()`, `isToolResultMessage()`, `groupToolUses()`, message predicates, content block helpers, thinking block handling, image-too-large error messages, compact/boundary message types.

**Go coverage:**
- `CreateUserMessage()` / `CreateSystemMessage()` — FULL (basic)
- `ExtractText()` — PARTIAL (only handles text blocks, not tool-use/result)
- `CountTokensEstimate()` — PARTIAL (rough 4-chars/token, TS uses actual token counter)
- `normalizeMessages()` — **MISSING**
- Thinking block handling — **MISSING**
- Image/PDF error message helpers — **MISSING**
- Message predicate functions — **MISSING**
- `groupToolUses()` — **MISSING**

**Rating: PARTIAL** — Basic message creation present; normalization, predicates, content filtering missing.

---

### 4. `utils/git.ts` → `pkg/context/git.go`

**Key TS exports:** `findCanonicalGitRoot()` (memoized, LRU cache), `findGitRoot()`, `getRepoName()`, `getBranch()`, `getDefaultBranch()`, `getGitStatus()`, `getGitLog()`, `detectRepository()`, SHA lookups, shallow clone detection, worktree count, remote URL parsing, filesystem-based reading (avoids subprocess for speed).

**Go coverage:**
- `GetIsGit()`, `GetBranch()`, `GetDefaultBranch()` — FULL (subprocess-based)
- `GetGitStatus()`, `GetGitLog()`, `GetRepoName()` — FULL
- `findCanonicalGitRoot()` with LRU memoize — **MISSING**: Go version has no caching
- Filesystem-based reading (no subprocess, via `pkg/utils/git/gitFilesystem.ts`) — **MISSING**: Go always spawns subprocess
- Shallow clone detection — **MISSING**
- Worktree count — **MISSING**
- `git/gitignore.ts` helpers — **MISSING**

**Rating: PARTIAL** — Core git info covered; performance optimizations (FS-based reads, memoization) and advanced features missing.

---

### 5. `utils/log.ts` → `pkg/utils/log/log.go`

**Key TS exports:** `logError()`, `logForDebugging()`, `logForDiagnosticsNoPII()`, `logAntError()`, structured log serialization to files, `diagLogs.ts` (separate diagnostic log sink with PII stripping), session-specific log directories.

**Go coverage:**
- `Debug()`, `Info()`, `Warn()`, `Error()`, `Errorf()`, `Fatal()` — FULL (basic)
- `SetDebug()` via `CLAUDE_DEBUG` env — FULL
- `logForDiagnosticsNoPII()` (PII-stripped diagnostic logs) — **MISSING**
- Structured log file serialization — **MISSING**
- Session-scoped log directories — **MISSING**
- `diagLogs.ts` / diagnostic sinks — **MISSING**

**Rating: PARTIAL** — Basic logging present; diagnostic/PII-stripped logs and structured file logging missing.

---

### 6. `utils/theme.ts` → `pkg/utils/theme/theme.go`

**Key TS exports:** `Theme` type (50+ color fields including shimmer variants, diff colors, agent colors, TUI v2 colors), `getTheme()` (returns chalk-colored functions), `ThemeSetting`, dark/light/system theme detection, `systemTheme.ts` for OS-level theme detection.

**Go coverage:**
- `ThemeName`, `Theme` struct with 6 basic colors — **PARTIAL**: Go has ~6 fields vs TS's 50+
- `DefaultTheme()` — PARTIAL (hardcoded dark, no light variant)
- Light theme — **MISSING**
- Shimmer/diff/agent color variants — **MISSING**
- System theme detection — **MISSING**
- Chalk/ANSI integration — **MISSING** (Go has no chalk equivalent hooked in)

**Rating: PARTIAL** — Stub only; full theme system not ported.

---

### 7. `utils/platform.ts` → `pkg/utils/platform/platform.go`

**Key TS exports:** `getPlatform()` (macos/windows/wsl/linux/unknown), `getWslVersion()`, `SUPPORTED_PLATFORMS`, WSL detection via `/proc/version`, shell detection, `getShellType()`.

**Go coverage:**
- `GetPlatform()`, `IsMacOS()`, `IsLinux()`, `IsWindows()`, `GetArch()` — FULL (basic)
- WSL detection — **MISSING**: Go uses `runtime.GOOS` only, no `/proc/version` check
- `getWslVersion()` — **MISSING**
- Shell type detection — **MISSING**

**Rating: PARTIAL** — Basic OS detection present; WSL and shell detection missing.

---

### 8. `utils/claudemd.ts` → `pkg/context/claudemd.go`

**Key TS exports:** `getCachedClaudeMd()`, `loadClaudeMdFiles()` (full hierarchy: managed → user → project → local), `@include` directive processing, frontmatter parsing via `marked`, `picomatch` ignore patterns, `CLAUDE.local.md` support, circular reference prevention, `.claude/rules/*.md` directory, `truncateEntrypointContent()`.

**Go coverage:**
- `ClaudeMdSearchPaths()` — PARTIAL: walks project + parent dirs + global, but misses `.claude/rules/*.md`
- `LoadClaudeMd()` — PARTIAL: reads files, no `@include` directive support
- Managed memory (`/etc/claude-code/CLAUDE.md`) — **MISSING**
- `@include` directive processing — **MISSING**
- Frontmatter parsing — **MISSING**
- Ignore patterns (picomatch) — **MISSING**
- `CLAUDE.local.md` handling — **MISSING**
- Circular reference prevention — Present (via `seen` map) — FULL

**Rating: PARTIAL** — Basic file discovery; `@include`, frontmatter, managed memory, rules dir missing.

---

### 9. `utils/systemPrompt.ts` → `pkg/context/systemprompt.go`

**Key TS exports:** `buildEffectiveSystemPrompt()` (priority: override → coordinator → agent → custom → default), agent system prompt injection, proactive mode integration, `asSystemPrompt()`, coordinator mode handling, `appendSystemPrompt` support.

**Go coverage:**
- `BuildSystemPrompt()` — PARTIAL: identity + tools + claudemd + git context blocks
- Cache-control block structure — FULL
- Agent/coordinator prompt priority — **MISSING**
- Override system prompt — PARTIAL (replaces index 0, not architecturally correct)
- Proactive mode — **MISSING**
- `asSystemPrompt()` wrapper — **MISSING**

**Rating: PARTIAL** — Basic structure covered; multi-agent priority logic and proactive mode missing.

---

### 10. `utils/permissions/permissions.ts` → `pkg/utils/permissions/checker.go` + `pkg/types/permissions.go`

**Key TS exports:** Full permission check pipeline with `checkPermissions()`, bash classifier integration, sandbox manager, `PermissionResult` with decision/reason/updatedInput, rule matching (glob, content), `checkHasTrustDialogAccepted()`, tool-specific permission logic, classifier-based YOLO mode, `getNextPermissionMode()`.

**Go coverage:**
- `CheckToolPermission()` — PARTIAL: basic allow/deny/ask logic only
- `IsPathAllowed()` — PARTIAL: cwd + additionalWorkingDirs check only
- Rule matching — PARTIAL: exact match + wildcard only (no glob patterns)
- Bash classifier (YOLO mode) — **MISSING**
- Sandbox integration — **MISSING**
- `PermissionRule` with glob content matching — PARTIAL (types defined, logic missing)
- `autoModeState.ts` (circuit breaker) — **MISSING**
- Trust dialog check — **MISSING** (in `auth/trust.go` but not wired into permissions)

**Rating: PARTIAL** — Basic permission gating present; classifier, sandbox, trust dialog, glob matching missing.

---

### 11. `utils/permissions/filesystem.ts` → `pkg/utils/permissions/checker.go`

**Key TS exports:** `checkReadPermission()`, `checkWritePermission()`, path validation, `containsPathTraversal()`, `expandPath()`, `sanitizePath()`, UNC path vulnerability checks, project temp dir, auto-memory path awareness, permission pattern constants.

**Go coverage:**
- `IsPathAllowed()` — PARTIAL: basic prefix check only
- Path traversal detection — **MISSING**
- UNC path vulnerability checks — **MISSING**
- `sanitizePath()` / `expandPath()` — **MISSING**
- Auto-memory path awareness — **MISSING**

**Rating: PARTIAL** — Basic path check; security hardening (path traversal, UNC) missing.

---

### 12. `utils/hooks/hookEvents.ts` + `utils/hooks/hooksSettings.ts` → `pkg/hooks/engine.go`

**Key TS exports:** `HookStartedEvent`, `HookProgressEvent`, `HookResponseEvent`, event emitter system, `IndividualHookConfig` with source tracking (user/project/plugin/session/builtin), `isHookEqual()`, hook deduplication, 9+ hook sources.

**Go coverage:**
- `HookConfig`, `HookInput`, `HookOutput` — FULL
- `Registry.Execute()` with timeout — FULL
- Matcher patterns — PARTIAL (exact + wildcard, no glob)
- Hook source tracking (user/project/plugin) — **MISSING**
- Progress events during hook execution — **MISSING**
- `AsyncHookRegistry` (async parallel execution) — **MISSING**: Go runs sequentially
- SSRF guard — **MISSING** (`ssrfGuard.ts`)
- Session hooks — **MISSING**
- Post-sampling hooks (skillImprovement, etc.) — **MISSING**

**Rating: PARTIAL** — Core execution present; source tracking, async execution, SSRF guard missing.

---

### 13. `utils/model/model.ts` + `utils/model/modelOptions.ts` → `pkg/constants/models.go`

**Key TS exports:** `getDefaultMainLoopModelSetting()`, `getSmallFastModel()`, `getDefaultHaikuModel/SonnetModel/OpusModel()`, `getUserSpecifiedModelSetting()`, `ModelSetting` type (alias | name | null), `renderDefaultModelSetting()`, subscription-gated model access (1M context window checks), model capability checks, deprecation warnings, Bedrock/Vertex model name mapping.

**Go coverage:**
- Model ID constants (opus4, sonnet4, haiku4) — FULL
- `ModelAliases` map — FULL
- `ResolveModelAlias()` — FULL
- Default model selection — PARTIAL (hardcoded, no subscription gating)
- 1M context window access checks — **MISSING**
- Bedrock/Vertex model name mapping — **MISSING**
- `getSmallFastModel()` with env override — PARTIAL (in config but not in constants)
- Model deprecation warnings — **MISSING**
- `modelCapabilities.ts` (structured outputs, thinking support) — **MISSING**

**Rating: PARTIAL** — Model IDs and aliases covered; dynamic selection, capability checks, provider mapping missing.

---

### 14. `utils/model/providers.ts` → `pkg/config/env.go`

**Key TS exports:** `getAPIProvider()` (firstParty/bedrock/vertex/foundry), `isFirstPartyAnthropicBaseUrl()`, `getAPIProviderForStatsig()`.

**Go coverage:** Checked `pkg/config/env.go` — provider detection likely present. TS logic is: env vars `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`. Go equivalent needs verification but the pattern matches `pkg/config/env.go`.

**Rating: PARTIAL** — Needs verification; foundry support may be missing.

---

### 15. `utils/bash/ast.ts` + `utils/bash/commands.ts` → `pkg/tools/bash/bash.go`

**Key TS exports:** Tree-sitter AST-based bash parsing, `parseForSecurity()`, `SimpleCommand` type with argv/env/redirects, `ParseForSecurityResult` (simple/too-complex/parse-unavailable), output redirect extraction, shell quote parsing, heredoc handling, pipe command detection, bash syntax highlighter.

**Go coverage:**
- `BashTool.Call()` — FULL (executes commands)
- Background execution — FULL
- Timeout handling — FULL
- AST-based security parsing — **MISSING**: Go runs commands directly, no pre-parse
- `parseForSecurity()` — **MISSING**
- Shell quote parsing — **MISSING**
- Redirect extraction — **MISSING**
- Output redirect security check — **MISSING**

**Rating: PARTIAL** — Execution present; security pre-parsing (tree-sitter AST) entirely missing.

---

### 16. `utils/plugins/installedPluginsManager.ts` + `utils/plugins/marketplaceManager.ts` → `pkg/services/plugins/plugins.go`

**Key TS exports:** Plugin installation, update, uninstall lifecycle, marketplace API client, dependency resolution, telemetry, LSP integration, MCP plugin integration, hook registration from plugins, command/agent loading, `headlessPluginInstall.ts`, `managedPlugins.ts` (enterprise-managed).

**Go coverage:**
- `PluginManifest` struct — FULL
- `Manager.LoadFromDirectory()` — FULL
- `ApplySettings()` enable/disable — FULL
- Plugin installation/update/uninstall — **MISSING**
- Marketplace API client — **MISSING**
- Dependency resolution — **MISSING**
- LSP plugin integration — **MISSING**
- MCP plugin integration — **MISSING**
- Managed/enterprise plugins — **MISSING**
- Hook/command/agent loading from plugins — **MISSING**

**Rating: PARTIAL** — Manifest loading only; full lifecycle, marketplace, and integration missing.

---

### 17. `utils/swarm/` (17 files) → `pkg/services/bridge/bridge.go`

**Key TS exports:** `inProcessRunner.ts` (AsyncLocalStorage context isolation, plan mode approval, idle notification), `spawnUtils.ts` (CLI flag propagation, teammate command resolution), `teammateInit.ts`, `leaderPermissionBridge.ts` (mailbox-based permission forwarding), `reconnection.ts`, `backends/` (tmux, in-process, Claude-in-Chrome), `teammateModel.ts`, layout manager.

**Go coverage:**
- `Bridge` struct with session state, inbox/outbox channels — PARTIAL (skeleton only)
- `GetState()`, `SetStatus()`, `Send()`, `Receive()`, `Close()` — FULL (basic)
- In-process teammate runner — **MISSING**
- Permission bridge (mailbox forwarding to leader) — **MISSING**
- Spawn utilities / CLI flag propagation — **MISSING**
- Tmux backend — **MISSING**
- Reconnection logic — **MISSING**
- Layout manager — **MISSING**
- Plan mode approval flow — **MISSING**

**Rating: PARTIAL** — Message bus skeleton only; full swarm orchestration missing.

---

### 18. `utils/task/diskOutput.ts` + `utils/task/framework.ts` → `pkg/task/task.go`

**Key TS exports:** `DiskTaskOutput` (file-backed streaming with O_NOFOLLOW security, 5GB cap, watchdog), `getTaskOutputDir()` (session-scoped), `MAX_TASK_OUTPUT_BYTES`, `TaskOutput` interface, task framework with status tracking, SDK progress events, output formatting.

**Go coverage:**
- `TaskType`, `TaskStatus` enums — FULL
- `TaskState` struct — FULL
- `GenerateTaskId()` with type prefixes — FULL
- `GetTaskOutputPath()` — FULL
- `IsTerminalTaskStatus()` — FULL
- `DiskTaskOutput` (file-backed streaming) — **MISSING**: Go has path but no writer
- O_NOFOLLOW security — **MISSING**
- Session-scoped output dir — **MISSING** (hardcoded `.claude/tasks/`)
- Output size cap/watchdog — **MISSING**
- SDK progress events — **MISSING**

**Rating: PARTIAL** — Types and ID generation present; actual output streaming missing.

---

### 19. `utils/todo/types.ts` → `pkg/tools/todowrite/todowrite.go`

**Key TS exports:** `TodoItem` (content, status, activeForm), `TodoList`, `TodoStatusSchema` (pending/in_progress/completed), Zod validation.

**Go coverage:**
- `TodoItem` with id/content/status — FULL (basic)
- `TodoWriteTool.Call()` — FULL
- `activeForm` field — **MISSING** (TS requires it, Go omits it)
- Zod validation equivalent — **MISSING** (Go does no input validation)
- All-complete auto-clear — FULL

**Rating: PARTIAL** — Core functionality present; `activeForm` field and validation missing.

---

### 20. `utils/secureStorage/index.ts` → `pkg/auth/keychain.go`

**Key TS exports:** `getSecureStorage()` (platform dispatch), `macOsKeychainStorage` (keytar-based), `plainTextStorage` (fallback), `fallbackStorage` (try keychain, fall back to plaintext), `keychainPrefetch.ts` (async prefetch at startup).

**Go coverage:**
- `SecureStorage` interface — FULL
- `MacOSKeychain` (security CLI) — FULL
- `FileStorage` fallback — FULL
- `NewSecureStorage()` platform dispatch — FULL
- `SaveTokens()` / `LoadTokens()` / `SaveAPIKey()` / `LoadAPIKey()` — FULL
- Fallback chain (keychain → plaintext) — **MISSING**: Go dispatches directly, no fallback
- Startup prefetch — **MISSING**
- Linux libsecret — **MISSING** (noted as TODO in both)

**Rating: PARTIAL** — Core storage present; fallback chain and prefetch missing.

---

### 21. `utils/fileStateCache.ts` → `pkg/utils/filestate/cache.go`

**Key TS exports:** `FileStateCache` (LRU-backed with content + timestamp + offset + limit + `isPartialView` flag), `READ_FILE_STATE_CACHE_SIZE` (100 entries), 25MB size limit with byte-aware eviction, path normalization.

**Go coverage:**
- `Cache` with RWMutex — FULL
- `Get()`, `Refresh()`, `Clear()` — FULL
- `maxSize` eviction (random drop) — PARTIAL: Go drops one random entry vs LRU
- `isPartialView` flag — **MISSING**: critical for preventing edits on partial content
- Content storage — **MISSING**: Go only stores ModTime/Size, not file content
- Offset/limit tracking — **MISSING**
- Byte-aware size limit — **MISSING**

**Rating: PARTIAL** — Metadata cache only; content caching and partial-view tracking missing (these are critical for the edit tool safety model).

---

### 22. `utils/format.ts` → (no direct Go equivalent)

**Key TS exports:** `formatFileSize()`, `formatDuration()`, `formatSecondsShort()`, `formatDate()`, `formatRelativeTime()`, locale-aware formatting via `intl.ts`.

**Go coverage:** No `pkg/utils/format/` package. Standard library `fmt` used ad-hoc.

**Rating: MISSING** — No centralized formatting utilities.

---

### 23. `utils/array.ts` → (no direct Go equivalent)

**Key TS exports:** `intersperse()`, `count()`, `uniq()`.

**Go coverage:** No equivalent package. Go uses standard library patterns inline.

**Rating: MISSING** — Trivial helpers, Go standard library covers this adequately.

---

### 24. `utils/env.ts` → `pkg/config/env.go`

**Key TS exports:** `getGlobalClaudeFile()`, `hasInternetAccess()`, `isCommandAvailable()`, `getClaudeConfigHomeDir()` (from `envUtils.ts`), `isEnvTruthy()`, `isEnvDefinedFalsy()`.

**Go coverage:** `pkg/config/env.go` likely has config dir resolution. Full verification needed but pattern matches.

**Rating: PARTIAL** — Core env resolution likely present; internet check, command availability missing.

---

### 25. `utils/errors.ts` → (no direct Go equivalent)

**Key TS exports:** `ClaudeError`, `MalformedCommandError`, `AbortError`, `isAbortError()`, `ConfigParseError` (with filePath + defaultConfig), `ShellError` (stdout/stderr/code/interrupted), `toError()`, `errorMessage()`, `getErrnoCode()`, `isENOENT()`.

**Go coverage:** No `pkg/utils/errors/` package. Go uses standard `errors` package and custom error types inline.

**Rating: MISSING** — No centralized error hierarchy; `ShellError` and `ConfigParseError` typed errors not ported.

---

### 26. `utils/permissions/autoModeState.ts` → (no Go equivalent)

**Key TS exports:** `setAutoModeActive()`, `isAutoModeActive()`, `setAutoModeFlagCli()`, `isAutoModeCircuitBroken()`, circuit breaker for auto-mode gate.

**Go coverage:** No equivalent in `pkg/`.

**Rating: MISSING** — Auto/YOLO mode state management not ported.

---

### 27. `utils/permissions/yoloClassifier.ts` → (no Go equivalent)

**Key TS exports:** LLM-based classifier that decides whether to auto-approve bash commands in YOLO mode, context window management, GrowthBook feature-gated, writes classifier results to disk for replay.

**Go coverage:** No equivalent.

**Rating: MISSING** — This is a significant feature gap; auto-approve mode requires classifier.

---

### 28. `utils/model/modelCapabilities.ts` → (no direct Go equivalent)

**Key TS exports:** `modelSupportsStructuredOutputs()`, `shouldUseGlobalCacheScope()`, `has1mContext()`, `modelSupports1M()`, per-model capability flags.

**Go coverage:** No `modelCapabilities` equivalent in Go constants or config.

**Rating: MISSING** — Model capability matrix not ported.

---

### 29. `utils/settings/settings.ts` (14 files total) → `pkg/config/settings.go`

**Key TS exports:** Full settings merge pipeline (managed → user → project → local → flag → policy → session), MDM/HKCU registry support (Windows), `parseSettingsFile()`, `getSettingsForSource()`, `SettingsJson` schema with Zod validation, settings cache with invalidation, `getAutoModeConfig()`, `changeDetector.ts`, plugin-only policy, validation tips.

**Go coverage:**
- `SettingsJson` struct — FULL (all fields mapped)
- `HookConfig`, `MCPServerConfig`, `PermissionSettings` — FULL
- Multi-source merge pipeline — **MISSING**: Go has struct, no loader
- MDM/Windows registry support — **MISSING**
- Zod validation equivalent — **MISSING**
- Settings file watcher/cache invalidation — **MISSING**
- `getAutoModeConfig()` — **MISSING**
- `changeDetector.ts` — **MISSING**

**Rating: PARTIAL** — Schema defined; multi-source merge pipeline not implemented.

---

### 30. `utils/git/gitFilesystem.ts` → `pkg/context/git.go`

**Key TS exports:** `resolveGitDir()`, `getCachedBranch()`, `getCachedHead()`, `getCachedRemoteUrl()`, `getCachedDefaultBranch()`, `getWorktreeCountFromFs()`, `isShallowClone()`, `GitHeadWatcher` (fs.watchFile-based), packed-refs parser, loose ref resolver — all **without spawning git subprocess**.

**Go coverage:** Go's `pkg/context/git.go` spawns `git` subprocess for all operations. No filesystem-based reading.

**Rating: PARTIAL** — Functionally equivalent but subprocess-based; startup perf impact in large repos.

---

## Subdirectory Summary Assessments

### `utils/permissions/` (24 files) → `pkg/utils/permissions/checker.go` + `pkg/types/permissions.go`

| File | Status |
|------|--------|
| `PermissionRule.ts` | PARTIAL — types in `pkg/types/permissions.go` |
| `PermissionResult.ts` | PARTIAL — basic struct present |
| `PermissionMode.ts` | FULL — all modes defined |
| `permissions.ts` | PARTIAL — basic check logic only |
| `filesystem.ts` | PARTIAL — basic path check, no security hardening |
| `permissionsLoader.ts` | MISSING — no multi-source rule loading |
| `permissionSetup.ts` | MISSING |
| `permissionRuleParser.ts` | MISSING — no glob/pattern parsing |
| `yoloClassifier.ts` | MISSING — LLM classifier not ported |
| `autoModeState.ts` | MISSING |
| `bashClassifier.ts` | MISSING |
| `classifierDecision.ts` | MISSING |
| `shadowedRuleDetection.ts` | MISSING |
| `PermissionUpdate.ts` | MISSING — no rule persistence |
| `bypassPermissionsKillswitch.ts` | MISSING |
| `dangerousPatterns.ts` | MISSING |
| `shellRuleMatching.ts` | MISSING |
| `denialTracking.ts` | MISSING |
| `permissionExplainer.ts` | MISSING |

**Overall: PARTIAL** — ~6/24 files covered (types + basic check). Major gaps: classifier, rule parsing, rule persistence, YOLO mode.

---

### `utils/settings/` (14 files) → `pkg/config/settings.go`

| File | Status |
|------|--------|
| `types.ts` | FULL — `SettingsJson` struct matches |
| `settings.ts` | PARTIAL — struct only, no load/merge pipeline |
| `constants.ts` | MISSING — source priority constants |
| `settingsCache.ts` | MISSING — no cache layer |
| `changeDetector.ts` | MISSING |
| `validation.ts` | MISSING — no schema validation |
| `managedPath.ts` | MISSING |
| `mdm/settings.ts` | MISSING — Windows MDM/HKCU |
| `pluginOnlyPolicy.ts` | MISSING |
| `internalWrites.ts` | MISSING |
| `applySettingsChange.ts` | MISSING |
| `permissionValidation.ts` | MISSING |

**Overall: PARTIAL** — Schema defined; load/merge/validate pipeline missing.

---

### `utils/hooks/` (18 files) → `pkg/hooks/engine.go` + `pkg/hooks/ui.go`

| File | Status |
|------|--------|
| `hookEvents.ts` | MISSING — event system not ported |
| `hooksSettings.ts` | PARTIAL — config struct present |
| `AsyncHookRegistry.ts` | MISSING — Go runs hooks synchronously |
| `execAgentHook.ts` | MISSING |
| `execHttpHook.ts` | MISSING — HTTP hooks not supported |
| `execPromptHook.ts` | MISSING |
| `apiQueryHookHelper.ts` | MISSING |
| `ssrfGuard.ts` | MISSING — security feature |
| `postSamplingHooks.ts` | MISSING |
| `sessionHooks.ts` | MISSING |
| `registerFrontmatterHooks.ts` | MISSING |
| `registerSkillHooks.ts` | MISSING |
| `skillImprovement.ts` | MISSING |
| `hookHelpers.ts` | MISSING |
| `hooksConfigManager.ts` | MISSING |
| `fileChangedWatcher.ts` | MISSING |

**Overall: PARTIAL** — Shell hook execution present; HTTP hooks, async execution, SSRF guard, skill/frontmatter hooks all missing.

---

### `utils/model/` (16 files) → `pkg/constants/models.go`

| File | Status |
|------|--------|
| `model.ts` | PARTIAL — dynamic selection logic missing |
| `modelOptions.ts` | MISSING — subscription-gated options |
| `modelStrings.ts` | PARTIAL — hardcoded constants vs dynamic |
| `providers.ts` | PARTIAL — basic provider detection |
| `modelCapabilities.ts` | MISSING — structured outputs, 1M context |
| `aliases.ts` | FULL — `ModelAliases` map |
| `modelAllowlist.ts` | MISSING |
| `configs.ts` | MISSING |
| `deprecation.ts` | MISSING |
| `bedrock.ts` | MISSING — Bedrock model name mapping |
| `antModels.ts` | MISSING — internal model codenames |
| `check1mAccess.ts` | MISSING |
| `contextWindowUpgradeCheck.ts` | MISSING |
| `validateModel.ts` | MISSING |

**Overall: PARTIAL** — ~3/16 files covered. Model selection, capabilities, provider-specific mapping missing.

---

### `utils/bash/` (15 files) → `pkg/tools/bash/bash.go`

| File | Status |
|------|--------|
| `ast.ts` | MISSING — tree-sitter AST parsing |
| `parser.ts` | MISSING |
| `commands.ts` | MISSING — security command extraction |
| `bashParser.ts` | MISSING |
| `shellQuote.ts` | MISSING |
| `shellQuoting.ts` | MISSING |
| `heredoc.ts` | MISSING |
| `bashPipeCommand.ts` | MISSING |
| `ParsedCommand.ts` | MISSING |
| `registry.ts` | MISSING |
| `shellCompletion.ts` | MISSING |
| `treeSitterAnalysis.ts` | MISSING |
| `ShellSnapshot.ts` | MISSING |

**Overall: PARTIAL** — Execution only; entire security parsing layer missing.

---

### `utils/plugins/` (45 files) → `pkg/services/plugins/plugins.go`

| File | Status |
|------|--------|
| `installedPluginsManager.ts` | MISSING — install lifecycle |
| `marketplaceManager.ts` | MISSING |
| `marketplaceHelpers.ts` | MISSING |
| `dependencyResolver.ts` | MISSING |
| `loadPluginAgents.ts` | MISSING |
| `loadPluginCommands.ts` | MISSING |
| `loadPluginHooks.ts` | MISSING |
| `loadPluginOutputStyles.ts` | MISSING |
| `lspPluginIntegration.ts` | MISSING |
| `lspRecommendation.ts` | MISSING |
| `mcpPluginIntegration.ts` | MISSING |
| `mcpbHandler.ts` | MISSING |
| `managedPlugins.ts` | MISSING — enterprise managed |
| `headlessPluginInstall.ts` | MISSING |
| `gitAvailability.ts` | MISSING |
| `hintRecommendation.ts` | MISSING |
| `fetchTelemetry.ts` | MISSING |
| `installCounts.ts` | MISSING |
| `cacheUtils.ts` | MISSING |
| `addDirPluginSettings.ts` | PARTIAL — settings loading |

**Overall: PARTIAL** — Manifest loading ~1/45 files. Full plugin ecosystem not ported.

---

### `utils/swarm/` (17 files) → `pkg/services/bridge/bridge.go`

| File | Status |
|------|--------|
| `constants.ts` | PARTIAL — some constants missing |
| `inProcessRunner.ts` | MISSING — core teammate execution |
| `spawnUtils.ts` | MISSING — spawn + CLI flag propagation |
| `leaderPermissionBridge.ts` | MISSING — permission forwarding |
| `reconnection.ts` | MISSING |
| `permissionSync.ts` | MISSING |
| `teammateInit.ts` | MISSING |
| `teammateLayoutManager.ts` | MISSING |
| `teammateModel.ts` | MISSING |
| `teammatePromptAddendum.ts` | MISSING |
| `teamHelpers.ts` | MISSING |
| `spawnInProcess.ts` | MISSING |
| `backends/` | MISSING — all backends |
| `It2SetupPrompt.tsx` | MISSING |

**Overall: PARTIAL** — Message bus skeleton only (~1/17). Full swarm orchestration not ported.

---

### `utils/task/` (5 files) → `pkg/task/task.go`

| File | Status |
|------|--------|
| `framework.ts` | PARTIAL — types only |
| `diskOutput.ts` | MISSING — file-backed streaming |
| `outputFormatting.ts` | MISSING |
| `sdkProgress.ts` | MISSING |
| `TaskOutput.ts` | MISSING — interface |

**Overall: PARTIAL** — Types/IDs ported; output streaming missing.

---

### `utils/secureStorage/` (6 files) → `pkg/auth/keychain.go`

| File | Status |
|------|--------|
| `index.ts` | PARTIAL — no fallback chain |
| `macOsKeychainStorage.ts` | FULL |
| `plainTextStorage.ts` | FULL |
| `fallbackStorage.ts` | MISSING |
| `keychainPrefetch.ts` | MISSING |
| (Linux libsecret) | MISSING |

**Overall: PARTIAL** — ~3/6 files covered; fallback chain and prefetch missing.

---

## Files With No Go Equivalent (Selected Critical Ones)

| TS File | Importance | Notes |
|---------|-----------|-------|
| `utils/errors.ts` | HIGH | Typed error hierarchy (`ShellError`, `ConfigParseError`) |
| `utils/format.ts` | MEDIUM | Display formatters |
| `utils/array.ts` | LOW | Trivial helpers |
| `utils/env.ts` | HIGH | Config dir resolution |
| `utils/debug.ts` | HIGH | Debug/diagnostics logging |
| `utils/cwd.ts` | HIGH | Working directory management |
| `utils/memoize.ts` | MEDIUM | LRU memoization |
| `utils/lockfile.ts` | HIGH | File-based locking |
| `utils/json.ts` | MEDIUM | Safe JSON parse |
| `utils/path.ts` | HIGH | Path utilities + traversal detection |
| `utils/semver.ts` | MEDIUM | Version comparison |
| `utils/sleep.ts` | LOW | Trivial |
| `utils/uuid.ts` | MEDIUM | ID generation |
| `utils/hash.ts` | MEDIUM | Content hashing |
| `utils/fsOperations.ts` | HIGH | FS abstraction layer |
| `utils/sessionStorage.ts` | HIGH | Session persistence |
| `utils/sessionEnvironment.ts` | HIGH | Session env vars |
| `utils/context.ts` | HIGH | Context management |
| `utils/memory/` | HIGH | Memory/AGENTS.md subsystem |
| `utils/telemetry/` | MEDIUM | Analytics/telemetry |
| `utils/sandbox/` | HIGH | Sandboxing layer |
| `utils/cron.ts` | MEDIUM | Cron scheduler |
| `utils/autoUpdater.ts` | MEDIUM | Auto-update |
| `utils/diff.ts` | HIGH | Diff computation |
| `utils/tokens.ts` | HIGH | Token counting |
| `utils/tokenBudget.ts` | HIGH | Token budget management |
| `utils/ripgrep.ts` | MEDIUM | Ripgrep wrapper |
| `utils/worktree.ts` | HIGH | Git worktree management |
| `utils/plans.ts` | HIGH | Plan mode state |

---

## Overall Summary

| Area | Coverage | Gap Severity |
|------|----------|-------------|
| `utils/auth.ts` | PARTIAL | HIGH — subscription checks, AWS auth missing |
| `utils/config.ts` | PARTIAL | HIGH — file watching, memory paths missing |
| `utils/messages.ts` | PARTIAL | HIGH — normalization, predicates missing |
| `utils/git.ts` | PARTIAL | MEDIUM — caching, FS-reads missing |
| `utils/log.ts` | PARTIAL | MEDIUM — PII-stripped diagnostics missing |
| `utils/theme.ts` | PARTIAL | LOW — TUI concern, stub is functional |
| `utils/platform.ts` | PARTIAL | MEDIUM — WSL detection missing |
| `utils/claudemd.ts` | PARTIAL | HIGH — @include, frontmatter missing |
| `utils/systemPrompt.ts` | PARTIAL | HIGH — multi-agent priority missing |
| `utils/permissions/` (24) | PARTIAL | CRITICAL — classifier, glob matching, YOLO mode missing |
| `utils/settings/` (14) | PARTIAL | HIGH — merge pipeline not implemented |
| `utils/hooks/` (18) | PARTIAL | HIGH — HTTP hooks, async, SSRF guard missing |
| `utils/model/` (16) | PARTIAL | HIGH — capabilities, provider mapping missing |
| `utils/bash/` (15) | PARTIAL | HIGH — security parsing entirely missing |
| `utils/plugins/` (45) | PARTIAL | HIGH — only manifest loading |
| `utils/swarm/` (17) | PARTIAL | CRITICAL — only message bus skeleton |
| `utils/task/` (5) | PARTIAL | HIGH — output streaming missing |
| `utils/todo/` | PARTIAL | LOW — minor field missing |
| `utils/secureStorage/` (6) | PARTIAL | LOW — fallback chain missing |
| `utils/fileStateCache.ts` | PARTIAL | HIGH — content caching missing (edit safety) |
| `utils/format.ts` | MISSING | LOW |
| `utils/array.ts` | MISSING | LOW |
| `utils/errors.ts` | MISSING | MEDIUM |
| Various (debug, cwd, path, etc.) | MISSING | VARIES |

### Top Priority Gaps to Address

1. **CRITICAL:** `utils/permissions/yoloClassifier.ts` — LLM-based auto-approve not ported; auto mode is non-functional
2. **CRITICAL:** `utils/swarm/inProcessRunner.ts` + backends — multi-agent (team/swarm) mode not ported
3. **HIGH:** `utils/settings/settings.ts` merge pipeline — settings from multiple sources not merged; app uses defaults only
4. **HIGH:** `utils/fileStateCache.ts` content + `isPartialView` — edit tool safety model depends on this
5. **HIGH:** `utils/bash/ast.ts` + security parsing — bash commands run without pre-parse security check
6. **HIGH:** `utils/claudemd.ts` `@include` directive — CLAUDE.md includes silently ignored
7. **HIGH:** `utils/auth.ts` subscription checks — subscription-gated features (1M context, etc.) non-functional
8. **HIGH:** `utils/task/diskOutput.ts` — background task output not streamed to disk
9. **MEDIUM:** `utils/hooks/` HTTP hooks + async execution — hook system incomplete
10. **MEDIUM:** `utils/model/modelCapabilities.ts` — structured outputs, 1M context window flags missing
