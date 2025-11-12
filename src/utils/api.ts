import axios from "axios";
import {
  Story,
  Quiz,
  QuizSubmission,
  QuizResult,
  StoriesResponse,
  StoryResponse,
  QuizResponse,
  QuizSubmissionResponse,
} from "@/types";

const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? "/.netlify/functions"
    : "http://localhost:8888/.netlify/functions";

// Stories API
export const fetchStories = async (): Promise<Story[]> => {
  try {
    const response = await axios.get<StoriesResponse>(
      `${API_BASE_URL}/stories`
    );
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.error || "Failed to fetch stories");
  } catch (error) {
    console.error("Error fetching stories:", error);
    throw error;
  }
};

export const fetchStoryById = async (storyId: string): Promise<Story> => {
  try {
    const response = await axios.get<StoryResponse>(
      `${API_BASE_URL}/stories/${storyId}`
    );
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.error || "Failed to fetch story");
  } catch (error) {
    console.error("Error fetching story:", error);
    throw error;
  }
};

// Quiz API
export const fetchQuizByStoryId = async (storyId: string): Promise<Quiz> => {
  try {
    const response = await axios.get<QuizResponse>(
      `${API_BASE_URL}/quiz/${storyId}`
    );
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.error || "Failed to fetch quiz");
  } catch (error) {
    console.error("Error fetching quiz:", error);
    throw error;
  }
};

export const submitQuiz = async (
  storyId: string,
  answers: number[]
): Promise<QuizResult> => {
  try {
    const response = await axios.post<QuizSubmissionResponse>(
      `${API_BASE_URL}/quiz/${storyId}/submit`,
      { answers }
    );
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.error || "Failed to submit quiz");
  } catch (error) {
    console.error("Error submitting quiz:", error);
    throw error;
  }
};

// Helper functions
export const formatLevel = (level: string): string => {
  const levelMap: Record<string, string> = {
    A1: "Beginner",
    A2: "Elementary",
    B1: "Intermediate",
    B2: "Upper Intermediate",
    C1: "Advanced",
    C2: "Proficient",
  };
  return levelMap[level] || level;
};

export const getLevelColor = (level: string): string => {
  const colorMap: Record<string, string> = {
    A1: "bg-green-100 text-green-800",
    A2: "bg-blue-100 text-blue-800",
    B1: "bg-yellow-100 text-yellow-800",
    B2: "bg-orange-100 text-orange-800",
    C1: "bg-red-100 text-red-800",
    C2: "bg-purple-100 text-purple-800",
  };
  return colorMap[level] || "bg-gray-100 text-gray-800";
};

export const calculateReadingTime = (text: string): number => {
  const wordsPerMinute = 200; // Average reading speed
  const wordCount = text.trim().split(/\s+/).length;
  return Math.ceil(wordCount / wordsPerMinute);
};

export const getQuizPassingScore = (): number => {
  return 70; // 70% to pass
};
