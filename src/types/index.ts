export interface VocabularyWord {
  german: string;
  english: string;
  wordType: string; // "noun", "verb", "adjective", etc.
  article?: string; // "der", "die", "das" for nouns
}

export interface Story {
  id: string;
  title: string;
  content: string;
  level: string; // A1, A2, B1, B2, C1, C2
  wordCount: number;
  readingTime: number; // Estimated reading time in minutes
  topic: string; // e.g., "Family", "Travel", "Food"
  vocabulary: VocabularyWord[];
  summary: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Question {
  id: string;
  storyId: string;
  question: string;
  questionType: string; // "comprehension", "vocabulary", "grammar"
  options: string[];
  correctAnswer: number; // Index of correct answer in options array
  explanation: string;
  points: number;
  order: number; // Question order in quiz (1-10)
  createdAt: string;
}

export interface Quiz {
  storyId: string;
  story?: Story;
  questions: Question[];
}

export interface QuizSubmission {
  id: string;
  storyId: string;
  answers: number[]; // Array of selected answer indices
  score: number; // Number of correct answers
  totalQuestions: number;
  submittedAt: string;
}

export interface QuizResult extends QuizSubmission {
  percentage: number;
  passed: boolean; // true if score >= 70%
  correctAnswers: boolean[]; // Array indicating which answers were correct
}

// API Response types
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface StoriesResponse extends ApiResponse<Story[]> {}
export interface StoryResponse extends ApiResponse<Story> {}
export interface QuizResponse extends ApiResponse<Quiz> {}
export interface QuizSubmissionResponse extends ApiResponse<QuizResult> {}

// Component Props types
export interface StoryCardProps {
  story: Story;
  className?: string;
}

export interface QuizQuestionProps {
  question: Question;
  selectedAnswer?: number;
  showResult?: boolean;
  onAnswerSelect: (answerIndex: number) => void;
}

export interface QuizResultsProps {
  result: QuizResult;
  story: Story;
  onRetakeQuiz: () => void;
  onReturnToStory: () => void;
}

// User and Authentication types
export type UserRole = 'creator' | 'reviewer' | 'admin';
export type UserStatus = 'active' | 'suspended' | 'pending';

export interface User {
  id: string;
  username: string;
  email: string;
  firstName: string;
  lastName: string;
  role: UserRole;
  status: UserStatus;
  createdAt: string;
  updatedAt: string;
  lastLoginAt?: string;
}

export interface LoginCredentials {
  username: string;
  password: string;
  rememberMe?: boolean;
}

export interface RegisterData {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
  firstName: string;
  lastName: string;
  preferredRole: 'creator' | 'reviewer';
}

export interface AuthResponse {
  token: string;
  expiresAt: string;
  user: User;
}

export interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginCredentials) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => void;
  hasRole: (role: UserRole) => boolean;
  hasAnyRole: (roles: UserRole[]) => boolean;
}
