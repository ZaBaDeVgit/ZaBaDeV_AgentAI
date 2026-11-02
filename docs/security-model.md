# Security model and permissions

ZaBaDeV configures third-party coding agents; it is not an operating-system sandbox. OpenCode, the selected model, MCP servers, plugins, shell commands and the operating system all participate in the trust boundary. A prompt or file-access rule cannot make arbitrary shell execution safe.

## Safer defaults

The generated OpenCode permission overlay asks before shell commands and denies reads/edits matching common secret paths (`.env`, `.envrc`, `.npmrc`, `secrets/`, credential JSON and private-key-like filenames) through the agent's read/edit tools. Other supported agents are configured to retain their normal approval behavior; automatic edit/approve modes are not enabled by this installer.

These are defense-in-depth measures, not a guarantee:

- A user can approve a risky command, and a command interpreter may access files outside read-tool rules.
- Secret filenames and conventions vary. Do not store credentials in a project workspace.
- MCP servers and plugins run with the privileges available to their process.
- Repository instructions, webpages, build scripts and issue content can contain hostile prompt-injection text. Treat them as data and do not follow requests to expose secrets or bypass approval.

## Recommended operating practice

1. Keep OpenCode's confirmation prompts enabled. Review each command's exact command line and working directory before approval.
2. Use OS-level sandboxing or a disposable development container for untrusted repositories or autonomous runs. Restrict filesystem and network access where practical.
3. Keep API keys in the provider's secure credential store or environment manager, not in source files. Never ask an agent to print them.
4. Review dependency install scripts, migrations, cleanup operations and commands that access paths outside the workspace.
5. Inspect `git diff` and run tests before committing. Require a human decision for commits, pushes, releases, deployments and destructive operations.
6. Re-run `zabadev install --dry-run` before changing the installed agent set or components. Back up agent configuration before manual edits.

## Permission configuration

The permissions component merges policy into the agent's configuration while preserving unrelated settings. Users who want a different policy should edit the generated agent configuration deliberately and understand that re-running installation may update the managed permission section. Do not weaken approval globally merely to reduce prompts.

## Reporting vulnerabilities

See [SECURITY.md](../SECURITY.md) for private vulnerability reporting.