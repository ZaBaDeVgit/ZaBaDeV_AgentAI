---
name: zustand-5
description: >
  Zustand 5 state management patterns. Trigger: When managing React state with Zustand.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Managing global React state
- Creating simple, performant state stores
- Needing middleware (persistence, devtools)
- Composing multiple stores

---

## Critical Patterns

### Pattern 1: Create with Selectors

```typescript
import { create } from 'zustand';

// Basic store
interface CounterStore {
  count: number;
  increment: () => void;
  decrement: () => void;
}

const useCounterStore = create<CounterStore>((set) => ({
  count: 0,
  increment: () => set((state) => ({ count: state.count + 1 })),
  decrement: () => set((state) => ({ count: state.count - 1 })),
}));

// Usage with selectors (prevents unnecessary re-renders)
function Counter() {
  const count = useCounterStore((state) => state.count);
  const increment = useCounterStore((state) => state.increment);
  
  return (
    <div>
      <p>Count: {count}</p>
      <button onClick={increment}>+</button>
    </div>
  );
}
```

### Pattern 2: Middleware

```typescript
import { create } from 'zustand';
import { persist, devtools, immer } from 'zustand/middleware';

// Persist - saves to localStorage
const useStore = create(
  persist(
    (set) => ({
      todos: [],
      addTodo: (todo) => set((state) => ({ 
        todos: [...state.todos, todo] 
      })),
    }),
    { name: 'my-store' } // localStorage key
  )
);

// Devtools - Redux DevTools integration
const useStore = create(
  devtools(
    (set) => ({
      // ...store logic
    }),
    { name: 'MyStore' }
  )
);

// Immer - mutable-style updates
const useStore = create(
  immer(
    (set) => ({
      items: [],
      addItem: (item) => set((state) => {
        state.items.push(item); // Mutate directly!
      }),
    })
  )
);
```

### Pattern 3: Store Composition

```typescript
// Base store
interface UserStore {
  user: User | null;
  setUser: (user: User | null) => void;
}

const createUserStore = () => create<UserStore>((set) => ({
  user: null,
  setUser: (user) => set({ user }),
}));

// Another store
interface CartStore {
  items: CartItem[];
  addItem: (item: CartItem) => void;
}

const createCartStore = () => create<CartStore>((set) => ({
  items: [],
  addItem: (item) => set((state) => ({ 
    items: [...state.items, item] 
  })),
}));

// Combined store
interface AppStore extends UserStore, CartStore {}

const useAppStore = create<AppStore>((set, get) => ({
  ...createUserStore()(set, get),
  ...createCartStore()(set, get),
}));
```

### Pattern 4: Async Actions

```typescript
interface AsyncStore {
  data: Data | null;
  loading: boolean;
  error: string | null;
  fetchData: () => Promise<void>;
}

const useAsyncStore = create<AsyncStore>((set) => ({
  data: null,
  loading: false,
  error: null,
  
  fetchData: async () => {
    set({ loading: true, error: null });
    try {
      const data = await api.fetchData();
      set({ data, loading: false });
    } catch (error) {
      set({ error: error.message, loading: false });
    }
  },
}));
```

---

## Decision Tree

```
State type?
├── Simple → Basic create()
├── Persistent → Add persist middleware
├── Complex mutations → Add immer middleware
└── Debugging → Add devtools middleware

Component needs?
├── Single value → Selector: useStore(s => s.value)
├── Multiple values → Multiple selectors or shallow
└── Actions only → Select actions without state

Store size?
├── Small → Single store
├── Large → Composed stores or slices pattern
```

---

## Anti-Patterns

- ❌ Not using selectors - causes unnecessary re-renders
- ❌ Selecting entire store - `useStore()` without selector
- ❌ Mutating state directly (without immer middleware)
- ❌ Over-splitting stores - keep related state together

---

## Resources

- **Zustand Docs**: https://zustand-demo.pmnd.rs/
- **Middleware**: https://zustand-demo.pmnd.rs/middleware
- **Cookbook**: https://github.com/pmndrs/zustand/wiki/Recipes