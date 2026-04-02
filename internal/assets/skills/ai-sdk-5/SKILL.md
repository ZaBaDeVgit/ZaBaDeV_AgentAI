---
name: ai-sdk-5
description: >
  Vercel AI SDK 5 patterns. Trigger: When building AI chat features - breaking changes from v4.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Building chat interfaces
- Implementing streaming AI responses
- Adding tool/function calling
- Creating completion/summarization features

---

## Critical Patterns

### Pattern 1: useChat Hook

```typescript
'use client';

import { useChat } from 'ai/react';

export function Chat() {
  const { messages, input, handleInputChange, handleSubmit, isLoading } = useChat({
    api: '/api/chat',
    onError: (error) => {
      console.error('Chat error:', error);
    },
  });

  return (
    <div>
      {messages.map((m) => (
        <div key={m.id} role={m.role}>
          {m.content}
        </div>
      ))}
      
      <form onSubmit={handleSubmit}>
        <input
          value={input}
          onChange={handleInputChange}
          placeholder="Type a message..."
          disabled={isLoading}
        />
        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Sending...' : 'Send'}
        </button>
      </form>
    </div>
  );
}
```

### Pattern 2: useCompletion Hook

```typescript
'use client';

import { useCompletion } from 'ai/react';

export function AutoComplete() {
  const { completion, input, handleInputChange, handleSubmit, isLoading } =
    useCompletion({
      api: '/api/complete',
      onFinish: (prompt, completion) => {
        console.log('Completed:', completion);
      },
    });

  return (
    <div>
      <p>Suggestion: {completion}</p>
      <form onSubmit={handleSubmit}>
        <input
          value={input}
          onChange={handleInputChange}
          placeholder="Start typing..."
        />
      </form>
    </div>
  );
}
```

### Pattern 3: Breaking Changes from v4

```typescript
// v4 style
import { useChat } from 'ai/react';
import { generateText } from 'ai';

// v5: Check for new imports and APIs
// - Some hooks may have different signatures
// - Stream handling may differ
// - Check actual version docs

// v5 typically keeps similar API but verify:
// - Message types
// - Callback signatures
// - Configuration options
```

### Pattern 4: Streaming Responses

```typescript
// API route (App Router)
import { streamText } from 'ai';

export async function POST(req: Request) {
  const { messages } = await req.json();
  
  const result = streamText({
    model: openai('gpt-4o'),
    messages,
  });
  
  return result.toDataStreamResponse();
}
```

### Pattern 5: Tool/Function Calling

```typescript
// Define tools
const tools = {
  getWeather: {
    description: 'Get weather for a location',
    parameters: z.object({
      location: z.string().describe('City name'),
    }),
  },
};

// Server-side with streamText
import { streamText, tool } from 'ai';

export async function POST(req: Request) {
  const result = streamText({
    model: openai('gpt-4o'),
    messages,
    tools: {
      getWeather: tool({
        execute: async ({ location }) => {
          return { weather: 'sunny', temp: 72 };
        },
      }),
    },
  });
  
  return result.toDataStreamResponse();
}

// Client-side - tools are auto-called
const { messages } = useChat({
  api: '/api/chat',
  // Tools called automatically when model requests them
});
```

### Pattern 6: Custom Model Provider

```typescript
import { createAI } from 'ai';

// Custom provider configuration
const ai = createAI({
  provider: 'openai', // or anthropic, google, etc.
  model: 'gpt-4o',
  // v5 config options
});
```

---

## Decision Tree

```
Feature type?
├── Chat UI → useChat
├── Autocomplete → useCompletion
├── Simple generation → generateText
└── Streaming → streamText

Data source?
├── Server Actions → Use in tools
├── API Route → streamText
└── Direct call → generateText

Need tools?
├── Define with Zod schema
├── Use tool() helper
└── Auto-execution on client
```

---

## Anti-Patterns

- ❌ Using v4 patterns without verifying v5 changes
- ❌ Not handling errors from AI responses
- ❌ Forgetting to handle loading states
- ❌ Not implementing proper message keys

---

## Resources

- **AI SDK Docs**: https://sdk.vercel.ai/
- **v5 Changelog**: Check vercel/ai changelog
- **Providers**: https://sdk.vercel.ai/docs/ai-core/providers