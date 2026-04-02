---
name: react-19
description: >
  React 19 patterns with React Compiler. Trigger: When writing React components.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Writing React components (any version 18+ with compiler)
- Using Server Components
- Managing async data with promises
- Handling form state with Actions

---

## Critical Patterns

### Pattern 1: No Manual Memoization

React Compiler handles optimization automatically. Don't add useMemo/useCallback unless:
- The dependency array is complex and causing issues
- You're passing callbacks to optimized child components

```jsx
// ❌ Don't do this - compiler handles it
const value = useMemo(() => computeExpensiveValue(a, b), [a, b]);
const handleClick = useCallback(() => doSomething(a), [a]);

// ✅ Do this - let the compiler optimize
const value = computeExpensiveValue(a, b);
const handleClick = () => doSomething(a);
```

### Pattern 2: Server Components

Use `'use server'` for server-side logic:

```jsx
// server-component.jsx
async function ServerComponent() {
  const data = await fetchData(); // Runs on server
  return <div>{data.name}</div>;
}

// server-action.js
'use server';
async function createTodo(formData) {
  const title = formData.get('title');
  await db.todos.create({ title });
  revalidatePath('/todos');
}
```

### Pattern 3: Actions (useActionState, useFormStatus)

```jsx
'use server';
async function submitForm(prevState, formData) {
  // Validate and process
  if (invalid) return { errors: [...] };
  return { success: true };
}

function Form() {
  const [state, action, isPending] = useActionState(submitForm, null);
  
  return (
    <form action={action}>
      <input name="title" />
      <button disabled={isPending}>Submit</button>
    </form>
  );
}
```

### Pattern 4: use() Hook

Unwrap promises and context directly:

```jsx
// Unwrap promise
function TodoList() {
  const todos = use(fetchTodos()); // Suspense handles loading
  return <ul>{todos.map(t => <li key={t.id}>{t.title}</li>)}</ul>;
}

// Unwrap context (replaces useContext)
function ThemeButton() {
  const theme = use(ThemeContext); // Works with useState too
  return <button className={theme}>{theme.name}</button>;
}
```

### Pattern 5: useOptimistic

Optimistic UI updates:

```jsx
function TodoItem({ todo, onToggle }) {
  const [optimisticTodo, setOptimistic] = useOptimistic(
    todo,
    (state, completed) => ({ ...state, completed })
  );

  return (
    <li className={optimisticTodo.completed ? 'done' : ''}>
      {optimisticTodo.title}
      <button onClick={async () => {
        setOptimistic(!todo.completed);
        await onToggle(todo.id);
      }}>Toggle</button>
    </li>
  );
}
```

---

## Decision Tree

```
Component type?
├── Server Component → Use 'use server' for data fetching
├── Client Component → Use useState, useEffect normally
└── Shared → Use use() for promises, 'use client' sparingly

Data fetching?
├── Server-side → Server Component with fetch()
├── Client-side → use() with Suspense
└── Streaming → Use use() + loading.js

Form handling?
├── Simple → useState
├── With validation → useActionState
└── Optimistic → useOptimistic

State management?
├── Local → useState/useReducer
├── Global → Zustand (see zustand-5 skill)
└── Server state → React Query / use() with fetch
```

---

## Anti-Patterns

- ❌ Adding useMemo/useCallback "just in case" - let compiler optimize
- ❌ Using 'use client' on everything - only when needed (event handlers, hooks)
- ❌ Mixing Server/Client components improperly
- ❌ Not handling Suspense boundaries for async components

---

## Resources

- **React 19 Docs**: https://react.dev/blog/2024/04/25/react-19
- **React Compiler**: https://github.com/reactjs/rfcs/blob/main/text/0000-react-compiler.md
- **Server Actions**: https://react.dev/reference/rsc/server-functions