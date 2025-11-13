import axios from "axios";
import {
  testUsers,
  testStories,
  testQuestions,
  testQuizzes,
  createTestDataSet,
} from "./testData";

const API_BASE_URL = "http://localhost:8888/.netlify/functions";

interface TestUser {
  id?: string;
  email: string;
  password: string;
  name: string;
  role: string;
  token?: string;
}

interface TestContext {
  users: Record<string, TestUser>;
  stories: Record<string, any>;
  questions: Record<string, any>;
  quizzes: Record<string, any>;
}

/**
 * Database setup utility for Playwright tests
 * Handles creating and cleaning up test data
 */
export class TestDatabaseSetup {
  private context: TestContext = {
    users: {},
    stories: {},
    questions: {},
    quizzes: {},
  };

  /**
   * Setup test database with all fixture data
   */
  async setupTestData(): Promise<TestContext> {
    console.log("Setting up test database...");

    try {
      // 1. Create test users and get authentication tokens
      await this.createTestUsers();

      // 2. Create test stories with different statuses
      await this.createTestStories();

      // 3. Create test questions linked to stories
      await this.createTestQuestions();

      // 4. Create test quizzes
      await this.createTestQuizzes();

      console.log("Test database setup completed successfully");
      return this.context;
    } catch (error) {
      console.error("Failed to setup test database:", error);
      throw error;
    }
  }

  /**
   * Clean up all test data
   */
  async cleanupTestData(): Promise<void> {
    console.log("Cleaning up test database...");

    try {
      // Delete in reverse order to handle dependencies
      await this.deleteTestQuizzes();
      await this.deleteTestQuestions();
      await this.deleteTestStories();
      await this.deleteTestUsers();

      console.log("Test database cleanup completed");
    } catch (error) {
      console.error("Failed to cleanup test database:", error);
      // Don't throw here, we want tests to continue even if cleanup fails
    }
  }

  /**
   * Create test users and authenticate them
   */
  private async createTestUsers(): Promise<void> {
    for (const [key, userData] of Object.entries(testUsers)) {
      try {
        // Register user
        const registerResponse = await axios.post(
          `${API_BASE_URL}/user-management/register`,
          {
            username: userData.email.split("@")[0], // Use email prefix as username
            email: userData.email,
            password: userData.password,
            firstName: userData.name.split(" ")[0] || userData.name,
            lastName: userData.name.split(" ")[1] || "User",
            preferredRole:
              userData.role === "admin" ? "creator" : userData.role,
          }
        );

        if (registerResponse.status === 201) {
          // Login to get token
          const loginResponse = await axios.post(
            `${API_BASE_URL}/user-management/login`,
            {
              username: userData.email.split("@")[0],
              password: userData.password,
            }
          );

          this.context.users[key] = {
            ...userData,
            id: registerResponse.data.id,
            token: loginResponse.data.token,
          };

          console.log(`Created test user: ${userData.email}`);
        }
      } catch (error: any) {
        if (error.response?.status === 409) {
          // User already exists, try to login
          try {
            const loginResponse = await axios.post(
              `${API_BASE_URL}/user-management/login`,
              {
                username: userData.email.split("@")[0],
                password: userData.password,
              }
            );

            this.context.users[key] = {
              ...userData,
              token: loginResponse.data.token,
            };

            console.log(`Using existing test user: ${userData.email}`);
          } catch (loginError) {
            console.error(
              `Failed to login existing user ${userData.email}:`,
              loginError
            );
          }
        } else {
          console.error(`Failed to create user ${userData.email}:`, error);
        }
      }
    }
  }

  /**
   * Create test stories with different statuses
   */
  private async createTestStories(): Promise<void> {
    if (!this.context.users.creator?.token) {
      console.log("No creator token available, skipping story creation");
      return;
    }

    const creatorToken = this.context.users.creator.token;

    for (const [key, storyData] of Object.entries(testStories)) {
      try {
        const response = await axios.post(
          `${API_BASE_URL}/stories-management`,
          storyData,
          {
            headers: {
              Authorization: `Bearer ${creatorToken}`,
              "Content-Type": "application/json",
            },
          }
        );

        this.context.stories[key] = {
          ...storyData,
          id: response.data.id,
        };

        console.log(`Created test story: ${storyData.title}`);
      } catch (error) {
        console.error(`Failed to create story ${storyData.title}:`, error);
      }
    }
  }

