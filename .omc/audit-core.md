# TypeScript-to-Go Port Audit

Audited: 2026-03-31
Source: `old-src/` (TypeScript) → `pkg/` (Go)

---

## types/ (TS) → pkg/types/

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| `types/command.ts` | `LocalCommandResult`, `PromptCommand`, `LocalCommandCall`, `LocalCommandModule`, `LocalJSXCommandContext` | MISSING | No Go equivalent in pkg/types/ or pkg/tool/. CLI command dispatch types not ported. |
| `types/hooks.ts` | `isHookEvent`, `PromptRequest`, `PromptResponse`, `syncHookResponseSchema`, `hookJSONOutputSchema`, `isSyncHookJSONOutput`, `isAsyncHookJSONOutput`, `HookCallbackContext`, `HookCallback`, `HookCallbackMatcher`, `HookProgress`, `HookBlockingError`, `PermissionRequestResult`, `HookResult`, `AggregatedHookResult` | PARTIAL | Go `pkg/types/hooks.go` has `HookEvent`, `HookProgress`, `PromptRequest`, `PromptResponse`, `HookResult`, `HookBlockingError`, `AggregatedHookResult`. Missing: `HookCallback`, `HookCallbackMatcher`, `HookCallbackContext`, `isSyncHookJSONOutput`/`isAsyncHookJSONOutput` validator fns, Zod schemas (N/A in Go but no equivalent runtime validation). |
| `types/ids.ts` | `SessionId`, `AgentId`, `asSessionId`, `asAgentId`, `toAgentId` | FULL | Go `pkg/types/ids.go` has `SessionId`, `AgentId`, `AsSessionId`, `AsAgentId`, `ToAgentId`, plus bonus `NewAgentId`. Pattern match identical. |
| `types/logs.ts` | `SerializedMessage`, `LogOption`, `SummaryMessage`, `CustomTitleMessage`, `AiTitleMessage`, `LastPromptMessage`, `TaskSummaryMessage`, `TagMessage`, `AgentNameMessage`, `AgentColorMessage`, `AgentSettingMessage`, `PRLinkMessage`, `ModeEntry`, `PersistedWorktreeSession`, `WorktreeStateEntry`, `ContentReplacementEntry`, `FileHistorySnapshotMessage`, `AttributionSnapshotMessage`, `TranscriptMessage`, `SpeculationAcceptMessage`, `ContextCollapseCommitEntry`, `ContextCollapseSnapshotEntry`, `Entry`, `sortLogs` | MISSING | No dedicated `logs.go` in pkg/types/. Session persistence / transcript log types entirely absent. |
| `types/permissions.ts` | `ExternalPermissionMode`, `PermissionMode`, `PermissionBehavior`, `PermissionRuleSource`, `PermissionRuleValue`, `PermissionRule`, `PermissionUpdateDestination`, `PermissionUpdate`, `AdditionalWorkingDirectory`, `PermissionResult`, `PermissionDecision`, `PermissionAllowDecision`, `PermissionAskDecision`, `PermissionDenyDecision`, `ClassifierResult`, `RiskLevel`, `ToolPermissionRulesBySource`, `ToolPermissionContext` | PARTIAL | Go `pkg/types/permissions.go` has modes, behavior, rule source/value/rule, update destination/type, `PermissionResult`. Missing: `AdditionalWorkingDirectory`, `PermissionCommandMetadata`, `PermissionMetadata`, `PendingClassifierCheck`, `ClassifierResult`, `ClassifierBehavior`, `ClassifierUsage`, `YoloClassifierResult`, `RiskLevel`, `PermissionExplanation`, `ToolPermissionRulesBySource`. `ToolPermissionContext` exists in `Tool.ts` (Go: `pkg/types/permissions.go` partial). |
| `types/plugin.ts` | `BuiltinPluginDefinition`, `PluginRepository`, `PluginConfig`, `LoadedPlugin`, `PluginComponent`, `PluginError`, `PluginAuthor`, `PluginManifest`, `CommandMetadata` | MISSING | No plugin types in Go at all. |
| `types/textInputTypes.ts` | `InlineGhostText`, `BaseTextInputProps`, `OrphanedPermission`, `AssistantMessageForRender` (and many UI prop types) | MISSING | UI/TUI-specific types. No Go equivalent expected (TUI is separate), but `OrphanedPermission` is used in QueryEngine and has no Go counterpart. |
| `schemas/hooks.ts` (1 file in schemas/) | `BashCommandHookSchema`, `PromptHookSchema`, `AgentHookSchema`, `HookCommandSchema`, `HooksSettings` (Zod schemas) | PARTIAL | Runtime schema validation is Zod-specific. Go has struct types for hook config but no equivalent runtime schema validation or `HooksSettings` type (map of event → hook array). |

