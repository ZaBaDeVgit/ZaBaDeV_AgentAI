# Engram Protocol

This section contains the Engram memory system integration protocol.

## Memory Integration

When session starts, call `engram mem` to load persistent context.

## Session Tracking

Engram tracks conversation history across sessions using HTTP API at localhost:2112

## Session End Protocol

Before ending or saying "done" / "listo" / "finished", call `mem_save` with structured content:
- title: Short searchable description
- topic_key: Stable key for the type of work (e.g., "sdd/{change-name}/proposal")
- type: One of bugfix | decision | discovery | pattern | config | preference
- scope: project (default) or personal  
- content: What was done, why, where, and any learned gotchas

Example:
```
mem_save(
  title: "Fixed N+1 query in UserList",
  type: "bugfix",
  project: "my-project",
  content:
    What: Added eager loading for user posts relationship
    Why: User reported slow load times on user list page
    Where: internal/database/users.go, internal/api/handlers.go
    Learned: Always use preload for to-many relationships in high-traffic handlers
)
```