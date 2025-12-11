# Design: Replace Frontend with shadcn-admin

## Context

The dental prosthesis laboratory platform requires a professional admin dashboard for managing laboratories, clients, orders, prostheses, and technicians. The current frontend is a minimal Next.js scaffold that would require extensive development to reach production quality.

### Stakeholders
- Developers: Faster development with pre-built components
- Laboratory staff: Professional, intuitive interface
- Business: Reduced time-to-market

### Constraints
- Must integrate with existing Go backend (REST API)
- Must use Clerk for authentication (already configured)
- Must support multi-tenant architecture (laboratory isolation)

## Goals / Non-Goals

### Goals
- Replace frontend with production-ready admin dashboard foundation
- Maintain Clerk authentication integration
- Enable rapid development of domain-specific features
- Provide modern, responsive, accessible UI

### Non-Goals
- SSR/SEO optimization (not needed for authenticated admin dashboard)
- Preserving Next.js ecosystem (clean break acceptable)
- Custom design system (leverage shadcn-admin's existing design)

## Decisions

### Decision 1: Use Vite instead of Next.js

**Chosen**: Vite + React (SPA)

**Rationale**:
- shadcn-admin is built on Vite, maintaining compatibility reduces effort
- Admin dashboards don't benefit significantly from SSR
- Faster development builds with Vite
- Simpler deployment (static files)

**Alternatives considered**:
- Port shadcn-admin to Next.js: High effort, defeats purpose of using template
- Keep Next.js, copy components only: Misses layout/routing benefits

### Decision 2: TanStack Router for routing

**Chosen**: TanStack Router (comes with template)

**Rationale**:
- Type-safe routing
- Already integrated in shadcn-admin
- Modern React patterns (suspense, loaders)
- File-based route generation available

**Alternatives considered**:
- React Router: Would require migration effort from template
- Wouter: Less features for complex admin app

### Decision 3: Keep Clerk authentication

**Chosen**: Retain Clerk (@clerk/clerk-react for Vite)

**Rationale**:
- Backend already validates Clerk JWTs
- shadcn-admin already uses Clerk
- No authentication migration needed
- Consistent user management

### Decision 4: API client approach

**Chosen**: Custom fetch-based API client with Clerk token injection

**Rationale**:
- Simple, no additional dependencies
- Full control over request/response handling
- Matches current frontend pattern

**Implementation**:
```typescript
// lib/api-client.ts
import { useAuth } from '@clerk/clerk-react';

export const createApiClient = (getToken: () => Promise<string | null>) => ({
  async fetch(endpoint: string, options?: RequestInit) {
    const token = await getToken();
    return fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        ...options?.headers,
      },
    });
  },
});
```

## Architecture

### Project Structure (adapted from shadcn-admin)

```
frontend/
├── public/
│   └── images/
├── src/
│   ├── components/
│   │   ├── ui/              # shadcn components
│   │   ├── layout/          # Sidebar, header, etc.
│   │   └── [feature]/       # Domain components
│   ├── routes/              # TanStack Router routes
│   │   ├── _authenticated/  # Protected routes
│   │   │   ├── dashboard/
│   │   │   ├── laboratories/
│   │   │   ├── clients/
│   │   │   ├── orders/
│   │   │   ├── prostheses/
│   │   │   └── technicians/
│   │   └── _public/         # Public routes (login)
│   ├── lib/
│   │   ├── api-client.ts    # Backend API client
│   │   └── utils.ts
│   ├── hooks/               # Custom React hooks
│   ├── types/               # TypeScript types
│   ├── services/            # API service layer
│   └── main.tsx
├── .env.example
├── package.json
├── vite.config.ts
├── tsconfig.json
└── tailwind.config.ts
```

### Route Structure

| Route | Purpose |
|-------|---------|
| `/` | Dashboard with lab metrics |
| `/laboratories` | Laboratory management |
| `/clients` | Client (dental clinic) management |
| `/orders` | Order tracking and management |
| `/prostheses` | Prosthesis catalog |
| `/technicians` | Technician management |
| `/settings` | User/lab settings |

## Risks / Trade-offs

### Risk 1: Learning curve for TanStack Router
- **Impact**: Medium
- **Mitigation**: Good documentation, type-safe API reduces errors

### Risk 2: Template divergence
- **Impact**: Low
- **Mitigation**: Fork template, don't track upstream closely

### Risk 3: Missing SSR capabilities
- **Impact**: Low (admin dashboard)
- **Mitigation**: Not needed for authenticated B2B app

## Migration Plan

### Phase 1: Setup (Day 1)
1. Clone/fork shadcn-admin template
2. Remove existing `frontend/` directory
3. Configure Clerk environment variables
4. Verify authentication works

### Phase 2: Integration (Day 1-2)
1. Set up API client for Go backend
2. Configure CORS on backend if needed
3. Test authenticated API calls

### Phase 3: Domain Pages (Day 2-5)
1. Create laboratory management pages
2. Create client management pages
3. Create order management pages
4. Create prosthesis catalog pages
5. Create technician management pages
6. Customize dashboard with relevant metrics

### Phase 4: Polish (Day 5-7)
1. Customize branding/theming
2. Add domain-specific components
3. Implement global search for entities
4. Testing and bug fixes

### Rollback Plan
- Keep backup of current frontend before deletion
- Git history preserves all previous code
- Can restore Next.js version if critical issues found

## Open Questions

1. **Branding**: Should we customize the color scheme immediately or keep default?
   - Recommendation: Keep default initially, customize after core features work

2. **API error handling**: Standardize error display patterns?
   - Recommendation: Use shadcn-admin's toast/sonner patterns

3. **Form validation**: Use Zod (already in template) or something else?
   - Recommendation: Use Zod, already integrated with shadcn forms
