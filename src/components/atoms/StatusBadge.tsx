import React from "react";
import { Badge } from "../ui/badge";
import { ContentStatus, Size } from "../../types/common";

interface StatusBadgeProps {
  status: ContentStatus | string;
  size?: Size;
}

const STATUS_COLORS: Record<string, string> = {
  draft: "bg-gray-100 text-gray-800",
  preview: "bg-yellow-100 text-yellow-800",
  ready: "bg-green-100 text-green-800",
  published: "bg-blue-100 text-blue-800",
};

const STATUS_LABELS: Record<string, string> = {
  draft: "Draft",
  preview: "Preview",
  ready: "Ready",
  published: "Published",
};

export const StatusBadge: React.FC<StatusBadgeProps> = ({ status, size = "md" }) => {
  const colorClass = STATUS_COLORS[status] || "bg-gray-100 text-gray-800";
  const label = STATUS_LABELS[status] || status;
  
  return (
    <Badge className={colorClass}>
      {label}
    </Badge>
  );
};
