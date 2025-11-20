import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios";

// Configuration
const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? "/.netlify/functions"
    : "http://localhost:8888/.netlify/functions";

// Storage keys
export const TOKEN_KEY = "egaldeutsch_auth_token";
export const USER_KEY = "egaldeutsch_user";

// Create axios instance with default config
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 60000, // 60 seconds to accommodate AI generation
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    // Only add token if it exists
    const token = localStorage.getItem(TOKEN_KEY);

    console.log("Request interceptor - token:", token ? "present" : "null");

    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle 401 errors (token expired/invalid)
    if (error.response?.status === 401) {
      // Clear invalid token
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem("egaldeutsch_user");

      // Redirect to login if not already there
      if (
        typeof window !== "undefined" &&
        !window.location.pathname.includes("/auth")
      ) {
        window.location.href = "/auth/login";
      }
    }
    return Promise.reject(error);
  }
);

// Public API methods (no auth required)
export const publicApi = {
  // Stories
  getStories: (status?: string) =>
    apiClient.get(status ? `/stories?status=${status}` : "/stories"),
  getStoryById: (id: string) => apiClient.get(`/stories?id=${id}`),
  getStoryBySlug: (slug: string) => apiClient.get(`/stories/${slug}`),

  // Quiz (public viewing)
  getQuizByStoryId: (storyId: string) =>
    apiClient.get(`/quiz?story_id=${storyId}`),
  getQuizByStorySlug: (slug: string) =>
    apiClient.get(`/quiz?story_slug=${slug}`),
  submitQuiz: (storyId: string, answers: number[]) =>
    apiClient.post(`/quiz?story_id=${storyId}&action=submit`, { answers }),

  // Questions
  getQuestionsByStoryId: (storyId: string) =>
    apiClient.get(`/questions?storyId=${storyId}`),

  // Auth endpoints
  login: (credentials: any) =>
    apiClient.post("/user-management/login", credentials),
  register: (data: any) => apiClient.post("/user-management/register", data),

  // Leaderboard (public)
  getLeaderboard: () => apiClient.get("/leaderboard"),

  // Health check
  health: () => apiClient.get("/health"),
};

// Protected API methods (auth required)
export const protectedApi = {
  // User management (admin only)
  getUsers: () => apiClient.get("/user-management"),
  updateUser: (userId: string, data: any) =>
    apiClient.put(`/user-management/${userId}`, data),
  deleteUser: (userId: string) =>
    apiClient.delete(`/user-management/${userId}`),

  // Content management (creators/reviewers/admins)
  createStory: (data: any) => apiClient.post("/stories-management", data),
  updateStory: (storyId: string, data: any) =>
    apiClient.put(`/stories-management/${storyId}`, data),
  deleteStory: (storyId: string) =>
    apiClient.delete(`/stories-management/${storyId}`),

  // AI generation (admins)
  generateQuestions: (storyId: string) =>
    apiClient.post(`/ai-generator?type=questions&story_id=${storyId}`),

  // Quiz management
  getQuizs: () => apiClient.get("/quiz-management"),
  createQuiz: (data: any) => apiClient.post("/quiz-management", data),
  updateQuiz: (quizId: string, data: any) =>
    apiClient.put(`/quiz-management/${quizId}`, data),

  // Questions management
  getQuestions: () => apiClient.get("/questions-management"),
  createQuestion: (data: any) => apiClient.post("/questions-management", data),
  updateQuestion: (questionId: string, data: any) =>
    apiClient.put(`/questions-management/${questionId}`, data),
  deleteQuestion: (questionId: string) =>
    apiClient.delete(`/questions-management/${questionId}`),

  // Stories with filtering (for admin use)
  getStories: (status?: string) =>
    apiClient.get(
      status ? `/stories-management?status=${status}` : "/stories-management"
    ),

  updateStoryStatus: (id: string, status: string) =>
    apiClient.patch(`/stories-management/${id}`, { status }),

  // Specific method for AI generation filtering
  getStoriesForAIGeneration: () =>
    apiClient.get(
      "/stories-management?status=published&aiQuestionsGenerated=false"
    ),
};

// Utility functions
export const isAuthenticated = (): boolean => {
  return !!localStorage.getItem(TOKEN_KEY);
};

export const getAuthToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY);
};

// Export the client for advanced usage
export { apiClient };
export default apiClient;
