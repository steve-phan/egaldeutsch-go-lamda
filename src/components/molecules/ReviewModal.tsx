import React, { useState } from "react";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "../ui/card";
import { Button } from "../ui/button";
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";
import { ContentItem, ReviewData } from "../../types/content";

interface ReviewModalProps {
  content: ContentItem | null;
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: ReviewData) => Promise<void>;
}

const CONTENT_STATUSES = [
  { value: "draft", label: "Draft" },
  { value: "preview", label: "Preview" },
  { value: "ready", label: "Ready" },
  { value: "published", label: "Published" },
];

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

export const ReviewModal: React.FC<ReviewModalProps> = ({
  content,
  isOpen,
  onClose,
  onSubmit,
}) => {
  const [reviewComment, setReviewComment] = useState("");
  const [newStatus, setNewStatus] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  React.useEffect(() => {
    if (content) {
      const nextActions = getNextActions(content.status);
      setNewStatus(nextActions[0] || "");
    }
  }, [content]);

  if (!isOpen || !content) return null;

  const handleSubmit = async () => {
    if (!newStatus) {
      setError("Please select an action");
      return;
    }

    try {
      setIsSubmitting(true);
      setError("");

      await onSubmit({
        status: newStatus,
        comment: reviewComment,
      });

      // Reset form
      setReviewComment("");
      setNewStatus("");
    } catch (err: any) {
      setError(err.message || "Failed to update status");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Review Content</CardTitle>
          <CardDescription>
            Update status for: {content.title || content.question}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {error && (
            <div className="text-sm text-red-600 bg-red-50 p-2 rounded">
              {error}
            </div>
          )}
          <div>
            <Label htmlFor="new-status">New Status</Label>
            <select
              id="new-status"
              value={newStatus}
              onChange={(e) => setNewStatus(e.target.value)}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              {getNextActions(content.status).map((status) => {
                const statusInfo = CONTENT_STATUSES.find(
                  (s) => s.value === status
                );
                return (
                  <option key={status} value={status}>
                    {statusInfo?.label || status}
                  </option>
                );
              })}
            </select>
          </div>
          <div>
            <Label htmlFor="review-comment">Review Comment</Label>
            <Textarea
              id="review-comment"
              value={reviewComment}
              onChange={(e) => setReviewComment(e.target.value)}
              placeholder="Add a comment about this review decision..."
              rows={3}
            />
          </div>
          <div className="flex gap-2">
            <Button
              onClick={handleSubmit}
              disabled={isSubmitting}
              className="flex-1"
            >
              {isSubmitting ? "Updating..." : "Update Status"}
            </Button>
            <Button
              variant="outline"
              onClick={() => {
                onClose();
                setReviewComment("");
                setNewStatus("");
                setError("");
              }}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
