## ADDED Requirements

### Requirement: Admin Dashboard UI

The frontend SHALL provide a professional admin dashboard interface for dental prosthesis laboratory management built on shadcn-admin template.

#### Scenario: Dashboard loads for authenticated user
- **WHEN** an authenticated user navigates to the root URL
- **THEN** the dashboard displays with laboratory metrics and recent activity

#### Scenario: Dashboard redirects unauthenticated user
- **WHEN** an unauthenticated user navigates to any protected route
- **THEN** they are redirected to the sign-in page

### Requirement: Sidebar Navigation

The frontend SHALL provide a sidebar navigation component with links to all domain sections.

#### Scenario: Navigation sections displayed
- **WHEN** a user views the sidebar
- **THEN** they see navigation links for Dashboard, Laboratories, Clients, Orders, Prostheses, and Technicians

#### Scenario: Active route highlighted
- **WHEN** a user is on a specific route
- **THEN** the corresponding sidebar item is visually highlighted

### Requirement: Laboratory Management Pages

The frontend SHALL provide pages for managing laboratories with full CRUD operations.

#### Scenario: List laboratories
- **WHEN** a user navigates to `/laboratories`
- **THEN** they see a data table with all laboratories they have access to

#### Scenario: Create laboratory
- **WHEN** a user submits the create laboratory form with valid data
- **THEN** the laboratory is created and appears in the list

#### Scenario: Edit laboratory
- **WHEN** a user edits a laboratory and saves changes
- **THEN** the laboratory is updated with the new information

#### Scenario: Delete laboratory
- **WHEN** a user confirms laboratory deletion
- **THEN** the laboratory is removed from the system

### Requirement: Client Management Pages

The frontend SHALL provide pages for managing clients (dental clinics) with full CRUD operations.

#### Scenario: List clients
- **WHEN** a user navigates to `/clients`
- **THEN** they see a data table with all clients for their laboratory

#### Scenario: Create client
- **WHEN** a user submits the create client form with valid data
- **THEN** the client is created and appears in the list

#### Scenario: Edit client
- **WHEN** a user edits a client and saves changes
- **THEN** the client is updated with the new information

#### Scenario: Delete client
- **WHEN** a user confirms client deletion
- **THEN** the client is removed from the system

### Requirement: Order Management Pages

The frontend SHALL provide pages for managing prosthesis orders with full CRUD operations and workflow status management.

#### Scenario: List orders
- **WHEN** a user navigates to `/orders`
- **THEN** they see a data table with all orders for their laboratory

#### Scenario: Create order
- **WHEN** a user submits the create order form with client and prosthesis selection
- **THEN** the order is created with status "received"

#### Scenario: View order details
- **WHEN** a user clicks on an order
- **THEN** they see the full order details including client, prosthesis, and status history

#### Scenario: Update order status
- **WHEN** a user changes an order's status (e.g., "received" to "in production")
- **THEN** the order status is updated and reflected in the UI

#### Scenario: Cancel order
- **WHEN** a user confirms order cancellation
- **THEN** the order is marked as cancelled

### Requirement: Prosthesis Catalog Pages

The frontend SHALL provide pages for managing the prosthesis catalog with full CRUD operations.

#### Scenario: List prostheses
- **WHEN** a user navigates to `/prostheses`
- **THEN** they see a data table with all prosthesis types in the catalog

#### Scenario: Create prosthesis
- **WHEN** a user submits the create prosthesis form with valid data
- **THEN** the prosthesis is added to the catalog

#### Scenario: Edit prosthesis
- **WHEN** a user edits a prosthesis and saves changes
- **THEN** the prosthesis is updated with the new information

#### Scenario: Delete prosthesis
- **WHEN** a user confirms prosthesis deletion
- **THEN** the prosthesis is removed from the catalog

### Requirement: Technician Management Pages

The frontend SHALL provide pages for managing laboratory technicians with full CRUD operations.

#### Scenario: List technicians
- **WHEN** a user navigates to `/technicians`
- **THEN** they see a data table with all technicians for their laboratory

#### Scenario: Create technician
- **WHEN** a user submits the create technician form with valid data
- **THEN** the technician is created and appears in the list

#### Scenario: Edit technician
- **WHEN** a user edits a technician and saves changes
- **THEN** the technician is updated with the new information

#### Scenario: Delete technician
- **WHEN** a user confirms technician deletion
- **THEN** the technician is removed from the system

### Requirement: Dark Mode Support

The frontend SHALL support light and dark color themes.

#### Scenario: Toggle dark mode
- **WHEN** a user toggles the theme setting
- **THEN** the UI switches between light and dark mode

#### Scenario: Persist theme preference
- **WHEN** a user sets a theme preference
- **THEN** the preference is persisted across sessions

### Requirement: Responsive Design

The frontend SHALL be responsive and usable on desktop, tablet, and mobile devices.

#### Scenario: Mobile navigation
- **WHEN** a user views the app on a mobile device
- **THEN** the sidebar collapses into a mobile-friendly menu

#### Scenario: Data tables on mobile
- **WHEN** a user views data tables on a mobile device
- **THEN** the tables are scrollable or adapt to the screen size

### Requirement: API Integration

The frontend SHALL communicate with the Go backend via authenticated REST API calls.

#### Scenario: API calls include authentication
- **WHEN** the frontend makes an API request
- **THEN** the request includes the Clerk JWT token in the Authorization header

#### Scenario: Handle API errors
- **WHEN** an API call fails
- **THEN** the user sees an appropriate error message (toast notification)

#### Scenario: Handle loading states
- **WHEN** data is being fetched
- **THEN** the UI displays appropriate loading indicators
