import { useState, useMemo } from "react";
import { User } from "../types/index";

export const useUserFilters = (users: User[]) => {
  const [roleFilter, setRoleFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");

  const filteredUsers = useMemo(() => {
    let filtered = users;

    if (roleFilter && roleFilter !== "all") {
      filtered = filtered.filter((user) => user.role === roleFilter);
    }

    if (statusFilter && statusFilter !== "all") {
      filtered = filtered.filter((user) => user.status === statusFilter);
    }

    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (user) =>
          user.username.toLowerCase().includes(query) ||
          user.email.toLowerCase().includes(query) ||
          user.firstName.toLowerCase().includes(query) ||
          user.lastName.toLowerCase().includes(query)
      );
    }

    return filtered;
  }, [users, roleFilter, statusFilter, searchQuery]);

  const clearFilters = () => {
    setRoleFilter("all");
    setStatusFilter("all");
    setSearchQuery("");
  };

  return {
    filteredUsers,
    roleFilter,
    statusFilter,
    searchQuery,
    setRoleFilter,
    setStatusFilter,
    setSearchQuery,
    clearFilters,
  };
};
