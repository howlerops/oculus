# UI Module Audit: TypeScript → Go Port

**Date:** 2026-03-31
**Auditor:** Executor agent
**Scope:** old-src/{hooks,ink,components,screens,keybindings,vim,context,constants} vs pkg/{hooks,tui,constants}

---

## Legend

| Rating | Meaning |
|--------|---------|
| **FULL** | Go equivalent exists and covers the same logic |
| **PARTIAL** | Go has something but it's incomplete or surface-level |
| **MISSING** | Logic should exist in Go but doesn't |
| **N/A** | React/Ink-specific plumbing with no Go equivalent needed (bubbletea handles differently) |

---

## 1. old-src/hooks/ (104 files) → pkg/hooks/

### Go coverage: pkg/hooks/engine.go + pkg/hooks/ui.go

The Go `hooks` package conflates two very different things: the **hook execution engine** (running external commands on events) in `engine.go`, and **UI utilities** (terminal size, blink, elapsed timer) in `ui.go`. In TS the hooks directory is entirely React hooks.

| TS File | Exports / Purpose | Go Equivalent | Rating |
|---------|-------------------|---------------|--------|
| `useTerminalSize.ts` | `useTerminalSize()` → reads from TerminalSizeContext | `pkg/hooks/ui.go: GetTerminalSize(), OnTerminalResize()` | **FULL** |
| `useBlink.ts` | `useBlink(enabled, intervalMs)` → synchronized cursor blink animation | `pkg/hooks/ui.go: BlinkState, NewBlinkState(), IsVisible(), Stop()` | **FULL** |
| `useElapsedTime.ts` | `useElapsedTime(startTime, isRunning, ms, pausedMs, endTime)` → formatted elapsed time | `pkg/hooks/ui.go: ElapsedTimer, NewElapsedTimer(), Elapsed(), Reset()` — note: Go returns `time.Duration`, not formatted string | **PARTIAL** — missing `FormatDuration()` wrapper and pause/endTime support |
| `useTextInput.ts` | Full text input state machine (cursor, kill ring, yank, modifiers, image paste) | None | **MISSING** — critical; bubbletea `textarea` is far simpler |
| `useVimInput.ts` | Vim input handler wrapping `useTextInput` with operator/motion/transition logic | `pkg/tui/vim.go: VimState.HandleKey()` | **PARTIAL** — Go vim lacks operators (d/c/y), text objects, find (f/F/t/T), count prefix, register, repeat (`.`), visual mode |
| `useHistorySearch.ts` | Search through conversation history by substring | None | **MISSING** |
| `useCommandQueue.ts` | Subscribe to unified command queue | None | **MISSING** — no equivalent queue manager |
| `useInputBuffer.ts` | Undo buffer with debounced snapshots for text input | None | **MISSING** |
| `useArrowKeyHistory.tsx` | Arrow-key navigation through prompt history with mode filtering | `pkg/tui/app.go: historyIdx + history[]` | **PARTIAL** — Go has basic up/down, no mode filtering or disk-backed history |
| `useDoublePress.ts` | `useDoublePress(setPending, onDoublePress, onFirstPress)` → timed double-press detection | None | **MISSING** |
| `useTimeout.ts` | `useTimeout(delay, resetTrigger)` → boolean that becomes true after delay | None | **MISSING** — trivial but absent |
| `useTypeahead.tsx` | Full autocomplete/suggestion engine (slash commands, paths, shell history, MCP, agents) | None | **MISSING** — large and complex |
| `useVirtualScroll.ts` | Virtualized list rendering with height estimation and overscan | None | **MISSING** — performance-critical for long transcripts |
| `useExitOnCtrlCD.ts` | Double-press Ctrl-C/D exit logic | `pkg/tui/app.go: tea.KeyCtrlC → tea.Quit` | **PARTIAL** — Go has single-press exit, no double-press pending state or display |
| `useExitOnCtrlCDWithKeybindings.ts` | Same but wired to keybinding system | None | **MISSING** |
| `useGlobalKeybindings.tsx` | Registers all global keybinding handlers | None | **MISSING** — Go has bindings config but no handler registration system |
| `useCommandKeybindings.tsx` | Keybindings for slash commands | None | **MISSING** |
| `useSettings.ts` | `useSettings()` → reactive settings from AppState | None (Go config is read at startup) | **N/A** — Go config loaded once; reactive updates not needed in bubbletea |
| `useDynamicConfig.ts` | `useDynamicConfig(name, default)` → GrowthBook feature flag | None | **N/A** — no GrowthBook in Go port |
| `useCanUseTool.tsx` | Full tool permission decision engine (coordinator, swarm worker, interactive handlers) | None | **MISSING** — critical permission logic |
| `useCancelRequest.ts` | Cancel current request, exit vim mode, dismiss overlays | `pkg/tui/app.go: m.cancel()` | **PARTIAL** — Go cancels but doesn't handle vim/overlay dismissal |
| `useLogMessages.ts` | Append-only transcript logger to disk | None | **MISSING** |
| `useMemoryUsage.ts` | Poll Node.js heap usage, return status ('normal'/'high'/'critical') | None | **N/A** — Go uses garbage collection; no heap monitoring needed (Go runtime handles this) |
| `usePasteHandler.ts` | Clipboard image paste detection, text paste debounce | None | **MISSING** |
| `useDiffData.ts` | Load and parse diff data for file changes | None | **MISSING** |
| `useDiffInIDE.ts` | Open diffs in IDE | None | **MISSING** |
| `useIdeConnectionStatus.ts` | Poll IDE connection state | None | **N/A** — Go port doesn't have IDE integration |
| `useIDEIntegration.tsx` | IDE bridge integration hook | None | **N/A** |
| `useIdeAtMentioned.ts` | Handle @-mentions from IDE selection | None | **N/A** |
| `useIdeSelection.ts` | Track IDE code selection | None | **N/A** |
| `useIdeLogging.ts` | Log to IDE channel | None | **N/A** |
| `useInboxPoller.ts` | Poll message inbox for background tasks | None | **MISSING** |
| `useMailboxBridge.ts` | Bridge between mailbox and REPL | None | **MISSING** |
| `useMainLoopModel.ts` | Track/switch main loop model | None | **MISSING** |
| `useManagePlugins.ts` | Plugin install/update/remove | None | **MISSING** — but plugin management is in `pkg/services/plugins/` (different layer) |
| `useMergedClients.ts` | Merge MCP clients from settings + plugins | None | **MISSING** |
| `useMergedCommands.ts` | Merge built-in + plugin commands | None | **MISSING** |
| `useMergedTools.ts` | Merge built-in + plugin tools | None | **MISSING** |
| `useMinDisplayTime.ts` | Enforce minimum display duration for messages | None | **MISSING** |
| `useNotifyAfterTimeout.ts` | Show notification if operation exceeds timeout | None | **MISSING** |
| `usePrStatus.ts` | Fetch GitHub PR status | None | **MISSING** |
| `usePromptSuggestion.ts` | First-run prompt suggestions | None | **MISSING** |
| `useQueueProcessor.ts` | Process queued commands sequentially | None | **MISSING** |
| `useRemoteSession.ts` | Remote session management | None | **MISSING** |
| `useReplBridge.tsx` | REPL ↔ remote bridge | None | **MISSING** |
| `useScheduledTasks.ts` | Background scheduled task management | None | **MISSING** |
| `useSearchInput.ts` | Transcript search input handling | None | **MISSING** |
| `useSessionBackgrounding.ts` | Background session management | None | **MISSING** |
| `useSettingsChange.ts` | Detect settings file changes on disk | None | **MISSING** |
| `useSkillsChange.ts` | Detect skills directory changes | None | **MISSING** |
| `useSSHSession.ts` | SSH session management | None | **MISSING** |
| `useSwarmInitialization.ts` | Initialize swarm/multi-agent mode | None | **MISSING** |
| `useSwarmPermissionPoller.ts` | Poll pending permissions from swarm workers | None | **MISSING** |
| `useTaskListWatcher.ts` | Watch TodoWrite task list for changes | None | **MISSING** |
| `useTasksV2.ts` | Background task state management | None | **MISSING** |
| `useTeammateViewAutoExit.ts` | Auto-exit teammate view when done | None | **MISSING** |
| `useTeleportResume.tsx` | Resume teleport sessions | None | **MISSING** |
| `useAssistantHistory.ts` | Assistant message history | None | **MISSING** |
| `useAwaySummary.ts` | Away/idle summary generation | None | **MISSING** |
| `useBackgroundTaskNavigation.ts` | Navigate between background tasks | None | **MISSING** |
| `useApiKeyVerification.ts` | Verify API key on startup | None | **MISSING** |
| `useAfterFirstRender.ts` | Run effect after first render | None | **N/A** — React lifecycle, no Go equivalent needed |
| `useCopyOnSelect.ts` | Terminal text selection copy | None | **N/A** — terminal-level feature |
| `useDeferredHookMessages.ts` | Deferred hook system messages | None | **MISSING** |
| `useDirectConnect.ts` | Direct connection mode | None | **MISSING** |
| `useFileHistorySnapshotInit.ts` | Initialize file history snapshots | None | **MISSING** |
| `useIssueFlagBanner.ts` | Issue flag banner display | None | **MISSING** |
| `useClipboardImageHint.ts` | Hint when clipboard has image | None | **MISSING** |
| `useClaudeCodeHintRecommendation.tsx` | Plugin hint menu recommendation | None | **MISSING** |
| `useLspPluginRecommendation.tsx` | LSP plugin recommendation | None | **MISSING** |
| `useOfficialMarketplaceNotification.tsx` | Marketplace notification | None | **MISSING** |
| `usePluginRecommendationBase.tsx` | Base plugin recommendation hook | None | **MISSING** |
| `usePromptsFromClaudeInChrome.tsx` | Claude-in-Chrome prompt integration | None | **N/A** |
| `useChromeExtensionNotification.tsx` | Chrome extension notification | None | **N/A** |
| `useVoice.ts` | Voice recording/transcription state | None | **N/A** — voice not planned for Go port |
| `useVoiceEnabled.ts` | Check if voice is enabled | None | **N/A** |
| `useVoiceIntegration.tsx` | Voice integration hook | None | **N/A** |
| `useTurnDiffs.ts` | Track file diffs per turn | None | **MISSING** |
| `useUpdateNotification.ts` | Auto-update notification | None | **MISSING** |
| `renderPlaceholder.ts` | Render ghost/placeholder text | None | **MISSING** |
| `fileSuggestions.ts` | Background file index for path completions | None | **MISSING** |
| `unifiedSuggestions.ts` | Unified suggestion pipeline | None | **MISSING** |
| `toolPermission/PermissionContext.ts` | Permission queue context | None | **MISSING** |
| `toolPermission/permissionLogging.ts` | Log permission decisions | None | **MISSING** |
| `toolPermission/handlers/coordinatorHandler.ts` | Coordinator permission handler | None | **MISSING** |
| `toolPermission/handlers/interactiveHandler.ts` | Interactive permission handler | None | **MISSING** |
| `toolPermission/handlers/swarmWorkerHandler.ts` | Swarm worker permission handler | None | **MISSING** |
| `notifs/useAutoModeUnavailableNotification.ts` | Auto-mode unavailable notification | None | **N/A** — notification system not ported |
| `notifs/useCanSwitchToExistingSubscription.tsx` | Subscription switch notification | None | **N/A** |
| `notifs/useDeprecationWarningNotification.tsx` | Deprecation warning | None | **N/A** |
| `notifs/useFastModeNotification.tsx` | Fast mode notification | None | **N/A** |
| `notifs/useIDEStatusIndicator.tsx` | IDE status indicator | None | **N/A** |
| `notifs/useInstallMessages.tsx` | Installation messages | None | **N/A** |
| `notifs/useLspInitializationNotification.tsx` | LSP init notification | None | **N/A** |
| `notifs/useMcpConnectivityStatus.tsx` | MCP connectivity notification | None | **MISSING** — MCP exists in Go |
| `notifs/useModelMigrationNotifications.tsx` | Model migration notification | None | **N/A** |
| `notifs/useNpmDeprecationNotification.tsx` | npm deprecation notification | None | **N/A** |
| `notifs/usePluginAutoupdateNotification.tsx` | Plugin autoupdate notification | None | **N/A** |
| `notifs/usePluginInstallationStatus.tsx` | Plugin install status | None | **N/A** |
| `notifs/useRateLimitWarningNotification.tsx` | Rate limit warning | None | **MISSING** — rate limiting exists in Go |
| `notifs/useSettingsErrors.tsx` | Settings error notifications | None | **MISSING** |
| `notifs/useStartupNotification.ts` | Startup notifications | None | **N/A** |
| `notifs/useTeammateShutdownNotification.ts` | Teammate shutdown notification | None | **N/A** |