  /**
   * Create test questions linked to stories
   */
  private async createTestQuestions(): Promise<void> {
    if (!this.context.users.creator?.token) {
      console.log("No creator token available, skipping question creation");
      return;
    }

    const creatorToken = this.context.users.creator.token;
    const storyIds = Object.values(this.context.stories).map(
      (story: any) => story.id
    );

    for (const [difficulty, questions] of Object.entries(testQuestions)) {
      for (const [index, questionData] of questions.entries()) {
        try {
          const response = await axios.post(
            `${API_BASE_URL}/questions-management`,
            {
              ...questionData,
              storyId: storyIds[index % storyIds.length], // Cycle through available stories
            },
            {
              headers: {
                Authorization: `Bearer ${creatorToken}`,
                "Content-Type": "application/json",
              },
            }
          );

          if (!this.context.questions[difficulty]) {
            this.context.questions[difficulty] = [];
          }

          this.context.questions[difficulty].push({
            ...questionData,
            id: response.data.id,
          });

          console.log(`Created test question: ${questionData.question}`);
        } catch (error) {
          console.error(
            `Failed to create question ${questionData.question}:`,
            error
          );
        }
      }
    }
  }

  /**
   * Create test quizzes
   */
  private async createTestQuizzes(): Promise<void> {
    if (!this.context.users.creator?.token) {
      console.log("No creator token available, skipping quiz creation");
      return;
    }

    const creatorToken = this.context.users.creator.token;

    for (const [key, quizData] of Object.entries(testQuizzes)) {
      try {
        const response = await axios.post(
          `${API_BASE_URL}/quiz-management`,
          quizData,
          {
            headers: {
              Authorization: `Bearer ${creatorToken}`,
              "Content-Type": "application/json",
            },
          }
        );

        this.context.quizzes[key] = {
          ...quizData,
          id: response.data.id,
        };

        console.log(`Created test quiz: ${quizData.title}`);
      } catch (error) {
        console.error(`Failed to create quiz ${quizData.title}:`, error);
      }
    }
  }

  // Cleanup methods
  private async deleteTestUsers(): Promise<void> {
    for (const [key, user] of Object.entries(this.context.users)) {
      if (user.id && user.token) {
        try {
          await axios.delete(`${API_BASE_URL}/user-management/${user.id}`, {
            headers: {
              Authorization: `Bearer ${user.token}`,
            },
          });
          console.log(`Deleted test user: ${user.email}`);
        } catch (error) {
          console.error(`Failed to delete user ${user.email}:`, error);
        }
      }
    }
  }

  private async deleteTestStories(): Promise<void> {
    if (!this.context.users.admin?.token) {
      console.log("No admin token available, skipping story cleanup");
      return;
    }

    const adminToken = this.context.users.admin.token;

    for (const [key, story] of Object.entries(this.context.stories)) {
      if ((story as any).id) {
        try {
          await axios.delete(
            `${API_BASE_URL}/stories-management/${(story as any).id}`,
            {
              headers: {
                Authorization: `Bearer ${adminToken}`,
              },
            }
          );
          console.log(`Deleted test story: ${(story as any).title}`);
        } catch (error) {
          console.error(
            `Failed to delete story ${(story as any).title}:`,
            error
          );
        }
      }
    }
  }

  private async deleteTestQuestions(): Promise<void> {
    if (!this.context.users.admin?.token) {
      console.log("No admin token available, skipping question cleanup");
      return;
    }

    const adminToken = this.context.users.admin.token;

    for (const [difficulty, questions] of Object.entries(
      this.context.questions
    )) {
      for (const question of questions as any[]) {
        if (question.id) {
          try {
            await axios.delete(
              `${API_BASE_URL}/questions-management/${question.id}`,
              {
                headers: {
                  Authorization: `Bearer ${adminToken}`,
                },
              }
            );
            console.log(`Deleted test question: ${question.question}`);
          } catch (error) {
            console.error(
              `Failed to delete question ${question.question}:`,
              error
            );
          }
        }
      }
    }
  }

  private async deleteTestQuizzes(): Promise<void> {
    if (!this.context.users.admin?.token) {
      console.log("No admin token available, skipping quiz cleanup");
      return;
    }

    const adminToken = this.context.users.admin.token;

    for (const [key, quiz] of Object.entries(this.context.quizzes)) {
      if ((quiz as any).id) {
        try {
          await axios.delete(
            `${API_BASE_URL}/quiz-management/${(quiz as any).id}`,
            {
              headers: {
                Authorization: `Bearer ${adminToken}`,
              },
            }
          );
          console.log(`Deleted test quiz: ${(quiz as any).title}`);
        } catch (error) {
          console.error(`Failed to delete quiz ${(quiz as any).title}:`, error);
        }
      }
    }
  }

  /**
   * Get the current test context
   */
  getContext(): TestContext {
    return this.context;
  }
}
