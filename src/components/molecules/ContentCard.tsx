import React from "react";
import { navigate } from "gatsby";
import { Card, CardContent } from "../ui/card";
import { Button } from "../ui/button";
import { Badge } from "../ui/badge";
import { StatusBadge } from "../atoms/StatusBadge";
import { ContentItem } from "../../types/content";

interface ContentCardProps {
  content: ContentItem;
  onReview: (content: ContentItem) => void;
  canReview: boolean;
  showReviewButton: boolean;
}

export const ContentCard: React.FC<ContentCardProps> = ({
  content,
  onReview,
  canReview,
  showReviewButton,
}) => {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardContent className="pt-6">
        <div className="flex justify-between items-start mb-4">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <Badge variant="outline" className="text-xs">
                {content.type.toUpperCase()}
              </Badge>
              <StatusBadge status={content.status} />
              {content.level && (
                <Badge variant="secondary" className="text-xs">
                  {content.level}
                </Badge>
              )}
            </div>
            <h3 className="text-lg font-semibold mb-2">
              {content.title || content.question || "Untitled"}
            </h3>
            {content.topics && (
              <div className="flex flex-wrap gap-1 mb-2">
                {content.topics.map((topic, index) => (
                  <Badge key={index} variant="outline" className="text-xs">
                    {topic}
                  </Badge>
                ))}
              </div>
            )}
            <div className="text-sm text-muted-foreground">
              Created: {new Date(content.createdAt).toLocaleDateString()} •
              Version: {content.version}
              {content.reviewedAt && (
                <>
                  {" "}
                  • Reviewed:{" "}
                  {new Date(content.reviewedAt).toLocaleDateString()}
                </>
              )}
            </div>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate(`/${content.type}/${content.id}`)}
            >
              View
            </Button>
            {canReview && showReviewButton && (
              <Button size="sm" onClick={() => onReview(content)}>
                Review
              </Button>
            )}
          </div>
        </div>

        {content.comments && content.comments.length > 0 && (
          <div className="mt-4 pt-4 border-t">
            <h4 className="text-sm font-medium mb-2">Recent Comments:</h4>
            <div className="space-y-2">
              {content.comments.slice(-2).map((comment) => (
                <div key={comment.id} className="bg-muted p-3 rounded-md">
                  <p className="text-sm">{comment.comment}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {new Date(comment.createdAt).toLocaleDateString()}
                  </p>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};
