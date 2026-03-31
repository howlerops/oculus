# TypeScript → Go Port Gap Analysis

## Summary
- **TS Source**: 1,884 files across 35 modules + 18 top-level files
- **Go Port**: 54 files across 40 packages
- **Coverage**: Core architecture ported, many advanced/niche modules not yet ported

## Module-by-Module Status

### FULLY PORTED (core functionality present)
| TS Module | Go Package | TS Files | Go Files | Status |
|-----------|-----------|----------|----------|--------|
| types/ | pkg/types/ | 11 | 5 | FULL - message, permissions, ids, tools, hooks |
| Tool.ts | pkg/tool/ | 1 | 1 | FULL - Tool interface + BaseTool |
| Task.ts | pkg/task/ | 1 | 2 | FULL - with tests |
| context.ts | pkg/context/ | 1+9 | 4 | FULL - git, claudemd, system, user |
| query.ts + QueryEngine.ts | pkg/query/ | 2 | 1 | FULL - streaming loop + tool dispatch |
| tools.ts | cmd/claude-go/ | 1 | 1 | FULL - tool registry in main |
| commands.ts | pkg/commands/ | 1 | 2 | FULL - registry + 12 builtins |
| state/ | pkg/state/ | 6 | 4 | FULL - AppState, Store, selectors, tests |
| main.tsx | cmd/claude-go/ | 1 | 2 | FULL - cobra CLI + integration tests |
| cost-tracker.ts | pkg/services/cost/ | 1 | 1 | FULL |
| history.ts | pkg/services/history/ | 1 | 1 | FULL |

### PARTIALLY PORTED (core present, advanced features missing)
| TS Module | Go Package | TS Files | Go Files | Missing |
|-----------|-----------|----------|----------|---------|
| tools/ | pkg/tools/ | 184 | 19 pkgs | 12 tool dirs missing (BriefTool, ConfigTool, MCPTool, etc.) |
| services/ | pkg/services/ | 130 | 5 | analytics, api details, compact advanced, diagnostics, mcp advanced |
| config | pkg/config/ | 4+ | 4 | settings validation, mdm, managed settings |
| ink/ | pkg/tui/ | 96 | 1 | Full Ink/React replacement via bubbletea (1 file covers core) |

### NOT YET PORTED (listed by criticality)

#### HIGH PRIORITY (needed for full functionality)
| TS Module | Files | What It Does | Go Equivalent Needed |
|-----------|-------|-------------|---------------------|
| hooks/ | 104 | React hooks for UI state | pkg/hooks/ (Go channels/callbacks) |
| utils/permissions/ | 24 | Full permission system | pkg/utils/permissions/ (has basic, needs full) |
| utils/messages/ | 3+ | Message normalization | pkg/utils/messages/ (has basic, needs full) |
| utils/hooks/ | 18 | Hook execution engine | pkg/hooks/engine.go |
| constants/ | 21 | Product constants, limits | pkg/constants/ |
| entrypoints/ | 8 | SDK, print, headless modes | pkg/entrypoints/ |

#### MEDIUM PRIORITY (advanced features)
| TS Module | Files | What It Does |
|-----------|-------|-------------|
| utils/model/ | 16 | Model selection, capabilities, aliases |
| utils/settings/ | 14 | Settings validation, schema |
| utils/bash/ | 15 | Bash parsing, safety checking |
| utils/plugins/ | 45 | Plugin system |
| utils/swarm/ | 17 | Team/swarm orchestration |
| bridge/ | 31 | Remote session bridge |
| cli/ | 19 | CLI argument parsing details |
| components/ | 389 | UI components (covered by bubbletea) |
| screens/ | 3 | Screen routing |

#### LOW PRIORITY (niche/platform-specific)
| TS Module | Files | What It Does |
|-----------|-------|-------------|
| vim/ | 5 | Vim keybindings |
| voice/ | 1 | Voice input |
| buddy/ | 6 | Companion mode |
| assistant/ | 1 | Kairos assistant |
| migrations/ | 11 | Data migrations |
| upstreamproxy/ | 2 | Proxy support |
| native-ts/ | 4 | Native TS bindings |
| moreright/ | 1 | Advanced permissions |
| coordinator/ | 1 | Coordinator mode |
| bootstrap/ | 1 | Bootstrap state |
| schemas/ | 1 | Schema definitions |
| outputStyles/ | 1 | Output formatting |
| memdir/ | 8 | Memory directory |
| keybindings/ | 14 | Key binding system |
| remote/ | 4 | Remote execution |
| server/ | 3 | Server mode |
| skills/ | 20 | Built-in skills |
| plugins/ | 2 | Plugin types |

## Missing Tools (12 of 40)
BriefTool, ConfigTool, ListMcpResourcesTool, McpAuthTool, MCPTool, PowerShellTool, ReadMcpResourceTool, RemoteTriggerTool, REPLTool, ScheduleCronTool, SleepTool, SyntheticOutputTool

Most are ant-only (REPLTool, BriefTool, SyntheticOutputTool), platform-specific (PowerShellTool), or advanced (ScheduleCronTool, RemoteTriggerTool).

## Critical Path for Working Binary
The binary already has: CLI → API Client → Streaming → Query Loop → Tool Dispatch → 6 Core Tools
What's needed: verify it works end-to-end with a real API call.
