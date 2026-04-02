---
name: gepeto
description: >
  Guide for building 1-click launchers and apps with launchers built-in using Pinokio.
  Trigger: When building Pinokio apps - package.json, index.js, app configuration.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Building a Pinokio app from scratch
- Creating a 1-click launcher
- Setting up app configuration
- Integrating launcher into an app

---

## Critical Patterns

### Pattern 1: Package.json Structure

```json
{
  "name": "my-app",
  "version": "1.0.0",
  "pinokio": {
    "runtime": "node",
    "command": "node index.js",
    "autoLaunch": true
  },
  "scripts": {
    "start": "node index.js"
  },
  "dependencies": {
    "@pinokio/client": "^1.0.0"
  }
}
```

### Pattern 2: Index.js with Built-in Launcher

```javascript
const { pterm } = require('@pinokio/client');

async function main() {
  // Your app logic
  console.log('App running...');
  
  // Optionally: integrate with Pinokio
  // This makes it a "launcher-aware" app
  const status = await pterm.list();
  console.log('Running apps:', status.length);
}

// Run with Pinokio integration
if (require.main === module) {
  main().catch(console.error);
}

module.exports = { main };
```

### Pattern 3: App Configuration

Create `pinokio.json`:

```json
{
  "name": "My App",
  "id": "my-app",
  "description": "A Pinokio-powered app",
  "runtime": "node",
  "entry": "index.js",
  "env": {
    "NODE_ENV": "production"
  },
  "autoLaunch": true,
  "icon": "assets/icon.png"
}
```

### Pattern 4: Building with Launchers

Structure for launcher-integrated app:

```
my-app/
├── package.json        # With Pinokio config
├── index.js           # Main app + launcher integration
├── pinokio.json       # App metadata
├── launcher/          # Built-in launcher components
│   ├── install.js
│   └── configure.js
└── assets/
    └── icon.png
```

```javascript
// index.js - Launcher-aware
const { launcher } = require('./launcher');

async function start() {
  // Check if launched via Pinokio or standalone
  const isPinokio = process.env.PINOKIO_APP_ID !== undefined;
  
  if (isPinokio) {
    await launcher.initialize();
  }
  
  // Your app
  await runApp();
}
```

---

## Decision Tree

```
Building approach?
├── Simple app → Basic package.json + index.js
├── Full launcher → Add launcher/ directory
└── Config-driven → Use pinokio.json

Launcher features?
├── Auto-install → Add installer script
├── Configuration → Add configure.js
└── UI → Add launcher components
```

---

## Anti-Patterns

- ❌ Not using package.json Pinokio config
- ❌ Missing entry point
- ❌ Not testing with Pinokio runtime
- ❌ Hardcoding paths

---

## Resources

- **Pinokio Client**: `@pinokio/client` npm package
- **Examples**: See `pinokio/prototype/system/examples/`
- **Launcher Pattern**: Learn from existing launchers