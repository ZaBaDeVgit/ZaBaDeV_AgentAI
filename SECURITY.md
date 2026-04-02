# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in ZaBaDeV AgentAI, please report it
privately by emailing the maintainers or opening a private security advisory
on GitHub.

**Do not open a public issue for security vulnerabilities.**

## Scope

This policy covers:
- The ZaBaDeV AgentAI CLI (`zabadev`)
- All components in the `internal/` directory
- Embedded assets and skills
- Installation scripts

## What We Consider a Security Issue

- Remote code execution vulnerabilities
- Path traversal or file inclusion bugs
- Credential or secret leakage
- Privilege escalation
- Denial of service in critical paths

## Response Timeline

We aim to acknowledge reports within 48 hours and provide a fix within 14 days
for confirmed critical vulnerabilities.
