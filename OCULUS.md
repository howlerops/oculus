# Oculus Project Instructions

## Language & Runtime
- Go 1.25+ required
- Module: github.com/howlerops/oculus
- Build: `go build ./cmd/oculus/`
- Test: `go test ./...`

## Code Style
- Run `gofmt` before committing
- Run `go vet ./...` — zero warnings allowed
- Prefer `Edit` over `Write` for existing files
- No `TODO` comments in committed code — use issues instead

## Architecture
- `pkg/` for library code, `cmd/` for binaries
- Each tool in its own package under `pkg/tools/`
- TUI split into focused files (model, init, update, view, messages, handlers)
- Orchestration (ralph, plan, ultrawork) in `pkg/orchestration/`
- Lens system routes through Focus/Scan/Craft

## Testing
- All new packages must have `_test.go` files
- Target 80%+ coverage on critical paths (query, lens, bridge)
- Use `t.TempDir()` for file-based tests
- Mock external services, don't hit real APIs in tests

## Git
- Conventional commits: feat/fix/refactor/test/docs/chore
- Don't force push to main
- Run `go build ./... && go test ./... && go vet ./...` before committing

## Config
- User config: `~/.oculus/settings.json`
- Project config: `OCULUS.md` (this file) or `.oculus/settings.json`
- Legacy `~/.claude/` and `CLAUDE.md` supported as fallback
