# Tasks: Replace Frontend with shadcn-admin

## 1. Setup

- [x] 1.1 Backup current frontend (create archive branch or tag)
- [x] 1.2 Remove existing `frontend/` directory
- [x] 1.3 Clone shadcn-admin template into `frontend/`
- [x] 1.4 Remove shadcn-admin's `.git` directory
- [x] 1.5 Update `package.json` name to `dental-prosthesis-frontend`
- [x] 1.6 Install dependencies with pnpm
- [x] 1.7 Verify dev server starts successfully

## 2. Authentication Integration

- [x] 2.1 Configure Clerk environment variables (`.env.local.example`)
- [x] 2.2 Add ClerkProvider to root (main.tsx)
- [x] 2.3 Update sign-in/sign-up pages to use Clerk components
- [x] 2.4 Implement auth guard for protected routes

## 3. API Client Setup

- [x] 3.1 Create API client with Clerk token injection (`src/lib/api-client.ts`)
- [x] 3.2 Configure API base URL environment variable
- [x] 3.3 Create `useApiAuth` hook to sync Clerk token with API client
- [ ] 3.4 Test authenticated API call to backend (e.g., list laboratories)
- [ ] 3.5 Configure CORS on backend if needed

## 4. Domain Routes & Pages

### 4.1 Laboratories
- [x] 4.1.1 Create `/laboratories` route
- [ ] 4.1.2 Implement laboratory list page with data table
- [ ] 4.1.3 Implement create laboratory form
- [ ] 4.1.4 Implement edit laboratory form
- [ ] 4.1.5 Implement delete laboratory confirmation

### 4.2 Clients
- [x] 4.2.1 Create `/clients` route
- [ ] 4.2.2 Implement client list page with data table
- [ ] 4.2.3 Implement create client form
- [ ] 4.2.4 Implement edit client form
- [ ] 4.2.5 Implement delete client confirmation

### 4.3 Orders
- [x] 4.3.1 Create `/orders` route
- [ ] 4.3.2 Implement order list page with data table
- [ ] 4.3.3 Implement create order form (with client/prosthesis selection)
- [ ] 4.3.4 Implement order detail view
- [ ] 4.3.5 Implement order status workflow (received → in production → ready → delivered)
- [ ] 4.3.6 Implement edit order form
- [ ] 4.3.7 Implement cancel order confirmation

### 4.4 Prostheses
- [x] 4.4.1 Create `/prostheses` route
- [ ] 4.4.2 Implement prosthesis catalog list page
- [ ] 4.4.3 Implement create prosthesis form
- [ ] 4.4.4 Implement edit prosthesis form
- [ ] 4.4.5 Implement delete prosthesis confirmation

### 4.5 Technicians
- [x] 4.5.1 Create `/technicians` route
- [ ] 4.5.2 Implement technician list page with data table
- [ ] 4.5.3 Implement create technician form
- [ ] 4.5.4 Implement edit technician form
- [ ] 4.5.5 Implement delete technician confirmation

## 5. Dashboard

- [x] 5.1 Customize dashboard layout for dental lab metrics
- [x] 5.2 Add order status summary widget (orders by status)
- [ ] 5.3 Add recent orders widget
- [x] 5.4 Add client count metric
- [ ] 5.5 Add production timeline/calendar widget (stretch goal)

## 6. Navigation & Layout

- [x] 6.1 Update sidebar navigation with domain sections
- [x] 6.2 Configure navigation icons (Lucide)
- [x] 6.3 Update application title/branding
- [ ] 6.4 Configure global search for domain entities (stretch goal)

## 7. Types & Services

- [x] 7.1 Create TypeScript types matching backend DTOs (`src/types/`)
- [x] 7.2 Create API service layer (`src/services/`)
  - [x] 7.2.1 Laboratory service
  - [x] 7.2.2 Client service
  - [x] 7.2.3 Order service
  - [x] 7.2.4 Prosthesis service
  - [x] 7.2.5 Technician service

## 8. Testing & Polish

- [ ] 8.1 Verify all CRUD operations work end-to-end
- [ ] 8.2 Test responsive design on mobile/tablet
- [ ] 8.3 Test dark mode
- [x] 8.4 Fix any console errors/warnings (build passes, lint passes)
- [ ] 8.5 Update project.md if conventions change

## 9. Documentation

- [x] 9.1 Update frontend README with new setup instructions
- [x] 9.2 Document environment variables
- [x] 9.3 Document development workflow