### Section Summary
- **FULL:** 2 (`useTerminalSize`, `useBlink`)
- **PARTIAL:** 5 (`useElapsedTime`, `useVimInput`, `useArrowKeyHistory`, `useExitOnCtrlCD`, `useCancelRequest`)
- **MISSING:** ~60 files — most application logic hooks have no Go equivalent
- **N/A:** ~17 files (React lifecycle, voice, Chrome extension, Node.js heap, IDE bridge)

---

## 2. old-src/ink/ (96 files) → pkg/tui/

The `ink/` directory is the entire custom Ink (React-ink) fork with its own React reconciler, layout engine (yoga), terminal I/O, event system, and rendering pipeline. **Bubbletea replaces the entire Ink framework** — there is no 1:1 mapping. Most files are N/A.

### Core Framework Files (N/A — replaced by bubbletea)

| TS File | Purpose | Rating |
|---------|---------|--------|
| `ink.tsx` | Main Ink entry point, `render()` function | **N/A** — `tea.NewProgram()` |
| `reconciler.ts` | React reconciler for terminal DOM | **N/A** |
| `renderer.ts` | Diff-based terminal renderer | **N/A** — bubbletea `View()` |
| `root.ts` | React root management | **N/A** |
| `screen.ts` | Screen buffer management | **N/A** |
| `output.ts` | Output stream management | **N/A** |
| `dom.ts` | Virtual DOM types | **N/A** |
| `frame.ts` | Frame rendering | **N/A** |
| `instances.ts` | Ink instance registry | **N/A** |
| `optimizer.ts` | Render optimization | **N/A** |
| `focus.ts` | Focus management | **N/A** |
| `node-cache.ts` | Node cache | **N/A** |
| `squash-text-nodes.ts` | Text node squashing | **N/A** |
| `log-update.ts` | Log-update integration | **N/A** |

