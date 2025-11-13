import { Page } from "@playwright/test";

/**
 * Helper functions for authentication workflows in E2E tests
 */

/**
 * Login as a specific user
 * @param page - Playwright page object
 * @param email - User email
 * @param password - User password
 */
export async function loginAs(
  page: Page,
  email: string,
  password: string
): Promise<void> {
  await page.goto("/auth/login");
  await page.fill('input[type="email"]', email);
  await page.fill('input[type="password"]', password);
  await page.click('button[type="submit"]');
  await page.waitForURL("/");
}

/**
 * Login as a user from the test context
 * @param page - Playwright page object
 * @param userKey - Key of the user in test context (e.g., "creator", "admin")
 * @param context - Test context containing user data
 */
export async function loginAsUser(
  page: Page,
  userKey: string,
  context: any
): Promise<void> {
  const user = context.users[userKey];
  if (!user) {
    throw new Error(`User ${userKey} not found in test context`);
  }
  await loginAs(page, user.email, user.password);
}

/**
 * Logout current user
 * @param page - Playwright page object
 */
export async function logout(page: Page): Promise<void> {
  await page.click('[data-testid="user-menu"]');
  await page.click("text=Logout");
}

/**
 * Verify user is logged in
 * @param page - Playwright page object
 * @param userName - Expected user name to verify
 * @returns boolean - True if user is logged in
 */
export async function isLoggedIn(
  page: Page,
  userName?: string
): Promise<boolean> {
  const userMenuVisible = await page
    .locator('[data-testid="user-menu"]')
    .isVisible();

  if (!userMenuVisible) {
    return false;
  }

  if (userName) {
    const userNameVisible = await page.locator(`text=${userName}`).isVisible();
    return userNameVisible;
  }

  return true;
}

/**
 * Verify user is logged out
 * @param page - Playwright page object
 * @returns boolean - True if user is logged out
 */
export async function isLoggedOut(page: Page): Promise<boolean> {
  const userMenuHidden = await page
    .locator('[data-testid="user-menu"]')
    .isHidden();
  return userMenuHidden;
}
