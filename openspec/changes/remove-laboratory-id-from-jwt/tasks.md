## 1. Authentication Middleware Changes
- [ ] 1.1 Remove `laboratory_id` extraction from JWT custom claims parsing
- [ ] 1.2 Remove `laboratory_id` validation check in Authenticate middleware
- [ ] 1.3 Remove `LaboratoryIDKey` from context injection
- [ ] 1.4 Update `Claims` struct to remove `LaboratoryID` field (or make optional)
- [ ] 1.5 Update `CustomClaims` struct to remove `LaboratoryID` field
- [ ] 1.6 Remove or deprecate `GetLaboratoryID` function (or make it return empty string)

## 2. Handler Updates
- [ ] 2.1 Update `order.go` handlers to extract `laboratory_id` from query parameter (`c.Query("laboratory_id")`)
- [ ] 2.2 Update `client.go` handlers to extract `laboratory_id` from query parameter
- [ ] 2.3 Update `laboratory.go` handlers to extract `laboratory_id` from query parameter (if needed)
- [ ] 2.4 Add validation in all handlers to ensure `laboratory_id` query parameter is present and non-empty
- [ ] 2.5 Return HTTP 400 Bad Request if `laboratory_id` query parameter is missing
- [ ] 2.6 Replace all `auth.GetLaboratoryID()` calls with query parameter extraction

## 3. Domain & Application Layer
- [ ] 3.1 Review domain entities (Client, Order) - determine if laboratory_id should come from different source
- [ ] 3.2 Update application services to handle laboratory context differently
- [ ] 3.3 Ensure multi-tenancy is still enforced through alternative mechanism

## 4. Testing Review & Updates
- [ ] 4.1 Review all existing test files to identify tests that need updates
  - [ ] 4.1.1 Review `backend/pkg/auth/clerk_test.go` - Remove laboratory_id validation tests
  - [ ] 4.1.2 Review `backend/internal/adapters/inbound/http/handler/client_test.go` - Update all test requests
  - [ ] 4.1.3 Review `backend/internal/adapters/inbound/http/handler/order_test.go` - Update all test requests
  - [ ] 4.1.4 Review `backend/internal/adapters/inbound/http/handler/laboratory_test.go` - Update if needed
  - [ ] 4.1.5 Review all integration tests in `backend/test/` directory
- [ ] 4.2 Update authentication middleware tests to remove laboratory_id validation
- [ ] 4.3 Update all handler tests to include `laboratory_id` query parameter in test requests
  - [ ] 4.3.1 Update test helper functions that create requests
  - [ ] 4.3.2 Update all HTTP test requests (GET, POST, PUT, PATCH, DELETE)
  - [ ] 4.3.3 Ensure test requests include `?laboratory_id=test-lab-id` or similar
- [ ] 4.4 Add new test cases for query parameter validation
  - [ ] 4.4.1 Add tests for missing `laboratory_id` query parameter (should return 400)
  - [ ] 4.4.2 Add tests for empty `laboratory_id` query parameter (should return 400)
  - [ ] 4.4.3 Add tests for invalid `laboratory_id` format (if applicable)
- [ ] 4.5 Review and update integration tests
  - [ ] 4.5.1 Update integration test requests to include `laboratory_id` query parameter
  - [ ] 4.5.2 Verify laboratory isolation still works correctly in integration tests
- [ ] 4.6 Verify test coverage
  - [ ] 4.6.1 Run test coverage report to ensure no regressions
  - [ ] 4.6.2 Ensure all test cases are updated and passing
- [ ] 4.7 Run full test suite to ensure all tests pass with new query parameter approach

## 5. Documentation
- [ ] 5.1 Update backend README to remove laboratory_id from JWT claims documentation
- [ ] 5.2 Update backend README to document `laboratory_id` query parameter requirement
- [ ] 5.3 Update API documentation (client.http files) to include `laboratory_id` query parameter in all examples
- [ ] 5.4 Update frontend documentation about JWT setup changes
- [ ] 5.5 Document breaking change: all API requests now require `laboratory_id` query parameter