### Layout Engine (N/A — bubbletea uses lipgloss)

| TS File | Purpose | Rating |
|---------|---------|--------|
| `layout/engine.ts` | Yoga layout engine integration | **N/A** — lipgloss |
| `layout/geometry.ts` | Geometry types | **N/A** |
| `layout/node.ts` | Layout node | **N/A** |
| `layout/yoga.ts` | Yoga WASM bindings | **N/A** |
| `styles.ts` | CSS-like style system | **N/A** — lipgloss styles |
| `get-max-width.ts` | Max width calculation | **N/A** |
| `measure-element.ts` | Element measurement | **N/A** |
| `hit-test.ts` | Mouse click hit testing | **N/A** |

### Text/String Utilities (PARTIAL/MISSING — ported logic needed)

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `measure-text.ts` | Measure text width+height with ANSI awareness | None in tui | **MISSING** — needed for layout |
| `wrap-text.ts` | Word-wrap text respecting ANSI codes | None | **MISSING** |
| `wrapAnsi.ts` | ANSI-aware string wrapping | None | **MISSING** |
| `stringWidth.ts` | Unicode-aware string width (CJK, emoji) | None | **MISSING** — `runewidth` library would cover this |
| `widest-line.ts` | Find widest line in multi-line string | None | **MISSING** |
| `line-width-cache.ts` | LRU cache for line width | None | **MISSING** |
| `searchHighlight.ts` | Apply search highlight spans to text | None | **MISSING** |
| `selection.ts` | Text selection range tracking | None | **MISSING** |
| `Ansi.tsx` | `<Ansi>` component: renders raw ANSI string | None | **N/A** — render via lipgloss |
| `render-node-to-output.ts` | Render DOM node to output buffer | **N/A** | N/A |
| `render-border.ts` | Render box borders | **N/A** — lipgloss borders | N/A |
| `render-to-screen.ts` | Full-screen render | **N/A** | N/A |

### Terminal I/O (MISSING — critical)

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `parse-keypress.ts` | Parse terminal escape sequences to key events | None in tui | **MISSING** — bubbletea handles some but custom sequences need this |
| `terminal.ts` | Terminal capabilities detection | None | **MISSING** |
| `terminal-focus-state.ts` | Terminal focus tracking | None | **MISSING** |
| `terminal-querier.ts` | Query terminal for capabilities (DA2, etc.) | None | **MISSING** |
| `termio.ts` | Terminal I/O main module | None | **MISSING** |
| `termio/ansi.ts` | ANSI C0/C1 constants | None | **MISSING** |
| `termio/csi.ts` | CSI escape sequence constants and parser | None | **MISSING** |
| `termio/dec.ts` | DEC private mode sequences | None | **MISSING** |
| `termio/esc.ts` | ESC sequence parser | None | **MISSING** |
| `termio/osc.ts` | OSC sequence parser (terminal title, clipboard, hyperlinks) | None | **MISSING** |
| `termio/sgr.ts` | SGR (color/style) parser | None | **MISSING** |
| `termio/tokenize.ts` | Streaming escape sequence tokenizer | None | **MISSING** |
| `termio/types.ts` | Termio type definitions | None | **MISSING** |
| `tabstops.ts` | Tab stop tracking | None | **MISSING** |
| `bidi.ts` | Bidirectional text support | None | **MISSING** |
| `colorize.ts` | Apply colors to text | None | **N/A** — lipgloss |
| `clearTerminal.ts` | Clear terminal screen | None | **MISSING** |
| `supports-hyperlinks.ts` | Detect OSC 8 hyperlink support | None | **MISSING** |
| `warn.ts` | Warning output helper | None | **N/A** — Go uses log package |
| `useTerminalNotification.ts` | Send OS/terminal notifications (OSC) | None | **MISSING** |

### Events System (N/A — replaced by bubbletea messages)

| TS File | Purpose | Rating |
|---------|---------|--------|
| `events/click-event.ts` | Mouse click event type | **N/A** |
| `events/dispatcher.ts` | Event dispatcher | **N/A** |
| `events/emitter.ts` | Event emitter | **N/A** |
| `events/event-handlers.ts` | Event handler registry | **N/A** |
| `events/event.ts` | Base event type | **N/A** |
| `events/focus-event.ts` | Focus event | **N/A** |
| `events/input-event.ts` | Input event | **N/A** |
| `events/keyboard-event.ts` | Keyboard event with modifier flags | **N/A** |
| `events/terminal-event.ts` | Terminal event | **N/A** |
| `events/terminal-focus-event.ts` | Terminal focus event | **N/A** |

### Ink Hooks (N/A — replaced by bubbletea `tea.Msg` pattern)

| TS File | Purpose | Rating |
|---------|---------|--------|
| `hooks/use-animation-frame.ts` | RAF-like animation tick | **N/A** — `spinner.Tick` |
| `hooks/use-app.ts` | `useApp()` → exit/clear methods | **N/A** |
| `hooks/use-declared-cursor.ts` | Cursor shape management | **N/A** |
| `hooks/use-input.ts` | `useInput(handler)` → key input hook | **N/A** — bubbletea `tea.KeyMsg` |
| `hooks/use-interval.ts` | `useInterval(callback, ms)` | **N/A** — `time.Ticker` |
| `hooks/use-search-highlight.ts` | Search highlight state | **N/A** |
| `hooks/use-selection.ts` | Text selection state | **N/A** |
| `hooks/use-stdin.ts` | Access stdin stream | **N/A** |
| `hooks/use-tab-status.ts` | Terminal tab status | **N/A** |
| `hooks/use-terminal-focus.ts` | Terminal focus state | **N/A** |
| `hooks/use-terminal-title.ts` | Set terminal window title | **MISSING** — useful functionality absent |
| `hooks/use-terminal-viewport.ts` | Viewport dimensions | **N/A** — `tea.WindowSizeMsg` |

