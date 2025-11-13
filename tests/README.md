# E2E Testing with Playwright

This project uses Playwright for comprehensive end-to-end testing of the German learning platform workflow.

## Test Structure

### Test Files

#### Core Workflow Tests
- `tests/e2e/creator-workflow.spec.ts` - Tests for content creator user journey
- `tests/e2e/admin-workflow.spec.ts` - Tests for admin review and management
- `tests/e2e/complete-lifecycle.spec.ts` - Full workflow integration tests

#### Enhanced E2E Tests
- `tests/e2e/auth-workflows.spec.ts` - Comprehensive authentication and session management tests
- `tests/e2e/admin-user-management.spec.ts` - Detailed admin user management scenarios
- `tests/e2e/story-workflows.spec.ts` - Complete story creation and review workflows
- `tests/e2e/question-quiz-workflows.spec.ts` - Question and quiz management tests
- `tests/e2e/error-handling.spec.ts` - Error handling and edge case tests

#### Role-Based CRUD Tests (NEW)
- `tests/e2e/role-based-crud-workflows.spec.ts` - Comprehensive CRUD testing for all roles (Creator, Reviewer, Admin)
  - Stories: Create, update, review, publish workflows
  - Questions: Create, update, link to stories, change difficulty
  - Quizzes: Create, update settings, preserve associations
  - Permission enforcement and security boundaries
- `tests/e2e/cross-role-scenarios.spec.ts` - Multi-role collaboration and edge cases
  - Complete lifecycle: Creator → Reviewer → Creator → Admin
  - Multiple reviewers, admin overrides, audit trails
  - Concurrent operations and conflict handling
  - Security: Role permission enforcement

### Fixtures

- `tests/fixtures/testData.ts` - Test data definitions (users, stories, questions, quizzes)
- `tests/fixtures/databaseSetup.ts` - Database setup and cleanup utilities
- `tests/fixtures/authHelpers.ts` - Authentication helper functions for tests
- `tests/fixtures/contentHelpers.ts` - Content creation and management helpers

## Content Status Workflow

The platform uses a simplified 4-state workflow:

1. **Draft** - Creator can edit, not visible to public
2. **Preview** - Creator submitted for review, ready for admin
3. **Ready** - Admin approved, ready to publish
4. **Published** - Live content, visible to public

### Status Transitions

- Creator: `draft` → `preview`
- Admin: `preview` → `ready` or `preview` → `draft` (rejection)
- Admin: `ready` → `published` or `ready` → `draft` (revision needed)
- Admin: `published` → `ready` (unpublish)

## Running Tests

### Prerequisites

1. Start the development servers:

   ```bash
   npm run develop  # Gatsby frontend (localhost:8000)
   netlify dev --functions-port 8888  # Backend functions
   ```

2. Ensure MongoDB is running and accessible

### Test Commands

```bash
# Run all E2E tests
npm run test:e2e

# Run tests with browser UI visible
npm run test:e2e:headed

# Run tests in debug mode
npm run test:e2e:debug

# Run specific test suites (original)
npm run test:e2e:creator
npm run test:e2e:admin
npm run test:e2e:lifecycle

# Run specific enhanced test suites
npm run test:e2e:auth           # Authentication workflows
npm run test:e2e:user-mgmt      # Admin user management
npm run test:e2e:stories        # Story workflows
npm run test:e2e:questions      # Question and quiz workflows
npm run test:e2e:errors         # Error handling
npm run test:e2e:roles          # Role-based CRUD operations (NEW)
npm run test:e2e:cross-role     # Cross-role collaboration scenarios (NEW)

# Run tests with Playwright UI mode
npm run test:e2e:ui

# View test report
npm run test:report
```

### Test Configuration

Tests are configured in `playwright.config.ts`:

- Multi-browser testing (Chrome, Firefox, Safari)
- Mobile viewport testing
- Automatic server startup
- Trace collection on failures
- HTML reporting

## Test Scenarios

### Authentication & Session Management Tests (`auth-workflows.spec.ts`)