### Generated types (types/generated/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| `types/generated/events_mono` | Anthropic API event stream types | MISSING | File not readable but corresponds to streaming event types. Go `pkg/types/message.go` has `StreamEvent` but full generated API types not verified. |
| `types/generated/google` | Google/Vertex API types | MISSING | No Google/Vertex provider types in Go. |

---

## state/ (6 TS files) → pkg/state/

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| `state/AppState.tsx` | Re-exports from `AppStateStore.ts`: `AppState`, `AppStateStore`, `CompletionBoundary`, `getDefaultAppState`, `IDLE_SPECULATION_STATE`, `SpeculationResult`, `SpeculationState`; plus `AppStoreContext`, `AppStateProvider` (React) | PARTIAL | `AppStateProvider`/`AppStoreContext` are React-specific (N/A in Go). `CompletionBoundary`, `SpeculationResult`, `SpeculationState`, `IDLE_SPECULATION_STATE` not in Go. |
| `state/AppStateStore.ts` | `AppState` (large struct), `AppStateStore`, `CompletionBoundary`, `SpeculationState`, `SpeculationResult`, `IDLE_SPECULATION_STATE`, `FooterItem`, `getDefaultAppState` | PARTIAL | Go `pkg/state/appstate.go` has `AppState` and `NewAppState`. Missing from Go `AppState`: `MainLoopModelForSession`, `ShowTeammateMessagePreview`, `SelectedIPAgentIndex`, `CoordinatorTaskIndex`, `ViewSelectionMode`, `FooterSelection`, `SpinnerTip`, `Agent`, `KairosEnabled`, `RemoteSessionUrl`, `RemoteConnectionStatus`, `RemoteBackgroundTaskCount`, `ReplBridge*` fields (12 fields), `IsUltraplanMode`, `ViewingAgentTaskId`, `Notifications`, `FileHistory`, `Plugins`, `SessionHooks`, `AttributionState`, `DenialTracking`, `TodoList` and many others. Go AppState is a minimal subset. |
| `state/onChangeAppState.ts` | `onChangeAppState`, `externalMetadataToAppState` | MISSING | No equivalent change-listener with permission-mode sync logic. |
| `state/selectors.ts` | `getViewedTeammateTask`, `ActiveAgentForInput`, `getActiveAgentForInput` | PARTIAL | Go `pkg/state/selectors.go` has `GetRunningTasks`, `GetPendingTasks`, `GetActiveTasks`, `HasInProgressTools`, `GetMessageCount`, `GetLastMessage`. Missing: `GetViewedTeammateTask`, `GetActiveAgentForInput`, `ActiveAgentForInput` (teammate routing selectors). |
| `state/store.ts` | `Store<T>` (generic), `createStore` | PARTIAL | Go `pkg/state/store.go` has `Store` (AppState-specific, not generic), `NewStore`, `Get`, `Set`, `Update`, `Subscribe`. Functionally equivalent but not generic — tied to `AppState` only. TS version is generic `Store<T>`. |
| `state/teammateViewHelpers.ts` | `enterTeammateView`, `exitTeammateView` | MISSING | No teammate view transition functions in Go. |

---

