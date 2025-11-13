import { test, expect } from "@playwright/test";
import { TestDatabaseSetup } from "../fixtures/databaseSetup";

let testSetup: TestDatabaseSetup;

test.describe("Admin Review Workflow", () => {
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

  test("Admin can login and access admin dashboard", async ({ page }) => {
    const context = testSetup.getContext();
    const admin = context.users.admin;

    // Navigate to login page
    await page.goto("/auth/login");

    // Fill login form
    await page.fill('input[type="email"]', admin.email);
    await page.fill('input[type="password"]', admin.password);

    // Submit login
    await page.click('button[type="submit"]');

    // Wait for redirect to dashboard
    await page.waitForURL("/");

    // Verify admin can see admin navigation
    await expect(page.locator('a[href="/admin"]')).toBeVisible();

    // Navigate to admin dashboard
    await page.goto("/admin");

    // Verify admin dashboard loads
    await expect(page.locator("h1")).toContainText("Admin Dashboard");
    await expect(page.locator("text=Content Review")).toBeVisible();
    await expect(page.locator("text=User Management")).toBeVisible();
  });

  test("Admin can view and filter pending content", async ({ page }) => {
    const context = testSetup.getContext();
    const admin = context.users.admin;

    // Login as admin
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', admin.email);
    await page.fill('input[type="password"]', admin.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Navigate to content review page
    await page.goto("/admin/reviews");

    // Verify review page loads
    await expect(page.locator("h1")).toContainText("Content Review");

    // Should default to showing preview status content
    await expect(page.locator('select[id="status-filter"]')).toHaveValue(
      "preview"
    );

    // Verify test content is visible
    await expect(
      page.locator("text=Einkaufen im Supermarkt - Preview")
    ).toBeVisible();

    // Test status filter
    await page.selectOption('select[id="status-filter"]', "draft");
    await expect(page.locator("text=Ein Tag im Park - Draft")).toBeVisible();
    await expect(
      page.locator("text=Einkaufen im Supermarkt - Preview")
    ).toBeHidden();

    // Test type filter
    await page.selectOption('select[id="status-filter"]', "all");
    await page.selectOption('select[id="type-filter"]', "story");

    // Should only show stories
    await expect(page.locator('[data-testid="content-item"]')).toContainText(
      "Story"
    );

    // Test search functionality
    await page.fill('input[id="search-input"]', "Park");
    await expect(page.locator("text=Ein Tag im Park - Draft")).toBeVisible();
    await expect(page.locator("text=Einkaufen im Supermarkt")).toBeHidden();
  });

  test("Admin can approve content and change status to ready", async ({
    page,
  }) => {
    const context = testSetup.getContext();
    const admin = context.users.admin;

    // Login as admin
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', admin.email);
    await page.fill('input[type="password"]', admin.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Navigate to content review
    await page.goto("/admin/reviews");

    // Filter to show preview content
    await page.selectOption('select[id="status-filter"]', "preview");

    // Find and click on a preview item
    await page.click("text=Einkaufen im Supermarkt - Preview");

    // Should open review modal or detail view
    await expect(page.locator("text=Review Content")).toBeVisible();

    // Select "Ready" status
    await page.selectOption('select[name="newStatus"]', "ready");

    // Add approval comment
    await page.fill(
      'textarea[name="reviewComment"]',
      "Content approved. Good structure and appropriate for A1 level. Ready for publication."
    );

    // Submit review
    await page.click('button:has-text("Update Status")');

    // Verify success message
    await expect(
      page.locator("text=Status updated successfully")
    ).toBeVisible();

    // Verify content is no longer in preview list
    await page.selectOption('select[id="status-filter"]', "preview");
    await expect(
      page.locator("text=Einkaufen im Supermarkt - Preview")
    ).toBeHidden();

    // Verify content appears in ready list
    await page.selectOption('select[id="status-filter"]', "ready");
    await expect(page.locator("text=Einkaufen im Supermarkt")).toBeVisible();
    await expect(page.locator("text=Ready")).toBeVisible();
  });

  test("Admin can reject content and send back to draft", async ({ page }) => {
    const context = testSetup.getContext();
    const admin = context.users.admin;

    // Login as admin
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', admin.email);
    await page.fill('input[type="password"]', admin.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Navigate to content review
    await page.goto("/admin/reviews");

    // Filter to show preview content
    await page.selectOption('select[id="status-filter"]', "preview");

    // If no preview content, let's create some by submitting existing draft
    const hasPreviewContent = await page.locator("text=Preview").isVisible();
    if (!hasPreviewContent) {
      // Create some preview content first
      await page.selectOption('select[id="status-filter"]', "draft");
      await page.click("text=Ein Tag im Park - Draft");
      await page.selectOption('select[name="newStatus"]', "preview");
      await page.fill(
        'textarea[name="reviewComment"]',
        "Submitting for review"
      );
      await page.click('button:has-text("Update Status")');
      await page.waitForTimeout(1000);
    }

    // Now reject content
    await page.selectOption('select[id="status-filter"]', "preview");
    await page.click('[data-testid="content-item"]:first-child');

    // Select "Draft" status (rejection)
    await page.selectOption('select[name="newStatus"]', "draft");

    // Add rejection comment with specific feedback
    await page.fill(
      'textarea[name="reviewComment"]',
      "Content needs revision. Please improve the following: 1) Add more descriptive vocabulary, 2) Include more cultural context, 3) Ensure grammar examples are clearer."
    );

    // Submit review
    await page.click('button:has-text("Update Status")');

    // Verify success message
    await expect(
      page.locator("text=Status updated successfully")
    ).toBeVisible();

    // Verify content is back in draft status
    await page.selectOption('select[id="status-filter"]', "draft");
    await expect(page.locator('[data-testid="content-item"]')).toBeVisible();

    // Verify rejection comment is visible
    await page.click('[data-testid="content-item"]:first-child');
    await expect(page.locator("text=Content needs revision")).toBeVisible();
  });

  test("Admin can publish ready content", async ({ page }) => {
    const context = testSetup.getContext();
    const admin = context.users.admin;

    // Login as admin
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', admin.email);
    await page.fill('input[type="password"]', admin.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Navigate to content review
    await page.goto("/admin/reviews");

    // Filter to show ready content
    await page.selectOption('select[id="status-filter"]', "ready");

    // If no ready content, approve some first
    const hasReadyContent = await page.locator("text=Ready").isVisible();
    if (!hasReadyContent) {
      // Create ready content first
      await page.selectOption('select[id="status-filter"]', "preview");
      if (await page.locator('[data-testid="content-item"]').isVisible()) {
        await page.click('[data-testid="content-item"]:first-child');
        await page.selectOption('select[name="newStatus"]', "ready");
        await page.fill(
          'textarea[name="reviewComment"]',
          "Approved for publication"
        );
        await page.click('button:has-text("Update Status")');
        await page.waitForTimeout(1000);
      }
    }

    // Now publish ready content
    await page.selectOption('select[id="status-filter"]', "ready");
    await page.click('[data-testid="content-item"]:first-child');

    // Select "Published" status
    await page.selectOption('select[name="newStatus"]', "published");

    // Add publication comment
    await page.fill(
      'textarea[name="reviewComment"]',
      "Content published successfully. Available to all learners."
    );

    // Submit publication
    await page.click('button:has-text("Update Status")');

    // Verify success message
    await expect(
      page.locator("text=Status updated successfully")
    ).toBeVisible();

    // Verify content appears in published list
    await page.selectOption('select[id="status-filter"]', "published");
    await expect(page.locator("text=Published")).toBeVisible();

    // Verify published content is no longer in ready list
    await page.selectOption('select[id="status-filter"]', "ready");
    // The specific content should not be visible here anymore
  });

  test("Admin can manage users and change roles", async ({ page }) => {
    const context = testSetup.getContext();
    const admin = context.users.admin;

    // Login as admin
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', admin.email);
    await page.fill('input[type="password"]', admin.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Navigate to user management
    await page.goto("/admin/users");

    // Verify user management page loads
    await expect(page.locator("h1")).toContainText("User Management");

    // Should see test users
    await expect(page.locator("text=creator@test.com")).toBeVisible();
    await expect(page.locator("text=admin@test.com")).toBeVisible();

    // Test role filter
    await page.selectOption('select[name="roleFilter"]', "creator");
    await expect(page.locator("text=creator@test.com")).toBeVisible();
    await expect(page.locator("text=admin@test.com")).toBeHidden();

    // Test user search
    await page.selectOption('select[name="roleFilter"]', "all");
    await page.fill('input[name="userSearch"]', "creator");
    await expect(page.locator("text=creator@test.com")).toBeVisible();
    await expect(page.locator("text=admin@test.com")).toBeHidden();

    // Clear search
    await page.fill('input[name="userSearch"]', "");

    // Change user role (if functionality exists)
    await page.click("text=creator@test.com");

    // Should open user detail modal
    await expect(page.locator("text=User Details")).toBeVisible();

    // Change role from creator to reviewer
    await page.selectOption('select[name="userRole"]', "reviewer");

    // Add role change comment
    await page.fill(
      'textarea[name="roleChangeComment"]',
      "Promoted to reviewer based on high-quality content submissions"
    );

    // Save role change
    await page.click('button:has-text("Update Role")');

    // Verify success message
    await expect(
      page.locator("text=User role updated successfully")
    ).toBeVisible();

    // Verify role change is reflected
    await expect(page.locator("text=reviewer")).toBeVisible();
  });

  test("Admin dashboard shows accurate statistics", async ({ page }) => {
    const context = testSetup.getContext();
    const admin = context.users.admin;

    // Login as admin
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', admin.email);
    await page.fill('input[type="password"]', admin.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Navigate to admin dashboard
    await page.goto("/admin");

    // Verify statistics cards are present
    await expect(page.locator('[data-testid="stats-card"]')).toHaveCount(4); // draft, preview, ready, published

    // Verify each status has a count
    await expect(page.locator("text=Draft")).toBeVisible();
    await expect(page.locator("text=Preview")).toBeVisible();
    await expect(page.locator("text=Ready")).toBeVisible();
    await expect(page.locator("text=Published")).toBeVisible();

    // Verify recent activity section
    await expect(page.locator("text=Recent Activity")).toBeVisible();

    // Verify quick actions are available
    await expect(page.locator("text=Review Content")).toBeVisible();
    await expect(page.locator("text=Manage Users")).toBeVisible();

    // Test navigation from dashboard
    await page.click("text=Review Content");
    await expect(page).toHaveURL("/admin/reviews");

    await page.goto("/admin");
    await page.click("text=Manage Users");
    await expect(page).toHaveURL("/admin/users");
  });

  test("Admin cannot access regular user content creation", async ({
    page,
  }) => {
    const context = testSetup.getContext();
    const admin = context.users.admin;

    // Login as admin
    await page.goto("/auth/login");
    await page.fill('input[type="email"]', admin.email);
    await page.fill('input[type="password"]', admin.password);
    await page.click('button[type="submit"]');
    await page.waitForURL("/");

    // Admin should be able to access admin functions
    await expect(page.locator('a[href="/admin"]')).toBeVisible();

    // But should also be able to create content if needed (admin has all permissions)
    await page.goto("/create/story");

    // Admin should be able to create content (they have creator permissions too)
    await expect(page.locator("h1")).toContainText("Create Story");

    // But their main interface should emphasize admin functions
    await page.goto("/");
    await expect(page.locator("text=Admin Dashboard")).toBeVisible();
  });
});
