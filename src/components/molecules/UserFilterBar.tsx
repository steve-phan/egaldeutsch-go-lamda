import React from "react";
import { Button } from "../ui/button";
import { FilterSelect } from "../atoms/FilterSelect";
import { SearchInput } from "../atoms/SearchInput";
import { Card, CardContent } from "../ui/card";

interface UserFilterBarProps {
  roleFilter: string;
  statusFilter: string;
  searchQuery: string;
  onRoleChange: (role: string) => void;
  onStatusChange: (status: string) => void;
  onSearchChange: (query: string) => void;
  onRefresh: () => void;
  loading?: boolean;
}

const ROLE_OPTIONS = [
  { value: "all", label: "All Roles" },
  { value: "creator", label: "Creator" },
  { value: "reviewer", label: "Reviewer" },
  { value: "admin", label: "Admin" },
];

const STATUS_OPTIONS = [
  { value: "all", label: "All Statuses" },
  { value: "active", label: "Active" },
  { value: "suspended", label: "Suspended" },
  { value: "pending", label: "Pending" },
];

export const UserFilterBar: React.FC<UserFilterBarProps> = ({
  roleFilter,
  statusFilter,
  searchQuery,
  onRoleChange,
  onStatusChange,
  onSearchChange,
  onRefresh,
  loading = false,
}) => {
  return (
    <Card className="mb-6">
      <CardContent className="pt-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <FilterSelect
            id="role-filter"
            label="Filter by Role"
            value={roleFilter}
            onChange={onRoleChange}
            options={ROLE_OPTIONS}
            placeholder="All Roles"
          />
          <FilterSelect
            id="status-filter"
            label="Filter by Status"
            value={statusFilter}
            onChange={onStatusChange}
            options={STATUS_OPTIONS}
            placeholder="All Statuses"
          />
          <SearchInput
            id="search-input"
            label="Search Users"
            placeholder="Search by name, email, username..."
            value={searchQuery}
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
