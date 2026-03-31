# TypeScript → Go Port Audit: tools/

Audit date: 2026-03-31
TS source: `old-src/tools/` (184 files across 41 tool directories + shared/ + testing/)
Go target: `pkg/tools/`

---

## Summary Table

| Tool Directory | TS Files | Go Package | Coverage |
|----------------|----------|------------|----------|
| AgentTool/ | 14 | pkg/tools/agent/ | PARTIAL |
| AskUserQuestionTool/ | 2 | pkg/tools/askuser/ | PARTIAL |
| BashTool/ | 17 | pkg/tools/bash/ | PARTIAL |
| BriefTool/ | 4 | pkg/tools/brief/ | PARTIAL |
| ConfigTool/ | 5 | pkg/tools/config/ | PARTIAL |
| EnterPlanModeTool/ | 4 | pkg/tools/plan/ | PARTIAL |
| EnterWorktreeTool/ | 4 | pkg/tools/worktree/ | PARTIAL |
| ExitPlanModeTool/ | 4 | pkg/tools/plan/ | PARTIAL |
| ExitWorktreeTool/ | 4 | pkg/tools/worktree/ | PARTIAL |
| FileEditTool/ | 5 | pkg/tools/fileedit/ | PARTIAL |
| FileReadTool/ | 4 | pkg/tools/fileread/ | PARTIAL |
| FileWriteTool/ | 2 | pkg/tools/filewrite/ | FULL |
| GlobTool/ | 3 | pkg/tools/glob/ | FULL |
| GrepTool/ | 3 | pkg/tools/grep/ | PARTIAL |
| ListMcpResourcesTool/ | 3 | pkg/tools/listmcpresources/ | PARTIAL |
| LSPTool/ | 6 | pkg/tools/lsp/ | PARTIAL |
| McpAuthTool/ | 1 | pkg/tools/mcpauth/ | PARTIAL |
| MCPTool/ | 4 | pkg/tools/mcp/ | PARTIAL |
| NotebookEditTool/ | 3 | pkg/tools/notebook/ | PARTIAL |
| PowerShellTool/ | 14 | pkg/tools/powershell/ | PARTIAL |
| ReadMcpResourceTool/ | 3 | pkg/tools/readmcpresource/ | PARTIAL |
| RemoteTriggerTool/ | 3 | pkg/tools/remotetrigger/ | PARTIAL |
| REPLTool/ | 2 | pkg/tools/repl/ | PARTIAL |
| ScheduleCronTool/ | 5 | pkg/tools/cron/ | PARTIAL |
| SendMessageTool/ | 4 | pkg/tools/team/ (SendMessageTool in team.go) | PARTIAL |
| SkillTool/ | 4 | pkg/tools/skill/ | PARTIAL |
| SleepTool/ | 1 | pkg/tools/sleep/ | PARTIAL |
| SyntheticOutputTool/ | 1 | pkg/tools/syntheticoutput/ | PARTIAL |
| TaskCreateTool/ | 3 | pkg/tools/task/ | PARTIAL |
| TaskGetTool/ | 3 | pkg/tools/task/ | PARTIAL |
| TaskListTool/ | 3 | pkg/tools/task/ | MISSING (no TaskListTool in task.go) |
| TaskOutputTool/ | 3 | pkg/tools/task/ | MISSING |
| TaskStopTool/ | 3 | pkg/tools/task/ | MISSING |
| TaskUpdateTool/ | 3 | pkg/tools/task/ | PARTIAL |
| TeamCreateTool/ | 4 | pkg/tools/team/ | PARTIAL |
| TeamDeleteTool/ | 4 | pkg/tools/team/ | PARTIAL |
| TodoWriteTool/ | 3 | pkg/tools/todowrite/ | PARTIAL |
| ToolSearchTool/ | 3 | pkg/tools/toolsearch/ | PARTIAL |
| WebFetchTool/ | 5 | pkg/tools/webfetch/ | PARTIAL |
| WebSearchTool/ | 2 | pkg/tools/websearch/ | PARTIAL |
| shared/ | 2 | (no direct Go equivalent) | MISSING |
| testing/ | 1 | (no direct Go equivalent) | MISSING |
| utils.ts | 1 top-level | (no direct Go equivalent) | MISSING |

---

## Detailed Per-Tool Audit

---