## tasks/ (12 TS files) → pkg/task/

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| `tasks/types.ts` | `TaskState` (union), `BackgroundTaskState` (union), `isBackgroundTask` | PARTIAL | Go `pkg/task/task.go` has `TaskState` as a single struct (not a discriminated union). `BackgroundTaskState` union and `IsBackgroundTask` predicate not ported. All task variants (LocalShell, LocalAgent, Remote, InProcessTeammate, LocalWorkflow, MonitorMCP, Dream) are collapsed into one `TaskState` struct — variant-specific fields are absent. |
| `tasks/pillLabel.ts` | `getPillLabel` | MISSING | UI display helper. No Go equivalent. |
| `tasks/stopTask.ts` | `StopTaskError`, `StopTaskContext`, `StopTaskResult`, `stopTask` | MISSING | No `stopTask` implementation in Go pkg/task/. |
| `tasks/LocalMainSessionTask.ts` | `LocalMainSessionTaskState`, `createLocalMainSessionTask` | MISSING | Main-session backgrounding not ported. |
| `tasks/DreamTask/DreamTask.ts` | `DreamTurn`, `DreamPhase`, `DreamTaskState`, `isDreamTask`, `registerDreamTask` | MISSING | Dream/memory consolidation task not ported. `TaskType` enum has `TaskTypeDream` but no `DreamTaskState` struct. |
| `tasks/InProcessTeammateTask/types.ts` | `TeammateIdentity`, `InProcessTeammateTaskState` | MISSING | In-process teammate task state not ported. |
| `tasks/InProcessTeammateTask/InProcessTeammateTask.tsx` | Teammate task execution (React component + logic) | MISSING | Not ported. |
| `tasks/LocalAgentTask/LocalAgentTask.tsx` | `ToolActivity`, `AgentProgress`, `ProgressTracker`, `LocalAgentTaskState`, `createProgressTracker`, `getTokenCountFromTracker`, `updateProgressFromMessage`, `ActivityDescriptionResolver` | MISSING | Agent task state and progress tracking not ported. |
| `tasks/LocalShellTask/guards.ts` | `BashTaskKind`, `LocalShellTaskState`, `isLocalShellTask` | MISSING | Shell task state not ported. |
| `tasks/LocalShellTask/killShellTasks.ts` | `killShellTasksForAgent` | MISSING | Not ported. |
| `tasks/RemoteAgentTask/RemoteAgentTask.tsx` | `RemoteAgentTaskState`, `RemoteTaskType`, `RemoteTaskMetadata`, `AutofixPrRemoteTaskMetadata`, `RemoteTaskCompletionChecker`, `registerRemoteAgentTask` | MISSING | Remote agent task state not ported. |
| `tasks/LocalMainSessionTask.ts` | `LocalMainSessionTaskState` | MISSING | See above. |

---

## query/ (4 TS files) → pkg/query/

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| `query/config.ts` | `QueryConfig`, `buildQueryConfig` | MISSING | No `QueryConfig` struct or `BuildQueryConfig` in `pkg/query/engine.go`. Feature gates (statsig, env flags) not represented. |
| `query/deps.ts` | `QueryDeps`, `productionDeps` | MISSING | No dependency injection pattern in Go engine. Go calls API client directly. |
| `query/stopHooks.ts` | `handleStopHooks` (async generator), `StopHookResult` | MISSING | Stop hook execution pipeline not in Go. |
| `query/tokenBudget.ts` | `BudgetTracker`, `createBudgetTracker`, `TokenBudgetDecision`, `checkTokenBudget` | MISSING | Token budget tracking/continuation logic not ported. |

---

## schemas/ (1 TS file) → pkg/types/

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| `schemas/hooks.ts` | `BashCommandHookSchema`, `PromptHookSchema`, `AgentHookSchema`, `HookCommandSchema` (Zod), `HooksSettings` (map of event→hooks), `buildHookSchemas` | MISSING | No `HooksSettings` type (map of HookEvent to hook config arrays) in Go. Hook configuration schema not ported. |

---

## Top-level files → various pkg/

