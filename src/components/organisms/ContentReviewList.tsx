import React from "react";
import { ContentCard } from "../molecules/ContentCard";
import { LoadingSpinner } from "../atoms/LoadingSpinner";
import { EmptyState } from "../atoms/EmptyState";
import { ContentItem, ContentStatus } from "../../types/content";

interface ContentReviewListProps {
  contents: ContentItem[];
  loading: boolean;
  error?: string;
  onReview: (content: ContentItem) => void;
  canReview: boolean;
}

const getNextActions = (currentStatus: string): string[] => {
  switch (currentStatus) {
    case ContentStatus.DRAFT:
      return [ContentStatus.PREVIEW];
    case ContentStatus.PREVIEW:
      return [ContentStatus.READY, ContentStatus.DRAFT];
    case ContentStatus.READY:
      return [ContentStatus.PUBLISHED, ContentStatus.DRAFT];
    case ContentStatus.PUBLISHED:
      return [ContentStatus.DRAFT];
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
    return (
      <EmptyState message="No content items found for the selected filters." />
    );
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
