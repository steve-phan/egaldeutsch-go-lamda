// Test fixture data for German learning platform
export const testUsers = {
  creator: {
    email: "creator@test.com",
    password: "password123",
    name: "Test Creator",
    role: "creator",
  },
  admin: {
    email: "admin@test.com",
    password: "password123",
    name: "Test Admin",
    role: "admin",
  },
  reviewer: {
    email: "reviewer@test.com",
    password: "password123",
    name: "Test Reviewer",
    role: "reviewer",
  },
};

export const testStories = {
  draft: {
    title: "Ein Tag im Park - Draft",
    content:
      "Es ist ein schöner Sommertag. Anna geht in den Park. Sie sieht viele Blumen und Bäume. Ein kleiner Hund läuft zu ihr. Anna streichelt den Hund. Der Hund ist sehr freundlich. Anna spielt mit dem Hund im Park. Sie werfen einen Ball. Der Hund bringt den Ball zurück. Anna lacht und ist glücklich. Am Ende des Tages geht Anna nach Hause. Sie denkt an den schönen Tag im Park.",
    level: "A1",
    topics: ["Natur", "Tiere", "Freizeit"],
    status: "draft",
  },
  preview: {
    title: "Einkaufen im Supermarkt - Preview",
    content:
      "Maria braucht Lebensmittel. Sie geht zum Supermarkt. Im Supermarkt gibt es viele Regale. Maria nimmt einen Einkaufswagen. Zuerst kauft sie Brot und Milch. Dann geht sie zur Obst- und Gemüseabteilung. Sie kauft Äpfel, Bananen und Tomaten. An der Kasse bezahlt sie mit ihrer Kreditkarte. Die Kassiererin ist sehr freundlich. Maria bekommt eine Quittung. Zu Hause packt sie die Lebensmittel aus. Sie ist zufrieden mit ihrem Einkauf.",
    level: "A1",
    topics: ["Einkaufen", "Lebensmittel", "Alltag"],
    status: "preview",
  },
  ready: {
    title: "Mein erster Schultag - Ready",
    content:
      "Heute ist Toms erster Schultag. Er ist sehr aufgeregt. Seine Mutter macht ihm ein leckeres Frühstück. Tom packt seinen neuen Schulranzen. In der Schule trifft er seine neue Lehrerin, Frau Schmidt. Sie ist sehr nett. Tom lernt viele neue Kinder kennen. In der Pause spielen sie zusammen auf dem Schulhof. Nach der Schule erzählt Tom seiner Familie von seinem ersten Tag. Er freut sich schon auf morgen.",
    level: "A1",
    topics: ["Schule", "Familie", "Kinder"],
    status: "ready",
  },
  published: {
    title: "Die vier Jahreszeiten - Published",
    content:
      "In Deutschland gibt es vier Jahreszeiten. Im Frühling werden die Tage länger und wärmer. Die Blumen blühen und die Bäume bekommen grüne Blätter. Im Sommer ist es sehr warm. Die Menschen gehen schwimmen oder machen Urlaub. Im Herbst werden die Blätter bunt und fallen von den Bäumen. Im Winter ist es kalt und manchmal schneit es. Jede Jahreszeit hat ihren eigenen Charme und ihre besonderen Aktivitäten.",
    level: "A2",
    topics: ["Natur", "Wetter", "Jahreszeiten"],
    status: "published",
  },
};

export const testQuestions = {
  easy: [
    {
      question: "Was macht Anna im Park?",
      questionType: "multiple_choice",
      options: [
        "Sie liest ein Buch",
        "Sie spielt mit einem Hund",
        "Sie schläft auf einer Bank",
        "Sie isst ein Eis",
      ],
      correctAnswer: "Sie spielt mit einem Hund",
      difficulty: "easy",
      status: "draft",
    },
    {
      question: "Welche Tiere sieht Anna im Park?",
      questionType: "multiple_choice",
      options: ["Katze", "Hund", "Vogel", "Fisch"],
      correctAnswer: "Hund",
      difficulty: "easy",
      status: "preview",
    },
  ],
  medium: [
    {
      question: "Wie fühlt sich Anna am Ende des Tages?",
      questionType: "multiple_choice",
      options: ["müde", "traurig", "glücklich", "hungrig"],
      correctAnswer: "glücklich",
      difficulty: "medium",
      status: "ready",
    },
    {
      question: "Was kauft Maria im Supermarkt? (Mehrere Antworten möglich)",
      questionType: "multiple_select",
      options: ["Brot", "Milch", "Äpfel", "Fleisch"],
      correctAnswers: ["Brot", "Milch", "Äpfel"],
      difficulty: "medium",
      status: "published",
    },
  ],
  hard: [
    {
      question: "Beschreibe Annas Tag im Park in eigenen Worten.",
      questionType: "open_text",
      difficulty: "hard",
      sampleAnswer:
        "Anna verbringt einen schönen Tag im Park, wo sie mit einem freundlichen Hund spielt und Ball wirft. Am Ende ist sie glücklich und denkt an den schönen Tag.",
      status: "draft",
    },
  ],
};

export const testQuizzes = {
  draft: {
    title: "Park und Natur Quiz - Draft",
    description: "Ein einfaches Quiz über Annas Tag im Park",
    difficulty: "easy",
    questions: [], // Will be populated with question IDs
    status: "draft",
  },
  preview: {
    title: "Einkaufen Quiz - Preview",
    description: "Quiz über Marias Einkaufserlebnis",
    difficulty: "medium",
    questions: [], // Will be populated with question IDs
    status: "preview",
  },
  ready: {
    title: "Alltag in Deutschland Quiz - Ready",
    description: "Quiz über alltägliche Situationen in Deutschland",
    difficulty: "medium",
    questions: [], // Will be populated with question IDs
    status: "ready",
  },
  published: {
    title: "Jahreszeiten Quiz - Published",
    description: "Quiz über die vier Jahreszeiten in Deutschland",
    difficulty: "easy",
    questions: [], // Will be populated with question IDs
    status: "published",
  },
};

// Helper function to create test data with proper relationships
export function createTestDataSet() {
  return {
    users: testUsers,
    stories: testStories,
    questions: testQuestions,
    quizzes: testQuizzes,
    // This would be used to establish relationships between stories, questions, and quizzes
    relationships: {
      "park-story": {
        story: "draft",
        questions: ["easy.0", "easy.1", "hard.0"],
        quiz: "draft",
      },
      "shopping-story": {
        story: "preview",
        questions: ["medium.1"],
        quiz: "preview",
      },
      "school-story": {
        story: "ready",
        questions: ["medium.0"],
        quiz: "ready",
      },
      "seasons-story": {
        story: "published",
        questions: ["medium.1"],
        quiz: "published",
      },
    },
  };
}