- ✅ Complete login flow with form validation
- ✅ Empty form validation
- ✅ Invalid email format handling
- ✅ Incorrect credentials error handling
- ✅ Session persistence across page refreshes
- ✅ Session expiration and auto-logout
- ✅ Logout clears session and redirects
- ✅ Multi-tab session synchronization
- ✅ Login with different user roles
- ✅ Concurrent login attempts handling
- ✅ Password visibility toggle

### Admin User Management Tests (`admin-user-management.spec.ts`)

- ✅ View all users with filtering and search
- ✅ Role filter functionality
- ✅ User search by name or email
- ✅ Change user roles with audit trail
- ✅ Role change comment tracking
- ✅ Activate/deactivate user accounts
- ✅ View user activity history
- ✅ User statistics on admin dashboard
- ✅ Non-admin access prevention

### Story Workflow Tests (`story-workflows.spec.ts`)

- ✅ Complete story creation with validation
- ✅ Form validation error handling
- ✅ Story submission for review
- ✅ Review cycle with detailed comments
- ✅ Story rejection and revision cycle
- ✅ Multiple creators working simultaneously
- ✅ Content length validation
- ✅ Level-appropriate vocabulary checking
- ✅ Topic management

### Question & Quiz Tests (`question-quiz-workflows.spec.ts`)

- ✅ Create multiple choice questions
- ✅ Create open text questions
- ✅ Create multiple select questions
- ✅ Link questions to stories
- ✅ Auto-generate quizzes from questions
- ✅ Configure quiz difficulty balance
- ✅ Admin review and publish quiz
- ✅ Quiz preview functionality
- ✅ Quiz takes and scoring
- ✅ Question validation

### Error Handling & Edge Cases (`error-handling.spec.ts`)

- ✅ API failure graceful handling
- ✅ Network timeout scenarios
- ✅ Concurrent status update conflicts
- ✅ Form submission error handling
- ✅ Invalid input data validation
- ✅ XSS attack prevention
- ✅ Session timeout during form filling
- ✅ Browser back/forward navigation
- ✅ Network disconnection handling
- ✅ Rapid consecutive submissions
- ✅ Malformed API response handling

### Creator Workflow Tests

- ✅ User authentication and dashboard access
- ✅ Story creation with validation
- ✅ Content submission for review (draft → preview)
- ✅ Status tracking and history
- ✅ Role-based access control
- ✅ Form validation and error handling

### Admin Workflow Tests

- ✅ Admin dashboard access and navigation
- ✅ Content filtering and search
- ✅ Content approval (preview → ready → published)
- ✅ Content rejection with feedback (preview → draft)
- ✅ User role management
- ✅ Dashboard statistics and analytics

### Complete Lifecycle Tests

- ✅ End-to-end content creation and publication
- ✅ Rejection and revision cycles
- ✅ Multi-user concurrent workflows
- ✅ Public content accessibility
- ✅ Audit trail and history tracking

## Test Data Management

### Fixture Data

The test suite uses comprehensive fixture data including:

- **Users**: Creator, Admin, Reviewer roles with authentication
- **Stories**: Content in all workflow states (draft, preview, ready, published)
- **Questions**: Multiple choice, multiple select, and open text questions
- **Quizzes**: Auto-generated quizzes with balanced difficulty

### Database Setup

- Tests use `TestDatabaseSetup` class for data lifecycle management
- Automatic creation of test data before test runs
- Automatic cleanup after test completion
- Isolated test environments prevent data conflicts

### Mock Data Structure

```typescript
{
  users: { creator, admin, reviewer },
  stories: { draft, preview, ready, published },
  questions: { easy: [...], medium: [...], hard: [...] },
  quizzes: { draft, preview, ready, published }
}
```

## Best Practices

### Test Isolation

- Each test suite creates fresh test data
- Browser contexts are isolated between user roles
- Database cleanup ensures no test pollution

### Real User Simulation

- Tests use actual form interactions, not direct API calls
- Multi-browser testing ensures cross-platform compatibility
- Mobile testing verifies responsive design

### Comprehensive Coverage

- Tests cover happy path and error scenarios
- Role-based access control verification
- Form validation and user feedback testing
- Concurrent user workflow testing

## Troubleshooting

### Common Issues

