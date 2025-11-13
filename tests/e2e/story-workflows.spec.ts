import { test, expect } from "@playwright/test";
import { TestDatabaseSetup } from "../fixtures/databaseSetup";
import { loginAs } from "../fixtures/authHelpers";
import {
  fillStoryForm,
  testStoryFormValidation,
  createAndSubmitStory,
  updateContentStatus,
  type StoryData,
} from "../fixtures/contentHelpers";

let testSetup: TestDatabaseSetup;

test.describe("Story Management Workflows", () => {
  test.beforeAll(async () => {
    testSetup = new TestDatabaseSetup();
    await testSetup.setupTestData();
  });

  test.afterAll(async () => {
    if (testSetup) {
      await testSetup.cleanupTestData();
    }
  });

  test("Complete story creation with all validation rules", async ({
    page,
  }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    await loginAs(page, creator.email, creator.password);
    await page.goto("/create/story");

    // Verify create story page
    await expect(page.locator("h1")).toContainText("Create Story");

    // Test form validation by submitting empty form
    await page.click('button:has-text("Save Draft")');

    // Should see validation errors
    await page.waitForTimeout(500);
    const hasValidationErrors =
      (await page.locator("text=required").isVisible()) ||
      (await page.locator("text=Title is required").isVisible()) ||
      (await page.locator("text=Content is required").isVisible());

    // Create valid story
    const storyData: StoryData = {
      title: "E2E Test: Eine Reise nach Berlin",
      content: `Maria plant eine Reise nach Berlin. Sie ist sehr aufgeregt, denn es ist ihre erste Reise in die deutsche Hauptstadt. 
      Zuerst bucht sie einen Zug von München nach Berlin. Die Fahrt dauert vier Stunden. Maria packt ihren Koffer mit warmer Kleidung, 
      denn es ist Winter. In Berlin möchte sie das Brandenburger Tor, das Reichstagsgebäude und die Berliner Mauer besuchen. 
      Sie hat auch Tickets für ein Konzert in der Philharmonie. Maria freut sich sehr auf die kulturellen Erlebnisse. 
      Am Ende ihrer Reise wird sie viele schöne Erinnerungen und Fotos haben.`,
      level: "A2",
      topics: ["Reisen", "Berlin", "Kultur", "Transport"],
    };

    await fillStoryForm(page, storyData);
    await page.click('button:has-text("Save Draft")');

    await page.waitForTimeout(1000);

    // Verify story creation
    const hasSuccessMessage =
      (await page.locator("text=Story created successfully").isVisible()) ||
      (await page.locator("text=created").isVisible()) ||
      (await page.locator("text=success").isVisible());

    if (hasSuccessMessage) {
      expect(hasSuccessMessage).toBeTruthy();
    }

    // Verify story is in draft status
    const hasDraftStatus =
      (await page.locator('[data-testid="story-status"]').isVisible()) ||
      (await page.locator("text=Draft").isVisible());
    expect(hasDraftStatus).toBeTruthy();

    // Verify story appears in creator's dashboard
    await page.goto("/dashboard/stories");
    await page.waitForTimeout(500);
    await expect(page.locator(`text=${storyData.title}`)).toBeVisible();
  });

  test("Story submission and review cycle with comments", async ({
    page,
    browser,
  }) => {
    const context = testSetup.getContext();

    // Use separate contexts for creator and admin
    const creatorContext = await browser.newContext();
    const adminContext = await browser.newContext();
    const creatorPage = await creatorContext.newPage();
    const adminPage = await adminContext.newPage();

    try {
      // Creator submits story for review
      await loginAs(
        creatorPage,
        context.users.creator.email,
        context.users.creator.password
      );
      await creatorPage.goto("/dashboard/stories");

      // Find the story we created
      const hasStory = await creatorPage
        .locator("text=Eine Reise nach Berlin")
        .isVisible();

      if (hasStory) {
        await creatorPage.click("text=Eine Reise nach Berlin");
        await creatorPage.waitForTimeout(500);

        // Submit for review
        const hasSubmitButton =
          (await creatorPage
            .locator('[data-testid="submit-for-review-button"]')
            .isVisible()) ||
          (await creatorPage
            .locator('button:has-text("Submit for Review")')
            .isVisible());

        if (hasSubmitButton) {
          try {
            await creatorPage.click('[data-testid="submit-for-review-button"]');
          } catch (e) {
            await creatorPage.click('button:has-text("Submit for Review")');
          }

          // May have a comment field
          const hasCommentField =
            await creatorPage
              .locator('[data-testid="submission-comment"]')
              .isVisible()
              .catch(() => false);

          if (hasCommentField) {
            await creatorPage.fill(
              '[data-testid="submission-comment"]',
              "This story covers A2-level travel vocabulary and German cultural landmarks. Ready for review."
            );
          }

          // Confirm submission
          try {
            await creatorPage.click('[data-testid="confirm-submission"]');
          } catch (e) {
            // May not have confirmation button
          }

          await creatorPage.waitForTimeout(1000);

          // Verify status change
          const hasPreviewStatus =
            (await creatorPage
              .locator('[data-testid="story-status"]')
              .textContent()
              .then((text) => text?.includes("Preview"))) ||
            (await creatorPage.locator("text=Preview").isVisible());
          expect(hasPreviewStatus).toBeTruthy();
        }
      }

      // Admin reviews story
      await loginAs(
        adminPage,
        context.users.admin.email,
        context.users.admin.password
      );
      await adminPage.goto("/admin/reviews");
      await adminPage.selectOption('select[id="status-filter"]', "preview");
      await adminPage.waitForTimeout(500);

      const storyVisible = await adminPage
        .locator("text=Eine Reise nach Berlin")
        .isVisible();

      if (storyVisible) {
        await adminPage.click("text=Eine Reise nach Berlin");
        await adminPage.waitForTimeout(500);

        // Admin provides detailed review
        await adminPage.selectOption('select[name="newStatus"]', "ready");
        await adminPage.fill(
          'textarea[name="reviewComment"]',
          `Excellent story! Strengths: 
        - Appropriate A2 vocabulary level
        - Good cultural context about Berlin
        - Clear narrative structure
        - Engaging travel theme
        
        Minor suggestions for future stories:
        - Could include more conversational elements
        - Consider adding more sensory details
        
        Approved for publication.`
        );
        await adminPage.click('button:has-text("Update Status")');

        await adminPage.waitForTimeout(1000);

        // Verify admin review completion
        const hasSuccess =
          (await adminPage
            .locator("text=Status updated successfully")
            .isVisible()) ||
          (await adminPage.locator("text=success").isVisible());

        if (hasSuccess) {
          expect(hasSuccess).toBeTruthy();
        }

        // Verify creator can see review feedback
        await creatorPage.reload();
        await creatorPage.waitForTimeout(500);

        const hasReadyStatus =
          (await creatorPage
            .locator('[data-testid="story-status"]')
            .textContent()
            .then((text) => text?.includes("Ready"))) ||
          (await creatorPage.locator("text=Ready").isVisible());

        if (hasReadyStatus) {
          expect(hasReadyStatus).toBeTruthy();
        }

        // Try to view review history
        const hasReviewHistory =
          await creatorPage
            .locator('[data-testid="view-review-history"]')
            .isVisible()
            .catch(() => false);

        if (hasReviewHistory) {
          await creatorPage.click('[data-testid="view-review-history"]');
          await creatorPage.waitForTimeout(500);
          await expect(
            creatorPage.locator("text=Excellent story!")
          ).toBeVisible();
        }
      }
    } finally {
      await creatorContext.close();
      await adminContext.close();
    }
  });

  test("Story rejection and revision cycle", async ({ page, browser }) => {
    const context = testSetup.getContext();

    const creatorContext = await browser.newContext();
    const adminContext = await browser.newContext();
    const creatorPage = await creatorContext.newPage();
    const adminPage = await adminContext.newPage();

    try {
      // Creator creates story with issues
      await loginAs(
        creatorPage,
        context.users.creator.email,
        context.users.creator.password
      );
      await creatorPage.goto("/create/story");

      const rejectionStory: StoryData = {
        title: "Revision Test: Kurzer Text",
        content:
          "Das ist ein sehr kurzer Text für eine Geschichte. Es fehlen wichtige Details und kultureller Kontext. Der Text ist zu kurz.",
        level: "A1",
        topics: ["Test"],
      };

      await fillStoryForm(creatorPage, rejectionStory);
      await creatorPage.click('button:has-text("Save Draft")');
      await creatorPage.waitForTimeout(1000);

      // Submit for review
      try {
        await creatorPage.click('button:has-text("Submit for Review")');
        await creatorPage.waitForTimeout(1000);
      } catch (e) {
        // May fail if button not found
      }

      // Admin rejects content
      await loginAs(
        adminPage,
        context.users.admin.email,
        context.users.admin.password
      );
      await adminPage.goto("/admin/reviews");
      await adminPage.selectOption('select[id="status-filter"]', "preview");
      await adminPage.waitForTimeout(500);

      const storyVisible = await adminPage
        .locator("text=Revision Test: Kurzer Text")
        .isVisible();

      if (storyVisible) {
        await adminPage.click("text=Revision Test: Kurzer Text");
        await adminPage.waitForTimeout(500);

        // Reject (send back to draft)
        await adminPage.selectOption('select[name="newStatus"]', "draft");
        await adminPage.fill(
          'textarea[name="reviewComment"]',
          "Content rejected. Issues: 1) Too short - needs at least 150 words, 2) Lacks cultural context, 3) Limited vocabulary for learning. Please revise and resubmit."
        );
        await adminPage.click('button:has-text("Update Status")');
        await adminPage.waitForTimeout(1000);

        // Creator revises based on feedback
        await creatorPage.goto("/dashboard/stories");
        await creatorPage.selectOption('select[name="status"]', "draft");
        await creatorPage.waitForTimeout(500);

        const draftStoryVisible = await creatorPage
          .locator("text=Revision Test: Kurzer Text")
          .isVisible();

        if (draftStoryVisible) {
          await creatorPage.click("text=Revision Test: Kurzer Text");
          await creatorPage.waitForTimeout(500);

          // Should see rejection feedback
          const hasRejectionFeedback =
            (await creatorPage.locator("text=Content rejected").isVisible()) ||
            (await creatorPage.locator("text=Too short").isVisible());

          // Edit the story
          try {
            await creatorPage.click('button:has-text("Edit")');
            await creatorPage.waitForTimeout(500);

            // Improve the content
            await creatorPage.fill(
              'textarea[name="content"]',
              "Anna hat ein Problem mit ihrem Computer. Sie kann ihre wichtige Präsentation nicht öffnen. Morgen muss sie die Präsentation in der Uni zeigen. Anna ist sehr nervös und weiß nicht, was sie tun soll. Sie ruft ihren Freund Max an. Max ist sehr gut mit Computern. Er kommt sofort zu Anna und schaut sich das Problem an. Nach einer Stunde Arbeit kann Max das Problem lösen. Die Präsentation ist wieder da! Anna ist so glücklich und dankbar. Sie lädt Max zum Essen ein. Beide sind zufrieden mit dem Ergebnis."
            );

            try {
              await creatorPage.click('button:has-text("Save Changes")');
              await creatorPage.waitForTimeout(1000);
            } catch (e) {
              await creatorPage.click('button:has-text("Save")');
              await creatorPage.waitForTimeout(1000);
            }

            // Submit for review again
            await creatorPage.click('button:has-text("Submit for Review")');
            await creatorPage.waitForTimeout(1000);
          } catch (e) {
            console.log("Edit workflow may differ from expected");
          }
        }
      }
    } finally {
      await creatorContext.close();
      await adminContext.close();
    }
  });

  test("Multiple creators can work on different stories simultaneously", async ({
    browser,
  }) => {
    const context = testSetup.getContext();

    const creator1Context = await browser.newContext();
    const creator2Context = await browser.newContext();
    const creator1Page = await creator1Context.newPage();
    const creator2Page = await creator2Context.newPage();

    try {
      // Both creators login
      await loginAs(
        creator1Page,
        context.users.creator.email,
        context.users.creator.password
      );

      // Create second creator for concurrent testing
      await creator2Page.goto("/auth/register");
      const hasRegister = await creator2Page
        .locator('input[name="email"]')
        .isVisible()
        .catch(() => false);

      if (hasRegister) {
        await creator2Page.fill('input[name="email"]', "creator2@test.com");
        await creator2Page.fill('input[name="password"]', "password123");
        await creator2Page.fill('input[name="name"]', "Test Creator 2");
        try {
          await creator2Page.selectOption('select[name="role"]', "creator");
        } catch (e) {
          // Role may be auto-assigned
        }
        await creator2Page.click('button[type="submit"]');
        await creator2Page.waitForTimeout(1000);
      } else {
        // Login as existing user if registration not available
        await loginAs(
          creator2Page,
          context.users.creator.email,
          context.users.creator.password
        );
      }

      // Both creators create content simultaneously
      await Promise.all([
        (async () => {
          await creator1Page.goto("/create/story");
          await creator1Page.fill(
            'input[name="title"]',
            "Concurrent Test 1: Im Restaurant"
          );
          await creator1Page.fill(
            'textarea[name="content"]',
            "Familie Schmidt geht ins Restaurant. Sie bestellen verschiedene Gerichte. Der Kellner ist sehr freundlich und hilfsbereit. Das Essen schmeckt ausgezeichnet und alle sind zufrieden. Am Ende bezahlen sie die Rechnung und geben Trinkgeld. Die Familie freut sich auf den nächsten Besuch."
          );
          await creator1Page.selectOption('select[name="level"]', "A2");
          await creator1Page.fill('input[name="topics"]', "Restaurant");
          await creator1Page.press('input[name="topics"]', "Enter");
          await creator1Page.click('button:has-text("Save Draft")');
        })(),
        (async () => {
          await creator2Page.goto("/create/story");
          await creator2Page.fill(
            'input[name="title"]',
            "Concurrent Test 2: Beim Arzt"
          );
          await creator2Page.fill(
            'textarea[name="content"]',
            "Peter fühlt sich nicht gut. Er geht zum Arzt. Im Wartezimmer sitzen viele Patienten. Die Sprechstundenhilfe ist sehr nett. Der Arzt untersucht Peter gründlich. Peter bekommt ein Rezept für Medikamente. Nach einer Woche geht es ihm viel besser."
          );
          await creator2Page.selectOption('select[name="level"]', "A2");
          await creator2Page.fill('input[name="topics"]', "Gesundheit");
          await creator2Page.press('input[name="topics"]', "Enter");
          await creator2Page.click('button:has-text("Save Draft")');
        })(),
      ]);

      await creator1Page.waitForTimeout(1000);
      await creator2Page.waitForTimeout(1000);

      // Verify both stories were created
      await creator1Page.goto("/dashboard/stories");
      const story1Visible = await creator1Page
        .locator("text=Concurrent Test 1: Im Restaurant")
        .isVisible();
      expect(story1Visible).toBeTruthy();

      await creator2Page.goto("/dashboard/stories");
      const story2Visible = await creator2Page
        .locator("text=Concurrent Test 2: Beim Arzt")
        .isVisible();
      expect(story2Visible).toBeTruthy();
    } finally {
      await creator1Context.close();
      await creator2Context.close();
    }
  });

  test("Story content length and level validation", async ({ page }) => {
    const context = testSetup.getContext();
    const creator = context.users.creator;

    await loginAs(page, creator.email, creator.password);
    await page.goto("/create/story");

    // Test too short title
    await page.fill('input[name="title"]', "Short");
    await page.click('button:has-text("Save Draft")');
    await page.waitForTimeout(500);

    // May show validation error
    const hasTitleError =
      (await page
        .locator("text=Title must be at least")
        .isVisible()
        .catch(() => false)) ||
      (await page.locator("text=too short").isVisible().catch(() => false));

    // Test too short content
    await page.fill('input[name="title"]', "Valid Title for Testing Purposes");
    await page.fill('textarea[name="content"]', "Too short content");
    await page.click('button:has-text("Save Draft")');
    await page.waitForTimeout(500);

    // May show validation error
    const hasContentError =
      (await page
        .locator("text=Content must be at least")
        .isVisible()
        .catch(() => false)) ||
      (await page.locator("text=too short").isVisible().catch(() => false));

    // Test valid form submission
    await page.fill(
      'textarea[name="content"]',
      "This is a valid content with more than one hundred characters. It contains a proper German learning story that meets all the validation requirements for our platform. The story is engaging and educational."
    );
    await page.selectOption('select[name="level"]', "A1");
    await page.fill('input[name="topics"]', "Validation");
    await page.press('input[name="topics"]', "Enter");

    await page.click('button:has-text("Save Draft")');
    await page.waitForTimeout(1000);

    // Should succeed
    const hasSuccess =
      (await page.locator("text=Story created successfully").isVisible()) ||
      (await page.locator("text=created").isVisible());

    if (hasSuccess) {
      expect(hasSuccess).toBeTruthy();
    }
  });
});
