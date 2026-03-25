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

<p><strong>One command. OpenCode fully configured with the complete ZaBaDeV ecosystem.</strong></p>

<p>
<a href="https://github.com/zabadev/agent-ai/releases"><img src="https://img.shields.io/github/v/release/zabadev/agent-ai" alt="Release"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
<img src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white" alt="Go 1.24+">
<img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey" alt="Platform">
</p>

</div>

---

## Captura de Pantalla

![ZaBaDeV Agent AI](Captura.png)

---

## Que es ZaBaDeV?

**Senior ZaBaDeV** es un ecosistema de desarrollo de IA que transforma tu editor/agent de IA en un asistente de desarrollo profesional con:

- **Memoria Persistente (Engram)** — Recuerda decisiones, bugs y convenciones entre sesiones
- **Workflow SDD** — Spec-Driven Development: planifica antes de codificar
- **Skills profesionales** — Patrones de codificación para React, TypeScript, Tailwind, testing y más
- **MCP Servers** — Context7 para documentación actualizada
- **Persona Teaching-First** — Un mentor arquitectónico que explica el "por qué" antes del "qué"
- **Review automático con GGA** — Guardian Angel revisa cada commit

---

## Características Principales

| Caracteristica | Descripcion |
|----------------|-------------|
| **Engram** | Sistema de memoria persistente que survive entre sesiones |
| **SDD Workflow** | 9 skills para Spec-Driven Development: init, explore, propose, spec, design, tasks, apply, verify, archive |
| **Skills** | 11+ skills profesionales para desarrollo moderno |
| **MCP Servers** | Context7, Notion, Jira para integracion de documentacion y proyectos |
| **GGA** | Guardian Angel — revision de codigo AI en cada commit |
| **Persona ZaBaDeV** | Modo mentor arquitectonico con estilo teaching-first |
| **Multi-plataforma** | macOS, Linux, Windows (WSL) |

---

## Instalacion

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/zabadev/agent-ai/main/scripts/install.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/zabadev/agent-ai/main/scripts/install.ps1 | iex
```

### Homebrew (macOS / Linux)

```bash
brew tap zabadev/homebrew-tap
brew install zabadev
```

### Go install (cualquier plataforma con Go 1.24+)

```bash
go install github.com/zabadev/agent-ai/cmd/zabadev@latest
```

### Desde Releases

Descarga el binario para tu plataforma desde [GitHub Releases](https://github.com/zabadev/agent-ai/releases).

---

## Uso

```bash
# Instalacion estandar en OpenCode (comportamiento por defecto)
zabadev

# Mostrar ayuda
zabadev --help

# Mostrar version
zabadev version
```

---

## Que se Instala

Cuando ejecutas `zabadev`, se instala en tu agente de IA:

| Componente | Descripcion |
|------------|-------------|
| **Engram** | Sistema de memoria persistente |
| **SDD Workflow** | Spec-Driven Development completo |
| **Skills** | 11+ skills profesionales |
| **Context7** | Servidor MCP para documentacion |
| **GGA** | Automatizacion de agente global |
| **Permisos** | Configuracion security-first |
| **Persona** | Senior ZaBaDeV modo ensenanza |

---

## Estructura del Proyecto

```
cmd/zabadev/             # Punto de entrada CLI
internal/
  app/                   # Dispatch de comandos + wiring
  model/                 # Tipos de dominio
  catalog/               # Definiciones de registro
  system/                # Deteccion OS/distro
  cli/                   # Flags de instalacion
  pipeline/              # Ejecucion por etapas
  backup/                # Snapshot de config
  components/            # Logica por componente
    engram/  sdd/  skills/  mcp/  persona/
  agents/                # Adaptadores de agente
    claude/  opencode/  gemini/  cursor/
  verify/                # Health checks
  tui/                   # Interfaz Bubbletea
scripts/                 # Scripts de instalacion
e2e/                     # Tests E2E en Docker
testdata/                # Fixtures golden
```

---

## Testing

```bash
# Tests unitarios
go test ./...

# Docker E2E (Ubuntu + Arch + Fedora, requiere Docker)
RUN_FULL_E2E=1 RUN_BACKUP_TESTS=1 ./e2e/docker-test.sh
```

---

## Documentacion Adicional

- [Arquitectura](docs/architecture.md) — Detalles tecnicos del proyecto
- [Usage](docs/usage.md) — Guia de uso avanzada
- [Non-Interactive](docs/non-interactive.md) — Modo no interactivo para CI
- [Quickstart](docs/quickstart.md) — Guia de inicio rapido
- [Platforms](docs/platforms.md) — Notas especificas por plataforma

---

## Relación con Gentleman.Dots

| | Gentleman.Dots | ZaBaDeV |
|--|---------------|---------|
| **Proposito** | Entorno de desarrollo (editores, shells, terminales) | Capa de desarrollo IA (agentes, memoria, skills) |
| **Instala** | Neovim, Fish/Zsh, Tmux/Zellij, Ghostty | Configura Claude Code, OpenCode, Gemini CLI, Cursor |
| **Superposicion** | Ninguna — complementario | Ninguna — diferente capa |

Instala Gentleman.Dots primero para tu entorno de desarrollo, luego ZaBaDeV para la capa de IA.

---

## Licencia

MIT License — consulta el archivo [LICENSE](LICENSE) para mas detalles.

---

## Agradecimientos

<div align="center">

**Dedicado a Gentleman, el creador original del ecosistema ZaBaDeV.**

</div>

Este proyecto esta inspirado en el trabajo visionario de **Gentleman Programming**, quien creo el concepto de un ecosistema de desarrollo de IA completo y profesional. Su contribucion al ecosistema de desarrollo con herramientas como Engram, SDD, y las skills de desarrollo ha sido fundamental para hacer posible este proyecto.

Para mas informacion sobre el ecosistema Gentleman original, visita: [gentleman.ai](https://gentleman.ai)

---

<div align="center">

_Con ❤️ desde la comunidad ZaBaDeV_

</div>