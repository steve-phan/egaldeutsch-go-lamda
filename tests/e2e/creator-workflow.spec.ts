import { test, expect } from "@playwright/test";
import { TestDatabaseSetup } from "../fixtures/databaseSetup";

let testSetup: TestDatabaseSetup;

test.describe("Creator Workflow", () => {
  test.beforeAll(async () => {
    // Setup test database with fixture data
    testSetup = new TestDatabaseSetup();
    await testSetup.setupTestData();
  });

  test.afterAll(async () => {
    // Cleanup test data after all tests
    if (testSetup) {
      await testSetup.cleanupTestData();
    }
  });

  test("Creator can login and access dashboard", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    // Navigate to login page
    await page.goto("/auth/login");

    // Verify login page loads
    await expect(page).toHaveTitle(/EgalDeutsch/);
    await expect(page.locator("h1")).toContainText("Login");

    // Fill login form
    await page.fill('input[type="email"]', creator.email);
    await page.fill('input[type="password"]', creator.password);

    // Submit login
    await page.click('button[type="submit"]');

    // Wait for redirect to dashboard
    await page.waitForURL("/");

    // Verify successful login
    await expect(page.locator("text=Welcome")).toBeVisible();
    await expect(page.locator(`text=${creator.name}`)).toBeVisible();
  });

  test("Creator can create a new story", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    // Login first
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', creator.email);
    await page.fill('input[type="password"]', creator.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Navigate to create story page (assuming it exists)
    await page.goto("/create/story");

    // Verify create story page
    await expect(page.locator("h1")).toContainText("Create Story");

    // Fill story form
    const testStory = {
      title: "E2E Test Story - Creator Workflow",
      content:
        "Dies ist eine Testgeschichte für unsere E2E-Tests. Anna geht spazieren und trifft einen Freund. Sie unterhalten sich über das Wetter und planen einen Ausflug. Am Ende sind beide glücklich und freuen sich auf das Wochenende. Diese Geschichte enthält mehr als hundert Zeichen und ist perfekt für unsere Tests geeignet.",
      level: "A1",
      topics: ["Test", "E2E", "Automatisierung"],
    };

    await page.fill('input[name="title"]', testStory.title);
    await page.fill('textarea[name="content"]', testStory.content);
    await page.selectOption('select[name="level"]', testStory.level);

    // Add topics (assuming multi-select or tag input)
    for (const topic of testStory.topics) {
      await page.fill('input[name="topics"]', topic);
      await page.press('input[name="topics"]', "Enter");
    }

    // Save as draft
    await page.click('button:has-text("Save Draft")');

    // Verify success message
    await expect(page.locator("text=Story created successfully")).toBeVisible();

    // Verify redirect to story list or story view
    await page.waitForURL(/\/story\/.+/);

    // Verify story content is displayed
    await expect(page.locator("h1")).toContainText(testStory.title);
    await expect(page.locator("text=Draft")).toBeVisible();
  });

  test("Creator can preview and submit story for review", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    // Login
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', creator.email);
    await page.fill('input[type="password"]', creator.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Go to creator's stories (assuming a dashboard or profile page exists)
    await page.goto("/dashboard/stories");

    // Find the draft story we created
    await expect(
      page.locator("text=E2E Test Story - Creator Workflow")
    ).toBeVisible();

    // Click on the story to view it
    await page.click("text=E2E Test Story - Creator Workflow");

    // Verify we're on the story page
    await expect(page.locator("h1")).toContainText(
      "E2E Test Story - Creator Workflow"
    );
    await expect(page.locator("text=Draft")).toBeVisible();

    // Click preview/edit button
    await page.click('button:has-text("Edit")');

    // Make a small edit to test the workflow
    await page
      .locator('textarea[name="content"]')
      .fill(
        "Dies ist eine Testgeschichte für unsere E2E-Tests. Anna geht spazieren und trifft einen Freund. Sie unterhalten sich über das Wetter und planen einen Ausflug. Am Ende sind beide glücklich und freuen sich auf das Wochenende. Diese Geschichte wurde aktualisiert und ist bereit für die Überprüfung."
      );

    // Submit for review (change status to preview)
    await page.click('button:has-text("Submit for Review")');

    // Verify success message
    await expect(page.locator("text=Story submitted for review")).toBeVisible();

    // Verify status changed to preview
    await expect(page.locator("text=Preview")).toBeVisible();

    // Verify creator can no longer edit (edit button should be disabled or hidden)
    await expect(page.locator('button:has-text("Edit")')).toBeDisabled();
  });

  test("Creator can view their content status history", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    // Login
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', creator.email);
    await page.fill('input[type="password"]', creator.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Navigate to creator dashboard
    await page.goto("/dashboard");

    // Verify dashboard shows content statistics
    await expect(page.locator("text=My Content")).toBeVisible();

    // Check for status counts
    await expect(page.locator("text=Draft")).toBeVisible();
    await expect(page.locator("text=Preview")).toBeVisible();

    // Verify creator can see their submitted content
    await page.goto("/dashboard/stories");

    // Should see at least the test story we created
    await expect(
      page.locator("text=E2E Test Story - Creator Workflow")
    ).toBeVisible();

    // Verify status filter works
    await page.selectOption('select[name="status"]', "preview");
    await expect(
      page.locator("text=E2E Test Story - Creator Workflow")
    ).toBeVisible();

    await page.selectOption('select[name="status"]', "draft");
    await expect(
      page.locator("text=E2E Test Story - Creator Workflow")
    ).toBeHidden();
  });

  test("Creator cannot access admin functions", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    // Login as creator
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', creator.email);
    await page.fill('input[type="password"]', creator.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Try to access admin pages - should be redirected or show error
    await page.goto("/admin/reviews");

    // Should either redirect to login/home or show access denied
    const currentUrl = page.url();
    const hasAccessDenied = await page
      .locator("text=Access Denied")
      .isVisible();
    const redirectedHome =
      currentUrl.includes("/auth/login") ||
      currentUrl === "http://localhost:8000/";

    expect(hasAccessDenied || redirectedHome).toBeTruthy();

    // Admin navigation should not be visible
    await page.goto("/");
    await expect(page.locator('a[href="/admin"]')).toBeHidden();
    await expect(page.locator("text=Admin Dashboard")).toBeHidden();
  });

  test("Form validation works correctly", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    // Login
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', creator.email);
    await page.fill('input[type="password"]', creator.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Go to create story page
    await page.goto("/create/story");

    // Try to submit empty form
    await page.click('button:has-text("Save Draft")');

    // Should see validation errors
    await expect(page.locator("text=Title is required")).toBeVisible();
    await expect(page.locator("text=Content is required")).toBeVisible();

    // Test title length validation
    await page.fill('input[name="title"]', "Short"); // Too short
    await page.click('button:has-text("Save Draft")');
    await expect(
      page.locator("text=Title must be at least 10 characters")
    ).toBeVisible();

    // Test content length validation
    await page.fill('input[name="title"]', "Valid Title for Testing Purposes");
    await page.fill('textarea[name="content"]', "Too short"); // Too short
    await page.click('button:has-text("Save Draft")');
    await expect(
      page.locator("text=Content must be at least 100 characters")
    ).toBeVisible();

    // Test valid form submission
    await page.fill(
      'textarea[name="content"]',
      "This is a valid content with more than one hundred characters. It contains a proper German learning story that meets all the validation requirements for our platform."
    );
    await page.selectOption('select[name="level"]', "A1");
    await page.fill('input[name="topics"]', "Validation");
    await page.press('input[name="topics"]', "Enter");

    await page.click('button:has-text("Save Draft")');

    // Should succeed now
    await expect(page.locator("text=Story created successfully")).toBeVisible();
  });
});
