import { Page } from "@playwright/test";
import { MockApiSetup } from "./mockApi";

/**
 * Helper functions for authentication workflows in E2E tests with Mock API
 */

/**
 * Login as a specific user using mock API
 * @param page - Playwright page object
 * @param username - Username
 * @param password - User password
 */
export async function mockLoginAs(
  page: Page,
  username: string,
  password: string
): Promise<void> {
  const mockApi = new MockApiSetup(page);
  await mockApi.setupMockResponses();

  await page.goto("/auth/login");
  await page.fill('input[name="username"]', username);
  await page.fill('input[name="password"]', password);
  await page.click('button[type="submit"]');
  await page.waitForTimeout(2000); // Wait for navigation
}

/**
 * Setup authenticated user with mock API
 * @param page - Playwright page object
 * @param role - User role ("creator", "admin", "reviewer")
 */
export async function setupMockUser(
  page: Page,
  role: "creator" | "admin" | "reviewer" = "creator"
): Promise<void> {
  const mockApi = new MockApiSetup(page);
  await mockApi.setupMockResponses();
  await page.goto("/");
  await mockApi.setupAuthenticatedUser(role);
}

/**
 * Mock logout current user
 * @param page - Playwright page object
 */
export async function mockLogout(page: Page): Promise<void> {
  // Navigate to a page first
  await page.goto("/");

  // Try to find logout button - implementation may vary
  const logoutSelectors = [
    '[data-testid="user-menu"]',
    'button:has-text("Logout")',
    'a:has-text("Logout")',
    ".user-menu",
  ];

  for (const selector of logoutSelectors) {
    try {
      const element = page.locator(selector);
      if (await element.isVisible({ timeout: 1000 })) {
        await element.click();
        if (selector === '[data-testid="user-menu"]') {
          await page.click("text=Logout");
        }
        break;
      }
    } catch (e) {
      // Continue to next selector
    }
  }

  // Clear localStorage as fallback
  await page.evaluate(() => {
    if (typeof window !== "undefined" && window.localStorage) {
      localStorage.removeItem("authToken");
      localStorage.removeItem("user");
    }
  });
}

/**
 * Check if user is logged in (mock API version)
 * @param page - Playwright page object
 */
export async function isMockLoggedIn(page: Page): Promise<boolean> {
  try {
    const hasToken = await page.evaluate(() => {
      return (
        typeof window !== "undefined" &&
        window.localStorage &&
        !!localStorage.getItem("authToken")
      );
    });
    return hasToken;
  } catch (e) {
    return false;
  }
}

/**
 * Check if user is logged out (mock API version)
 * @param page - Playwright page object
 */
export async function isMockLoggedOut(page: Page): Promise<boolean> {
  return !(await isMockLoggedIn(page));
}

// Keep original functions for backward compatibility
export async function loginAs(
  page: Page,
  email: string,
  password: string
): Promise<void> {
  // Convert email to username for mock API compatibility
  const username = email.split("@")[0];
  return mockLoginAs(page, username, password);
}

export async function logout(page: Page): Promise<void> {
  return mockLogout(page);
}

export async function isLoggedIn(page: Page): Promise<boolean> {
  return isMockLoggedIn(page);
}

export async function isLoggedOut(page: Page): Promise<boolean> {
  return isMockLoggedOut(page);
}
