---
name: sdd-new
description: >
  The missing SDD phase - explore + propose combined. Trigger: /sdd-new <change> - quick explore + propose workflow.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- User wants a quick start: explore + propose in one command
- Change is small enough that separate explore/propose adds overhead
- You need to quickly validate the change is worth pursuing

**Do NOT use when:**
- Change is complex and needs deep exploration
- Requirements are unclear and need investigation
- Design decisions need separate analysis

---

## Workflow

### Step 1: Investigate

Quickly explore the codebase to understand:
- What already exists related to this change
- Where the change would need to be made
- Potential risks or complexities

```bash
# Use grep/glob to find relevant files
grep -r "related-term" --include="*.go"
glob "**/relevant/*.go"
```

### Step 2: Create Proposal

Write a proposal with:
- **Intent**: What problem does this solve?
- **Scope**: What's included/excluded?
- **Approach**: How will you tackle it?

### Step 3: Save to Memory

Store in engram with key: `sdd/{change}/proposal`

```markdown
# Proposal: {change-name}

## Intent
{One paragraph on why this change matters}

## Scope
- ✅ Included: {list}
- ❌ Excluded: {list}

## Approach
- {High-level approach}
- {Key technical decisions}

## Risks
- {Potential risk 1}
- {Potential risk 2}

## Next Steps
- Create specs (sdd-spec)
- Create design (sdd-design)
- Break into tasks (sdd-tasks)
```

---

## Artifact Store

| Mode | When to Use |
|------|-------------|
| `engram` | **Default** - Persistent across sessions |
| `openspec` | File-based - when explicitly requested |
| `hybrid` | Both backends |
| `none` | Inline only - NOT RECOMMENDED |

---

## Decision: sdd-new vs Separate Phases

```
Use sdd-new when:
├── Change is small (1-3 files)
├── Requirements are clear
├── Approach is straightforward
└── User wants speed

Use separate explore + propose when:
├── Change is large/complex
├── Requirements unclear
├── Multiple approaches possible
├── Needs deep investigation
└── Design needs separate thinking
```

---

## Relationship to Other SDD Phases

```
sdd-new
   │
   ├──→ sdd-spec (write specs)
   ├──→ sdd-design (design decisions)
   ├──→ sdd-tasks (break into tasks)
   └──→ sdd-apply (implementation)

Alternative chain:
sdd-explore → sdd-propose → sdd-spec → ...
```

---

## Commands

| Command | Description |
|---------|-------------|
| `/sdd-new <change>` | Quick explore + propose |
| `/sdd-explore <topic>` | Just investigation |
| `/sdd-propose` | Create/update proposal only |

---

## Resources

- **SDD Phases**: See individual skill files
- **Engram Protocol**: Use mem_search/mem_save
- **Artifact Keys**: `sdd/{change}/proposal`