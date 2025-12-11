# Change: Replace Frontend with shadcn-admin Template

## Why

The current Next.js frontend is a minimal scaffold with only placeholder pages and basic Clerk authentication. Building a professional admin dashboard from scratch would require significant effort. The [shadcn-admin](https://github.com/satnaing/shadcn-admin) template provides a production-ready admin dashboard UI that aligns with our existing tech choices (shadcn/ui, Clerk, TypeScript, Tailwind) and accelerates development.

## What Changes

- **BREAKING**: Remove existing Next.js frontend completely
- Replace with Vite + React SPA based on shadcn-admin template
- Switch routing from Next.js App Router to TanStack Router
- Retain Clerk authentication (shadcn-admin already uses Clerk)
- Configure API client to communicate with existing Go backend
- Adapt dashboard for dental prosthesis laboratory domain

## Impact

- **Affected specs**: None currently (no frontend specs exist)
- **New specs**: `frontend` capability spec will be created
- **Affected code**:
  - `frontend/` - Complete replacement
  - Backend remains unchanged (REST API compatible with any client)
- **Breaking changes**:
  - No SSR capabilities (acceptable for B2B admin dashboard)
  - Different routing paradigm (TanStack Router vs Next.js)
  - Development workflow changes (Vite vs Next.js dev server)

## Benefits

- 10+ ready-made pages (dashboard, settings, error pages)
- Global search command (⌘K)
- Built-in sidebar navigation
- Light/dark mode
- Responsive and accessible design
- RTL support
- Form components, dialogs, tables pre-built

## Risks

- Learning curve for TanStack Router
- Loss of SSR (mitigated: not needed for admin dashboard)
- Template customization effort
