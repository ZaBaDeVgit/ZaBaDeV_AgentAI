# Supported Agents

← [Back to README](../README.md)

---

| Agent | ID | Skills | MCP | Sub-agents | Output Styles | Slash Commands | Config Path |
|-------|-----|--------|-----|------------|---------------|----------------|-------------|
| Claude Code | `claude-code` | Yes | Yes | Yes | Yes | No | `~/.claude` |
| OpenCode | `opencode` | Yes | Yes | Yes | No | Yes | `~/.config/opencode` |
| Gemini CLI | `gemini-cli` | Yes | Yes | Yes (experimental) | No | No | `~/.gemini` |
| Cursor | `cursor` | Yes | Yes | Yes | No | No | `~/.cursor` |
| VS Code Copilot | `vscode-copilot` | Yes | Yes | Yes | No | No | `~/.copilot` + VS Code User profile |

Agent capabilities differ by product and version. ZaBaDeV writes the available instructions, skills and configuration for each adapter; it cannot guarantee that a third-party agent loaded them or that delegation is enabled. For OpenCode, multi-agent delegation depends on the installed background-agents plugin and its runtime dependencies. Verify behavior in the actual agent after installation.

## Notes

- **OpenCode** now exposes a single public named agent in `~/.config/opencode/opencode.json`: `senior-zabadev`. In multi-agent SDD mode, hidden `sdd-*` sub-agents are still installed for orchestration, but the old public `gentleman` and `sdd-orchestrator` tabs are no longer written.
- **Gemini CLI** sub-agents are experimental and require `experimental.enableAgents: true` in `settings.json`. Custom sub-agents are defined as markdown files in `~/.gemini/agents/`.
- **Cursor** supports async sub-agents (v2.5+) that can run in background and spawn nested sub-agent trees.
- **VS Code Copilot** uses the `runSubagent` tool with support for parallel execution and custom agent definitions.
- **Output Styles** are currently a Claude Code exclusive feature (`~/.claude/output-styles/`).
- **Slash Commands** are currently supported by OpenCode only.
- **VS Code Copilot** stores skills under `~/.copilot/skills/` (global), system prompt under `Code/User/prompts/zabadev.instructions.md`, and MCP config under `Code/User/mcp.json`.
- **Codex** currently gets shared instructions via `~/.codex/agents.md`, but this repo does not yet manage tab-switchable named agent profiles for Codex the way it does for OpenCode.