### Components (N/A — no direct Go equivalent; bubbletea uses View())

| TS File | Purpose | Rating |
|---------|---------|--------|
| `components/AlternateScreen.tsx` | Alternate screen buffer | **N/A** — `tea.WithAltScreen()` |
| `components/App.tsx` | App root component | **N/A** |
| `components/AppContext.ts` | App React context | **N/A** |
| `components/Box.tsx` | Flex layout box | **N/A** — lipgloss |
| `components/Button.tsx` | Button component | **N/A** |
| `components/ClockContext.tsx` | Clock context for animations | **N/A** |
| `components/CursorDeclarationContext.ts` | Cursor shape context | **N/A** |
| `components/ErrorOverview.tsx` | Error display | **N/A** |
| `components/Link.tsx` | Hyperlink (OSC 8) | **MISSING** — OSC 8 links useful |
| `components/Newline.tsx` | Newline component | **N/A** |
| `components/NoSelect.tsx` | Disable text selection | **N/A** |
| `components/RawAnsi.tsx` | Raw ANSI pass-through | **N/A** |
| `components/ScrollBox.tsx` | Scrollable container | **MISSING** — complex scroll logic needed |
| `components/Spacer.tsx` | Flexible spacer | **N/A** — lipgloss |
| `components/StdinContext.ts` | Stdin React context | **N/A** |
| `components/TerminalFocusContext.tsx` | Terminal focus context | **N/A** |
| `components/TerminalSizeContext.tsx` | Terminal size context | **N/A** — `tea.WindowSizeMsg` |
| `components/Text.tsx` | Styled text component | **N/A** — lipgloss |
| `constants.ts` | Ink internal constants | **N/A** |

### Section Summary
- **FULL:** 0
- **PARTIAL:** 0
- **MISSING:** ~25 files (text measurement, terminal I/O, escape sequences, scroll, title)
- **N/A:** ~70 files (entire React reconciler, layout engine, event system, Ink hooks, base components — all replaced by bubbletea+lipgloss)

---

## 3. old-src/components/ (389 files) → pkg/tui/components/

The Go `pkg/tui/components/components.go` is a single 112-line file. It implements only: `PermissionDialog`, `ProgressBar`, `TokenCounter`, `SpinnerFrames`, `MessageBubble`, and basic lipgloss styles.

**Important:** In bubbletea, components are rendered as strings from `View()` methods. The TS components are React JSX trees. Most are **N/A** as JSX components, but the **logic** they contain (permission decisions, diff rendering, etc.) should be ported.

### Design System

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `design-system/color.ts` | Color palette constants | lipgloss colors in components.go | **PARTIAL** — no named palette |
| `design-system/Dialog.tsx` | Modal dialog frame | None | **MISSING** |
| `design-system/Divider.tsx` | Horizontal divider | None | **MISSING** |
| `design-system/FuzzyPicker.tsx` | Fuzzy search picker | None | **MISSING** |
| `design-system/KeyboardShortcutHint.tsx` | Display keyboard shortcut | None | **MISSING** |
| `design-system/ListItem.tsx` | List item display | None | **MISSING** |
| `design-system/LoadingState.tsx` | Loading spinner display | `components.go: SpinnerFrames` | **PARTIAL** |
| `design-system/Pane.tsx` | Content pane with borders | None | **MISSING** |
| `design-system/ProgressBar.tsx` | Progress bar | `components.go: ProgressBar()` | **FULL** |
| `design-system/Ratchet.tsx` | Minimum-height expanding area | None | **MISSING** |
| `design-system/StatusIcon.tsx` | Status indicator icon | None | **MISSING** |
| `design-system/Tabs.tsx` | Tabbed interface | None | **MISSING** |
| `design-system/ThemedBox.tsx` | Theme-aware box | None | **N/A** |
| `design-system/ThemedText.tsx` | Theme-aware text | None | **N/A** |
| `design-system/ThemeProvider.tsx` | Theme React context provider | None | **N/A** |
| `design-system/Byline.tsx` | Byline display | None | **N/A** |

### Permissions

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `permissions/PermissionDialog.tsx` | Main permission dialog frame | `components.go: PermissionDialog()` | **PARTIAL** — Go has basic dialog; TS has full decision flow |
| `permissions/PermissionRequest.tsx` | Permission request dispatching | None | **MISSING** |
| `permissions/PermissionPrompt.tsx` | Permission prompt UI | None | **MISSING** |
| `permissions/BashPermissionRequest/BashPermissionRequest.tsx` | Bash tool permission | None | **MISSING** |
| `permissions/BashPermissionRequest/bashToolUseOptions.tsx` | Bash permission options | None | **MISSING** |
| `permissions/FileEditPermissionRequest/` | File edit permission | None | **MISSING** |
| `permissions/FileWritePermissionRequest/` | File write permission + diff | None | **MISSING** |
| `permissions/FilesystemPermissionRequest/` | Filesystem permission | None | **MISSING** |
| `permissions/FilePermissionDialog/` | File permission dialog with IDE diff | None | **MISSING** |
| `permissions/PowerShellPermissionRequest/` | PowerShell permission | None | **N/A** — no Windows in Go port |
| `permissions/NotebookEditPermissionRequest/` | Notebook edit permission | None | **MISSING** |
| `permissions/AskUserQuestionPermissionRequest/` | Ask user question UI | None | **MISSING** |
| `permissions/EnterPlanModePermissionRequest/` | Enter plan mode | None | **MISSING** |
| `permissions/ExitPlanModePermissionRequest/` | Exit plan mode | None | **MISSING** |
| `permissions/ComputerUseApproval/` | Computer use approval | None | **N/A** |
| `permissions/SandboxPermissionRequest.tsx` | Sandbox permission | None | **MISSING** |
| `permissions/SedEditPermissionRequest/` | Sed edit permission | None | **MISSING** |
| `permissions/SkillPermissionRequest/` | Skill permission | None | **MISSING** |
| `permissions/WebFetchPermissionRequest/` | Web fetch permission | None | **MISSING** |
| `permissions/FallbackPermissionRequest.tsx` | Fallback permission UI | None | **MISSING** |
| `permissions/PermissionDecisionDebugInfo.tsx` | Debug info display | None | **N/A** |
| `permissions/PermissionExplanation.tsx` | Explain permission decision | None | **MISSING** |
| `permissions/PermissionRuleExplanation.tsx` | Permission rule explanation | None | **MISSING** |
| `permissions/PermissionRequestTitle.tsx` | Permission dialog title | None | **MISSING** |
| `permissions/rules/` (8 files) | Permission rule management UI | None | **MISSING** |
| `permissions/WorkerBadge.tsx` | Worker agent badge | None | **MISSING** |
| `permissions/WorkerPendingPermission.tsx` | Worker pending permission | None | **MISSING** |
| `permissions/hooks.ts` | Permission helper hooks | None | **MISSING** |
| `permissions/shellPermissionHelpers.tsx` | Shell permission helpers | None | **MISSING** |
| `permissions/useShellPermissionFeedback.ts` | Shell permission feedback | None | **MISSING** |
| `permissions/utils.ts` | Permission utilities | None | **MISSING** |
| `permissions/BypassPermissionsModeDialog.tsx` | Bypass mode dialog | None | **MISSING** |

