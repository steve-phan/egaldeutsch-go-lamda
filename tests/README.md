# E2E Testing with Playwright

This project uses Playwright for comprehensive end-to-end testing of the German learning platform workflow.

## Test Structure

### Test Files

- `tests/e2e/creator-workflow.spec.ts` - Tests for content creator user journey
- `tests/e2e/admin-workflow.spec.ts` - Tests for admin review and management
- `tests/e2e/complete-lifecycle.spec.ts` - Full workflow integration tests

### Fixtures

- `tests/fixtures/testData.ts` - Test data definitions (users, stories, questions, quizzes)
- `tests/fixtures/databaseSetup.ts` - Database setup and cleanup utilities

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

# Run specific test suites
npm run test:e2e:creator
npm run test:e2e:admin
npm run test:e2e:lifecycle

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
