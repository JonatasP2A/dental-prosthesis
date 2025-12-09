# How to Access Clerk JWT Token in Browser

This guide explains different ways to access the Clerk JWT session token for making authenticated API calls to your backend.

## Methods to Access JWT Token

### 1. **Server Components** (Recommended for Server-Side)

Use `auth()` from `@clerk/nextjs/server`:

```typescript
import { auth } from '@clerk/nextjs/server';

export default async function MyServerComponent() {
  const { getToken } = await auth();
  const token = await getToken();
  
  // Use token in API calls
  const response = await fetch('http://localhost:8080/api/v1/laboratories', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}
```

**Location:** `app/lib/api-client.ts` provides `authenticatedFetch()` helper function.

### 2. **Client Components** (For Client-Side)

Use `useAuth()` hook from `@clerk/nextjs`:

```typescript
'use client';

import { useAuth } from '@clerk/nextjs';

export default function MyClientComponent() {
  const { getToken } = useAuth();
  
  const handleApiCall = async () => {
    const token = await getToken();
    
    const response = await fetch('http://localhost:8080/api/v1/laboratories', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  };
  
  return <button onClick={handleApiCall}>Call API</button>;
}
```

**Location:** `components/client-api-example.tsx` shows a complete example.

### 3. **Browser Console** (For Debugging Only)

You can inspect the token in browser DevTools:

1. Open DevTools (F12)
2. Go to **Application** tab → **Cookies** → `http://localhost:3000`
3. Look for the `__session` cookie
4. The cookie value contains the session token (but it's encoded)

**⚠️ Note:** Don't manually extract tokens from cookies. Always use Clerk's API (`getToken()`).

### 4. **API Routes** (Next.js API Routes)

```typescript
import { auth } from '@clerk/nextjs/server';
import { NextResponse } from 'next/server';

export async function GET() {
  const { getToken } = await auth();
  const token = await getToken();
  
  // Use token to call your backend
  const response = await fetch('http://localhost:8080/api/v1/laboratories', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  
  const data = await response.json();
  return NextResponse.json(data);
}
```

## Helper Functions Provided

### `lib/api-client.ts`

- **`authenticatedFetch(endpoint, options)`** - Automatically adds JWT token to requests
- **`getAuthToken()`** - Returns the JWT token string

### Usage Example:

```typescript
import { authenticatedFetch } from '@/lib/api-client';

// In a Server Component
const response = await authenticatedFetch('/api/v1/laboratories');
const data = await response.json();
```

## Example Pages

1. **`/api-example`** - Server Component example showing token access
2. **`components/client-api-example.tsx`** - Client Component example

## Important Notes

1. **Token Format**: The token returned by `getToken()` is a JWT that starts with `eyJ` (base64 encoded)
2. **Token Expiration**: Tokens expire automatically. Clerk handles refresh.
3. **Security**: Never log or expose tokens in client-side code or console
4. **Backend**: Your backend expects the token in the `Authorization: Bearer <token>` header

## Testing in Browser Console

If you want to test getting the token in the browser console:

```javascript
// This only works if you're on a page that uses Clerk
// And you need to be signed in
// Note: This is for debugging only!

// In browser console (after signing in):
document.cookie.split('; ').find(row => row.startsWith('__session='))?.split('=')[1]
```

But again, **always prefer using Clerk's `getToken()` API** instead of manually accessing cookies.

## Backend Integration

Your backend expects the JWT token in this format:

```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Important Breaking Change**: The backend no longer extracts `laboratory_id` from JWT tokens. Instead, you must include `laboratory_id` as a query parameter in all API requests that require laboratory context.

The backend will:
1. Extract the token from the Authorization header
2. Verify it using Clerk SDK
3. Extract user ID (`sub` claim) from the token
4. Require `laboratory_id` as a query parameter for client and order endpoints

### Example API Calls

```typescript
// List clients for a laboratory
const response = await authenticatedFetch(
  `/api/v1/clients?laboratory_id=${laboratoryId}`
);

// Create an order
const response = await authenticatedFetch(
  `/api/v1/orders?laboratory_id=${laboratoryId}`,
  {
    method: 'POST',
    body: JSON.stringify({
      client_id: 'client-123',
      prosthesis: [...]
    })
  }
);
```

**Note**: Laboratory endpoints (`/api/v1/laboratories`) do not require the `laboratory_id` query parameter.
