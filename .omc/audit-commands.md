# commands/ Audit: TypeScript -> Go Port

Audited: 2026-03-31
Source: `old-src/commands/` (TS) vs `pkg/commands/` (Go)

## Go Port Summary

The Go port has **4 files** in `pkg/commands/`:
- `command.go` — Registry, Command struct, helpers (FULL)
- `builtins.go` — Registers 13 built-in commands (stubs/placeholders)
- `login.go` — Registers `login` + `logout` (calls real `auth` package)
- `resume.go` — Registers `resume` + `session` (calls real `sessions` package)

---

### commands/ (189 TS files -> pkg/commands/)

| TS Command | Description | Go Coverage | Notes / Missing |
|---|---|---|---|
| `/add-dir` | Add a new working directory | MISSING | No Go equivalent; adds dirs to context scope |
| `/agents` | Manage agent configurations | MISSING | Agent config management not ported |
| `/branch` | Create a branch of the current conversation | MISSING | Conversation branching not ported |
| `/bridge` (`/remote-control`) | Connect terminal for remote-control sessions | MISSING | Bridge/daemon mode not ported |
| `/brief` | Toggle brief-only mode | MISSING | Feature-flagged in TS (`KAIROS`); not ported |
| `/btw` | Background task annotation | MISSING | btw command not ported |
| `/chrome` | Claude in Chrome (Beta) settings | MISSING | Browser integration not ported |
| `/clear` | Clear conversation history and free up context | REGISTERED | Go stub returns static string; TS clears message history + caches |
| `/color` | Set the prompt bar color for this session | MISSING | TUI color theming not ported |
| `/commit` | Create a git commit | MISSING | Git commit workflow not ported |
| `/commit-push-pr` | Commit, push, and open a PR | MISSING | Full git workflow not ported |
| `/compact` | Summarize and compact the conversation | REGISTERED | Go stub returns static string; TS triggers actual LLM summarization |
| `/config` | Show or edit configuration | REGISTERED | Go stub points to settings.json; TS opens interactive config panel |
| `/context` (+ noninteractive) | Visualize/show current context usage | MISSING | Context visualization not ported |
| `/copy` | Copy last response to clipboard | MISSING | Clipboard integration not ported |
| `/cost` | Show token usage and cost | REGISTERED | Go stub returns placeholder; TS reads live cost tracking state |
| `/desktop` | Continue session in Claude Desktop | MISSING | Desktop handoff not ported |
| `/diff` | View uncommitted changes and per-turn diffs | MISSING | Git diff integration not ported |
| `/doctor` | Diagnose Claude Code installation | REGISTERED | Go has basic checks (API key, git, ripgrep); TS checks many more (node, permissions, config validity, etc.) |
| `/effort` | Set effort level for model usage | MISSING | Model effort/budget not ported |
| `/exit` | Exit the REPL | REGISTERED | Go registers as alias of `quit`; TS is its own command |
| `/export` | Export conversation to file or clipboard | MISSING | Export not ported |
| `/extra-usage` | Configure extra usage for rate limits | MISSING | Rate-limit handling not ported |
| `/fast` | Toggle fast/lightweight model mode | MISSING | Model switching shortcut not ported |
| `/feedback` | Submit feedback about Claude Code | MISSING | Feedback submission not ported |
| `/files` | List all files currently in context | MISSING | Context file listing not ported |
| `/heapdump` | Dump the JS heap to ~/Desktop | MISSING | JS-specific; N/A for Go (no direct equivalent needed) |
| `/help` | Show available commands | FULL | Go registry formats all registered commands |
| `/hooks` | View hook configurations for tool events | MISSING | Hook config viewer not ported |
| `/ide` | Manage IDE integrations | MISSING | IDE integration not ported |
| `/init` | Initialize CLAUDE.md for this project | REGISTERED | Go creates basic CLAUDE.md; TS has richer template + skill scaffolding |
| `/init-verifiers` | Initialize verifier agents | MISSING | Verifier agent setup not ported |
| `/insights` | Session insights / backfill | MISSING | Analytics/insights not ported |
| `/install-github-app` | Set up Claude GitHub Actions | MISSING | GitHub Actions setup not ported |
| `/install-slack-app` | Install the Claude Slack app | MISSING | Slack integration not ported |
| `/keybindings` | Open/create keybindings config | MISSING | Keybindings editor not ported |
| `/login` | Sign in with Anthropic account | PARTIAL | Go calls `auth.GetAuthToken`; TS handles OAuth browser flow, API key auth, multiple auth modes |
| `/logout` | Sign out and clear credentials | PARTIAL | Go calls `auth.Logout`; TS has multi-auth-mode logout + credential cleanup |
| `/mcp` | Manage MCP servers | MISSING | MCP server management not ported |
| `/memory` | Edit Claude memory files | REGISTERED | Go stub redirects to /context; TS opens interactive CLAUDE.md editor |
| `/mobile` | Show QR code for Claude mobile app | MISSING | QR code display not ported |
| `/model` | Switch/show model | MISSING | Model selection not ported |
| `/onboarding` | Onboarding flow | MISSING | Onboarding not ported |
| `/output-style` | Deprecated: use /config | MISSING | Deprecated TS shim; not worth porting |
| `/passes` | Configure model passes | MISSING | Multi-pass reasoning not ported |
| `/permissions` | Manage tool permission rules | MISSING | Permission allow/deny rules not ported |
| `/plan` | Enable plan mode / view plan | MISSING | Plan mode not ported |
| `/pr-comments` | Get comments from a GitHub PR | MISSING | GitHub PR integration not ported |
| `/privacy-settings` | View/update privacy settings | MISSING | Privacy settings not ported |
| `/rate-limit-options` | Show options when rate limited | MISSING | Rate limit handling not ported |
| `/release-notes` | View release notes | MISSING | Release notes not ported |
| `/reload-plugins` | Activate pending plugin changes | MISSING | Plugin system not ported |
| `/remote-env` | Configure default remote environment | MISSING | Remote env config not ported |
| `/remote-control` (bridge) | Remote control server | MISSING | Bridge/daemon mode not ported |
| `/rename` | Rename the current conversation | MISSING | Session rename not ported |
| `/resume` | Resume a previous conversation | PARTIAL | Go loads session by ID and lists recent; TS has interactive fuzzy-picker UI |
| `/review` | Review a pull request | MISSING | PR review workflow not ported |
| `/review` (ultrareview) | Deep bug-finding review (web) | MISSING | Ultra-review not ported |
| `/rewind` | Restore code/conversation to previous point | MISSING | Checkpoint/rewind not ported |
| `/sandbox` | Toggle sandbox mode | MISSING | Sandbox toggling not ported |
| `/security-review` | Security review of pending changes | MISSING | Security review workflow not ported |
| `/session` | Show current session info | REGISTERED | Go returns static string; TS shows full session metadata + actions |
| `/share` | Show remote session URL / QR code | MISSING | Remote sharing not ported |
| `/skills` | List available skills | MISSING | Skills system not ported |
| `/stats` | Show Claude Code usage statistics | MISSING | Usage stats/analytics not ported |
| `/status` | Show session status | REGISTERED | Go shows runtime/memory stats; TS shows API status, model, context usage |
| `/stickers` | Order Claude Code stickers | MISSING | Stickers URL launcher not ported |
| `/tag` | Toggle a tag on the current session | MISSING | Session tagging not ported |
| `/tasks` | List and manage background tasks | MISSING | Background task management not ported |
| `/teleport` | Teleport session to another machine | MISSING | Teleport functionality not ported |
| `/terminal-setup` | Set up terminal integrations | MISSING | Terminal setup not ported |
| `/theme` | Change color theme | REGISTERED | Go sets theme name only; TS applies theme to TUI rendering |
| `/think-back` / `/thinkback-play` | Thinkback animation | MISSING | Easter egg/animation not ported |
| `/upgrade` | Upgrade to Max plan | MISSING | Upgrade flow not ported |
| `/usage` | Show plan usage limits | MISSING | Usage/quota display not ported |
| `/version` | Show version information | REGISTERED | Go returns hardcoded string; TS reads build metadata |
| `/vim` | Toggle vim keybindings | REGISTERED | Go returns stub message; TS actually toggles vim mode in TUI |
| `/voice` | Toggle voice mode | MISSING | Voice mode not ported |
| `/web-setup` | Web-based remote setup | MISSING | Feature-flagged (`CCR_REMOTE_SETUP`); not ported |