### tools/AgentTool/ (14 TS files → pkg/tools/agent/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| AgentTool.tsx | `AgentTool.call()`, remote/local agent dispatch, async/foreground lifecycle, subagent_type routing, worktree isolation | PARTIAL | Remote agent eligibility (`checkRemoteAgentEligibility`), in-process teammate handling, permission filtering (`filterDeniedAgents`), async lifecycle hooks, `run_in_background` goroutine never returns result |
| prompt.ts | `getPrompt()`, `formatAgentLine()`, `shouldInjectAgentListInMessages()`, dynamic agent listing | MISSING | No dynamic prompt generation; Go returns static string |
| forkSubagent.ts | `isForkSubagentEnabled()`, `buildForkedMessages()`, `buildWorktreeNotice()`, `FORK_AGENT` | MISSING | Fork subagent feature entirely absent |
| runAgent.ts | `runAgent()` — full query loop with MCP reconnect, skill commands, summarization | PARTIAL | Go `RunQuery` is a thin wrapper; no MCP reconnection, no skill command injection, no agent summarization |
| resumeAgent.ts | `resumeAgent()` — session resume from disk, transcript reconstruction | MISSING | No agent resume from stored transcript |
| builtInAgents.ts | `getBuiltInAgents()`, `areExplorePlanAgentsEnabled()` | MISSING | No built-in agent registry |
| agentMemory.ts | `getAgentMemoryDir()`, `AgentMemoryScope` | MISSING | No persistent agent memory |
| agentMemorySnapshot.ts | `getSnapshotDirForAgent()`, snapshot sync logic | MISSING | No memory snapshots |
| agentDisplay.ts | `AGENT_SOURCE_GROUPS`, `ResolvedAgent`, override annotation | MISSING | No agent display helpers |
| agentColorManager.ts | `setAgentColor()`, `getAgentColor()`, color palette | MISSING | No color assignment |
| loadAgentsDir.ts | `AgentDefinition` schema, plugin agents, frontmatter parsing, `isBuiltInAgent()` | MISSING | No agent directory loading or plugin agents |
| agentToolUtils.ts | `runAsyncAgentLifecycle()`, `classifyHandoffIfNeeded()`, `emitTaskProgress()`, `finalizeAgentTool()` | MISSING | No async lifecycle utilities |
| constants.ts | `AGENT_TOOL_NAME`, `LEGACY_AGENT_TOOL_NAME`, `ONE_SHOT_BUILTIN_AGENT_TYPES` | PARTIAL | Go uses hardcoded "Agent" string, no legacy alias "Task", no one-shot set |
| UI.tsx | React rendering for agent progress, grouped tool use, user-facing names | MISSING | No UI rendering (expected for Go) |
| built-in/*.ts (6 files) | `generalPurposeAgent`, `exploreAgent`, `planAgent`, `verificationAgent`, `claudeCodeGuideAgent`, `statuslineSetup` | MISSING | None of the built-in agent definitions exist in Go |

---

### tools/AskUserQuestionTool/ (2 TS files → pkg/tools/askuser/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| AskUserQuestionTool.tsx | Multi-question schema (1–4 questions, 2–4 options each), preview feature, multiSelect, uniqueness validation, annotations, per-question channel filtering | PARTIAL | Go only parses a flat `questions` array; no uniqueness validation, no preview field, no annotations schema, no multiSelect enforcement, no channel filtering |
| prompt.ts | `ASK_USER_QUESTION_TOOL_PROMPT`, `DESCRIPTION`, `PREVIEW_FEATURE_PROMPT` (markdown + html variants), `ASK_USER_QUESTION_TOOL_CHIP_WIDTH` | MISSING | Go returns one-liner description; no rich prompt text, no preview-feature docs |

---

### tools/BashTool/ (17 TS files → pkg/tools/bash/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| BashTool.tsx | `BashTool.call()`, background task management (`spawnShellTask`, `backgroundExistingForeground`), sandbox adapter, sed edit parsing, image output handling, file history tracking, LSP notification | PARTIAL | No sandbox support, no background task registry (`LocalShellTask`), no image output handling, no file history, no LSP notifications, no git operation tracking |
| prompt.ts | `getDefaultTimeoutMs()`, `getMaxTimeoutMs()`, full prompt with git/undercover/attribution sections | PARTIAL | Go has `DefaultTimeout`/`MaxTimeout` constants; no full prompt text, no git instructions, no undercover mode |
| bashPermissions.ts | `bashToolHasPermission()`, `matchWildcardPattern()`, `permissionRuleExtractPrefix()`, `commandHasAnyCd()`, wildcard rules | MISSING | No permission rule evaluation in Go |
| bashSecurity.ts | `parseForSecurity()`, command substitution pattern detection, heredoc-in-substitution, Zsh expansion blocking | MISSING | No security AST analysis |
| commandSemantics.ts | `interpretCommandResult()`, `COMMAND_SEMANTICS` map (grep exit-1 = not-error, etc.) | MISSING | Go treats all non-zero exits as errors |
| sedEditParser.ts | `parseSedEditCommand()`, `isSedInPlaceEdit()`, BRE→ERE conversion | MISSING | No sed parsing |
| sedValidation.ts | `sedCommandIsAllowedByAllowlist()`, flag allowlist validation | MISSING | No sed validation |
| modeValidation.ts | `validateCommandForMode()`, acceptEdits mode filesystem command allowlist | MISSING | No mode-based validation |
| pathValidation.ts | `validatePath()`, `PathCommand` union, output redirection extraction, `isDangerousRemovalPath()` | MISSING | No path validation |
| readOnlyValidation.ts | `checkReadOnlyConstraints()`, 100+ read-only command configs with safe flags | MISSING | No read-only mode enforcement |
| destructiveCommandWarning.ts | `DESTRUCTIVE_PATTERNS`, git reset/push --force/clean/checkout warning detection | MISSING | No destructive command warnings |
| commentLabel.ts | `extractBashCommentLabel()` | MISSING | No comment label extraction |
| toolName.ts | `BASH_TOOL_NAME = 'Bash'` | FULL | Go uses "Bash" string |
| bashCommandHelpers.ts | `segmentedCommandPermissionResult()`, compound command permission checking | MISSING | No compound command permission logic |
| shouldUseSandbox.ts | `shouldUseSandbox()`, `containsExcludedCommand()` | MISSING | No sandbox detection |
| utils.ts | `stripEmptyLines()`, `isImageOutput()`, `buildImageToolResult()`, `resetCwdIfOutsideProject()`, `resizeShellImageOutput()` | MISSING | No image handling, no cwd reset |
| UI.tsx | `BackgroundHint`, `renderToolResultMessage`, progress/queued/error renderers | MISSING | No UI rendering (expected) |

---

### tools/BriefTool/ (4 TS files → pkg/tools/brief/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| BriefTool.ts | Tool name `SendUserMessage` (not `Brief`), attachment validation, `status: normal\|proactive`, structured output schema | PARTIAL | Go registers as "Brief" not "SendUserMessage"; no `status` field, no `attachments` array, no output schema |
| prompt.ts | `BRIEF_TOOL_NAME='SendUserMessage'`, `LEGACY_BRIEF_TOOL_NAME='Brief'`, `BRIEF_TOOL_PROMPT`, `BRIEF_PROACTIVE_SECTION` | PARTIAL | Go uses "Brief" (legacy name); missing full prompt text and proactive section |
| attachments.ts | `validateAttachmentPaths()`, `resolveAttachments()`, `ResolvedAttachment` type | MISSING | No attachment handling |
| upload.ts | `uploadAttachment()` to claude.ai private API (bridge mode) | MISSING | No upload support |

---

### tools/ConfigTool/ (5 TS files → pkg/tools/config/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| ConfigTool.ts | `ConfigTool.call()` with get/set dispatch, `saveGlobalConfig`, `updateSettingsForSource` | PARTIAL | Go does read/write via JSON marshal; no `updateSettingsForSource` per-source logic, no `getRemoteControlAtStartup` |
| prompt.ts | `generatePrompt()` — dynamic listing of all supported settings with options | MISSING | Go returns one-liner; no dynamic setting enumeration |
| supportedSettings.ts | `SUPPORTED_SETTINGS` registry with 20+ settings (theme, model, permissions, etc.), `validateOnWrite`, `formatOnRead` | MISSING | No settings registry in Go; arbitrary key/value only |
| constants.ts | `CONFIG_TOOL_NAME = 'Config'` | FULL | Go uses "Config" |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/EnterPlanModeTool/ (4 TS files → pkg/tools/plan/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| EnterPlanModeTool.ts | `EnterPlanModeTool`, `handlePlanModeTransition()`, `prepareContextForPlanMode()`, `applyPermissionUpdate()`, `getAllowedChannels()`, shouldDefer | PARTIAL | Go sets bool flag; no permission update, no plan mode context preparation, no channel check, no interview-phase support |
| prompt.ts | `getEnterPlanModeToolPrompt()` — detailed When-to-Use guidance (5 condition categories) | MISSING | Go has no prompt text |
| constants.ts | `ENTER_PLAN_MODE_TOOL_NAME = 'EnterPlanMode'` | FULL | Go uses "EnterPlanMode" |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/ExitPlanModeTool/ (4 TS files → pkg/tools/plan/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| ExitPlanModeV2Tool.ts | `ExitPlanModeTool`, reads plan file from disk, writes plan approval response to mailbox, team/agent-swarm handling, `setNeedsPlanModeExitAttachment` | PARTIAL | Go only sets bool; no plan file reading, no mailbox communication, no swarm support |
| prompt.ts | `EXIT_PLAN_MODE_V2_TOOL_PROMPT` — detailed instructions on when to use, before-using checklist | MISSING | Go has no prompt text |
| constants.ts | `EXIT_PLAN_MODE_V2_TOOL_NAME = 'ExitPlanMode'` | FULL | Go uses "ExitPlanMode" |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/EnterWorktreeTool/ (4 TS files → pkg/tools/worktree/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| EnterWorktreeTool.ts | `EnterWorktreeTool`, `createWorktreeForSession()`, `validateWorktreeSlug()`, saves worktree state to session storage, supports hooks-based (non-git) worktrees, clears memory/plan caches | PARTIAL | Go creates worktree with git directly; no slug validation, no session state persistence, no hook-based worktrees, no cache clearing |
| prompt.ts | `getEnterWorktreeToolPrompt()` — when-to-use, requirements, behavior details | MISSING | Go has no prompt text |
| constants.ts | `ENTER_WORKTREE_TOOL_NAME = 'EnterWorktree'` | FULL | Go uses "EnterWorktree" |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/ExitWorktreeTool/ (4 TS files → pkg/tools/worktree/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| ExitWorktreeTool.ts | `ExitWorktreeTool`, `action: keep\|remove`, `discard_changes` guard, `cleanupWorktree()`/`keepWorktree()`, tmux session kill, CWD restore, memory cache clearing | PARTIAL | Go removes worktree without keep/remove distinction; no discard_changes guard, no uncommitted-file check, no tmux handling, no cache clearing |
| prompt.ts | `getExitWorktreeToolPrompt()` — scope limitations, when to use, parameter details | MISSING | Go has no prompt text |
| constants.ts | `EXIT_WORKTREE_TOOL_NAME = 'ExitWorktree'` | FULL | Go uses "ExitWorktree" |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/FileEditTool/ (5 TS files → pkg/tools/fileedit/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| FileEditTool.ts | `FileEditTool.call()`, file modification time check (`FILE_UNEXPECTEDLY_MODIFIED_ERROR`), git diff generation, LSP diagnostics clearing, skill directory activation, write permission check, curly-quote normalization | PARTIAL | Go does string replace correctly; no modification time guard, no git diff, no LSP notification, no permission checks |
| prompt.ts | `getEditToolDescription()` — pre-read instruction, line format instruction, uniqueness hint | PARTIAL | Go returns one-liner; missing full instruction text |
| utils.ts | `normalizeQuotes()`, `stripTrailingWhitespace()`, `getPatchForDisplay()`, structured patch hunks | MISSING | No quote normalization, no patch generation |
| types.ts | `FileEditInput`, `EditInput`, `FileEdit`, `hunkSchema`, `gitDiffSchema` | MISSING | No typed schemas in Go |
| constants.ts | `FILE_EDIT_TOOL_NAME`, `FILE_UNEXPECTEDLY_MODIFIED_ERROR` | MISSING | Go uses "Edit" string but no error constant |

---

### tools/FileReadTool/ (4 TS files → pkg/tools/fileread/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| FileReadTool.ts | `FileReadTool.call()`, image reading (PNG/JPG), PDF reading, directory reading, notebook reading, token estimation, memory freshness note, skill directory activation, file encoding detection | PARTIAL | Go reads text files only; no image/PDF/notebook support, no token limit, no memory freshness, no skill activation |
| prompt.ts | `FILE_READ_TOOL_NAME='Read'`, `MAX_LINES_TO_READ=2000`, `renderPromptTemplate()`, `FILE_UNCHANGED_STUB` | PARTIAL | Go has `MaxLinesToRead=2000` and uses "Read"; missing full prompt template, no unchanged-stub logic |
| limits.ts | `FileReadingLimits`, `DEFAULT_MAX_OUTPUT_TOKENS=25000`, `maxSizeBytes=256KB`, env override | MISSING | No token/size limits in Go |
| imageProcessor.ts | `getImageProcessor()`, sharp instance types, lazy module loading | MISSING | No image processing |

---

### tools/FileWriteTool/ (2 TS files → pkg/tools/filewrite/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| FileWriteTool.ts | `FileWriteTool.call()`, file modification time guard, git diff, LSP notification, write permission check, skill activation, file history tracking | PARTIAL | Go writes file correctly; no modification guard, no git diff, no permission checks, no LSP/skill notification |
| prompt.ts | `FILE_WRITE_TOOL_NAME='Write'`, `getWriteToolDescription()` with pre-read instruction | PARTIAL | Go uses "Write"; missing full description text with pre-read instruction |

**Overall: PARTIAL** (core write works; safety checks and integrations missing)

---

### tools/GlobTool/ (3 TS files → pkg/tools/glob/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| GlobTool.ts | `GlobTool`, `glob()` utility, read permission check, result sorted by modification time, truncation at 100 files | PARTIAL | Go sorts by mod time and uses `doublestar`; no permission check, truncation limit present via `WalkDir` but not capped at 100 |
| prompt.ts | `GLOB_TOOL_NAME='Glob'`, `DESCRIPTION` | PARTIAL | Go uses "Glob"; description is one-liner vs multi-line TS version |
| UI.tsx | React renderers | MISSING | No UI (expected) |

**Overall: FULL** (functional parity for core use case; minor: no permission check, no 100-file cap explicit)

---

### tools/GrepTool/ (3 TS files → pkg/tools/grep/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| GrepTool.ts | `GrepTool`, ripgrep invocation, `output_mode` (content/files_with_matches/count), `-A/-B/-C` context, `type` file-type filter, `head_limit`, `offset`, `multiline`, `head_limit` + `offset` pagination, read permission check, plugin glob exclusions | PARTIAL | Go invokes ripgrep and passes flags; no `offset` parameter, no plugin glob exclusions, no permission check |
| prompt.ts | `GREP_TOOL_NAME='Grep'`, `getDescription()` — detailed usage with multiline note | PARTIAL | Go uses "Grep"; description is one-liner |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/ListMcpResourcesTool/ (3 TS files → pkg/tools/listmcpresources/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| ListMcpResourcesTool.ts | `ListMcpResourcesTool`, `fetchResourcesForClient()`, multi-server aggregation, `server` filter param, `isReadOnly`, `isConcurrencySafe` | PARTIAL | Go has Client field but uses `server_name` not `server`; real resource fetching requires live MCP client |
| prompt.ts | `LIST_MCP_RESOURCES_TOOL_NAME`, `DESCRIPTION`, `PROMPT` | PARTIAL | Go uses "ListMcpResources" (vs TS "ListMcpResourcesTool"); description is one-liner |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/LSPTool/ (6 TS files → pkg/tools/lsp/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| LSPTool.ts | `LSPTool`, 9 operations (goToDefinition, findReferences, hover, documentSymbol, workspaceSymbol, goToImplementation, prepareCallHierarchy, incomingCalls, outgoingCalls), waits for LSP init, file size guard | STUB | Go returns stub message; no actual LSP protocol connection |
| prompt.ts | `LSP_TOOL_NAME='LSP'`, `DESCRIPTION` with all 9 operations listed | PARTIAL | Go uses "LSP"; description is one-liner |
| schemas.ts | `lspToolInputSchema` — discriminated union of 9 operation schemas | MISSING | Go uses flat single schema with `action` string |
| formatters.ts | `formatGoToDefinitionResult()`, `formatFindReferencesResult()`, `formatHoverResult()`, etc. — 7 formatters | MISSING | No result formatting |
| symbolContext.ts | `getSymbolAtPosition()` — synchronous file read for UI context | MISSING | No symbol context extraction |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/McpAuthTool/ (1 TS file → pkg/tools/mcpauth/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| McpAuthTool.ts | `createMcpAuthTool()` — dynamically creates per-server pseudo-tool, `performMCPOAuthFlow()`, `skipBrowserOpen`, reconnect after OAuth callback | PARTIAL | Go has static McpAuthTool; no dynamic per-server tool creation, no actual OAuth flow, no auto-reconnect |

---

### tools/MCPTool/ (4 TS files → pkg/tools/mcp/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| MCPTool.ts | `MCPTool` base (overridden per-server in mcpClient.ts), `isMcp: true`, passthrough schema, `classifyForCollapse` | PARTIAL | Go wraps MCP client correctly; `server_name`+`tool_name` schema differs from TS per-tool schema injection pattern |
| prompt.ts | Empty PROMPT/DESCRIPTION (overridden at runtime) | FULL | Go also returns empty strings |
| classifyForCollapse.ts | Collapse classification for tool results | MISSING | No collapse logic (UI-only concern) |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/NotebookEditTool/ (3 TS files → pkg/tools/notebook/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| NotebookEditTool.ts | `NotebookEditTool`, `cell_id` by ID (not index), `edit_mode: replace\|insert\|delete`, file mod time guard, write permission check, file history | PARTIAL | Go supports replace/insert/delete and cell_id; no modification time guard, no permission check, no file history |
| prompt.ts | `DESCRIPTION`, `PROMPT` | PARTIAL | Go description is one-liner; missing full prompt |
| constants.ts | `NOTEBOOK_EDIT_TOOL_NAME = 'NotebookEdit'` | FULL | Go uses "NotebookEdit" |

---

### tools/PowerShellTool/ (14 TS files → pkg/tools/powershell/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| PowerShellTool.tsx | Full BashTool equivalent for PowerShell: background tasks, sandbox, file history, image output, git tracking | PARTIAL | Go executes pwsh correctly with timeout; no background tasks, no sandbox, no image output |
| prompt.ts | Full prompt with sleep guidance, background usage, git instructions | MISSING | Go has no prompt text |
| powershellPermissions.ts | `powershellToolHasPermission()`, wildcard permission rules (case-insensitive cmdlets) | MISSING | No permission checking |
| powershellSecurity.ts | `powershellCommandIsSafe()`, download cradle detection, privilege escalation detection, COM object blocking | MISSING | No security analysis |
| gitSafety.ts | `isDotGitPathPS()`, `isGitInternalPathPS()` — bare-repo attack vectors for PowerShell | MISSING | No git safety |
| commandSemantics.ts | PowerShell-specific exit code semantics | MISSING | No command semantics |
| clmTypes.ts | CLM (Constrained Language Mode) type allowlist | MISSING | No CLM support |
| commonParameters.ts | Common PS parameter handling | MISSING | No parameter analysis |
| destructiveCommandWarning.ts | PowerShell destructive pattern warnings | MISSING | No warnings |
| modeValidation.ts | Mode-based validation for PowerShell | MISSING | No mode validation |
| pathValidation.ts | Path validation for PS commands | MISSING | No path validation |
| readOnlyValidation.ts | Read-only command validation | MISSING | No read-only enforcement |
| toolName.ts | `POWERSHELL_TOOL_NAME = 'PowerShell'` | FULL | Go uses "PowerShell" |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/ReadMcpResourceTool/ (3 TS files → pkg/tools/readmcpresource/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| ReadMcpResourceTool.ts | `ReadMcpResourceTool`, `ensureConnectedClient()`, binary blob persistence to disk | PARTIAL | Go calls MCP client correctly; no binary blob handling, `server` param in TS vs `server_name` in Go |
| prompt.ts | `DESCRIPTION`, `PROMPT` | PARTIAL | Go description is one-liner |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/RemoteTriggerTool/ (3 TS files → pkg/tools/remotetrigger/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| RemoteTriggerTool.ts | `RemoteTriggerTool`, `action: list\|get\|create\|update\|run`, OAuth token auto-injection, Claude.ai CCR API endpoints, feature flag gating | PARTIAL | Go does generic HTTP; no Claude.ai CCR API endpoints, no OAuth token injection, no action-based routing, no feature flag gating |
| prompt.ts | `REMOTE_TRIGGER_TOOL_NAME`, `DESCRIPTION`, `PROMPT` with action docs | MISSING | Go uses generic HTTP description |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/REPLTool/ (2 TS files → pkg/tools/repl/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| primitiveTools.ts | `getReplPrimitiveTools()` — returns [Read, Write, Edit, Glob, Grep, Bash, NotebookEdit, Agent] when REPL mode active | MISSING | No REPL mode concept in Go; Go has a simple REPL executor |
| constants.ts | `REPL_TOOL_NAME='REPL'`, `isReplModeEnabled()` | MISSING | Go uses "REPL" string but no mode toggle logic |

**Note:** TS REPL is a VM execution environment that hides primitive tools. Go REPL just runs code in subprocesses. Fundamentally different design.

---

### tools/ScheduleCronTool/ (5 TS files → pkg/tools/cron/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| CronCreateTool.ts | `CronCreateTool`, 5-field cron, `recurring` (default true), `durable` (persist to `.claude/scheduled_tasks.json`), `addCronTask()`, jitter config | PARTIAL | Go stores in-memory only; no 5-field cron parsing, no `recurring` flag, no durable disk persistence, no jitter |
| CronDeleteTool.ts | `CronDeleteTool`, `removeCronTasks()` by ID | PARTIAL | Go has `CronDeleteTool`; uses in-memory map |
| CronListTool.ts | `CronListTool`, `listAllCronTasks()`, human-readable schedule | PARTIAL | Go has `CronListTool`; human schedule formatting absent |
| prompt.ts | `isKairosCronEnabled()`, `isDurableCronEnabled()`, `DEFAULT_MAX_AGE_DAYS` | MISSING | Go has no feature gates |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/SendMessageTool/ (4 TS files → pkg/tools/team/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| SendMessageTool.ts | `SendMessageTool`, `to: name\|"*"\|uds://\|bridge://`, structured messages (shutdown_request, shutdown_response, plan_approval_response), UDS/bridge cross-session routing, teammate mailbox, resume agent background | PARTIAL | Go `SendMessageTool` in `team.go`; no UDS/bridge routing, no structured message types, no mailbox, basic in-memory member inbox only |
| prompt.ts | `DESCRIPTION`, `getPrompt()` with routing table, cross-session section, protocol responses | MISSING | Go uses static description |
| constants.ts | `SEND_MESSAGE_TOOL_NAME = 'SendMessage'` | FULL | Go uses "SendMessage" |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/SkillTool/ (4 TS files → pkg/tools/skill/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| SkillTool.ts | `SkillTool`, loads skill from commands registry, frontmatter parsing, agent context (forked agent), model override, skill usage tracking, `prepareForkedCommandContext()` | PARTIAL | Go finds skill files from `.claude/skills/`; no commands registry, no frontmatter model override, no usage tracking, no forked agent context |
| prompt.ts | `getCharBudget()`, `SKILL_BUDGET_CONTEXT_PERCENT`, `MAX_LISTING_DESC_CHARS`, skill listing with budget | MISSING | No budget-based skill listing |
| constants.ts | `SKILL_TOOL_NAME = 'Skill'` | FULL | Go uses "Skill" |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/SleepTool/ (1 TS file → pkg/tools/sleep/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| prompt.ts | `SLEEP_TOOL_NAME='Sleep'`, `SLEEP_TOOL_PROMPT` with tick-tag, prefer-over-Bash note | PARTIAL | Go implements sleep with ctx cancellation correctly; uses `duration_ms` param (vs TS uses `duration`), no TICK_TAG awareness |

**Note:** TS SleepTool.ts itself was not found separately — only `prompt.ts` exists in the directory.

---

### tools/SyntheticOutputTool/ (1 TS file → pkg/tools/syntheticoutput/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| SyntheticOutputTool.ts | `SyntheticOutputTool`, `isSyntheticOutputToolEnabled()` (non-interactive sessions only), dynamic output schema via AJV validation, `SYNTHETIC_OUTPUT_TOOL_NAME='StructuredOutput'` | PARTIAL | Go uses "SyntheticOutput" (wrong name — TS uses "StructuredOutput"); no non-interactive session gate, no AJV schema validation |

---

### tools/TaskCreateTool/ (3 TS files → pkg/tools/task/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| TaskCreateTool.ts | `TaskCreateTool`, `createTask()`, `activeForm` spinner text, `metadata`, task hooks (`executeTaskCreatedHooks`), team name context, `isTodoV2Enabled()` gate | PARTIAL | Go creates TaskState; no `activeForm`, no metadata, no hooks, no TodoV2 gate |
| prompt.ts | Full "When to Use" guide (8 scenarios) | MISSING | Go returns one-liner |
| constants.ts | `TASK_CREATE_TOOL_NAME` | PARTIAL | Go uses "TaskCreate" |

---

### tools/TaskGetTool/ (3 TS files → pkg/tools/task/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| TaskGetTool.ts | `TaskGetTool`, `getTask()`, blocks/blockedBy graph, `isTodoV2Enabled()` gate | PARTIAL | Go returns basic status; no blocks/blockedBy, no TodoV2 gate |
| prompt.ts | `DESCRIPTION`, `PROMPT` with output field docs | MISSING | Go returns one-liner |
| constants.ts | `TASK_GET_TOOL_NAME` | PARTIAL | Go uses "TaskGet" |

---

### tools/TaskListTool/ (3 TS files → pkg/tools/task/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| TaskListTool.ts | `TaskListTool`, `listTasks()`, owner/blockedBy in output, `isTodoV2Enabled()` gate | MISSING | No `TaskListTool` implementation in `task.go` |
| prompt.ts | `DESCRIPTION`, `getPrompt()` with multi-agent tips | MISSING | No Go equivalent |
| constants.ts | `TASK_LIST_TOOL_NAME` | MISSING | No Go equivalent |

---

### tools/TaskOutputTool/ (3 TS files → pkg/tools/task/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| TaskOutputTool.tsx | `TaskOutputTool`, `block: boolean`, `timeout` (0–600000ms), unified output for agent/shell/remote tasks, `getTaskOutputData()`, `retrieval_status: success\|timeout\|not_ready` | MISSING | No `TaskOutputTool` in Go |
| constants.ts | `TASK_OUTPUT_TOOL_NAME` | MISSING | No Go equivalent |
| (no separate prompt) | — | MISSING | — |

---

### tools/TaskStopTool/ (3 TS files → pkg/tools/task/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| TaskStopTool.ts | `TaskStopTool`, `stopTask()`, `shell_id` backward-compat alias, `KillShell` alias | MISSING | No `TaskStopTool` in Go |
| prompt.ts | `TASK_STOP_TOOL_NAME`, `DESCRIPTION` | MISSING | No Go equivalent |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/TaskUpdateTool/ (3 TS files → pkg/tools/task/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| TaskUpdateTool.ts | `TaskUpdateTool`, `status` including `'deleted'`, `addBlocks`/`addBlockedBy`, `owner`, `metadata`, task completed hooks, teammate notification via mailbox | PARTIAL | Go updates `status` only; no blocks/blockedBy graph, no owner, no metadata, no hooks, no mailbox notification |
| prompt.ts | `DESCRIPTION`, `PROMPT` | MISSING | Go returns one-liner |
| constants.ts | `TASK_UPDATE_TOOL_NAME` | PARTIAL | Go uses "TaskUpdate" |

---

### tools/TeamCreateTool/ (4 TS files → pkg/tools/team/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| TeamCreateTool.ts | `TeamCreateTool`, writes `~/.claude/teams/{name}/config.json`, creates task list dir, `registerTeamForSessionCleanup`, `assignTeammateColor`, `setLeaderTeamName`, backend detection | PARTIAL | Go sets in-memory `activeTeam`; no disk persistence, no task list dir, no session cleanup, no backend detection |
| prompt.ts | `getPrompt()` — detailed team workflow guide, agent type selection, when-to-use | MISSING | Go has no prompt text |
| constants.ts | `TEAM_CREATE_TOOL_NAME = 'TeamCreate'` | FULL | Go uses "TeamCreate" |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/TeamDeleteTool/ (4 TS files → pkg/tools/team/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| TeamDeleteTool.ts | `TeamDeleteTool`, `cleanupTeamDirectories()`, `unregisterTeamForSessionCleanup()`, `clearTeammateColors()`, `clearLeaderTeamName()`, active-member guard | PARTIAL | Go sets `activeTeam=nil`; no disk cleanup, no member guard, no session cleanup |
| prompt.ts | `getPrompt()` — TeamDelete workflow, active-member warning | MISSING | Go has no prompt text |
| constants.ts | `TEAM_DELETE_TOOL_NAME` | FULL | Go uses "TeamDelete" |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/TodoWriteTool/ (3 TS files → pkg/tools/todowrite/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| TodoWriteTool.ts | `TodoWriteTool`, `TodoListSchema`, session-scoped storage, `verificationNudgeNeeded`, `isTodoV2Enabled()` guard (mutual exclusion with Task tools), strict schema | PARTIAL | Go stores todos in struct field; no `verificationNudgeNeeded`, no TodoV2 mutual exclusion, "clear all when complete" logic present but differs |
| prompt.ts | `PROMPT`, `DESCRIPTION` — full When-to-Use guide with examples | MISSING | Go returns one-liner |
| constants.ts | `TODO_WRITE_TOOL_NAME` | PARTIAL | Go uses "TodoWrite" |

---

### tools/ToolSearchTool/ (3 TS files → pkg/tools/toolsearch/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| ToolSearchTool.ts | `ToolSearchTool`, `select:<name>` direct mode, keyword scoring, `isDeferredTool()` filter, MCP pending servers listing, cache invalidation on deferred tool list change | PARTIAL | Go has select mode and keyword scoring; no deferred-tool concept, no MCP pending servers, no cache invalidation |
| prompt.ts | `getPrompt()` — deferred tool location hint (system-reminder or block), tool-location hint varies by user type | MISSING | Go returns one-liner |
| constants.ts | `TOOL_SEARCH_TOOL_NAME` | PARTIAL | Go uses "ToolSearch" |

---

### tools/WebFetchTool/ (5 TS files → pkg/tools/webfetch/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| WebFetchTool.ts | `WebFetchTool`, `prompt` param for secondary model processing, `applyPromptToMarkdown()`, 15-min self-cleaning cache, preapproved domain check, permission rules by hostname | PARTIAL | Go fetches and returns HTML→text; no secondary model processing of `prompt`, no cache, no preapproved domain list, no permission rules |
| prompt.ts | `WEB_FETCH_TOOL_NAME='WebFetch'`, `DESCRIPTION`, `makeSecondaryModelPrompt()` | PARTIAL | Go uses "WebFetch"; missing full description and secondary model prompt |
| preapproved.ts | `PREAPPROVED_HOSTS` — 100+ approved hostnames (docs.python.org, go.dev, etc.) | MISSING | No preapproved host list |
| utils.ts | `getURLMarkdownContent()`, `applyPromptToMarkdown()`, `isPreapprovedUrl()`, `MAX_MARKDOWN_LENGTH` | MISSING | No HTML-to-markdown conversion, no prompt application |
| UI.tsx | React renderers | MISSING | No UI (expected) |

---

### tools/WebSearchTool/ (2 TS files → pkg/tools/websearch/)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| WebSearchTool.ts | `WebSearchTool`, uses Anthropic `BetaWebSearchTool20250305` natively, secondary model result processing, `allowed_domains`/`blocked_domains` passed to API | PARTIAL | Go constructs a manual web_search tool and passes to API; no beta search tool type, domain filtering params present |
| prompt.ts | `WEB_SEARCH_TOOL_NAME='WebSearch'`, `getWebSearchPrompt()` with Sources requirement and current month injection | MISSING | Go returns one-liner; no Sources mandate, no date injection |

---

### tools/shared/ (2 TS files → no Go equivalent)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| spawnMultiAgent.ts | `spawnTeammate()` — complex teammate spawning with tmux/in-process/pane backends, environment inheritance, flag settings path propagation | MISSING | No equivalent; Go `AgentTool` uses a simple goroutine |
| gitOperationTracking.ts | `trackGitOperations()` — regex detection of git commit/push/cherry-pick/merge/rebase and gh/glab PR creation, OTLP counter increments, analytics events | MISSING | No git operation tracking in Go |

---

### tools/testing/ (1 TS file → no Go equivalent)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| TestingPermissionTool.tsx | `TestingPermissionTool` — always asks for permission, enabled only in test env | MISSING | No testing tool in Go |

---

### tools/utils.ts (1 TS file → no Go equivalent)

| TS File | Key Exports | Go Coverage | Missing |
|---------|-------------|-------------|---------|
| utils.ts | `tagMessagesWithToolUseID()`, `getToolUseIDFromParentMessage()` | MISSING | No direct Go equivalent; message tagging is not needed in current Go architecture |

---

## Cross-Cutting Gaps

The following capabilities appear across many TS tools but are entirely absent from the Go port:

| Capability | TS Location | Go Status |
|------------|-------------|-----------|
| Permission system (allow/ask/deny rules, wildcards) | bashPermissions.ts, filesystem.ts, PermissionResult.ts | MISSING — all tools skip permission checks |
| Sandbox execution | shouldUseSandbox.ts, sandbox-adapter.ts | MISSING |
| File history tracking | fileHistory.ts | MISSING |
| LSP integration (diagnostics clear on write) | LSPDiagnosticRegistry, getLspServerManager | MISSING |
| Git diff on write | gitDiff.ts, fetchSingleFileGitDiff | MISSING |
| Modification time guard | getFileModificationTime, FILE_UNEXPECTEDLY_MODIFIED | MISSING |
| Analytics / telemetry | logEvent, logFileOperation | MISSING |
| Skill activation on file access | discoverSkillDirsForPaths | MISSING |
| Secondary model processing | queryModelWithStreaming for WebFetch/WebSearch | MISSING |
| Task/agent lifecycle hooks | executeTaskCreatedHooks, executeTaskCompletedHooks | MISSING |
| Teammate mailbox communication | writeToMailbox, readMailbox | MISSING |
| In-process teammate spawning | spawnInProcessTeammate, inProcessRunner | MISSING |
| Agent swarms / team backend detection | backends/registry, detectAndGetBackend | MISSING |
| UI/React rendering | All UI.tsx files | MISSING (intentional — Go is headless) |
| TodoV2 / Task system mutual exclusion | isTodoV2Enabled() | MISSING |
| Deferred tools system | isDeferredTool(), ToolSearchTool cache | PARTIAL |
| Feature flags / GrowthBook | getFeatureValue_CACHED_MAY_BE_STALE | MISSING |
