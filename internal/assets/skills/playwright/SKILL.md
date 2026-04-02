---
name: playwright
description: >
  Playwright E2E testing patterns. Trigger: When writing E2E tests - Page Objects, selectors, MCP workflow.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Writing end-to-end tests for web applications
- Testing complex user flows
- Creating visual regression tests
- Automating browser interactions

---

## Critical Patterns

### Pattern 1: Page Objects

```typescript
// pages/LoginPage.ts
import { Page, Locator } from '@playwright/test';

export class LoginPage {
  readonly page: Page;
  readonly usernameInput: Locator;
  readonly passwordInput: Locator;
  readonly submitButton: Locator;
  readonly errorMessage: Locator;

  constructor(page: Page) {
    this.page = page;
    this.usernameInput = page.getByLabel('Username');
    this.passwordInput = page.getByLabel('Password');
    this.submitButton = page.getByRole('button', { name: 'Login' });
    this.errorMessage = page.getByRole('alert');
  }

  async goto() {
    await this.page.goto('/login');
  }

  async login(username: string, password: string) {
    await this.usernameInput.fill(username);
    await this.passwordInput.fill(password);
    await this.submitButton.click();
  }
}
```

```typescript
// tests/login.spec.ts
import { test, expect } from '@playwright/test';
import { LoginPage } from '../pages/LoginPage';

test('successful login', async ({ page }) => {
  const loginPage = new LoginPage(page);
  await loginPage.goto();
  await loginPage.login('admin', 'password123');
  
  await expect(page).toHaveURL('/dashboard');
});
```

### Pattern 2: Modern Selectors

Prefer accessibility-aware selectors:

```typescript
// Best - getByRole (accessibility)
await page.getByRole('button', { name: 'Submit' });
await page.getByRole('textbox', { name: 'Email' });
await page.getByRole('link', { name: 'Learn more' });

// Good - getByText for non-interactive elements
await page.getByText('Welcome back');

// Good - getByLabel for form fields
await page.getByLabel('Email').fill('test@example.com');
await page.getByLabel('Password').fill('secret');

// Fallback - getByTestId when no semantic option
await page.getByTestId('submit-button');

// Avoid - CSS selectors (fragile)
await page.locator('.btn-primary').click();
```

### Pattern 3: MCP Workflow for AI Testing

Use MCP server for AI-assisted test generation:

```typescript
// playwright.config.ts
import { defineConfig } from '@playwright/test';

export default defineConfig({
  use: {
    // AI test generation via MCP
    launchOptions: {
      // Setup MCP connection
    },
  },
});
```

### Pattern 4: Visual Regression Testing

```typescript
import { test, expect } from '@playwright/test';

test('visual regression - home page', async ({ page }) => {
  await page.goto('/');
  
  // Full page screenshot
  await expect(page).toHaveScreenshot('homepage.png');
  
  // Specific element
  await expect(page.locator('.hero')).toHaveScreenshot('hero.png');
});
```

### Pattern 5: Fixtures and Custom Test

```typescript
// fixtures.ts
import { test as base } from '@playwright/test';

export const test = base.extend({
  authenticatedPage: async ({ page }, use) => {
    // Setup: login
    await page.goto('/login');
    await page.getByLabel('Email').fill('admin@test.com');
    await page.getByLabel('Password').fill('password');
    await page.getByRole('button', { name: 'Login' }).click();
    
    await use(page);
    
    // Teardown
    await page.context().clearCookies();
  },
});

// Use in tests
test('dashboard loads', async ({ authenticatedPage }) => {
  await authenticatedPage.goto('/dashboard');
  await expect(authenticatedPage.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
});
```

---

## Decision Tree

```
Writing a test?
├── Page Object needed? → Create pages/ directory
├── New flow? → Add to existing Page Object
└── Unique page? → New Page Object

Selector priority?
├── Interactive element? → getByRole
├── Form field? → getByLabel
├── Text content? → getByText
└── Last resort? → getByTestId

Test scope?
├── Single feature → Feature file
├── Full flow → Multiple tests in file
└── Cross-browser → Use @cross-browser tag
```

---

## Anti-Patterns

- ❌ Using CSS selectors (`.class`, `#id`) - fragile, breaks on changes
- ❌ Hard-coding waits (sleep) - use expect with timeout
- ❌ Not using Page Objects - leads to duplication
- ❌ Testing implementation details - test behavior, not code

---

## Commands

```bash
npx playwright test                    # Run all tests
npx playwright test --ui               # Interactive UI
npx playwright test --headed           # Visible browser
npx playwright test --debug            # Debug mode
npx playwright test --grep "login"      # Run matching tests
npx playwright show-report             # View HTML report
```

---

## Resources

- **Playwright Docs**: https://playwright.dev/
- **Selectors Guide**: https://playwright.dev/docs/locators
- **Page Objects**: https://playwright.dev/docs/test-pom