### Messages

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `messages/AssistantMessage.tsx` | Render assistant message | `pkg/tui/app.go: View() assistant case` | **PARTIAL** — Go has minimal rendering |
| `messages/UserPromptMessage.tsx` | Render user message | `pkg/tui/app.go: View() user case` | **PARTIAL** |
| `messages/ToolUseMessage.tsx` | Render tool use block | None | **MISSING** |
| `messages/ToolResultMessage.tsx` | Render tool result | None | **MISSING** |
| `messages/SystemTextMessage.tsx` | Render system message | None | **MISSING** |
| `messages/SystemAPIErrorMessage.tsx` | API error display | None | **MISSING** |
| `messages/RateLimitMessage.tsx` | Rate limit display | None | **MISSING** |
| `messages/ShutdownMessage.tsx` | Shutdown message | None | **MISSING** |
| `messages/UserBashInputMessage.tsx` | Bash input display | None | **MISSING** |
| `messages/UserBashOutputMessage.tsx` | Bash output display | None | **MISSING** |
| `messages/UserImageMessage.tsx` | Image display | None | **MISSING** |
| `messages/UserToolResultMessage/` (8 files) | Tool result variants | None | **MISSING** |
| All other message types (15+) | Various message renderers | None | **MISSING** |

### Prompt Input

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `PromptInput/PromptInput.tsx` | Main prompt input component | `pkg/tui/app.go: textarea.Model` | **PARTIAL** — Go has basic textarea |
| `PromptInput/PromptInputFooter.tsx` | Footer with suggestions | None | **MISSING** |
| `PromptInput/PromptInputFooterSuggestions.tsx` | Autocomplete suggestions | None | **MISSING** |
| `PromptInput/PromptInputFooterLeftSide.tsx` | Footer left side (mode, hints) | None | **MISSING** |
| `PromptInput/PromptInputHelpMenu.tsx` | Help menu | None | **MISSING** |
| `PromptInput/PromptInputModeIndicator.tsx` | Input mode indicator | None | **MISSING** |
| `PromptInput/PromptInputQueuedCommands.tsx` | Queued commands display | None | **MISSING** |
| `PromptInput/PromptInputStashNotice.tsx` | Stash notice | None | **MISSING** |
| `PromptInput/HistorySearchInput.tsx` | History search input | None | **MISSING** |
| `PromptInput/Notifications.tsx` | Input area notifications | None | **MISSING** |
| `PromptInput/ShimmeredInput.tsx` | Shimmering loading input | None | **MISSING** |
| `PromptInput/VoiceIndicator.tsx` | Voice recording indicator | None | **N/A** |
| `PromptInput/SandboxPromptFooterHint.tsx` | Sandbox hint | None | **MISSING** |
| `PromptInput/IssueFlagBanner.tsx` | Issue flag banner | None | **MISSING** |
| `PromptInput/inputModes.ts` | Input mode logic (vim/normal/etc) | None | **MISSING** |
| `PromptInput/inputPaste.ts` | Paste handling logic | None | **MISSING** |
| `PromptInput/utils.ts` | Prompt input utilities | None | **MISSING** |
| `PromptInput/useMaybeTruncateInput.ts` | Truncate long input | None | **MISSING** |
| `PromptInput/usePromptInputPlaceholder.ts` | Placeholder text logic | None | **MISSING** |
| `PromptInput/useShowFastIconHint.ts` | Fast mode icon hint | None | **MISSING** |
| `PromptInput/useSwarmBanner.ts` | Swarm mode banner | None | **MISSING** |

### Diff Rendering

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `StructuredDiff.tsx` | Structured diff rendering | None | **MISSING** |
| `StructuredDiff/colorDiff.ts` | Color diff algorithm | None | **MISSING** |
| `StructuredDiff/Fallback.tsx` | Diff fallback renderer | None | **MISSING** |
| `StructuredDiffList.tsx` | List of diffs | None | **MISSING** |
| `diff/DiffDetailView.tsx` | Diff detail view | None | **MISSING** |
| `diff/DiffDialog.tsx` | Diff dialog | None | **MISSING** |
| `diff/DiffFileList.tsx` | Diff file list | None | **MISSING** |

### Spinner / Progress

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `Spinner/index.ts` | Spinner entry point | `pkg/tui/app.go: spinner.Model` | **PARTIAL** |
| `Spinner/SpinnerGlyph.tsx` | Spinner animation glyph | `components.go: SpinnerFrames` | **PARTIAL** |
| `Spinner/SpinnerAnimationRow.tsx` | Spinner with verb | None | **MISSING** |
| `Spinner/GlimmerMessage.tsx` | Shimmer loading message | None | **MISSING** |
| `Spinner/ShimmerChar.tsx` | Shimmer character | None | **MISSING** |
| `Spinner/FlashingChar.tsx` | Flashing character | None | **MISSING** |
| `Spinner/TeammateSpinnerLine.tsx` | Teammate spinner | None | **MISSING** |
| `Spinner/TeammateSpinnerTree.tsx` | Teammate spinner tree | None | **MISSING** |
| `Spinner/useShimmerAnimation.ts` | Shimmer animation hook | None | **MISSING** |
| `Spinner/useStalledAnimation.ts` | Stalled animation hook | None | **MISSING** |
| `Spinner/utils.ts` | Spinner utilities | None | **MISSING** |
| `Spinner/teammateSelectHint.ts` | Teammate select hint | None | **MISSING** |
| `BashModeProgress.tsx` | Bash mode progress display | None | **MISSING** |
| `AgentProgressLine.tsx` | Agent progress line | None | **MISSING** |

