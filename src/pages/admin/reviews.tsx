import React, { useState, useEffect } from "react";
import Layout from "../../components/layout";
import ProtectedRoute from "../../components/auth/ProtectedRoute";
import { useAuth } from "../../contexts/AuthContext";
import { ErrorAlert } from "../../components/atoms/ErrorAlert";
import { FilterBar } from "../../components/molecules/FilterBar";
import { ReviewModal } from "../../components/molecules/ReviewModal";
import { ContentReviewList } from "../../components/organisms/ContentReviewList";
import { useContentReview } from "../../hooks/useContentReview";
import { useContentFilters } from "../../hooks/useContentFilters";
import { ContentItem, ReviewData } from "../../types/content";

const ReviewsPage: React.FC = () => {
  const { user } = useAuth();
  const {
    contentItems,
    loading,
    error,
    loadContentItems,
    updateContentStatus,
  } = useContentReview();
  const {
    filteredItems,
    statusFilter,
    typeFilter,
    searchTerm,
    setStatusFilter,
    setTypeFilter,
    setSearchTerm,
  } = useContentFilters(contentItems);
  const [selectedContent, setSelectedContent] = useState<ContentItem | null>(null);

  useEffect(() => {
    loadContentItems();
  }, [loadContentItems]);

  const handleReview = (content: ContentItem) => {
    setSelectedContent(content);
  };

  const handleStatusUpdate = async (data: ReviewData) => {
    if (!selectedContent) return;

    await updateContentStatus(
      selectedContent.id,
      selectedContent.type,
      data.status,
      data.comment
    );
    setSelectedContent(null);
    await loadContentItems();
  };

  const canReview = user?.role === "admin" || user?.role === "reviewer";

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

        <ErrorAlert message={error} />

        <FilterBar
          statusFilter={statusFilter}
          typeFilter={typeFilter}
          searchTerm={searchTerm}
          onStatusChange={setStatusFilter}
          onTypeChange={setTypeFilter}
          onSearchChange={setSearchTerm}
          onRefresh={loadContentItems}
          loading={loading}
        />

        <ContentReviewList
          contents={filteredItems}
          loading={loading}
          error={error}
          onReview={handleReview}
          canReview={canReview}
        />

        <ReviewModal
          content={selectedContent}
          isOpen={!!selectedContent}
          onClose={() => setSelectedContent(null)}
          onSubmit={handleStatusUpdate}
        />
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
