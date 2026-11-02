---
name: project-engineering
description: >
  Detect project stacks and safely build, test, debug and deliver changes. Trigger: Before implementing in an unfamiliar repository or building a complete application.
license: MIT
metadata:
  author: ZaBaDeV
  version: "1.0"
---

# Project Engineering Workflow

Use this skill when starting in an unfamiliar repository, creating an application, or when the project's build and test workflow is unclear.

## 1. Identify the project from evidence

Inspect the project root and relevant subdirectories for manifests, lockfiles, CI workflows, README instructions and existing source. Common signals include:

- JavaScript/TypeScript: `package.json`, lockfiles, `tsconfig.json`, framework config.
- Python: `pyproject.toml`, `requirements*.txt`, `pytest.ini`, `manage.py`.
- Go: `go.mod`, `Makefile`, `*_test.go`.
- Rust: `Cargo.toml`, `Cargo.lock`.
- .NET: `.sln`, `.csproj`, `global.json`.
- Java/Kotlin: Gradle or Maven manifests.
- Mobile: platform manifests and project files (Android Gradle, Xcode project, Flutter `pubspec.yaml`).
- Monorepos: workspace manifests, package manager config, and per-package scripts.

Do not infer the whole stack from a single file. Check the package manager, runtime/toolchain version, scripts, CI and existing code. If multiple applications exist, identify which one the user means.

## 2. Discover commands before running them

Prefer commands already documented in the README, Makefile, package scripts, task runner or CI. Report when a required tool or dependency is missing. Do not install packages, download browser binaries, run migrations, or alter global machine state without user approval.

Use the narrowest relevant checks first, then broader checks if practical:

1. formatter / lint / static analysis;
2. targeted unit or integration tests;
3. full test suite;
4. production build;
5. browser or device-based end-to-end checks for user-facing flows.

Never invent a passing result. Distinguish checks that passed, failed, were skipped, or could not run.

## 3. Implement complete user-facing behavior

For an application feature, identify the user journey and cover the states it needs: loading, success, empty, validation errors, server errors and recovery. Respect existing design conventions. Consider semantic HTML, keyboard navigation, contrast, responsive layouts, input validation, authorization, privacy, observability and graceful failure where relevant.

For API/backend work, inspect existing contracts and add or update tests for validation, error cases, authorization and persistence. For data migrations, explain compatibility and rollback considerations. Avoid adding dependencies when existing project capabilities are sufficient.

## 4. Verify web UI changes

If the repository already has Playwright or another browser test tool, use the installed project workflow to exercise the changed user journey. Prefer accessible role/label selectors and assert visible behavior, navigation and relevant console/network errors. If a browser tool is not installed, do not claim visual verification; explain the missing prerequisite and provide the exact next check.

Do not create placeholder or fictitious browser/MCP configuration. Use only integrations that are actually installed and documented in the current project.

## 5. Finish responsibly

Review the diff for unrelated edits, secrets, generated files and accidental dependency churn. Summarize the implementation, checks and remaining limitations. Ask before committing, pushing, deploying, deleting user data or changing files outside the project.