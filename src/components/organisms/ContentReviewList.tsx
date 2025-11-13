import React from "react";
import { ContentCard } from "../molecules/ContentCard";
import { LoadingSpinner } from "../atoms/LoadingSpinner";
import { EmptyState } from "../atoms/EmptyState";
import { ContentItem } from "../../types/content";

interface ContentReviewListProps {
  contents: ContentItem[];
  loading: boolean;
  error?: string;
  onReview: (content: ContentItem) => void;
  canReview: boolean;
}

const getNextActions = (currentStatus: string): string[] => {
  switch (currentStatus) {
    case "draft":
      return ["pending_review"];
    case "pending_review":
      return ["approved", "rejected"];
    case "approved":
      return ["active", "rejected"];
    case "active":
      return ["inactive", "archived"];
    case "inactive":
      return ["active", "archived"];
    case "rejected":
      return ["pending_review", "archived"];
    default:
      return [];
  }
};

export const ContentReviewList: React.FC<ContentReviewListProps> = ({
  contents,
  loading,
  error,
  onReview,
  canReview,
}) => {
  if (loading) {
    return <LoadingSpinner message="Loading content..." />;
  }

  if (contents.length === 0) {
    return <EmptyState message="No content items found for the selected filters." />;
  }

  return (
    <div className="grid gap-4">
      {contents.map((item) => (
        <ContentCard
          key={item.id}
          content={item}
          onReview={onReview}
          canReview={canReview}
          showReviewButton={getNextActions(item.status).length > 0}
        />
      ))}
    </div>
  );
};
