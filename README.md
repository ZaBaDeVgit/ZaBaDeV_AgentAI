# ZaBaDeV-AgentAI

<div align="center">

<pre>
██████╗
╚══███╔╝
  ███╔╝
 ███╔╝
██████╗
╚══════╝
</pre>

<h1>Senior ZaBaDeV — AI Agent Ecosystem</h1>

<p><strong>One command. OpenCode fully configured with the complete ZaBaDeV ecosystem Full.</strong></p>

<p>
<a href="https://github.com/zabadev/agent-ai/releases"><img src="https://img.shields.io/github/v/release/zabadev/agent-ai" alt="Release"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
<img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white" alt="Go 1.21+">
<img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey" alt="Platform">
<a href="https://github.com/zabadev/agent-ai/actions"><img src="https://github.com/zabadev/agent-ai/workflows/CI/badge.svg" alt="CI"></a>
</p>

</div>

---

## Screenshot

![ZaBaDeV Agent AI](Captura.png)

---

## What is ZaBaDeV?

**Senior ZaBaDeV** is an AI development ecosystem that transforms your AI editor/agent into a professional development assistant with:

- **Persistent Memory (Engram)** — Remembers decisions, bugs, and conventions across sessions
- **SDD Workflow** — Spec-Driven Development: plan before you code
- **Professional Skills** — Coding patterns for React, TypeScript, Tailwind, testing, and more
- **MCP Servers** — Context7 for up-to-date documentation
- **Teaching-First Persona** — An architectural mentor that explains the "why" before the "what"
- **Automated Review with GGA** — Guardian Angel reviews every commit

---

## Key Features

| Feature                     | Description                                       |
| --------------------------- | ------------------------------------------------- |
| **One-Click Setup**         | Single command installs the complete ecosystem    |
| **Cross-Platform**          | Works on macOS, Linux, and Windows                |
| **Multi-Agent Support**     | OpenCode, Claude Code, Cursor, Gemini CLI, VSCode |
| **Persistent Learning**     | Remembers your preferences and project patterns   |
| **Spec-Driven Development** | Structured planning and implementation workflow   |
| **Automated Code Review**   | GGA checks code quality on every commit           |
| **Context7 Integration**    | Real-time documentation and examples              |
| **Extensible Skills**       | Modular skills system for custom workflows        |

---

## Quick Start

### Prerequisites

- Go 1.21 or later
- Git
- One of: OpenCode, Claude Code, Cursor, Gemini CLI, or VSCode

### Installation

```bash
# Clone the repository
git clone https://github.com/zabadev/agent-ai.git
cd agent-ai

# Build and install
make build
make install

# Or install directly
go install github.com/zabadev/agent-ai/cmd/zabadev@latest
```

### Setup

```bash
# Initialize ZaBaDeV ecosystem
zabadev install

# For development with hot reload
make dev
make run-dev
```

---

## Usage

### Basic Commands

```bash
# Install agents and tools
zabadev install

# Sync configurations
zabadev sync

# Update tools
zabadev upgrade

# Restore from backup
zabadev restore --list
zabadev restore <backup-id>

# Show version
zabadev version
```

### SDD Workflow

```bash
# Initialize project
zabadev sdd-init

# Explore feature
zabadev sdd-explore "add user authentication"

# Create change proposal
zabadev sdd-new "implement-user-auth"

# Write specifications
zabadev sdd-spec

# Design implementation
zabadev sdd-design

# Plan tasks
zabadev sdd-tasks

# Implement
zabadev sdd-apply

# Verify
zabadev sdd-verify

# Archive
zabadev sdd-archive
```

---

## Architecture

```
ZaBaDeV-AgentAI/
├── cmd/zabadev/          # Main CLI application
├── internal/
│   ├── agents/            # Agent-specific implementations
│   ├── app/               # Core application logic
│   ├── assets/            # Embedded assets and templates
│   ├── backup/            # Backup and restore functionality
│   ├── cli/               # CLI commands and execution
│   ├── components/        # UI components and rendering
│   ├── model/             # Data models and types
│   └── system/            # System detection and utilities
├── pkg/                   # Public packages
└── testdata/              # Test fixtures and golden files
```

---

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage

# Lint code
make lint

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests with race detection
go test ./... -race

# Run specific package tests
go test ./internal/app -v
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Write tests for your changes
4. Ensure CI passes: `make ci`
5. Submit a pull request

---

## Configuration

### Agent Configuration

ZaBaDeV supports multiple AI agents:

- **OpenCode** — Primary development environment
- **Claude Code** — Anthropic's coding assistant
- **Cursor** — AI-first code editor
- **Gemini CLI** — Google's AI assistant
- **VSCode** — With GitHub Copilot

### Skills System

Extend functionality with skills:

```bash
# List available skills
zabadev sync --skills

# Install specific skills
zabadev sync --skills react,typescript,testing
```

### Backup & Restore

```bash
# Create backup
zabadev install  # Creates automatic backup

# List available backups
zabadev restore --list

# Restore specific backup
zabadev restore backup-2024-01-15-14-30-00
```

---

## Documentation

- [Architecture Overview](docs/architecture.md)
- [SDD Workflow Guide](docs/sdd-workflow.md)
- [Agent Configuration](docs/agent-config.md)
- [Skills Development](docs/skills.md)
- [Troubleshooting](docs/troubleshooting.md)

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Support

- **Issues**: [GitHub Issues](https://github.com/zabadev/agent-ai/issues)
- **Discussions**: [GitHub Discussions](https://github.com/zabadev/agent-ai/discussions)
- **Documentation**: [docs/](docs/)

---

<div align="center">

**Built with ❤️ for the AI development community**

[⭐ Star us on GitHub](https://github.com/zabadev/agent-ai) • [📖 Read the docs](docs/) • [🐛 Report bugs](https://github.com/zabadev/agent-ai/issues)

</div>

