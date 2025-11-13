import { useState, useCallback } from "react";
import axios from "axios";
import { ContentItem } from "../types/content";

const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? "/.netlify/functions"
    : "http://localhost:8888/.netlify/functions";

export const useContentReview = () => {
  const [contentItems, setContentItems] = useState<ContentItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>("");

  const loadContentItems = useCallback(async () => {
    try {
      setLoading(true);
      setError("");

      // Load all content types with proper error handling
      const [storiesRes, questionsRes, quizzesRes] = await Promise.all([
        axios
          .get(`${API_BASE_URL}/stories-management`)
          .catch(() => ({ data: { stories: [] } })),
        axios
          .get(`${API_BASE_URL}/questions-management`)
          .catch(() => ({ data: { questions: [] } })),
        axios
          .get(`${API_BASE_URL}/quiz-management`)
          .catch(() => ({ data: { quizzes: [] } })),
      ]);

      const allItems: ContentItem[] = [];

      // Process stories
      if (storiesRes.data?.stories) {
        const stories = storiesRes.data.stories.map((item: any) => ({
          ...item,
          type: "story" as const,
        }));
        allItems.push(...stories);
      }

      // Process questions
      if (questionsRes.data?.questions) {
        const questions = questionsRes.data.questions.map((item: any) => ({
          ...item,
          type: "question" as const,
        }));
        allItems.push(...questions);
      }

      // Process quizzes
      if (quizzesRes.data?.quizzes) {
        const quizzes = quizzesRes.data.quizzes.map((item: any) => ({
          ...item,
          type: "quiz" as const,
        }));
        allItems.push(...quizzes);
      }

      setContentItems(allItems);
    } catch (err) {
      console.error("Error loading content items:", err);
      setError("Failed to load content items");
    } finally {
      setLoading(false);
    }
  }, []);

  const updateContentStatus = useCallback(
    async (id: string, type: string, status: string, comment?: string) => {
      try {
        const endpoint =
          type === "story"
            ? "stories-management"
            : type === "question"
            ? "questions-management"
            : "quiz-management";

        await axios.patch(`${API_BASE_URL}/${endpoint}/${id}/status`, {
          status,
          comment,
        });

        return true;
      } catch (err: any) {
        console.error("Error updating status:", err);
        throw new Error(err.response?.data?.error || "Failed to update status");
      }
    },
    []
  );

  return {
    contentItems,
    loading,
    error,
    loadContentItems,
    updateContentStatus,
    refetch: loadContentItems,
  };
};
