# Dental Lab Pro - Frontend

Admin dashboard for dental prosthesis laboratory management, built on [shadcn-admin](https://github.com/satnaing/shadcn-admin).

## Features

- 🔐 **Clerk Authentication** - Secure sign-in/sign-up with Clerk
- 🏠 **Dashboard** - Overview of laboratory metrics and order status
- 🏢 **Laboratory Management** - Manage dental prosthesis laboratories
- 👥 **Client Management** - Track dental clinic clients
- 📋 **Order Management** - Track prosthesis orders through production workflow
- 🦷 **Prosthesis Catalog** - Manage prosthesis types and materials
- 👷 **Technician Management** - Manage laboratory staff
- 🌙 **Dark Mode** - Light and dark theme support
- 📱 **Responsive** - Works on desktop, tablet, and mobile

## Tech Stack

- **Framework**: [Vite](https://vitejs.dev/) + [React 19](https://react.dev/)
- **Routing**: [TanStack Router](https://tanstack.com/router)
- **Data Fetching**: [TanStack Query](https://tanstack.com/query)
- **UI Components**: [shadcn/ui](https://ui.shadcn.com/)
- **Styling**: [Tailwind CSS v4](https://tailwindcss.com/)
- **Authentication**: [Clerk](https://clerk.com/)
- **Form Handling**: [React Hook Form](https://react-hook-form.com/) + [Zod](https://zod.dev/)
- **HTTP Client**: [Axios](https://axios-http.com/)

## Getting Started

### Prerequisites

- Node.js 20.19+ or 22.12+
- pnpm (recommended)

### Installation

1. Install dependencies:

```bash
pnpm install
```

2. Create environment file:

```bash
cp .env.local.example .env.local
```

3. Configure environment variables in `.env.local`:

```env
# Clerk Authentication (required)
VITE_CLERK_PUBLISHABLE_KEY=pk_test_your_clerk_publishable_key

# Backend API (required)
VITE_API_BASE_URL=http://localhost:8080
```

4. Start the development server:

```bash
pnpm dev
```

5. Open [http://localhost:5173](http://localhost:5173) in your browser.

## Project Structure

```
src/
├── assets/          # Static assets and icons
├── components/
│   ├── ui/          # shadcn/ui components
│   ├── layout/      # Layout components (sidebar, header)
│   └── ...          # Shared components
├── context/         # React context providers
├── features/        # Feature modules
│   ├── dashboard/
│   ├── laboratories/
│   ├── clients/
│   ├── orders/
│   ├── prostheses/
│   └── technicians/
├── hooks/           # Custom React hooks
├── lib/             # Utilities and API client
├── routes/          # TanStack Router routes
├── services/        # API service layer
├── stores/          # Zustand stores
├── styles/          # Global styles
└── types/           # TypeScript type definitions
```

## Available Scripts

| Command | Description |
|---------|-------------|
| `pnpm dev` | Start development server |
| `pnpm build` | Build for production |
| `pnpm preview` | Preview production build |
| `pnpm lint` | Run ESLint |
| `pnpm format` | Format code with Prettier |
| `pnpm format:check` | Check code formatting |

## API Integration

The frontend communicates with the Go backend via REST API. Authentication tokens from Clerk are automatically injected into API requests.

### Services

- `laboratoryService` - Laboratory CRUD operations
- `clientService` - Client CRUD operations  
- `orderService` - Order management and status updates
- `prosthesisService` - Prosthesis catalog operations
- `technicianService` - Technician CRUD operations

## Customization

### Adding New Routes

1. Create a new route file in `src/routes/_authenticated/`
2. Create a feature component in `src/features/`
3. Add navigation item in `src/components/layout/data/sidebar-data.ts`
4. Run `pnpm dev` to regenerate the route tree

### Theming

Theme configuration is in `src/styles/theme.css`. The app uses Tailwind CSS with CSS custom properties for theming.

## Credits

Based on [shadcn-admin](https://github.com/satnaing/shadcn-admin) by [@satnaing](https://github.com/satnaing).

## License

MIT
