# Engram Protocol

This document contains the complete Engram memory system integration protocol. All agents MUST follow this protocol to ensure consistent, hallucination-free operation with persistent memory across sessions.

---

## 1. Anti-Hallucination Mandate

**CRITICAL: ZERO tolerance for hallucinations. Memory is non-negotiable.**

Before ANY statement about code, APIs, or project behavior, you MUST:
1. **CHECK MEMORY FIRST**: Call `mem_search(query: "relevant-keywords", project: "{project}")`
2. **VERIFY WITH CODE**: Read actual files when discussing implementation details
3. **IF UNSURE**: Say "Necesito verificar esto" and investigate before proceeding

**The cost of asking is always lower than the cost of hallucinating.**

---

## 2. Memory Operations

Engram provides the following memory operations:

### mem_save
Saves structured observations to persistent storage.

**When to call:**
- Bug fix completed
- Architecture or design decision made
- Non-obvious discovery about the codebase
- Configuration change or environment setup
- Pattern established (naming, structure, convention)
- User preference or constraint learned

**Format:**
```python
mem_save(
  title: "Verb + what — short, searchable (e.g., 'Fixed N+1 query in UserList')",
  type: "bugfix | decision | architecture | discovery | pattern | config | preference",
  scope: "project" | "personal",  # default is project
  topic_key: "Stable key for evolving topics (e.g., 'architecture/auth-model')",
  content:
    What: One sentence — what was done
    Why: What motivated it (user request, bug, performance, etc.)
    Where: Files or paths affected
    Learned: Gotchas, edge cases, things that surprised you (omit if none)
)
```

### mem_search
Full-text search across all observations using FTS5.

**Returns:** 300-character PREVIEWS with observation IDs.

**Always follow with mem_get_observation** to get full content.

### mem_get_observation
Retrieves full content of a specific observation by ID.

**Required after mem_search** — previews are insufficient for decision-making.

### mem_update
Updates an existing observation when you have the exact observation ID.

### mem_context
Checks recent session history — fast and cheap, used at session start.

### mem_session_summary
Structured session end summary (REQUIRED before ending any session).

---

## 3. Two-Step Recovery Protocol

**CRITICAL: mem_search returns 300-char PREVIEWS. ALWAYS call mem_get_observation for full content.**

```
Step 1: mem_search(query: "{topic_key}", project: "{project}") → observation ID
Step 2: mem_get_observation(id: {id}) → full content
```

**Anti-pattern:** Making decisions based on preview snippets without retrieving full context.

---

## 4. Topic Key Format

Use stable, consistent topic keys for evolving topics. Different topics must NOT overwrite each other.

### SDD Topic Keys (Spec-Driven Development)

| Artifact | Topic Key |
|----------|-----------|
| Project context | `sdd-init/{project}` |
| Exploration | `sdd/{change}/explore` |
| Proposal | `sdd/{change}/proposal` |
| Spec | `sdd/{change}/spec` |
| Design | `sdd/{change}/design` |
| Tasks | `sdd/{change}/tasks` |
| Apply progress | `sdd/{change}/apply-progress` |
| Verify report | `sdd/{change}/verify-report` |
| Archive report | `sdd/{change}/archive-report` |

### Architecture Topic Keys

- `architecture/auth-model` — Authentication/authorization decisions
- `architecture/database-schema` — Schema changes
- `architecture/api-contract` — API design decisions
- `architecture/frontend-state` — State management patterns
- `architecture/deployment` — Infrastructure decisions

**Rule:** Reuse the same `topic_key` to update an evolving topic instead of creating new observations. If unsure about the key, call `mem_suggest_topic_key` first.

---

## 5. Session Start Protocol

At session start (or after compaction), you MUST:

1. **mem_context** — Check recent session history (fast, cheap)
2. **mem_search** — Search relevant topics for current work
3. **mem_get_observation** — Retrieve full content for any relevant hits
4. **Load project context** — If working on a specific project, load its context

```python
# Session start sequence
mem_context()
mem_search(query: "{keywords from user's message}", project: "{project}")
# For each relevant hit:
mem_get_observation(id: {id})
# Continue with verified information only
```

---

## 6. Session End Protocol (MANDATORY)

Before ending or saying "done" / "listo" / "finished", call `mem_session_summary`:

```markdown
## Goal
[What we were working on this session]

## Instructions
[User preferences or constraints discovered — skip if none]

## Discoveries
- [Technical findings, gotchas, non-obvious learnings]

## Accomplished
- [Completed items with key details]

## Next Steps
- [What remains to be done — for the next session]

## Relevant Files
- path/to/file — [what it does or what changed]
```

**This is NOT optional.** If skipped, the next session starts blind.

---

## 7. Sub-Agent Context Rules