1. **Server not starting**: Ensure both Gatsby and Netlify dev servers are running
2. **Database connection**: Verify MongoDB is accessible and credentials are correct
3. **Test timeouts**: Increase timeout values in playwright.config.ts if needed
4. **Browser issues**: Run `npx playwright install` to update browsers

### Debug Mode

Use `npm run test:e2e:debug` to:

- Step through tests interactively
- Inspect page state at each step
- Debug failing assertions
- View network requests and responses

### Visual Testing

Use `npm run test:e2e:ui` for:

- Visual test runner interface
- Real-time test execution viewing
- Easy test debugging and analysis
- Test result comparison

## CI/CD Integration

The test suite is designed for CI/CD integration:

- Headless execution by default
- HTML report generation
- Screenshots and videos on failure
- Parallel test execution
- Cross-browser compatibility testing

For CI environments, use:

```bash
npm run test:e2e --reporter=junit
```

## Extending Tests

### Adding New Test Cases

1. Create new spec files in `tests/e2e/`
2. Use existing fixtures or create new ones in `tests/fixtures/`
3. Follow existing patterns for browser context management
4. Add test scripts to package.json

### Custom Fixtures

1. Add test data to `testData.ts`
2. Update `databaseSetup.ts` for data creation/cleanup
3. Ensure proper relationships between entities
4. Maintain data consistency across test runs

## Enhanced Test Coverage Summary

### New Test Files Added

1. **`auth-workflows.spec.ts`** (8 tests)
   - Complete authentication flows with validation
   - Session management and persistence
   - Multi-tab synchronization
   - Concurrent login handling

2. **`admin-user-management.spec.ts`** (8 tests)
   - User filtering and search
   - Role management with audit trails
   - Account activation/deactivation
   - Access control verification

3. **`story-workflows.spec.ts`** (5 tests)
   - End-to-end story creation
   - Review cycles with detailed feedback
   - Rejection and revision workflows
   - Concurrent creator workflows

4. **`question-quiz-workflows.spec.ts`** (6 tests)
   - Multiple question types (multiple choice, open text, multiple select)
   - Quiz auto-generation
   - Quiz review and publication
   - Quiz taking and scoring

5. **`error-handling.spec.ts`** (10 tests)
   - API failure recovery
   - Network timeout handling
   - Concurrent operations
   - Form validation edge cases
   - Security testing (XSS prevention)

### Helper Functions Added

- **`authHelpers.ts`**: Reusable authentication functions
  - `loginAs()`, `logout()`, `isLoggedIn()`, `isLoggedOut()`
  
- **`contentHelpers.ts`**: Content creation utilities
  - `fillStoryForm()`, `fillQuestionForm()`
  - `createAndSubmitStory()`, `updateContentStatus()`
  - `navigateToReviews()`

### Test Coverage Statistics

- **Total Enhanced Tests**: 37 new test scenarios
- **Original Tests**: 16 test scenarios (creator, admin, lifecycle)
- **Combined Total**: 53+ comprehensive test scenarios
- **Test Files**: 8 test specification files
- **Helper Modules**: 2 utility files

### Key Features Tested

✅ **Authentication**: Login, logout, session management, concurrent access  
✅ **User Management**: Role changes, filtering, search, audit trails  
✅ **Content Creation**: Stories, questions, quizzes with full validation  
✅ **Review Workflows**: Approval, rejection, revision cycles  
✅ **Error Handling**: API failures, network issues, race conditions  
✅ **Security**: XSS prevention, session timeout, malformed responses  
✅ **Concurrent Operations**: Multiple users, race conditions  
✅ **Edge Cases**: Invalid input, rapid submissions, offline mode  

### Browser & Mobile Coverage

All tests run across:
- ✅ Desktop Chrome
- ✅ Desktop Firefox
- ✅ Desktop Safari (WebKit)
- ✅ Mobile Chrome (Pixel 5)
- ✅ Mobile Safari (iPhone 12)

### Next Steps

To run the enhanced test suite:

1. Ensure development servers are running
2. Run `npm run test:e2e` for all tests
3. Run specific suites with `npm run test:e2e:auth`, etc.
4. View detailed reports with `npm run test:report`

For CI/CD integration, tests are configured for:
- Headless execution
- Parallel test execution
- Screenshot and video capture on failure
- HTML and JUnit reporting
