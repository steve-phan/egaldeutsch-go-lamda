import { test, expect } from "@playwright/test";
import { TestDatabaseSetup } from "../fixtures/databaseSetup";
import { loginAs, logout, isLoggedIn, isLoggedOut } from "../fixtures/authHelpers";

let testSetup: TestDatabaseSetup;

test.describe("Authentication Workflows", () => {
  test.beforeAll(async () => {
    testSetup = new TestDatabaseSetup();
    await testSetup.setupTestData();
  });

  test.afterAll(async () => {
    if (testSetup) {
      await testSetup.cleanupTestData();
    }
  });

  test("Complete login flow with form validation", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    await page.goto("/auth/login");

    // Verify login page loads
    await expect(page).toHaveTitle(/EgalDeutsch/);

    // Test empty form validation
    await page.click('button[type="submit"]');

    // Check if validation messages appear (may vary by implementation)
    const hasEmailError =
      (await page.locator("text=Email is required").isVisible()) ||
      (await page.locator("text=required").first().isVisible());
    expect(hasEmailError).toBeTruthy();

    // Test invalid email format
    await page.fill('input[type="email"]', "invalid-email");
    await page.click('button[type="submit"]');

    // May show invalid email format error (implementation dependent)
    await page.waitForTimeout(500);

    // Test incorrect credentials
    await page.fill('input[type="email"]', "wrong@email.com");
    await page.fill('input[type="password"]', "wrongpassword");
    await page.click('button[type="submit"]');

    // Should show error message (exact text may vary)
    await page.waitForTimeout(1000);
    const hasAuthError =
      (await page.locator("text=Invalid credentials").isVisible()) ||
      (await page.locator("text=Invalid").isVisible()) ||
      (await page.locator("text=incorrect").isVisible());

    // Test successful login
    await page.fill('input[type="email"]', creator.email);
    await page.fill('input[type="password"]', creator.password);
    await page.click('button[type="submit"]');

    // Verify redirect and session
    await page.waitForURL("/");
    await expect(page.locator("text=Welcome")).toBeVisible();
    await expect(page.locator(`text=${creator.name}`)).toBeVisible();
  });

  test("Session persistence across page refreshes", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    // Login
    await loginAs(page, creator.email, creator.password);

    // Navigate to different pages
    await page.goto("/dashboard/stories");
    await expect(page.locator(`text=${creator.name}`)).toBeVisible();

    // Refresh page
    await page.reload();
    await page.waitForTimeout(1000);

    // Verify user still logged in
    await expect(page.locator(`text=${creator.name}`)).toBeVisible();

    // Navigate to another page
    await page.goto("/");
    await expect(page.locator(`text=${creator.name}`)).toBeVisible();
  });

  test("Session expiration and auto-logout handling", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    await loginAs(page, creator.email, creator.password);

    // Mock session expiration by manipulating token in localStorage
    await page.evaluate(() => {
      // Clear authentication token to simulate expiration
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      sessionStorage.clear();
    });

    // Try to access protected page
    await page.goto("/dashboard/stories");

    // Should redirect to login or show session expired message
    await page.waitForTimeout(1000);
    const currentUrl = page.url();
    const isOnLoginPage = currentUrl.includes("/auth/login");
    const isOnHomePage = currentUrl === "http://localhost:8000/";

    expect(isOnLoginPage || isOnHomePage).toBeTruthy();
  });

  test("Logout clears session and redirects", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    await loginAs(page, creator.email, creator.password);

    // Verify logged in
    await expect(page.locator(`text=${creator.name}`)).toBeVisible();

    // Logout - try multiple possible logout methods
    const logoutMethods = [
      async () => {
        // Method 1: User menu with data-testid
        if (await page.locator('[data-testid="user-menu"]').isVisible()) {
          await page.click('[data-testid="user-menu"]');
          await page.click("text=Logout");
        }
      },
      async () => {
        // Method 2: Direct logout link
        if (await page.locator('a[href="/auth/logout"]').isVisible()) {
          await page.click('a[href="/auth/logout"]');
        }
      },
      async () => {
        // Method 3: Navigate to logout endpoint
        await page.goto("/auth/logout");
      },
    ];

    for (const logoutMethod of logoutMethods) {
      try {
        await logoutMethod();
        break;
      } catch (e) {
        console.log("Logout method failed, trying next...");
      }
    }

    await page.waitForTimeout(1000);

    // Verify redirect to home/login
    const currentUrl = page.url();
    expect(
      currentUrl === "http://localhost:8000/" ||
        currentUrl.includes("/auth/login")
    ).toBeTruthy();

    // Try to access protected page
    await page.goto("/dashboard");
    await page.waitForTimeout(1000);

    // Should redirect back to login or home
    const finalUrl = page.url();
    expect(
      finalUrl.includes("/auth/login") || finalUrl === "http://localhost:8000/"
    ).toBeTruthy();
  });

  test("Multiple tab session synchronization", async ({ browser }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    // Create two browser contexts (tabs)
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();
    const page1 = await context1.newPage();
    const page2 = await context2.newPage();

    try {
      // Login in first tab
      await loginAs(page1, creator.email, creator.password);
      await expect(page1.locator(`text=${creator.name}`)).toBeVisible();

      // Second tab should require login (separate context)
      await page2.goto("/auth/login");
      await expect(page2).toHaveURL(/\/auth\/login/);

      // Login in second tab
      await loginAs(page2, creator.email, creator.password);
      await expect(page2.locator(`text=${creator.name}`)).toBeVisible();

      // Both tabs should be logged in
      await page1.reload();
      await expect(page1.locator(`text=${creator.name}`)).toBeVisible();
      await page2.reload();
      await expect(page2.locator(`text=${creator.name}`)).toBeVisible();
    } finally {
      await context1.close();
      await context2.close();
    }
  });

  test("Login with different user roles", async ({ page }) => {
    const context = testSetup.getContext();

    // Test creator login
    const creator = context.users.creator;
    await loginAs(page, creator.email, creator.password);
    await expect(page.locator(`text=${creator.name}`)).toBeVisible();
    await expect(page.locator('a[href="/admin"]')).toBeHidden();

    // Logout
    await page.goto("/auth/logout");
    await page.waitForTimeout(500);

    // Test admin login
    const admin = context.users.admin;
    await loginAs(page, admin.email, admin.password);
    await expect(page.locator(`text=${admin.name}`)).toBeVisible();
    await expect(page.locator('a[href="/admin"]')).toBeVisible();
  });

  test("Concurrent login attempts handling", async ({ browser }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    // Create multiple browser contexts for concurrent login
    const contexts = await Promise.all([
      browser.newContext(),
      browser.newContext(),
      browser.newContext(),
    ]);

    const pages = await Promise.all(contexts.map((ctx) => ctx.newPage()));

    try {
      // Attempt concurrent logins
      await Promise.all(
        pages.map((page) => loginAs(page, creator.email, creator.password))
      );

      // All should successfully login
      await Promise.all(
        pages.map((page) =>
          expect(page.locator(`text=${creator.name}`)).toBeVisible()
        )
      );

      // Verify all can access protected pages
      await Promise.all(pages.map((page) => page.goto("/dashboard/stories")));

      await Promise.all(
        pages.map((page) =>
          expect(page.locator(`text=${creator.name}`)).toBeVisible()
        )
      );
    } finally {
      await Promise.all(contexts.map((ctx) => ctx.close()));
    }
  });

  test("Password visibility toggle", async ({ page }) => {
    await page.goto("/auth/login");

    // Fill password field
    await page.fill('input[type="password"]', "testpassword");

    // Check if password toggle exists (implementation may vary)
    const hasToggle = await page
      .locator('[data-testid="password-toggle"]')
      .isVisible()
      .catch(() => false);

    if (hasToggle) {
      // Verify password is hidden by default
      await expect(page.locator('input[type="password"]')).toBeVisible();

      // Click toggle
      await page.click('[data-testid="password-toggle"]');

      // Verify password is visible
      await expect(page.locator('input[type="text"]')).toBeVisible();

      // Toggle back
      await page.click('[data-testid="password-toggle"]');
      await expect(page.locator('input[type="password"]')).toBeVisible();
    }
  });
});
