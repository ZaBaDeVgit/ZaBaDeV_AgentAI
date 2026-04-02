---
name: zod-4
description: >
  Zod 4 schema validation patterns. Trigger: When using Zod for validation - breaking changes from v3.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Validating API data
- Type-safe form handling
- Schema definition
- Runtime type checking

---

## Critical Patterns

### Pattern 1: Basic Schema

```typescript
import { z } from 'zod';

// Simple schema
const UserSchema = z.object({
  id: z.string().uuid(),
  name: z.string().min(2).max(100),
  email: z.string().email(),
  age: z.number().min(0).optional(),
  role: z.enum(['admin', 'user', 'guest']),
});

// Infer TypeScript type
type User = z.infer<typeof UserSchema>;
```

### Pattern 2: Breaking Changes from v3

```typescript
// v3 style (deprecated in v4)
const SchemaV3 = z.object({
  // .nullable() - still works
  value: z.string().nullable(),
  
  // .optional() with default
  count: z.number().optional().default(0),
});

// v4: Check for new APIs and imports
import { z } from 'zod/';

// v4 may have different default behaviors
// Always verify with actual version
```

### Pattern 3: Schema Composition

```typescript
// Nested schemas
const AddressSchema = z.object({
  street: z.string(),
  city: z.string(),
  zip: z.string(),
});

const UserWithAddressSchema = z.object({
  ...UserSchema.shape,
  address: AddressSchema,
});

// Merge schemas
const FullSchema = z.object({
  ...BaseSchema.shape,
  ...ExtraSchema.shape,
});

// Partial / Partial
const PartialUser = UserSchema.partial();
const RequiredEmail = UserSchema.pick({ email: true });
```

### Pattern 4: Refinement & Transformation

```typescript
// Refinement - custom validation
const ValidatedSchema = z.object({
  password: z.string().min(8),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ['confirmPassword'],
});

// Transform - change value
const ProcessedSchema = z.string().transform((val) => val.trim().toLowerCase());

// Parse and transform
const result = ProcessedSchema.parse('  Hello  '); // 'hello'
```

### Pattern 5: Error Formatting

```typescript
import { z } from 'zod';

const schema = z.object({
  name: z.string().min(2),
  email: z.string().email(),
});

try {
  schema.parse(invalidData);
} catch (error) {
  if (error instanceof z.ZodError) {
    // Format errors
    console.log(error.format());
    // { name: { errors: [...], ... }, email: { ... } }
    
    // Flat errors
    error.errors.forEach((err) => {
      console.log(`${err.path.join('.')}: ${err.message}`);
    });
  }
}
```

### Pattern 6: Union & Discriminated Union

```typescript
// Simple union
const StringOrNumber = z.string().or(z.number());

// Discriminated union
const EventSchema = z.discriminatedUnion('type', [
  z.object({ type: z.literal('click'), x: z.number(), y: z.number() }),
  z.object({ type: z.literal('keydown'), key: z.string() }),
  z.object({ type: z.literal('scroll'), scrollY: z.number() }),
]);

type Event = z.infer<typeof EventSchema>;
```

---

## Decision Tree

```
Validation type?
├── Object → z.object()
├── Array → z.array()
├── Primitive → z.string(), z.number(), etc.
├── Union → z.union() or discriminatedUnion
└── Optional → .optional()

Need custom logic?
├── Simple check → .refine()
├── Transform value → .transform()
└── Complex → Separate function + .superRefine()

Need TypeScript?
├── Extract type → z.infer<typeof Schema>
└── Partial type → Schema.partial()
```

---

## Anti-Patterns

- ❌ Not using z.infer - loses TypeScript benefits
- ❌ Over-validating - trust some data sources
- ❌ Not handling errors properly - catch ZodError
- ❌ Using any type - use inferred types

---

## Commands

```bash
npm install zod@latest    # Install latest version
npx tsd                  # Type testing
```

---

## Resources

- **Zod Docs**: https://zod.dev/
- **v3→v4 Migration**: https://zod.dev/v4/changelog
- **Pattern Examples**: https://github.com/colinhacks/zod