---

## Coverage Summary

| Status | Count | Commands |
|---|---|---|
| FULL | 1 | `/help` |
| PARTIAL | 3 | `/login`, `/logout`, `/resume` |
| REGISTERED (stub only) | 11 | `/clear`, `/compact`, `/config`, `/cost`, `/doctor`, `/exit`, `/init`, `/memory`, `/session`, `/status`, `/theme`, `/version`, `/vim` |
| MISSING | 56 | All others |

> Note: `REGISTERED` means the command name exists in the Go registry with a description but the implementation is a static string or placeholder — the real TS logic (LLM calls, UI, file I/O, API calls) has not been ported.

---

## Top-Priority Gaps (Core UX)

These are commands used in normal daily operation that are completely absent from the Go port:

1. `/model` — Users need to switch models
2. `/diff` — View changes per turn
3. `/commit` — Git commit workflow
4. `/mcp` — MCP server management (critical for tool use)
5. `/permissions` — Tool allow/deny rules
6. `/plan` — Plan mode
7. `/context` — Context usage visibility
8. `/files` — See what's in context
9. `/tasks` — Background task management
10. `/rewind` — Undo/checkpoint recovery

## Commands Not Worth Porting (JS-Specific or Deprecated)

- `/heapdump` — Node.js heap dump; N/A in Go
- `/output-style` — Deprecated in TS, redirects to `/config`
- `/stickers` — Opens a URL; trivial if needed
- `/thinkback` / `/thinkback-play` — Easter egg animations
- `/btw` — Internal annotation command
- `/backfill-sessions` / `/break-cache` — Internal/maintenance commands
- `/insights` — Internal analytics backfill

## Feature-Flagged TS Commands (Low Priority)

These only activate under feature flags and can be deferred:
- `/brief` (`KAIROS` flag)
- `/voice` (`VOICE_MODE` flag)
- `/bridge` / `/remote-control` (`BRIDGE_MODE` + `DAEMON` flags)
- `/web-setup` (`CCR_REMOTE_SETUP` flag)
- `/workflows` (`WORKFLOW_SCRIPTS` flag)
