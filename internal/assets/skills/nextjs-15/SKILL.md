---
name: nextjs-15
description: >
  Next.js 15 App Router patterns. Trigger: When working with Next.js - routing, Server Actions, data fetching.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Building Next.js applications with App Router
- Creating route handlers and API endpoints
- Implementing Server Actions
- Setting up data fetching with caching

---

## Critical Patterns

### Pattern 1: Partial Prerendering (PPR)

Mix static and dynamic content:

```jsx
// app/page.jsx
import { Suspense } from 'react';

export default function Page() {
  return (
    <div>
      {/* Static - served immediately */}
      <header>
        <h1>My App</h1>
      </header>
      
      {/* Dynamic - streams in */}
      <Suspense fallback={<Skeleton />}>
        <UserFeed />
      </Suspense>
    </div>
  );
}
```

Enable in `next.config.js`:
```js
module.exports = {
  experimental: {
    ppr: true,
  },
};
```

### Pattern 2: Server Actions

```jsx
// app/actions.ts
'use server';

export async function createPost(formData: FormData) {
  const title = formData.get('title');
  
  // Direct DB call - no API route needed
  await db.posts.create({ title });
  
  // Revalidate
  revalidatePath('/posts');
  return { success: true };
}
```

```jsx
// app/posts/page.tsx
'use client';
import { createPost } from '@/app/actions';

export default function CreatePost() {
  return (
    <form action={createPost}>
      <input name="title" />
      <button type="submit">Create</button>
    </form>
  );
}
```

### Pattern 3: Data Fetching with Caching

```jsx
// Fetch with cache options
async function getData() {
  const res = await fetch('https://api.example.com/data', {
    next: { 
      revalidate: 3600,        // Revalidate every hour
      tags: ['posts'],         // For on-demand revalidation
    },
  });
  return res.json();
}

// Force dynamic (no cache)
export const dynamic = 'force-dynamic';
```

### Pattern 4: Route Handlers

```tsx
// app/api/posts/route.ts
import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
  const posts = await db.posts.findMany();
  return NextResponse.json(posts);
}

export async function POST(request: NextRequest) {
  const data = await request.json();
  const post = await db.posts.create(data);
  return NextResponse.json(post, { status: 201 });
}
```

### Pattern 5: Turbopack

Use `--turbopack` flag for dev server:
```bash
npm run dev -- --turbopack
# or
next dev --turbopack
```

---

## Decision Tree

```
Need to fetch data?
├── Server Component → fetch() directly with caching
├── Client Component → use() with Suspense or React Query
└── External API → Route handler for complex logic

Form submission?
├── Simple → Server Action (no API route needed)
├── Complex validation → Server Action with useActionState
└── Multi-step → Use useActionState + useOptimistic

Caching strategy?
├── Static (rarely changes) → default cache
├── Semi-static → revalidate: 3600
├── Dynamic → dynamic = 'force-dynamic'
└── On-demand → revalidateTag('tag')

Route type?
├── Page → app/page.tsx
├── API → app/api/route.ts
├── Layout → app/layout.tsx
└── Middleware → middleware.ts
```

---

## Anti-Patterns

- ❌ Using Pages Router in new projects - App Router is the future
- ❌ Making fetch calls in useEffect - use Server Components
- ❌ Not handling Suspense - breaks streaming
- ❌ Overusing 'use client' - only when needed

---

## Resources

- **Next.js 15 Docs**: https://nextjs.org/blog/next-15
- **PPR RFC**: https://github.com/vercel/next.js/discussions/64465
- **Server Actions**: https://nextjs.org/docs/app/building-your-application/data-fetching/server-actions-and-mutations