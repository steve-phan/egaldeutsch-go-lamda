import { test, expect } from "@playwright/test";
import { TestDatabaseSetup } from "../fixtures/databaseSetup";

let testSetup: TestDatabaseSetup;

test.describe("Complete Content Lifecycle", () => {
  test.beforeAll(async ({ browser }) => {
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

  test("Complete workflow: Creator creates → Admin reviews → Content goes live → Users access", async ({
    browser,
  }) => {
    const context = testSetup.getContext();

    // Use separate browser contexts for different user roles
    const creatorContext = await browser.newContext();
    const adminContext = await browser.newContext();
    const publicContext = await browser.newContext();

    const creatorPage = await creatorContext.newPage();
    const adminPage = await adminContext.newPage();
    const publicPage = await publicContext.newPage();

    try {
      // PHASE 1: Creator creates content
      console.log("Phase 1: Creator creates content");

      // Creator login
      await creatorPage.goto("/auth/login");
      await creatorPage.fill(
        'input[type="email"]',
        context.users.creator.email
      );
      await creatorPage.fill(
        'input[type="password"]',
        context.users.creator.password
      );
      await creatorPage.click('button[type="submit"]');
      await creatorPage.waitForURL("/");

      // Create a new story
      await creatorPage.goto("/create/story");

      const storyData = {
        title: "Lifecycle Test: Die Bibliothek",
        content:
          "Lisa geht zur Bibliothek. Sie möchte ein Buch über Geschichte lesen. Die Bibliothekarin hilft ihr dabei, das richtige Buch zu finden. Lisa setzt sich an einen ruhigen Tisch am Fenster. Sie liest zwei Stunden lang. Das Buch ist sehr interessant und lehrreich. Am Ende leiht sie das Buch aus und nimmt es mit nach Hause. Lisa freut sich darauf, heute Abend weiterzulesen.",
        level: "A2",
        topics: ["Bildung", "Lesen", "Bibliothek"],
      };

      await creatorPage.fill('input[name="title"]', storyData.title);
      await creatorPage.fill('textarea[name="content"]', storyData.content);
      await creatorPage.selectOption('select[name="level"]', storyData.level);

      // Add topics
      for (const topic of storyData.topics) {
        await creatorPage.fill('input[name="topics"]', topic);
        await creatorPage.press('input[name="topics"]', "Enter");
      }

      // Save as draft
      await creatorPage.click('button:has-text("Save Draft")');
      await expect(
        creatorPage.locator("text=Story created successfully")
      ).toBeVisible();

      // Verify story is in draft status
      await expect(creatorPage.locator("text=Draft")).toBeVisible();

      // Submit for review (change to preview status)
      await creatorPage.click('button:has-text("Submit for Review")');
      await expect(
        creatorPage.locator("text=Story submitted for review")
      ).toBeVisible();
      await expect(creatorPage.locator("text=Preview")).toBeVisible();

      // PHASE 2: Admin reviews and approves content
      console.log("Phase 2: Admin reviews content");

      // Admin login
      await adminPage.goto("/auth/login");
      await adminPage.fill('input[type="email"]', context.users.admin.email);
      await adminPage.fill(
        'input[type="password"]',
        context.users.admin.password
      );
      await adminPage.click('button[type="submit"]');
      await adminPage.waitForURL("/");

      // Navigate to admin review dashboard
      await adminPage.goto("/admin/reviews");

      // Filter to preview content
      await adminPage.selectOption('select[id="status-filter"]', "preview");

      // Find and review the story
      await adminPage.click("text=Lifecycle Test: Die Bibliothek");

      // Approve story (change to ready status)
      await adminPage.selectOption('select[name="newStatus"]', "ready");
      await adminPage.fill(
        'textarea[name="reviewComment"]',
        "Excellent story! Good vocabulary for A2 level. Cultural context is appropriate. Approved for publication."
      );
      await adminPage.click('button:has-text("Update Status")');
      await expect(
        adminPage.locator("text=Status updated successfully")
      ).toBeVisible();

      // Verify story is now in ready status
      await adminPage.selectOption('select[id="status-filter"]', "ready");
      await expect(
        adminPage.locator("text=Lifecycle Test: Die Bibliothek")
      ).toBeVisible();
      await expect(adminPage.locator("text=Ready")).toBeVisible();

      // PHASE 3: Admin publishes content
      console.log("Phase 3: Admin publishes content");

      // Click on the ready story
      await adminPage.click("text=Lifecycle Test: Die Bibliothek");

      // Publish the story
      await adminPage.selectOption('select[name="newStatus"]', "published");
      await adminPage.fill(
        'textarea[name="reviewComment"]',
        "Story published and now available to all learners. Excellent addition to our A2 content library."
      );
      await adminPage.click('button:has-text("Update Status")');
      await expect(
        adminPage.locator("text=Status updated successfully")
      ).toBeVisible();

      // Verify story is now published
      await adminPage.selectOption('select[id="status-filter"]', "published");
      await expect(
        adminPage.locator("text=Lifecycle Test: Die Bibliothek")
      ).toBeVisible();
      await expect(adminPage.locator("text=Published")).toBeVisible();

      // PHASE 4: Public users can access published content
      console.log("Phase 4: Public users access content");

      // Test with public/guest user
      await publicPage.goto("/");

      // Navigate to stories section (assuming public stories page exists)
      await publicPage.goto("/stories");

      // Should see published story
      await expect(
        publicPage.locator("text=Lifecycle Test: Die Bibliothek")
      ).toBeVisible();

      // Click to read the story
      await publicPage.click("text=Lifecycle Test: Die Bibliothek");

      // Verify story content is displayed
      await expect(publicPage.locator("h1")).toContainText(
        "Lifecycle Test: Die Bibliothek"
      );
      await expect(
        publicPage.locator("text=Lisa geht zur Bibliothek")
      ).toBeVisible();
      await expect(publicPage.locator("text=A2")).toBeVisible();

      // Verify story metadata is shown
      await expect(publicPage.locator("text=Bildung")).toBeVisible();
      await expect(publicPage.locator("text=Lesen")).toBeVisible();
      await expect(publicPage.locator("text=Bibliothek")).toBeVisible();

      // PHASE 5: Verify content is accessible in learning context
      console.log("Phase 5: Verify learning integration");

      // Test that story appears in level-filtered content
      await publicPage.goto("/stories?level=A2");
      await expect(
        publicPage.locator("text=Lifecycle Test: Die Bibliothek")
      ).toBeVisible();

      // Test that story appears in topic-filtered content
      await publicPage.goto("/stories?topic=Bildung");
      await expect(
        publicPage.locator("text=Lifecycle Test: Die Bibliothek")
      ).toBeVisible();

      // PHASE 6: Verify workflow history and audit trail
      console.log("Phase 6: Verify workflow history");

      // Admin can see full history
      await adminPage.goto("/admin/reviews");
      await adminPage.selectOption('select[id="status-filter"]', "published");
      await adminPage.click("text=Lifecycle Test: Die Bibliothek");

      // Should show complete comment history
      await expect(
        adminPage.locator("text=Approved for publication")
      ).toBeVisible();
      await expect(
        adminPage.locator("text=Story published and now available")
      ).toBeVisible();

      // Creator can still see their published content
      await creatorPage.goto("/dashboard/stories");
      await creatorPage.selectOption('select[name="status"]', "published");
      await expect(
        creatorPage.locator("text=Lifecycle Test: Die Bibliothek")
      ).toBeVisible();

      console.log("Complete lifecycle test passed successfully!");
    } finally {
      // Clean up browser contexts
      await creatorContext.close();
      await adminContext.close();
      await publicContext.close();
    }
  });

  test("Workflow handles rejection and revision cycle", async ({ browser }) => {
    const context = testSetup.getContext();

    const creatorContext = await browser.newContext();
    const adminContext = await browser.newContext();

    const creatorPage = await creatorContext.newPage();
    const adminPage = await adminContext.newPage();

    try {
      // PHASE 1: Creator creates content
      await creatorPage.goto("/auth/login");
      await creatorPage.fill(
        'input[type="email"]',
        context.users.creator.email
      );
      await creatorPage.fill(
        'input[type="password"]',
        context.users.creator.password
      );
      await creatorPage.click('button[type="submit"]');
      await creatorPage.waitForURL("/");

      // Create story with intentional issues for rejection
      await creatorPage.goto("/create/story");

      const rejectionStory = {
        title: "Rejection Test: Das Problem",
        content:
          "Das ist ein sehr kurzer Text. Es hat nicht genug Inhalt für eine richtige Geschichte. Außerdem fehlen wichtige Details und kultureller Kontext.", // Too short and lacks depth
        level: "A1",
        topics: ["Test"],
      };

      await creatorPage.fill('input[name="title"]', rejectionStory.title);
      await creatorPage.fill(
        'textarea[name="content"]',
        rejectionStory.content
      );
      await creatorPage.selectOption(
        'select[name="level"]',
        rejectionStory.level
      );
      await creatorPage.fill('input[name="topics"]', rejectionStory.topics[0]);
      await creatorPage.press('input[name="topics"]', "Enter");

      await creatorPage.click('button:has-text("Save Draft")');
      await creatorPage.click('button:has-text("Submit for Review")');

      // PHASE 2: Admin rejects content
      await adminPage.goto("/auth/login");
      await adminPage.fill('input[type="email"]', context.users.admin.email);
      await adminPage.fill(
        'input[type="password"]',
        context.users.admin.password
      );
      await adminPage.click('button[type="submit"]');
      await adminPage.waitForURL("/");

      await adminPage.goto("/admin/reviews");
      await adminPage.selectOption('select[id="status-filter"]', "preview");
      await adminPage.click("text=Rejection Test: Das Problem");

      // Reject (send back to draft)
      await adminPage.selectOption('select[name="newStatus"]', "draft");
      await adminPage.fill(
        'textarea[name="reviewComment"]',
        "Content rejected. Issues: 1) Too short - needs at least 150 words, 2) Lacks cultural context, 3) Limited vocabulary for learning. Please revise and resubmit."
      );
      await adminPage.click('button:has-text("Update Status")');

      // PHASE 3: Creator revises based on feedback
      await creatorPage.goto("/dashboard/stories");
      await creatorPage.selectOption('select[name="status"]', "draft");
      await creatorPage.click("text=Rejection Test: Das Problem");

      // Should see rejection feedback
      await expect(creatorPage.locator("text=Content rejected")).toBeVisible();
      await expect(
        creatorPage.locator("text=Too short - needs at least 150 words")
      ).toBeVisible();

      // Edit the story
      await creatorPage.click('button:has-text("Edit")');

      // Improve the content based on feedback
      await creatorPage.fill(
        'textarea[name="content"]',
        "Anna hat ein Problem mit ihrem Computer. Sie kann ihre wichtige Präsentation nicht öffnen. Morgen muss sie die Präsentation in der Uni zeigen. Anna ist sehr nervös und weiß nicht, was sie tun soll. Sie ruft ihren Freund Max an. Max ist sehr gut mit Computern. Er kommt sofort zu Anna und schaut sich das Problem an. Nach einer Stunde Arbeit kann Max das Problem lösen. Die Präsentation ist wieder da! Anna ist so glücklich und dankbar. Sie lädt Max zum Essen ein. Beide sind zufrieden mit dem Ergebnis."
      );

      await creatorPage.fill('input[name="topics"]', "Problem lösung");
      await creatorPage.press('input[name="topics"]', "Enter");
      await creatorPage.fill('input[name="topics"]', "Freundschaft");
      await creatorPage.press('input[name="topics"]', "Enter");

      await creatorPage.click('button:has-text("Save Changes")');
      await creatorPage.click('button:has-text("Submit for Review")');

      // PHASE 4: Admin approves revised content
      await adminPage.selectOption('select[id="status-filter"]', "preview");
      await adminPage.click("text=Rejection Test: Das Problem");

      // Approve the revised version
      await adminPage.selectOption('select[name="newStatus"]', "ready");
      await adminPage.fill(
        'textarea[name="reviewComment"]',
        "Much improved! Good length, cultural context added, appropriate vocabulary. Approved for publication."
      );
      await adminPage.click('button:has-text("Update Status")');

      // Verify successful revision cycle
      await adminPage.selectOption('select[id="status-filter"]', "ready");
      await expect(
        adminPage.locator("text=Rejection Test: Das Problem")
      ).toBeVisible();
      await expect(adminPage.locator("text=Ready")).toBeVisible();

      console.log("Rejection and revision cycle test passed successfully!");
    } finally {
      await creatorContext.close();
      await adminContext.close();
    }
  });

  test("Multiple users can work simultaneously without conflicts", async ({
    browser,
  }) => {
    const context = testSetup.getContext();

    // Create multiple browser contexts for concurrent work
    const creator1Context = await browser.newContext();
    const creator2Context = await browser.newContext();
    const adminContext = await browser.newContext();

    const creator1Page = await creator1Context.newPage();
    const creator2Page = await creator2Context.newPage();
    const adminPage = await adminContext.newPage();

    try {
      // Login all users simultaneously
      const loginPromises = [
        (async () => {
          await creator1Page.goto("/auth/login");
          await creator1Page.fill(
            'input[type="email"]',
            context.users.creator.email
          );
          await creator1Page.fill(
            'input[type="password"]',
            context.users.creator.password
          );
          await creator1Page.click('button[type="submit"]');
          await creator1Page.waitForURL("/");
        })(),
        (async () => {
          // Create a second creator account for this test
          await creator2Page.goto("/auth/register");
          await creator2Page.fill('input[name="email"]', "creator2@test.com");
          await creator2Page.fill('input[name="password"]', "password123");
          await creator2Page.fill('input[name="name"]', "Test Creator 2");
          await creator2Page.selectOption('select[name="role"]', "creator");
          await creator2Page.click('button[type="submit"]');
          await creator2Page.waitForURL("/");
        })(),
        (async () => {
          await adminPage.goto("/auth/login");
          await adminPage.fill(
            'input[type="email"]',
            context.users.admin.email
          );
          await adminPage.fill(
            'input[type="password"]',
            context.users.admin.password
          );
          await adminPage.click('button[type="submit"]');
          await adminPage.waitForURL("/");
        })(),
      ];

      await Promise.all(loginPromises);

      // Both creators create content simultaneously
      const creationPromises = [
        (async () => {
          await creator1Page.goto("/create/story");
          await creator1Page.fill(
            'input[name="title"]',
            "Concurrent Test 1: Im Restaurant"
          );
          await creator1Page.fill(
            'textarea[name="content"]',
            "Familie Schmidt geht ins Restaurant. Sie bestellen verschiedene Gerichte. Der Kellner ist sehr freundlich und hilfsbereit. Das Essen schmeckt ausgezeichnet und alle sind zufrieden. Am Ende bezahlen sie die Rechnung und geben Trinkgeld."
          );
          await creator1Page.selectOption('select[name="level"]', "A2");
          await creator1Page.fill('input[name="topics"]', "Restaurant");
          await creator1Page.press('input[name="topics"]', "Enter");
          await creator1Page.click('button:has-text("Save Draft")');
          await creator1Page.click('button:has-text("Submit for Review")');
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
          await creator2Page.click('button:has-text("Submit for Review")');
        })(),
      ];

      await Promise.all(creationPromises);

      // Admin reviews both pieces of content
      await adminPage.goto("/admin/reviews");
      await adminPage.selectOption('select[id="status-filter"]', "preview");

      // Should see both submitted stories
      await expect(
        adminPage.locator("text=Concurrent Test 1: Im Restaurant")
      ).toBeVisible();
      await expect(
        adminPage.locator("text=Concurrent Test 2: Beim Arzt")
      ).toBeVisible();

      // Process both reviews simultaneously (approve both)
      const reviewPromises = [
        (async () => {
          await adminPage.click("text=Concurrent Test 1: Im Restaurant");
          await adminPage.selectOption('select[name="newStatus"]', "ready");
          await adminPage.fill(
            'textarea[name="reviewComment"]',
            "Restaurant story approved - good cultural context"
          );
          await adminPage.click('button:has-text("Update Status")');
        })(),
        (async () => {
          // Wait a bit to avoid UI conflicts, then process second story
          await adminPage.waitForTimeout(2000);
          await adminPage.selectOption('select[id="status-filter"]', "preview");
          await adminPage.click("text=Concurrent Test 2: Beim Arzt");
          await adminPage.selectOption('select[name="newStatus"]', "ready");
          await adminPage.fill(
            'textarea[name="reviewComment"]',
            "Medical story approved - appropriate vocabulary"
          );
          await adminPage.click('button:has-text("Update Status")');
        })(),
      ];

      await Promise.all(reviewPromises);

      // Verify both stories are now in ready status
      await adminPage.selectOption('select[id="status-filter"]', "ready");
      await expect(
        adminPage.locator("text=Concurrent Test 1: Im Restaurant")
      ).toBeVisible();
      await expect(
        adminPage.locator("text=Concurrent Test 2: Beim Arzt")
      ).toBeVisible();

      console.log("Concurrent workflow test passed successfully!");
    } finally {
      await creator1Context.close();
      await creator2Context.close();
      await adminContext.close();
    }
  });
});
