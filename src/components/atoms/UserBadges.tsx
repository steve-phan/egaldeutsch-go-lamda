import React from "react";
import { Badge } from "../ui/badge";
import { UserRole, UserStatus } from "../../types/index";

interface RoleBadgeProps {
  role: UserRole;
}

const ROLE_COLORS: Record<UserRole, string> = {
  creator: "bg-blue-100 text-blue-800",
  reviewer: "bg-green-100 text-green-800",
  admin: "bg-purple-100 text-purple-800",
};

const ROLE_LABELS: Record<UserRole, string> = {
  creator: "Creator",
  reviewer: "Reviewer",
  admin: "Admin",
};

export const RoleBadge: React.FC<RoleBadgeProps> = ({ role }) => {
  const colorClass = ROLE_COLORS[role] || "bg-gray-100 text-gray-800";
  const label = ROLE_LABELS[role] || role;

  return <Badge className={colorClass}>{label}</Badge>;
};

interface UserStatusBadgeProps {
  status: UserStatus;
}

const STATUS_COLORS: Record<UserStatus, string> = {
  active: "bg-green-100 text-green-800",
  suspended: "bg-red-100 text-red-800",
  pending: "bg-yellow-100 text-yellow-800",
};

const STATUS_LABELS: Record<UserStatus, string> = {
  active: "Active",
  suspended: "Suspended",
  pending: "Pending",
};

export const UserStatusBadge: React.FC<UserStatusBadgeProps> = ({ status }) => {
  const colorClass = STATUS_COLORS[status] || "bg-gray-100 text-gray-800";
  const label = STATUS_LABELS[status] || status;

  return <Badge className={colorClass}>{label}</Badge>;
};