### Settings

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `Settings/Settings.tsx` | Settings screen | None | **MISSING** |
| `Settings/Config.tsx` | Config display | None | **MISSING** |
| `Settings/Status.tsx` | Status display | None | **MISSING** |
| `Settings/Usage.tsx` | Usage display | None | **MISSING** |

### Tasks

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `tasks/BackgroundTask.tsx` | Background task display | None | **MISSING** |
| `tasks/BackgroundTasksDialog.tsx` | Background tasks dialog | None | **MISSING** |
| `tasks/BackgroundTaskStatus.tsx` | Task status display | None | **MISSING** |
| `tasks/ShellProgress.tsx` | Shell task progress | None | **MISSING** |
| `tasks/RemoteSessionProgress.tsx` | Remote session progress | None | **MISSING** |
| `tasks/taskStatusUtils.tsx` | Task status utilities | None | **MISSING** |
| All other task components (6) | Various task UIs | None | **MISSING** |
| `TaskListV2.tsx` | Todo/task list display | None | **MISSING** |

### Agents

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `agents/AgentsList.tsx` | Agents list UI | None | **MISSING** |
| `agents/AgentsMenu.tsx` | Agents menu | None | **MISSING** |
| `agents/AgentDetail.tsx` | Agent detail view | None | **MISSING** |
| `agents/AgentEditor.tsx` | Agent editor | None | **MISSING** |
| `agents/ColorPicker.tsx` | Color picker for agents | None | **MISSING** |
| `agents/ModelSelector.tsx` | Model selection | None | **MISSING** |
| `agents/new-agent-creation/` (11 files) | Agent creation wizard | None | **MISSING** |
| `agents/agentFileUtils.ts` | Agent file utilities | None | **MISSING** |
| `agents/generateAgent.ts` | Agent generation logic | None | **MISSING** |
| `agents/types.ts` | Agent types | None | **MISSING** |
| `agents/utils.ts` | Agent utilities | None | **MISSING** |
| `agents/validateAgent.ts` | Agent validation | None | **MISSING** |

### CustomSelect

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `CustomSelect/select.tsx` | Custom select component | None | **MISSING** |
| `CustomSelect/SelectMulti.tsx` | Multi-select | None | **MISSING** |
| `CustomSelect/use-select-state.ts` | Select state machine | None | **MISSING** |
| `CustomSelect/use-multi-select-state.ts` | Multi-select state | None | **MISSING** |
| `CustomSelect/use-select-navigation.ts` | Select navigation | None | **MISSING** |
| `CustomSelect/use-select-input.ts` | Select input handling | None | **MISSING** |
| `CustomSelect/select-option.tsx` | Select option rendering | None | **MISSING** |
| `CustomSelect/select-input-option.tsx` | Select input option | None | **MISSING** |
| `CustomSelect/option-map.ts` | Option mapping | None | **MISSING** |

### Other Notable Components

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `App.tsx` | Root application component | `pkg/tui/app.go: Model` | **PARTIAL** — Go has minimal model |
| `VirtualMessageList.tsx` | Virtualized message list | None | **MISSING** — critical for performance |
| `SearchBox.tsx` | Search interface | None | **MISSING** |
| `ModelPicker.tsx` | Model selection UI | None | **MISSING** |
| `TextInput.tsx` | Single-line text input | `pkg/tui/app.go: textarea.Model` | **PARTIAL** |
| `VimTextInput.tsx` | Vim-mode text input | None | **MISSING** |
| `BaseTextInput.tsx` | Base text input | None | **MISSING** |
| `Spinner.tsx` | Spinner (legacy) | `pkg/tui/app.go: spinner.Model` | **PARTIAL** |
| `StatusLine.tsx` | Status line display | None | **MISSING** |
| `Stats.tsx` | Stats display | `components.go: TokenCounter()` | **PARTIAL** |
| `ThemePicker.tsx` | Theme selection | None | **MISSING** |
| `DiagnosticsDisplay.tsx` | LSP diagnostics display | None | **MISSING** |
| `ContextVisualization.tsx` | Context window visualization | None | **MISSING** |
| `CompactSummary.tsx` | Compact summary display | None | **MISSING** |
| `Onboarding.tsx` | Onboarding flow | None | **MISSING** |
| `ExitFlow.tsx` | Exit confirmation flow | None | **MISSING** |
| `EffortIndicator.ts` | Effort level logic | None | **MISSING** |
| `shell/ShellProgressMessage.tsx` | Shell progress | None | **MISSING** |
| `shell/OutputLine.tsx` | Shell output line | None | **MISSING** |
| `wizard/` (5 files) | Wizard framework | None | **MISSING** |
| `ui/TreeSelect.tsx` | Tree selection | None | **MISSING** |
| `AutoUpdater.tsx` + related | Auto-update UI | None | **MISSING** |
| `FeedbackSurvey/` (7 files) | Feedback survey | None | **N/A** |
| `sandbox/` (5 files) | Sandbox configuration | None | **MISSING** |
| All Teleport components (5) | Teleport session UI | None | **N/A** — not in Go scope |
| All teams/swarm components (5+) | Team/swarm UI | None | **MISSING** |

### Section Summary
- **FULL:** 1 (`ProgressBar`)
- **PARTIAL:** ~10 (basic message rendering, spinner, textarea, token counter)
- **MISSING:** ~300+ files (nearly all component logic)
- **N/A:** ~30 files (JSX structure, theme providers, React contexts, voice, Teleport, Chrome)

---

## 4. old-src/screens/ (3 files) → pkg/tui/screens.go

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `REPL.tsx` | Main REPL screen: full chat interface, search, keybindings, task management, token budget display | `pkg/tui/screens.go: ScreenRouter` (type enum + navigation only) | **PARTIAL** — Go has screen routing but zero REPL logic (search, transcript, task navigation, token budget, etc.) |
| `Doctor.tsx` | Doctor diagnostic screen: settings errors, keybinding warnings, sandbox, LSP, plugin errors, update checks | `pkg/tui/screens.go: ScreenType "settings"` only | **MISSING** — no diagnostic implementation |
| `ResumeConversation.tsx` | Resume conversation picker: session list, search, preview | `pkg/commands/resume.go` (headless only) | **PARTIAL** — Go has resume logic but no interactive picker UI |

### Section Summary
- **PARTIAL:** 2 (`REPL.tsx`, `ResumeConversation.tsx`)
- **MISSING:** 1 (`Doctor.tsx`)

---

