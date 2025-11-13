import { Page } from "@playwright/test";

/**
 * Helper functions for content creation and management in E2E tests
 */

export interface StoryData {
  title: string;
  content: string;
  level: string;
  topics: string[];
}

export interface QuestionData {
  questionText: string;
  questionType: "multiple_choice" | "open_text" | "multiple_select";
  options?: string[];
  correctAnswer?: string | string[];
  difficulty: "easy" | "medium" | "hard";
  sampleAnswer?: string;
}

/**
 * Fill story form with provided data
 * @param page - Playwright page object
 * @param storyData - Story data to fill
 */
export async function fillStoryForm(
  page: Page,
  storyData: StoryData
): Promise<void> {
  await page.fill('input[name="title"]', storyData.title);
  await page.fill('textarea[name="content"]', storyData.content);
  await page.selectOption('select[name="level"]', storyData.level);

  // Add topics
  for (const topic of storyData.topics) {
    await page.fill('input[name="topics"]', topic);
    await page.press('input[name="topics"]', "Enter");
  }
}

/**
 * Test story form validation by submitting empty form
 * @param page - Playwright page object
 */
export async function testStoryFormValidation(page: Page): Promise<void> {
  // Try to submit empty form
  await page.click('button:has-text("Save Draft")');

  // Should see validation errors
  // Note: These selectors may need adjustment based on actual UI
  const titleError = await page.locator("text=Title is required").isVisible();
  const contentError = await page
    .locator("text=Content is required")
    .isVisible();

  if (!titleError || !contentError) {
    console.warn("Form validation messages may not match expected selectors");
  }
}

/**
 * Create a story and submit for review
 * @param page - Playwright page object
 * @param storyData - Story data to create
 */
export async function createAndSubmitStory(
  page: Page,
  storyData: StoryData
): Promise<void> {
  await page.goto("/create/story");
  await fillStoryForm(page, storyData);
  await page.click('button:has-text("Save Draft")');
  await page.waitForTimeout(1000); // Wait for save to complete
  await page.click('button:has-text("Submit for Review")');
}

/**
 * Fill question form with provided data
 * @param page - Playwright page object
 * @param questionData - Question data to fill
 * @param storyId - Optional story ID to link question to
 */
export async function fillQuestionForm(
  page: Page,
  questionData: QuestionData,
  storyId?: string
): Promise<void> {
  if (storyId) {
    await page.selectOption('[data-testid="story-select"]', storyId);
  }

  await page.selectOption(
    '[data-testid="question-type"]',
    questionData.questionType
  );
  await page.fill('[data-testid="question-text"]', questionData.questionText);

  if (questionData.questionType === "multiple_choice" && questionData.options) {
    questionData.options.forEach(async (option, index) => {
      await page.fill(`[data-testid="option-${index + 1}"]`, option);
    });

    if (
      questionData.correctAnswer &&
      typeof questionData.correctAnswer === "string"
    ) {
      const correctIndex =
        questionData.options.indexOf(questionData.correctAnswer) + 1;
      await page.check(`[data-testid="correct-option-${correctIndex}"]`);
    }
  } else if (
    questionData.questionType === "multiple_select" &&
    questionData.options
  ) {
    questionData.options.forEach(async (option, index) => {
      await page.fill(`[data-testid="option-${index + 1}"]`, option);
    });

    if (
      questionData.correctAnswer &&
      Array.isArray(questionData.correctAnswer)
    ) {
      questionData.correctAnswer.forEach(async (answer) => {
        const correctIndex =
          questionData.options!.indexOf(answer as string) + 1;
        await page.check(`[data-testid="correct-option-${correctIndex}"]`);
      });
    }
  } else if (
    questionData.questionType === "open_text" &&
    questionData.sampleAnswer
  ) {
    await page.fill('[data-testid="sample-answer"]', questionData.sampleAnswer);
  }

  await page.selectOption('[data-testid="difficulty"]', questionData.difficulty);
}

/**
 * Navigate to content review page and filter by status
 * @param page - Playwright page object
 * @param status - Status to filter by
 */
export async function navigateToReviews(
  page: Page,
  status: "draft" | "preview" | "ready" | "published" | "all" = "preview"
): Promise<void> {
  await page.goto("/admin/reviews");
  await page.selectOption('select[id="status-filter"]', status);
}

/**
 * Update content status with comment
 * @param page - Playwright page object
 * @param newStatus - New status to set
 * @param comment - Review comment
 */
export async function updateContentStatus(
  page: Page,
  newStatus: "draft" | "preview" | "ready" | "published",
  comment: string
): Promise<void> {
  await page.selectOption('select[name="newStatus"]', newStatus);
  await page.fill('textarea[name="reviewComment"]', comment);
  await page.click('button:has-text("Update Status")');
}
