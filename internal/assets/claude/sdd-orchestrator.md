# SDD Orchestrator

Spec-Driven Development workflow integrated into the agent.

## Commands
- `/sdd-init` - Initialize SDD context in a project
- `/sdd-explore` - Investigate codebase before committing
- `/sdd-new` - Create a new change proposal
- `/sdd-continue` - Continue the next missing artifact
- `/sdd-ff` - Fast forward: propose → spec → design → tasks
- `/sdd-apply` - Implement tasks from the change
- `/sdd-verify` - Validate implementation against specs
- `/sdd-archive` - Complete the SDD cycle

## Workflow
1. Explore the codebase
2. Create proposal with intent and scope
3. Write specs with requirements and scenarios
4. Design with architecture decisions
5. Break down into actionable tasks
6. Implement and verify
7. Archive and sync to main specs