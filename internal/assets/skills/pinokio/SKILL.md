---
name: pinokio
description: >
  Discover, launch, and use apps and tools for the current task.
  Trigger: When working with Pinokio - runtime control, log checking, examples.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Working with Pinokio runtime
- Launching and controlling Pinokio apps
- Debugging Pinokio-related issues
- Following examples in the Pinokio ecosystem

---

## Critical Patterns

### Pattern 1: pterm for Pinokio Operations

Use `pterm` for all Pinokio control:

```typescript
import { pterm } from '@pinokio/client';

// List running apps
const apps = await pterm.list();

// Launch an app
const app = await pterm.launch({
  id: 'my-app',
  repo: 'github.com/user/repo',
  branch: 'main',
});

// Stop an app
await pterm.stop('my-app');
```

### Pattern 2: Log Checking

When debugging, check Pinokio logs:

```
pinokio/logs/
├── app-YYYY-MM-DD.log    # App-specific logs
├── runtime.log          # Runtime logs
└── error.log           # Error logs
```

```bash
# View recent logs
tail -f pinokio/logs/app.log

# Search for errors
grep -i error pinokio/logs/*.log
```

### Pattern 3: Examples

Follow examples in:
```
C:\pinokio\prototype\system\examples
# or Linux/Mac:
~/pinokio/prototype/system/examples
```

Key examples:
- `basic-launch/` - Simple app launch
- `runtime-control/` - Controlling running apps
- `file-operations/` - Reading/writing files

### Pattern 4: API Reference

See `PINOKIO.md` for full API:

```typescript
// Full API methods
pterm.launch(config)      // Launch app
pterm.stop(appId)         // Stop app
pterm.list()              // List running apps
pterm.exec(appId, cmd)   // Run command in app
pterm.files(appId, path) // Access files
```

---

## Decision Tree

```
Pinokio task?
├── Launch app → pterm.launch()
├── Stop app → pterm.stop()
├── Get status → pterm.list()
├── Run command → pterm.exec()
└── Debug → Check logs in pinokio/logs/

Problem debugging?
├── App crash → Check app-YYYY-MM-DD.log
├── Runtime issue → Check runtime.log
└── General error → Check error.log
```

---

## Anti-Patterns

- ❌ Using raw child_process instead of pterm
- ❌ Not checking logs when debugging
- ❌ Ignoring Pinokio examples
- ❌ Not reading PINOKIO.md API reference

---

## Resources

- **Pinokio Docs**: See PINOKIO.md in project
- **Examples**: `pinokio/prototype/system/examples/`
- **API**: `PINOKIO.md` in root