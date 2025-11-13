import { useState, useMemo } from "react";
import { ContentItem } from "../types/content";

export const useContentFilters = (items: ContentItem[]) => {
  const [statusFilter, setStatusFilter] = useState("preview");
  const [typeFilter, setTypeFilter] = useState("all");
  const [searchTerm, setSearchTerm] = useState("");

  const filteredItems = useMemo(() => {
    let filtered = items;

    // Filter by status
    if (statusFilter && statusFilter !== "all") {
      filtered = filtered.filter((item) => item.status === statusFilter);
    }

    // Filter by type
    if (typeFilter && typeFilter !== "all") {
      filtered = filtered.filter((item) => item.type === typeFilter);
    }

    // Filter by search term
    if (searchTerm) {
      const searchLower = searchTerm.toLowerCase();
      filtered = filtered.filter((item) => {
        // Handle different content types with their specific fields
        if (item.type === "story") {
          return (
            item.title?.toLowerCase().includes(searchLower) ||
            item.content?.toLowerCase().includes(searchLower) ||
            item.topics?.some((topic: string) =>
              topic.toLowerCase().includes(searchLower)
            )
          );
        } else if (item.type === "question") {
          return (
            item.question?.toLowerCase().includes(searchLower) ||
            item.questionType?.toLowerCase().includes(searchLower) ||
            item.difficulty?.toLowerCase().includes(searchLower)
          );
        } else if (item.type === "quiz") {
          return (
            item.title?.toLowerCase().includes(searchLower) ||
            item.description?.toLowerCase().includes(searchLower)
          );
        }
        return false;
      });
    }

    return filtered;
  }, [items, statusFilter, typeFilter, searchTerm]);

  const clearFilters = () => {
    setStatusFilter("all");
    setTypeFilter("all");
    setSearchTerm("");
  };

  return {
    filteredItems,
    statusFilter,
    typeFilter,
    searchTerm,
    setStatusFilter,
    setTypeFilter,
    setSearchTerm,
    clearFilters,
  };
};
