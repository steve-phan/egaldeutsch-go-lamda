import React from "react";
import { Button } from "../ui/button";
import { FilterSelect } from "../atoms/FilterSelect";
import { SearchInput } from "../atoms/SearchInput";
import { Card, CardContent } from "../ui/card";

interface FilterBarProps {
  statusFilter: string;
  typeFilter: string;
  searchTerm: string;
  onStatusChange: (status: string) => void;
  onTypeChange: (type: string) => void;
  onSearchChange: (term: string) => void;
  onRefresh: () => void;
  loading?: boolean;
}

const STATUS_OPTIONS = [
  { value: "all", label: "All Statuses" },
  { value: "draft", label: "Draft" },
  { value: "preview", label: "Preview" },
  { value: "ready", label: "Ready" },
  { value: "published", label: "Published" },
];

const TYPE_OPTIONS = [
  { value: "all", label: "All Types" },
  { value: "story", label: "Stories" },
  { value: "question", label: "Questions" },
  { value: "quiz", label: "Quizzes" },
];

export const FilterBar: React.FC<FilterBarProps> = ({
  statusFilter,
  typeFilter,
  searchTerm,
  onStatusChange,
  onTypeChange,
  onSearchChange,
  onRefresh,
  loading = false,
}) => {
  return (
    <Card className="mb-6">
      <CardContent className="pt-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <FilterSelect
            id="status-filter"
            label="Filter by Status"
            value={statusFilter}
            onChange={onStatusChange}
            options={STATUS_OPTIONS}
            placeholder="All Statuses"
          />
          <FilterSelect
            id="type-filter"
            label="Filter by Type"
            value={typeFilter}
            onChange={onTypeChange}
            options={TYPE_OPTIONS}
            placeholder="All Types"
          />
          <SearchInput
            id="search-input"
            label="Search Content"
            placeholder="Search by title, content, topics..."
            value={searchTerm}
            onChange={onSearchChange}
          />
          <div className="flex items-end">
            <Button onClick={onRefresh} disabled={loading}>
              {loading ? "Loading..." : "Refresh"}
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