### SDD Phase Agents
- **Who reads:** Main orchestrator reads all memory
- **Who writes:** Phase-specific agents (sdd-spec, sdd-design, sdd-apply, etc.)
- **Pattern:** Orchestrator delegates to sub-agent → sub-agent writes to memory → returns result

### Non-SDD Tasks
- **Who reads:** Any agent can read memory
- **Who writes:** The agent performing the work
- **Pattern:** Each agent saves its own observations for continuity

### Hard Stop Rule

Before Read/Edit/Write/Grep on source files:
1. **STOP** — "Is this orchestration or execution?"
2. If execution → **delegate to sub-agent**
3. **ONLY files you read directly:** git status/log, engram results, todo state

---

## 8. Memory Anti-Patterns (NEVER do)

- ❌ "Based on my knowledge..." → ALWAYS verify with memory/code
- ❌ "I believe the API is..." → ALWAYS check first
- ❌ "This pattern was used before..." → ALWAYS verify with `mem_search`
- ❌ Inventing support, APIs, or behaviors → VERIFY FIRST
- ❌ Making decisions based on mem_search previews without calling mem_get_observation
- ❌ Skipping session summary before ending

---

## 9. DAG State Persistence

The SDD workflow is a DAG (Directed Acyclic Graph) where:
- Each phase produces artifacts (proposal → specs → design → tasks → apply → verify → archive)
- Each artifact is stored in Engram with a stable topic key
- Compaction destroys conversation context but NOT Engram observations

### Recovery After Compaction

```
# Step 1: Recover what was done before compaction
mem_session_summary()  # Called by the main agent after compaction notice

# Step 2: Restore SDD state by searching for the change
mem_search(query: "sdd/{change-name}", project: "{project}")
# For each artifact found:
mem_get_observation(id: {id})

# Step 3: Continue from where left off
```

**This ensures SDD state survives context compaction.**

---

## 10. MCP Configuration

Engram connects as an MCP (Model Context Protocol) server to provide memory tools to agents.

### Connection Details

- **Server:** Engram MCP server
- **Endpoint:** localhost:2112
- **Command:** `engram mcp --tools=agent`

### MCP Config Examples

**Claude Code (separate-file strategy):**
```json
{
  "command": "engram",
  "args": ["mcp", "--tools=agent"]
}
```

**OpenCode (new MCP format):**
```json
{
  "mcp": {
    "engram": {
      "command": ["engram", "mcp", "--tools=agent"],
      "enabled": true,
      "type": "local"
    }
  }
}
```

**VS Code (servers key):**
```json
{
  "servers": {
    "engram": {
      "command": "engram",
      "args": ["mcp", "--tools=agent"]
    }
  }
}
```

**Codex (TOML):**
```toml
[mcp_servers.engram]
command = "engram"
args = ["mcp", "--tools=agent"]
```

### Agent Adapter Integration

The agent-ai system automatically injects the appropriate MCP configuration based on the detected adapter. This is handled by the `Inject()` function in `internal/components/engram/inject.go`, which:
1. Writes MCP server config using the adapter's strategy
2. Injects the Engram memory protocol into system prompts (if supported)

---

## 11. Memory Categories

Use these types when saving observations:

| Type | When to Use |
|------|-------------|
| `bugfix` | Fixed a bug — record what was wrong and how it was fixed |
| `decision` | Architectural or design decision — record rationale |
| `architecture` | Structural decisions about the codebase |
| `discovery` | Non-obvious findings about the codebase or tools |
| `pattern` | Established patterns (naming, structure, conventions) |
| `config` | Configuration changes or environment setup |
| `preference` | User preferences or constraints learned |

**Scope:** Defaults to `project` (shared across sessions). Use `personal` only for agent-specific notes not relevant to the project.

---

## 12. Scope Rules

| Scope | Description | Use Case |
|-------|-------------|----------|
| `project` (default) | Shared memory across all sessions for a project | Code decisions, bugs, patterns, conventions |
| `personal` | Agent-specific notes | Agent preferences, personal reminders |

**Rule:** Default to `project` scope unless the information is truly agent-specific and not relevant to the project.

---

## Summary: Mandatory Memory Flow

```
SESSION START:
  mem_context → mem_search → mem_search → mem_get_observation (for each hit)

DURING SESSION:
  When bug fixed → mem_save (type: bugfix)
  When decision made → mem_save (type: decision)
  When discovery made → mem_save (type: discovery)
  Before any statement about code → mem_search (verify first!)

SESSION END:
  mem_session_summary (REQUIRED)

AFTER COMPACTION:
  mem_session_summary → mem_context → restore SDD state via mem_search
```

---

**Engram ensures the agent remembers everything, eliminating hallucinations and ensuring continuity across sessions.**
