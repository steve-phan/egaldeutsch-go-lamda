import React, { useState, useEffect } from "react";
import { navigate } from "gatsby";
import Layout from "../../components/layout";
import ProtectedRoute from "../../components/auth/ProtectedRoute";
import { useAuth } from "../../contexts/AuthContext";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { Label } from "../../components/ui/label";
import { Textarea } from "../../components/ui/textarea";
import { Badge } from "../../components/ui/badge";
import { Alert, AlertDescription } from "../../components/ui/alert";
import { Separator } from "../../components/ui/separator";
import axios from "axios";

const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? "/.netlify/functions"
    : "http://localhost:8888/.netlify/functions";

interface ContentItem {
  id: string;
  title: string;
  content?: string;
  question?: string;
  type: "story" | "question" | "quiz";
  level?: string;
  topics?: string[];
  status: string;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
  reviewedBy?: string;
  reviewedAt?: string;
  comments?: Array<{
    id: string;
    reviewerId: string;
    comment: string;
    createdAt: string;
    type: string;
  }>;
  version: number;
}

const CONTENT_STATUSES = [
  { value: "draft", label: "Draft", color: "bg-gray-100 text-gray-800" },
  {
    value: "pending_review",
    label: "Pending Review",
    color: "bg-yellow-100 text-yellow-800",
  },
  {
    value: "approved",
    label: "Approved",
    color: "bg-green-100 text-green-800",
  },
  { value: "active", label: "Active", color: "bg-blue-100 text-blue-800" },
  { value: "inactive", label: "Inactive", color: "bg-gray-100 text-gray-600" },
  { value: "rejected", label: "Rejected", color: "bg-red-100 text-red-800" },
  {
    value: "archived",
    label: "Archived",
    color: "bg-purple-100 text-purple-800",
  },
];

