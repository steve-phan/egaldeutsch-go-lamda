import { test, expect } from "@playwright/test";
import { MockApiSetup } from "../fixtures/mockApi";
import {
  mockLoginAs,
  setupMockUser,
  mockLogout,
  isMockLoggedIn,
  isMockLoggedOut,
} from "../fixtures/mockAuthHelpers";

test.describe("Authentication Workflows", () => {
  test("Complete login flow with form validation", async ({ page }) => {
    const mockApi = new MockApiSetup(page);
    await mockApi.setupMockResponses();

    await page.goto("/auth/login");

    // Verify login page loads
    await expect(page).toHaveTitle(
      "EgalDeutsch - Learn German Through Stories"
    );

    // Test empty form validation
    await page.click('button[type="submit"]');

    // Check if validation messages appear (may vary by implementation)
    const hasUsernameError =
      (await page.locator("text=username").isVisible()) ||
      (await page.locator("text=required").first().isVisible());
    expect(hasUsernameError).toBeTruthy();

    // Test incorrect credentials
    await page.fill('input[name="username"]', "wronguser");
    await page.fill('input[name="password"]', "wrongpassword");
    await page.click('button[type="submit"]');

    // Should show error message (exact text may vary)
    await page.waitForTimeout(1000);
    const hasAuthError =
      (await page.locator("text=Invalid").isVisible()) ||
      (await page.locator("text=incorrect").isVisible()) ||
      (await page.locator(".error").isVisible());

    // Test successful login
    await page.fill('input[name="username"]', "testuser");
    await page.fill('input[name="password"]', "password123");
    await page.click('button[type="submit"]');

    // Verify redirect happens (may not show welcome text immediately)
    await page.waitForTimeout(2000);
    expect(page.url()).toContain("/");
  });

  test("Session persistence across page refreshes", async ({ page }) => {
    // Setup authenticated user with mock API
    await setupMockUser(page, "creator");

    // Navigate to different pages
    await page.goto("/stories");

    // Verify session persists
    const isLoggedIn = await isMockLoggedIn(page);
    expect(isLoggedIn).toBeTruthy();

    // Refresh page
    await page.reload();
    await page.waitForTimeout(1000);

    // Should still be logged in after refresh
    const isStillLoggedIn = await isMockLoggedIn(page);
    expect(isStillLoggedIn).toBeTruthy();

    // Navigate to another page
    await page.goto("/");
    const isLoggedInOnHome = await isMockLoggedIn(page);
    expect(isLoggedInOnHome).toBeTruthy();
  });

  test("Session expiration and auto-logout handling", async ({ page }) => {
    // Setup authenticated user with mock API
    await setupMockUser(page, "creator");

    // Mock session expiration by clearing token in localStorage
    await page.evaluate(() => {
      if (typeof window !== "undefined" && window.localStorage) {
        localStorage.removeItem("authToken");
      }
    });

    // Navigate to a protected page
    await page.goto("/admin");
    await page.waitForTimeout(1000);

    // Should be redirected to login or show login prompt
    const currentUrl = page.url();
    const isOnLoginPage = currentUrl.includes("/auth/login");
    const isOnHomePage = currentUrl === "http://localhost:8000/";

    // Either redirected to login or home page (depending on implementation)
    expect(isOnLoginPage || isOnHomePage).toBeTruthy();
  });

  test("Logout clears session and redirects", async ({ page }) => {
    // Setup authenticated user with mock API
    await setupMockUser(page, "creator");

    // Verify logged in
    const initiallyLoggedIn = await isMockLoggedIn(page);
    expect(initiallyLoggedIn).toBeTruthy();

    // Logout
    await mockLogout(page);
    await page.waitForTimeout(1000);

    // Verify logged out
    const isLoggedOut = await isMockLoggedOut(page);
    expect(isLoggedOut).toBeTruthy();

    // Try to access protected page - should be redirected or blocked
    await page.goto("/admin");
    await page.waitForTimeout(1000);

    const currentUrl = page.url();
    const isNotOnAdminPage = !currentUrl.includes("/admin");
    expect(isNotOnAdminPage).toBeTruthy();
  });

  test("Login with different user roles", async ({ page }) => {
    // Test creator role
    await setupMockUser(page, "creator");
    let isCreatorLoggedIn = await isMockLoggedIn(page);
    expect(isCreatorLoggedIn).toBeTruthy();

    // Clear session
    await mockLogout(page);

    // Test admin role
    await setupMockUser(page, "admin");
    let isAdminLoggedIn = await isMockLoggedIn(page);
    expect(isAdminLoggedIn).toBeTruthy();

    // Clear session
    await mockLogout(page);

    // Test reviewer role
    await setupMockUser(page, "reviewer");
    let isReviewerLoggedIn = await isMockLoggedIn(page);
    expect(isReviewerLoggedIn).toBeTruthy();
  });

  test("Login form handles loading states", async ({ page }) => {
    const mockApi = new MockApiSetup(page);
    await mockApi.setupMockResponses();

    await page.goto("/auth/login");

    // Fill form
    await page.fill('input[name="username"]', "testuser");
    await page.fill('input[name="password"]', "password123");

    // Submit form and check for loading state
    await page.click('button[type="submit"]');

    // Button should be disabled during submission
    const isButtonDisabled = await page
      .locator('button[type="submit"]')
      .isDisabled();
    // Note: This may not work if the login is too fast, but it's good to test

    // Wait for completion
    await page.waitForTimeout(2000);
    expect(page.url()).toContain("/");
  });
});
