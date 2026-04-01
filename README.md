# ◉ Oculus

**Go-native AI coding CLI with multi-model lens architecture.**

Built by [HowlerOps](https://github.com/howlerops) — achieving feature parity with leading AI coding tools while delivering Go's memory efficiency and startup speed.

[![CI](https://github.com/howlerops/oculus/actions/workflows/ci.yml/badge.svg)](https://github.com/howlerops/oculus/actions/workflows/ci.yml)
[![Release](https://github.com/howlerops/oculus/releases/latest/badge.svg)](https://github.com/howlerops/oculus/releases/latest)

## Install

```bash
# Go (recommended)
go install github.com/howlerops/oculus/cmd/oculus@latest

# npm
npm install -g @howlerops/oculus

# Binary (macOS/Linux/Windows)
# Download from https://github.com/howlerops/oculus/releases
```

## Quick Start

```bash
# Interactive REPL (launches TUI)
oculus

# Single prompt (non-interactive)
oculus -p "Explain this codebase"

# With specific model
oculus -m opus -p "Review this PR"
```

On first run, Oculus will open your browser for OAuth authentication with Anthropic. Alternatively, set `ANTHROPIC_API_KEY` directly.

## Features

### 🔭 Three-Lens Architecture
Oculus routes work through specialized lenses:
- **Focus** — Reasoning, planning, orchestration
- **Scan** — File search, codebase exploration, research
- **Craft** — Code writing, editing, command execution

Each lens can target a different model and provider.

### 🌐 6 Provider Bridges
- **Anthropic** — Claude Opus, Sonnet, Haiku
- **OpenAI** — GPT-4, GPT-4o (any OpenAI-compatible endpoint)
- **Ollama** — Local models (Llama, Mistral, CodeLlama)
- **Claude CLI** — Use your Claude Pro/Max subscription
- **Codex CLI** — Use your OpenAI subscription
- **Gemini CLI** — Use your Google subscription

### 🔄 Orchestration Engine
- **Ralph Loop** (`/ralph` or `--ralph`) — PRD-driven persistence until every story passes
- **Consensus Planning** (`/plan` or `--plan`) — Planner → Architect → Critic review loop
- **Ultrawork** — Parallel dispatch with dependency DAG and model tier routing
- **5 Agent Personas** — Architect, Critic, Executor, Explorer, Planner

### 🚀 Multi-Provider Onboarding
Auto-detects installed CLIs, API keys, and local models. Works without Anthropic — any provider can be primary.

### 🧠 Context Management
Episode-based compaction with LCM dual-threshold engine:
- Below 70%: Zero overhead
- 70-90%: Async compaction between turns
- Above 90%: Blocking compaction with TF-IDF keyword extraction

### 🛠 40 Built-in Tools
Bash, File Read/Write/Edit, Glob, Grep, Agent spawning, WebSearch, WebFetch, Notebook editing, MCP integration, and 30+ more.

### 💬 54 Slash Commands
`/model`, `/diff`, `/commit`, `/compact`, `/permissions`, `/resume`, `/agents`, `/skills`, `/mcp`, and more.

### 🖥 Full Terminal UI
- Permission dialog (approve/deny/always-allow)
- Markdown rendering with syntax highlighting
- Tool call badges with per-tool progress spinners
- Multi-line input with Ctrl+R history search
- Scrollable viewport, task panel, status bar
- Vim mode

### 🔒 Security
- Bash command safety analysis (24 dangerous patterns)
- Full permission system with glob-matching rules
- OAuth PKCE authentication with keychain storage
- SSRF guard on HTTP hooks

## Configuration

Config directory: `~/.oculus/` (auto-migrates from `~/.claude/`)

### OCULUS.md
Create project-specific instructions:

```markdown
# Project Instructions

- Use pnpm, not npm
- Run tests with `go test ./...`
- Follow conventional commits
```

### settings.json
```json
{
  "model": "claude-sonnet-4-20250514",
  "defaultMode": "default",
  "lenses": {
    "focus": { "model": "claude-opus-4-20250514" },
    "scan": { "model": "claude-sonnet-4-20250514" },
    "craft": { "model": "claude-sonnet-4-20250514", "provider": "ollama" }
  }
}
```

## Architecture

```
cmd/oculus/     CLI entry point
pkg/
├── api/        Anthropic streaming client
├── auth/       OAuth PKCE + keychain
├── bridge/     Multi-provider abstraction
├── lens/       Focus/Scan/Craft routing
├── query/      Conversation loop + tool dispatch
├── tools/      40 tool implementations
├── tui/        Bubbletea terminal UI (16 files)
├── services/   Episodes, MCP, analytics, compact
└── utils/      Permissions, git, bash parsing
```

## Documentation

Full docs: [howlerops.github.io/oculus](https://howlerops.github.io/oculus)

## License

MIT © [HowlerOps](https://github.com/howlerops)