const ReviewsPage: React.FC = () => {
  const { user } = useAuth();
  const [contentItems, setContentItems] = useState<ContentItem[]>([]);
  const [filteredItems, setFilteredItems] = useState<ContentItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selectedStatus, setSelectedStatus] = useState("pending_review");
  const [selectedType, setSelectedType] = useState("all");
  const [selectedItem, setSelectedItem] = useState<ContentItem | null>(null);
  const [reviewComment, setReviewComment] = useState("");
  const [newStatus, setNewStatus] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    loadContentItems();
  }, []);

  useEffect(() => {
    filterItems();
  }, [contentItems, selectedStatus, selectedType]);

  const loadContentItems = async () => {
    try {
      setLoading(true);

      // Load all content types
      const [storiesRes, questionsRes, quizzesRes] = await Promise.all([
        axios.get(`${API_BASE_URL}/stories-management`),
        axios.get(`${API_BASE_URL}/questions-management`),
        axios.get(`${API_BASE_URL}/quiz-management`),
      ]);

      const allItems: ContentItem[] = [
        ...(storiesRes.data.stories?.map((item: any) => ({
          ...item,
          type: "story" as const,
        })) || []),
        ...(questionsRes.data.questions?.map((item: any) => ({
          ...item,
          type: "question" as const,
        })) || []),
        ...(quizzesRes.data.quizzes?.map((item: any) => ({
          ...item,
          type: "quiz" as const,
        })) || []),
      ];

      setContentItems(allItems);
    } catch (err) {
      console.error("Error loading content items:", err);
      setError("Failed to load content items");
    } finally {
      setLoading(false);
    }
  };

  const filterItems = () => {
    let filtered = contentItems;

    if (selectedStatus && selectedStatus !== "all") {
      filtered = filtered.filter((item) => item.status === selectedStatus);
    }

    if (selectedType && selectedType !== "all") {
      filtered = filtered.filter((item) => item.type === selectedType);
    }

    setFilteredItems(filtered);
  };

  const handleStatusUpdate = async () => {
    if (!selectedItem || !newStatus) {
      setError("Please select an action");
      return;
    }

    try {
      setIsSubmitting(true);
      setError("");

      const endpoint =
        selectedItem.type === "story"
          ? "stories-management"
          : selectedItem.type === "question"
          ? "questions-management"
          : "quiz-management";

      await axios.patch(
        `${API_BASE_URL}/${endpoint}/${selectedItem.id}/status`,
        {
          status: newStatus,
          comment: reviewComment,
        }
      );

      // Reload content items
      await loadContentItems();

      // Close modal
      setSelectedItem(null);
      setReviewComment("");
      setNewStatus("");
    } catch (err: any) {
      console.error("Error updating status:", err);
      setError(err.response?.data?.error || "Failed to update status");
    } finally {
      setIsSubmitting(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const statusInfo = CONTENT_STATUSES.find((s) => s.value === status);
    return (
      <Badge className={statusInfo?.color || "bg-gray-100 text-gray-800"}>
        {statusInfo?.label || status}
      </Badge>
    );
  };

  const getNextActions = (currentStatus: string) => {
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

  return (
    <Layout>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">
            Content Review Dashboard
          </h1>
          <p className="text-muted-foreground">
            Review and manage content status across stories, questions, and
            quizzes
          </p>
        </div>

        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Filters */}
        <Card className="mb-6">
          <CardContent className="pt-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <Label htmlFor="status-filter">Filter by Status</Label>
                <select
                  id="status-filter"
                  value={selectedStatus}
                  onChange={(e) => setSelectedStatus(e.target.value)}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                >
                  <option value="all">All Statuses</option>
                  {CONTENT_STATUSES.map((status) => (
                    <option key={status.value} value={status.value}>
                      {status.label}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <Label htmlFor="type-filter">Filter by Type</Label>
                <select
                  id="type-filter"
                  value={selectedType}
                  onChange={(e) => setSelectedType(e.target.value)}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                >
                  <option value="all">All Types</option>
                  <option value="story">Stories</option>
                  <option value="question">Questions</option>
                  <option value="quiz">Quizzes</option>
                </select>
              </div>
              <div className="flex items-end">
                <Button onClick={loadContentItems} disabled={loading}>
                  {loading ? "Loading..." : "Refresh"}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Content Items */}
        {loading ? (
          <div className="text-center py-8">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            <p className="mt-2 text-muted-foreground">Loading content...</p>
          </div>
        ) : (
          <div className="grid gap-4">
            {filteredItems.length === 0 ? (
              <Card>
                <CardContent className="pt-6 text-center">
                  <p className="text-muted-foreground">
                    No content items found for the selected filters.
                  </p>
                </CardContent>
              </Card>
            ) : (
              filteredItems.map((item) => (
                <Card
                  key={item.id}
                  className="hover:shadow-md transition-shadow"
                >
                  <CardContent className="pt-6">
                    <div className="flex justify-between items-start mb-4">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          <Badge variant="outline" className="text-xs">
                            {item.type.toUpperCase()}
                          </Badge>
                          {getStatusBadge(item.status)}
                          {item.level && (
                            <Badge variant="secondary" className="text-xs">
                              {item.level}
                            </Badge>
                          )}
                        </div>
                        <h3 className="text-lg font-semibold mb-2">
                          {item.title || item.question || "Untitled"}
                        </h3>
                        {item.topics && (
                          <div className="flex flex-wrap gap-1 mb-2">
                            {item.topics.map((topic, index) => (
                              <Badge
                                key={index}
                                variant="outline"
                                className="text-xs"
                              >
                                {topic}
                              </Badge>
                            ))}
                          </div>
                        )}
                        <div className="text-sm text-muted-foreground">
                          Created:{" "}
                          {new Date(item.createdAt).toLocaleDateString()} •
                          Version: {item.version}
                          {item.reviewedAt && (
                            <>
                              {" "}
                              • Reviewed:{" "}
                              {new Date(item.reviewedAt).toLocaleDateString()}
                            </>
                          )}
                        </div>
                      </div>
                      <div className="flex gap-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => navigate(`/${item.type}/${item.id}`)}
                        >
                          View
                        </Button>
                        {(user?.role === "admin" ||
                          user?.role === "reviewer") &&
                          getNextActions(item.status).length > 0 && (
                            <Button
                              size="sm"
                              onClick={() => {
                                setSelectedItem(item);
                                setNewStatus(getNextActions(item.status)[0]);
                              }}
                            >
                              Review
                            </Button>
                          )}
                      </div>
                    </div>

                    {item.comments && item.comments.length > 0 && (
                      <div className="mt-4 pt-4 border-t">
                        <h4 className="text-sm font-medium mb-2">
                          Recent Comments:
                        </h4>
                        <div className="space-y-2">
                          {item.comments.slice(-2).map((comment) => (
                            <div
                              key={comment.id}
                              className="bg-muted p-3 rounded-md"
                            >
                              <p className="text-sm">{comment.comment}</p>
                              <p className="text-xs text-muted-foreground mt-1">
                                {new Date(
                                  comment.createdAt
                                ).toLocaleDateString()}
                              </p>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                  </CardContent>
                </Card>
              ))
            )}
          </div>
        )}

        {/* Review Modal */}
        {selectedItem && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <Card className="w-full max-w-md">
              <CardHeader>
                <CardTitle>Review Content</CardTitle>
                <CardDescription>
                  Update status for:{" "}
                  {selectedItem.title || selectedItem.question}
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="new-status">New Status</Label>
                  <select
                    id="new-status"
                    value={newStatus}
                    onChange={(e) => setNewStatus(e.target.value)}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  >
                    {getNextActions(selectedItem.status).map((status) => {
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
                    onClick={handleStatusUpdate}
                    disabled={isSubmitting}
                    className="flex-1"
                  >
                    {isSubmitting ? "Updating..." : "Update Status"}
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => {
                      setSelectedItem(null);
                      setReviewComment("");
                      setNewStatus("");
                    }}
                    disabled={isSubmitting}
                  >
                    Cancel
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </Layout>
  );
};

const AdminReviewsPage: React.FC = () => {
  return (
    <ProtectedRoute allowedRoles={["admin", "reviewer"]}>
      <ReviewsPage />
    </ProtectedRoute>
  );
};

export default AdminReviewsPage;