| TS File | Key Exports | Go Target | Go Coverage | Missing |
|---------|-------------|-----------|-------------|---------|
| `Tool.ts` | `ToolInputJSONSchema`, `ToolPermissionContext`, `getEmptyToolPermissionContext`, `CompactProgressEvent`, `ToolUseContext`, `QueryChainTracking`, `ValidationResult`, `SetToolJSXFn`, `ToolCallProgress`, `Tool<>` (generic interface), `Tools`, `toolMatchesName`, `findToolByName`, `buildTool`, `filterToolProgressMessages`, `ToolResult`, `Progress` | `pkg/tool/` | PARTIAL | Go `pkg/tool/tool.go` has `Tool` interface, `InputSchema`, `Result`, `ValidationResult`, `SearchOrReadInfo`, `InterruptBehavior`, `Tools`. Missing: `ToolUseContext` (large context object with options, app state, tools, etc.), `QueryChainTracking`, `CompactProgressEvent`, `SetToolJSXFn`, `toolMatchesName`/`findToolByName` helper functions, `buildTool` factory. `ToolPermissionContext` exists in `pkg/types/permissions.go` (partial). |
| `Task.ts` | `TaskType`, `TaskStatus`, `isTerminalTaskStatus`, `TaskHandle`, `SetAppState`, `TaskContext`, `TaskStateBase`, `LocalShellSpawnInput`, `Task` (interface), `generateTaskId`, `createTaskStateBase` | `pkg/task/` | PARTIAL | Go `pkg/task/task.go` has `TaskType`, `TaskStatus`, `IsTerminalTaskStatus`, `TaskHandle`, `TaskState` (base fields), `GenerateTaskId`, `NewTaskState`, `GetTaskOutputPath`. Missing: `SetAppState`, `TaskContext` (with abort controller), `LocalShellSpawnInput`, `Task` interface (kill method), `createTaskStateBase` equivalent. |
| `query.ts` | `QueryParams`, `query` (async generator — main loop) | `pkg/query/` | PARTIAL | Go `pkg/query/engine.go` has `Engine`, `RunQuery`, `StreamHandler`. Missing: `QueryParams` struct (full set of options including `verbose`, `dangerouslySkipPermissions`, `promptSuggestion`, `speculation`, `onMessage`, `hooks`, etc.), `query` as a standalone generator function. The Go engine is structurally similar but lacks the full parameter surface and does not implement stop hooks, token budget, speculation, or compaction. |
| `QueryEngine.ts` | `QueryEngineConfig`, `QueryEngine` (class), `ask` function | `pkg/query/` | PARTIAL | Go `Engine` covers basic conversation loop. Missing: `QueryEngineConfig` with full options, `ask` standalone function, headless profiler checkpoints, fast mode, file history snapshots, plugin loading, user input processing pipeline (`processUserInput`), session storage (`recordTranscript`), cost tracking integration, structured output enforcement. |
| `context.ts` | `getSystemPromptInjection`, `setSystemPromptInjection`, `getGitStatus`, `getSystemContext`, `getUserContext` | `pkg/context/` | PARTIAL | Go `pkg/context/` has `GetSystemContext`, `GetUserContext`, `LoadClaudeMd`, `BuildSystemPrompt`, git helpers. Missing: `GetSystemPromptInjection`/`SetSystemPromptInjection` (debug injection mechanism), memoization cache clearing on injection change. `getGitStatus` in TS returns rich formatted string; Go equivalent in `system.go` is similar but structured differently. |

---

## Summary

| Area | Total TS Files | FULL | PARTIAL | MISSING |
|------|---------------|------|---------|---------|
| types/ | 8 (+ 2 generated) | 1 (ids.ts) | 3 (hooks, permissions, schemas) | 4+ (command, logs, plugin, textInputTypes, generated) |
| state/ | 6 | 0 | 4 (AppState, AppStateStore, selectors, store) | 2 (onChangeAppState, teammateViewHelpers) |
| tasks/ | 12 | 0 | 1 (types.ts — TaskType/Status only) | 11 |
| query/ | 4 | 0 | 0 | 4 |
| schemas/ | 1 | 0 | 0 | 1 |
| Top-level | 5 | 0 | 5 (Tool, Task, query, QueryEngine, context) | 0 |
| **Total** | **36** | **1** | **13** | **22** |

### Critical gaps (blocking full port)

1. **Task variant states** — All 7 concrete task types (LocalShell, LocalAgent, Remote, InProcessTeammate, LocalWorkflow, MonitorMCP, Dream) have no Go struct definitions. Only the base `TaskState` exists.
2. **QueryParams / full query loop** — Token budget, stop hooks, speculation, compaction, file history, plugin loading, and user input processing pipeline are all absent.
3. **Teammate infrastructure** — `InProcessTeammateTaskState`, `TeammateIdentity`, `enterTeammateView`, `exitTeammateView`, `getActiveAgentForInput` not ported.
4. **AppState completeness** — Go `AppState` has ~15 fields vs ~50+ in TS. Remote bridge state, UI selection state, speculation state, attribution, denial tracking all missing.
5. **Session persistence (logs.ts)** — `SerializedMessage`, `LogOption`, transcript entry types, `sortLogs` — no Go equivalent at all.
6. **Permission richness** — Classifier types, `AdditionalWorkingDirectory`, `PermissionAllowDecision`/`AskDecision`/`DenyDecision`, `ToolPermissionRulesBySource` not ported.