## 5. old-src/keybindings/ (14 files) → pkg/tui/keybindings.go

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `defaultBindings.ts` | Default keybinding map for all actions (100+ bindings, platform-aware) | `pkg/tui/keybindings.go: defaultBindings` (8 bindings) | **PARTIAL** — Go has a tiny subset |
| `loadUserBindings.ts` | Load + parse `~/.claude/keybindings.json` | `pkg/tui/keybindings.go: LoadKeyBindings()` | **FULL** |
| `match.ts` | Match key event against parsed binding (modifier-aware) | None | **MISSING** — Go has no key matching logic |
| `parser.ts` | Parse keystroke strings like `"ctrl+shift+k"` to `ParsedKeystroke` | None | **MISSING** — Go uses raw strings |
| `resolver.ts` | Resolve key event → action with chord support, context filtering | None | **MISSING** — critical for proper keybinding dispatch |
| `schema.ts` | Zod schema for keybindings.json validation; defines all valid contexts and actions | None | **MISSING** — Go has no validation |
| `shortcutFormat.ts` | Format binding for display (e.g., "⌃C" on macOS) | None | **MISSING** |
| `template.ts` | Generate documented template keybindings.json | None | **MISSING** |
| `validate.ts` | Validate user keybindings for conflicts, reserved shortcuts, invalid actions | None | **MISSING** |
| `reservedShortcuts.ts` | Defines non-rebindable shortcuts (Ctrl-C, Ctrl-D, Ctrl-M) | None | **MISSING** |
| `useKeybinding.ts` | React hook: subscribe to an action's key handler | None | **N/A** — React hook; Go uses `Update()` |
| `useShortcutDisplay.ts` | React hook: get display text for an action's binding | None | **N/A** — React hook |
| `KeybindingContext.tsx` | React context for active keybinding contexts | None | **N/A** — React context |
| `KeybindingProviderSetup.tsx` | Sets up keybinding provider in component tree | None | **N/A** — React component |

### Section Summary
- **FULL:** 1 (`loadUserBindings.ts` / `LoadKeyBindings()`)
- **PARTIAL:** 1 (`defaultBindings.ts`)
- **MISSING:** 7 (match, parser, resolver, schema, shortcutFormat, template, validate, reservedShortcuts)
- **N/A:** 4 (React hooks/contexts)

---

## 6. old-src/vim/ (5 files) → pkg/tui/vim.go

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `types.ts` | Complete vim state machine types: VimState, CommandState (idle/count/operator/find/g/replace), PersistentState, RecordedChange | `pkg/tui/vim.go: VimState` (Mode, Enabled, PendingOperator, Count, Register, CommandBuffer, LastSearch) | **PARTIAL** — Go lacks CommandState sub-states, PersistentState (undo/redo), RecordedChange (dot-repeat), FindType |
| `transitions.ts` | State transition table for all vim key sequences; handles all normal-mode commands | `pkg/tui/vim.go: handleNormal()` | **PARTIAL** — Go handles ~12 keys; TS handles 50+. Missing: `D`, `C`, `Y`, `p`, `P`, `r`, `R`, `~`, `g`/`G`, `%`, `f`/`F`/`t`/`T`, `n`/`N`, count prefix, visual mode transitions, `>>`, `<<`, `J`, `.` (dot-repeat) |
| `operators.ts` | Execute vim operators: delete, change, yank, paste, toggle case, indent, join, replace, find, open line | None | **MISSING** — Go has no operator execution logic |
| `motions.ts` | Resolve vim motions to cursor positions: h/j/k/l, w/b/e, W/B/E, 0/^/$/gg/G, inclusive/linewise flags | `pkg/tui/vim.go: handleNormal()` (basic h/j/k/l/w/b) | **PARTIAL** — Go missing: E/W/B, ^, gg/G, `%`, inclusive/linewise, count-multiplied motion |
| `textObjects.ts` | Find text object boundaries: `iw`/`aw` (word), `i"`/`a"` (quotes), `i(`/`a(` (parens), all bracket pairs | None | **MISSING** — text objects entirely absent |

### Section Summary
- **FULL:** 0
- **PARTIAL:** 3 (`types.ts`, `transitions.ts`, `motions.ts`)
- **MISSING:** 2 (`operators.ts`, `textObjects.ts`)

---

## 7. old-src/context/ (9 files) → pkg/tui/ or pkg/state/

React contexts in TS provide global state accessible throughout the component tree. In Go, this state lives in the bubbletea `Model` struct or `pkg/state/`.

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `mailbox.tsx` | `Mailbox` queue for inter-component messages | None | **MISSING** — needed for queue-based message passing |
| `notifications.tsx` | Notification queue with priority, timeout, fold, invalidation | None | **MISSING** — no notification system in Go |
| `modalContext.tsx` | Modal context with rows/columns/scrollRef | None | **MISSING** — needed for modal sizing |
| `overlayContext.tsx` | Overlay registry for Escape key coordination (autocomplete, selects) | None | **MISSING** — needed so Escape dismisses overlays before canceling |
| `stats.tsx` | In-memory stats store (counters, histograms, percentiles) for performance metrics | None | **MISSING** |
| `voice.tsx` | Voice state: recording/processing/idle, audio levels, transcript | None | **N/A** — voice not in Go scope |
| `fpsMetrics.tsx` | FPS tracker context | None | **N/A** — React render performance, not relevant to Go |
| `QueuedMessageContext.tsx` | Context for queued message display state (isQueued, isFirst, paddingWidth) | None | **N/A** — React layout context |
| `promptOverlayContext.tsx` | Portal for prompt overlay dialogs that float above clipped layout areas | None | **N/A** — React layout portal |

### Section Summary
- **MISSING:** 4 (`mailbox`, `notifications`, `modalContext`, `overlayContext`, `stats`)
- **N/A:** 5 (voice, fpsMetrics, QueuedMessageContext, promptOverlayContext — React-specific or descoped)

---

## 8. old-src/constants/ (21 files) → pkg/constants/

