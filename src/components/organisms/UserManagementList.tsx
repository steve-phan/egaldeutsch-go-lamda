import React from "react";
import { UserCard } from "../molecules/UserCard";
import { LoadingSpinner } from "../atoms/LoadingSpinner";
import { EmptyState } from "../atoms/EmptyState";
import { User } from "../../types/index";
import { UserAction } from "../../hooks/useUserManagement";

interface UserManagementListProps {
  users: User[];
  loading: boolean;
  currentUserId?: string;
  onUserAction: (userId: string, action: UserAction) => void;
  isAdmin: boolean;
}

export const UserManagementList: React.FC<UserManagementListProps> = ({
  users,
  loading,
  currentUserId,
  onUserAction,
  isAdmin,
}) => {
  if (loading) {
    return <LoadingSpinner message="Loading users..." />;
  }

  if (users.length === 0) {
    return <EmptyState message="No users found for the selected filters." />;
  }

  return (
    <div className="grid gap-4">
      {users.map((user) => (
        <UserCard
          key={user.id}
          user={user}
          currentUserId={currentUserId}
          onAction={onUserAction}
          isAdmin={isAdmin}
        />
      ))}
    </div>
  );
};
