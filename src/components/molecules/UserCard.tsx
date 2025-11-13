import React from "react";
import { Card, CardContent } from "../ui/card";
import { Button } from "../ui/button";
import { RoleBadge, UserStatusBadge } from "../atoms/UserBadges";
import { User } from "../../types/index";
import { UserAction } from "../../hooks/useUserManagement";

interface UserCardProps {
  user: User;
  currentUserId?: string;
  onAction: (userId: string, action: UserAction) => void;
  isAdmin: boolean;
}

interface UserActionConfig {
  action: UserAction;
  label: string;
  variant: "default" | "destructive" | "outline";
}

const getAvailableActions = (user: User, currentUserId?: string): UserActionConfig[] => {
  if (user.id === currentUserId) return [];

  const actions: UserActionConfig[] = [];

  if (user.status === "active") {
    actions.push({
      action: "suspend",
      label: "Suspend",
      variant: "destructive",
    });
  } else if (user.status === "suspended") {
    actions.push({
      action: "activate",
      label: "Activate",
      variant: "default",
    });
  }

  if (user.role === "creator") {
    actions.push({
      action: "promote",
      label: "Promote to Reviewer",
      variant: "outline",
    });
  } else if (user.role === "reviewer") {
    actions.push({
      action: "promote",
      label: "Promote to Admin",
      variant: "outline",
    });
    actions.push({
      action: "demote",
      label: "Demote to Creator",
      variant: "outline",
    });
  } else if (user.role === "admin" && user.id !== currentUserId) {
    actions.push({
      action: "demote",
      label: "Demote to Reviewer",
      variant: "outline",
    });
  }

  actions.push({
    action: "delete",
    label: "Delete User",
    variant: "destructive",
  });

  return actions;
};

export const UserCard: React.FC<UserCardProps> = ({
  user,
  currentUserId,
  onAction,
  isAdmin,
}) => {
  const actions = getAvailableActions(user, currentUserId);
  const isCurrentUser = user.id === currentUserId;

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardContent className="pt-6">
        <div className="flex justify-between items-start">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <RoleBadge role={user.role} />
              <UserStatusBadge status={user.status} />
              {isCurrentUser && (
                <span className="text-xs text-muted-foreground">(You)</span>
              )}
            </div>
            <h3 className="text-lg font-semibold mb-1">
              {user.firstName} {user.lastName}
            </h3>
            <p className="text-sm text-muted-foreground mb-1">
              @{user.username}
            </p>
            <p className="text-sm text-muted-foreground mb-2">{user.email}</p>
            <div className="text-xs text-muted-foreground">
              Joined: {new Date(user.createdAt).toLocaleDateString()}
              {user.lastLoginAt && (
                <>
                  {" • Last login: "}
                  {new Date(user.lastLoginAt).toLocaleDateString()}
                </>
              )}
            </div>
          </div>
          {isAdmin && actions.length > 0 && (
            <div className="flex flex-col gap-2">
              {actions.map((action) => (
                <Button
                  key={action.action}
                  variant={action.variant}
                  size="sm"
                  onClick={() => onAction(user.id, action.action)}
                >
                  {action.label}
                </Button>
              ))}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
};
