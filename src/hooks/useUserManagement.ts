import { useState, useCallback } from "react";
import axios from "axios";
import { User } from "../types/index";

const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? "/.netlify/functions"
    : "http://localhost:8888/.netlify/functions";

export type UserAction = "suspend" | "activate" | "delete" | "promote" | "demote";

export const useUserManagement = (currentUserId?: string) => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>("");

  const loadUsers = useCallback(async () => {
    try {
      setLoading(true);
      setError("");
      const response = await axios.get(`${API_BASE_URL}/user-management`);
      setUsers(response.data);
    } catch (err: any) {
      console.error("Error loading users:", err);
      setError(err.response?.data?.error || "Failed to load users");
    } finally {
      setLoading(false);
    }
  }, []);

  const performUserAction = useCallback(
    async (userId: string, action: UserAction) => {
      if (userId === currentUserId && (action === "suspend" || action === "delete")) {
        throw new Error("Cannot perform this action on your own account");
      }

      try {
        switch (action) {
          case "suspend":
            await axios.put(`${API_BASE_URL}/user-management/${userId}`, {
              status: "suspended",
            });
            break;
          case "activate":
            await axios.put(`${API_BASE_URL}/user-management/${userId}`, {
              status: "active",
            });
            break;
          case "delete":
            if (
              window.confirm(
                "Are you sure you want to delete this user? This action cannot be undone."
              )
            ) {
              await axios.delete(`${API_BASE_URL}/user-management/${userId}`);
            } else {
              return false;
            }
            break;
          case "promote":
            const targetUser = users.find((u) => u.id === userId);
            if (targetUser) {
              const newRole =
                targetUser.role === "creator"
                  ? "reviewer"
                  : targetUser.role === "reviewer"
                  ? "admin"
                  : "admin";
              await axios.put(`${API_BASE_URL}/user-management/${userId}`, {
                role: newRole,
              });
            }
            break;
          case "demote":
            const targetUser2 = users.find((u) => u.id === userId);
            if (targetUser2) {
              const newRole =
                targetUser2.role === "admin"
                  ? "reviewer"
                  : targetUser2.role === "reviewer"
                  ? "creator"
                  : "creator";
              await axios.put(`${API_BASE_URL}/user-management/${userId}`, {
                role: newRole,
              });
            }
            break;
        }
        return true;
      } catch (err: any) {
        console.error("Error performing user action:", err);
        throw new Error(err.response?.data?.error || "Failed to perform action");
      }
    },
    [users, currentUserId]
  );

  return {
    users,
    loading,
    error,
    loadUsers,
    performUserAction,
    refetch: loadUsers,
  };
};