| TS File | Purpose | Go Equivalent | Rating |
|---------|---------|---------------|--------|
| `apiLimits.ts` | Image size limits (base64 max, raw target, PDF page limit, video/doc limits) | `pkg/constants/limits.go: MaxImageSizeBytes=20MB, MaxPDFPages=20` | **PARTIAL** — Go missing: `API_IMAGE_MAX_BASE64_SIZE` (5MB base64), `IMAGE_TARGET_RAW_SIZE` (3.75MB raw), document/video limits |
| `betas.ts` | Beta feature header strings (20+ headers for interleaved thinking, context-1m, structured outputs, etc.) | None | **MISSING** — all beta headers absent |
| `common.ts` | `getLocalISODate()`, `getSessionStartDate()` (memoized), `getLocalMonthYear()` | None | **MISSING** — date utilities for system prompt caching |
| `cyberRiskInstruction.ts` | `CYBER_RISK_INSTRUCTION` safety prompt constant | None | **MISSING** — but may not be needed if using hosted Claude |
| `errorIds.ts` | Numeric error IDs for production tracing (346 IDs) | None | **N/A** — Go uses error wrapping/types |
| `figures.ts` | Unicode glyphs: `BLACK_CIRCLE`, `LIGHTNING_BOLT`, `EFFORT_*`, `DIAMOND_*`, media icons, MCP indicators | None | **MISSING** — needed for consistent UI symbols |
| `files.ts` | `BINARY_EXTENSIONS` set for skipping binary files | None | **MISSING** — needed by file tools |
| `github-app.ts` | GitHub Actions workflow template content, PR title | None | **N/A** — GitHub integration not in Go scope |
| `keys.ts` | GrowthBook client keys | None | **N/A** — GrowthBook not in Go |
| `messages.ts` | `NO_CONTENT_MESSAGE = '(no content)'` | `pkg/constants/constants.go: NoContentMessage` | **FULL** |
| `oauth.ts` | OAuth config with staging/local/prod variants, `fileSuffixForOauthConfig()` | `pkg/constants/oauth.go` (prod only) | **PARTIAL** — Go missing staging/local variants and `fileSuffixForOauthConfig()` |
| `outputStyles.ts` | `OutputStyleConfig`, `OutputStyles` type, built-in output style configs (default, explanatory, learning) | None | **MISSING** |
| `product.ts` | `PRODUCT_URL`, `CLAUDE_AI_BASE_URL`, staging/local URL variants, `isRemoteSessionStaging()` | `pkg/constants/product.go` (basic constants only) | **PARTIAL** — Go missing `CLAUDE_AI_BASE_URL`, staging/local URLs, `isRemoteSessionStaging()`, `isRemoteSessionLocal()` |
| `prompts.ts` | System prompt prefix constants (`DEFAULT_PREFIX`, `AGENT_SDK_PREFIX`), `getCLISyspromptPrefix()` | None | **MISSING** — critical for system prompt construction |
| `spinnerVerbs.ts` | `SPINNER_VERBS` list, `getSpinnerVerbs()` with settings override | None | **MISSING** |
| `system.ts` | System prompt section framework (`systemPromptSection()`, volatile sections) | None | **MISSING** — system prompt composition |
| `systemPromptSections.ts` | Registers all system prompt sections with their compute functions | None | **MISSING** |
| `toolLimits.ts` | `DEFAULT_MAX_RESULT_SIZE_CHARS=50000`, `MAX_TOOL_RESULT_TOKENS=100000`, `BYTES_PER_TOKEN=4` | `pkg/constants/limits.go: MaxToolResultChars=100000` | **PARTIAL** — Go has `MaxToolResultChars` (100k) but missing `DEFAULT_MAX_RESULT_SIZE_CHARS` (50k) and `BYTES_PER_TOKEN` |
| `tools.ts` | All tool name constants and tool feature flags | None | **MISSING** — tool names are inline strings in Go |
| `turnCompletionVerbs.ts` | `TURN_COMPLETION_VERBS` list | None | **MISSING** |
| `xml.ts` | All XML tag name constants (command tags, bash tags, task tags, tick, fork glyph) | `pkg/constants/constants.go` (partial: command/system-reminder tags only) | **PARTIAL** — Go missing: bash I/O tags, task notification tags, tick/fork tags |

### Section Summary
- **FULL:** 1 (`messages.ts`)
- **PARTIAL:** 6 (`apiLimits`, `oauth`, `product`, `toolLimits`, `xml`, `models` implied)
- **MISSING:** 10 (`betas`, `common`, `figures`, `files`, `outputStyles`, `prompts`, `spinnerVerbs`, `system`, `systemPromptSections`, `tools`, `turnCompletionVerbs`)
- **N/A:** 4 (`errorIds`, `github-app`, `keys`, `cyberRiskInstruction`)

---

## Overall Summary

| Directory | Files | FULL | PARTIAL | MISSING | N/A |
|-----------|-------|------|---------|---------|-----|
| hooks/ | 104 | 2 | 5 | ~60 | ~17 |
| ink/ | 96 | 0 | 0 | ~25 | ~70 |
| components/ | 389 | 1 | ~10 | ~300 | ~30 |
| screens/ | 3 | 0 | 2 | 1 | 0 |
| keybindings/ | 14 | 1 | 1 | 7 | 4 |
| vim/ | 5 | 0 | 3 | 2 | 0 |
| context/ | 9 | 0 | 0 | 4 | 5 |
| constants/ | 21 | 1 | 6 | 10 | 4 |
| **Total** | **641** | **5** | **~27** | **~409** | **~130** |

---

## Critical Missing Logic (Highest Priority to Port)

These are not React-specific — they contain pure logic that must exist in Go:

1. **Keybinding resolver** (`keybindings/resolver.ts`, `match.ts`, `parser.ts`) — chord-based key dispatch with context filtering
2. **Vim operators + text objects** (`vim/operators.ts`, `vim/textObjects.ts`) — delete/change/yank/paste, `iw`/`aw`/`i"`/`a(` etc.
3. **Tool permission handlers** (`hooks/useCanUseTool.tsx`, `toolPermission/handlers/`) — the interactive/coordinator/worker permission decision logic
4. **Text input state machine** (`hooks/useTextInput.ts`) — kill ring, yank, cursor ops, image paste
5. **Beta header constants** (`constants/betas.ts`) — required for API calls to use newer features
6. **System prompt framework** (`constants/system.ts`, `systemPromptSections.ts`, `constants/prompts.ts`) — how system prompts are constructed
7. **Tool name constants** (`constants/tools.ts`) — scattered inline strings should be centralized
8. **Unicode figures** (`constants/figures.ts`) — UI symbols used throughout
9. **Notifications system** (`context/notifications.tsx`) — priority queue with fold/invalidate
10. **Overlay/escape coordination** (`context/overlayContext.tsx`) — prevents Escape from canceling when autocomplete is open

## What the Go Port Has

The current Go TUI is a basic proof-of-concept:
- `pkg/tui/app.go`: Single bubbletea Model with textarea + spinner + message list (minimal)
- `pkg/tui/vim.go`: Basic vim normal mode (~12 keys, no operators)
- `pkg/tui/keybindings.go`: Load/save keybindings from JSON (no matching/resolution)
- `pkg/tui/screens.go`: Screen type enum + router (no screen implementations)
- `pkg/tui/components/components.go`: 5 render functions (permission dialog, progress bar, token counter, spinner frames, message bubble)
- `pkg/hooks/engine.go`: Full hook execution engine (external command runner)
- `pkg/hooks/ui.go`: Terminal size, blink state, elapsed timer
- `pkg/constants/`: Partial constants (limits, product, oauth, XML tags, models)
