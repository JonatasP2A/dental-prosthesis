# Dental Prosthesis Platform - Frontend

Next.js frontend application with Clerk authentication for the Dental Prosthesis Platform.

## Getting Started

### Prerequisites

- Node.js 18+ installed
- Clerk account (sign up at https://clerk.com)

### Installation

1. Install dependencies:

```bash
npm install
```

2. Set up environment variables:

Copy `.env.local.example` to `.env.local`:

```bash
cp .env.local.example .env.local
```

3. Get your Clerk keys:

- Go to [Clerk Dashboard API Keys](https://dashboard.clerk.com/last-active?path=api-keys)
- Copy your **Publishable Key** and **Secret Key**
- Paste them into `.env.local`:

```bash
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_...
CLERK_SECRET_KEY=sk_test_...
```

4. Run the development server:

```bash
npm run dev
```

5. Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

```
frontend/
├── app/                    # Next.js App Router pages
│   ├── layout.tsx         # Root layout with ClerkProvider
│   ├── page.tsx           # Home page
│   └── globals.css        # Global styles
├── proxy.ts               # Clerk middleware
├── package.json
├── tsconfig.json
└── tailwind.config.ts     # Tailwind CSS configuration
```

## Features

- ✅ Clerk authentication integration
- ✅ Sign in/Sign up functionality
- ✅ Protected routes
- ✅ User profile management
- ✅ Responsive design with Tailwind CSS

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm start` - Start production server
- `npm run lint` - Run ESLint

## Clerk Integration

This project uses Clerk's Next.js SDK with the App Router approach:

- **Middleware**: `proxy.ts` uses `clerkMiddleware()` from `@clerk/nextjs/server`
- **Provider**: `<ClerkProvider>` wraps the app in `app/layout.tsx`
- **Components**: Uses Clerk's React components (`SignInButton`, `SignUpButton`, `UserButton`, etc.)

## Environment Variables

Required environment variables (in `.env.local`):

- `NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY` - Your Clerk publishable key
- `CLERK_SECRET_KEY` - Your Clerk secret key
- `NEXT_PUBLIC_API_URL` - Backend API URL (optional, defaults to http://localhost:8080)

## Learn More

- [Next.js Documentation](https://nextjs.org/docs)
- [Clerk Documentation](https://clerk.com/docs)
- [Clerk Next.js Quickstart](https://clerk.com/docs/quickstarts/nextjs)
