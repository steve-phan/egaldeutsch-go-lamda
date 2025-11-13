import { test, expect } from "@playwright/test";

test.describe("Basic Smoke Tests", () => {
  test("can load the homepage", async ({ page }) => {
    await page.goto("/");

    // Check if the page loads successfully
    await expect(page).toHaveTitle(
      "EgalDeutsch - Learn German Through Stories"
    );

    // Check for main heading
    await expect(page.locator("h1")).toContainText(
      "Learn German Through Stories"
    );
  });

  test("can navigate to login page", async ({ page }) => {
    await page.goto("/auth/login");

    // Check if login form loads
    await expect(page.locator('input[name="username"]')).toBeVisible();
    await expect(page.locator('input[name="password"]')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
  });

  test("can navigate to stories page", async ({ page }) => {
    await page.goto("/stories");

    // Check if stories page loads
    await expect(page.locator("h1")).toContainText("Stories");
  });
});
