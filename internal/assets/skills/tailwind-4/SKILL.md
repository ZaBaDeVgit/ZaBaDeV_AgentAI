---
name: tailwind-4
description: >
  Tailwind CSS 4 patterns and best practices. Trigger: When styling with Tailwind - cn(), theme variables, no var() in className.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Styling React/Next.js components
- Building responsive layouts
- Creating design systems
- Working with Tailwind v4

---

## Critical Patterns

### Pattern 1: cn() Utility for Class Merging

Use `clsx` or `tailwind-merge` for conditional classes:

```typescript
// lib/utils.ts
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// Usage
cn(
  "base-styles",
  isActive && "active-styles",
  variant === "primary" ? "primary-colors" : "secondary-colors"
);
```

### Pattern 2: Theme Variables (CSS @theme)

Define design tokens in CSS:

```css
/* app.css */
@theme {
  /* Colors */
  --color-primary: oklch(70% 0.2 250);
  --color-secondary: oklch(60% 0.15 200);
  
  /* Spacing */
  --spacing-container: 1280px;
  
  /* Typography */
  --font-display: "Inter", system-ui;
  --font-body: "Roboto", sans-serif;
  
  /* Animations */
  --animate-fade-in: fade-in 0.3s ease-out;
}

@keyframes fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}
```

Use in className:
```jsx
// ❌ Don't do this
<div className="text-[color:--color-primary]">

// ✅ Do this - use the CSS variable directly
<div className="text-primary">
```

### Pattern 3: No var() in className

Tailwind v4 generates utility classes from theme:

```jsx
// ❌ Old v3 way
<div className="text-[var(--color-primary)]">

// ✅ New v4 way - define in @theme, use directly
<div className="text-primary bg-secondary">
```

### Pattern 4: Container Queries

Responsive based on container, not viewport:

```jsx
// CSS
@theme {
  --spacing-container-sm: 640px;
  --spacing-container-md: 768px;
}

.card-container {
  container-type: inline-size;
}

// Component
function Card({ children }) {
  return (
    <div className="card-container bg-white rounded-lg p-4">
      <h2 className="@lg:text-xl @xl:text-2xl font-bold">
        Dynamic sizing
      </h2>
      <p className="@sm:text-sm @md:text-base">
        {children}
      </p>
    </div>
  );
}
```

### Pattern 5: Compose with @apply Sparingly

Prefer utility classes, use @apply for repeated groups:

```css
/* Good - repeated utilities in CSS */
.btn-primary {
  @apply bg-primary text-white px-4 py-2 rounded-lg;
}

/* Also good - simple utility chains in JSX */
<div className="flex items-center justify-between">
```

---

## Decision Tree

```
Styling approach?
├── Conditional classes → cn() utility
├── Design tokens → Define in @theme
├── Responsive container → Container queries
├── Repeated styles → CSS @apply
└── Custom animation → Define in @theme + @keyframes

Component type?
├── Simple → Utility classes directly
├── Reusable → CSS with @apply
└── Complex system → @theme variables
```

---

## Anti-Patterns

- ❌ Using `var()` in className - use theme variables directly
- ❌ Overusing @apply - defeats utility-first purpose
- ❌ Hardcoding colors in JSX - use theme tokens
- ❌ Not using cn() - breaks conditional class merging

---

## Resources

- **Tailwind v4 Docs**: https://tailwindcss.com/docs/upgrade-guide
- **CSS Theme**: https://tailwindcss.com/docs/theme
- **Container Queries**: https://tailwindcss.com/docs/container-queries