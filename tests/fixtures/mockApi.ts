import { Page } from "@playwright/test";

/**
 * Mock API responses for Playwright tests
 * This allows tests to run without requiring the actual server
 */
export class MockApiSetup {
  constructor(private page: Page) {}

  async setupMockResponses() {
    // Mock user-management endpoints
    await this.page.route("**/user-management/register", async (route) => {
      await route.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify({
          id: "mock-user-id",
          username: "testuser",
          email: "test@example.com",
          role: "creator",
        }),
      });
    });

    await this.page.route("**/user-management/login", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          token: "mock-jwt-token",
          user: {
            id: "mock-user-id",
            username: "testuser",
            email: "test@example.com",
            role: "creator",
          },
        }),
      });
    });

    // Mock stories-management endpoints
    await this.page.route("**/stories-management", async (route) => {
      if (route.request().method() === "GET") {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            data: [
              {
                id: "story-1",
                title: "Test Story",
                content: "This is a test story content.",
                level: "A1",
                status: "published",
                wordCount: 50,
                readingTime: 2,
                topics: ["Test"],
                vocabulary: [],
                summary: "A test story",
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString(),
              },
            ],
          }),
        });
      } else if (route.request().method() === "POST") {
        await route.fulfill({
          status: 201,
          contentType: "application/json",
          body: JSON.stringify({
            id: "new-story-id",
            title: "New Test Story",
            status: "draft",
          }),
        });
      }
    });

    // Mock questions-management endpoints
    await this.page.route("**/questions-management", async (route) => {
      if (route.request().method() === "GET") {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            questions: [
              {
                id: "question-1",
                storyId: "story-1",
                question: "What is the main topic?",
                questionType: "comprehension",
                options: ["A", "B", "C", "D"],
                correctAnswer: 0,
                explanation: "Test explanation",
                points: 10,
                order: 1,
                difficulty: "medium",
                status: "published",
              },
            ],
            total: 1,
            page: 1,
            limit: 10,
          }),
        });
      } else if (route.request().method() === "POST") {
        await route.fulfill({
          status: 201,
          contentType: "application/json",
          body: JSON.stringify({
            id: "new-question-id",
            status: "draft",
          }),
        });
      }
    });

    // Mock quiz-management endpoints
    await this.page.route("**/quiz-management", async (route) => {
      if (route.request().method() === "POST") {
        await route.fulfill({
          status: 201,
          contentType: "application/json",
          body: JSON.stringify({
            id: "new-quiz-id",
            status: "draft",
          }),
        });
      }
    });

    // Mock health endpoint
    await this.page.route("**/health", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ status: "ok" }),
      });
    });
  }

  async setupAuthenticatedUser(
    role: "creator" | "admin" | "reviewer" = "creator"
  ) {
    // Set up localStorage with mock authentication
    try {
      await this.page.evaluate((userRole) => {
        if (typeof window !== "undefined" && window.localStorage) {
          localStorage.setItem("authToken", "mock-jwt-token");
          localStorage.setItem(
            "user",
            JSON.stringify({
              id: "mock-user-id",
              username: "testuser",
              email: "test@example.com",
              role: userRole,
              firstName: "Test",
              lastName: "User",
            })
          );
        }
      }, role);
    } catch (error) {
      console.warn(
        "Could not set localStorage, likely due to security context:",
        error
      );
      // Alternative: Set via cookies or other means if needed
    }
  }
}
