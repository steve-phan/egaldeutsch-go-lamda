export type ContentType = "story" | "question" | "quiz";
export type ContentStatus = "draft" | "preview" | "ready" | "published";

export interface ContentItem {
  id: string;
  title?: string;
  question?: string;
  description?: string;
  content?: string;
  type: ContentType;
  level?: string;
  topics?: string[];
  storyId?: string;
  questionType?: string;
  difficulty?: string;
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

export interface ReviewData {
  status: string;
  comment?: string;
}

export interface DashboardStats {
  totalContent: number;
  pendingReviews: number;
  publishedContent: number;
  draftContent: number;
}

export interface Activity {
  id: string;
  type: string;
  description: string;
  timestamp: string;
}
