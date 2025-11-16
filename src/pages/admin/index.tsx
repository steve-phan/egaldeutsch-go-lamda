import React, { useState, useEffect } from "react";
import { Link, navigate } from "gatsby";
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
import { Badge } from "../../components/ui/badge";
import { Alert, AlertDescription } from "../../components/ui/alert";
import { ContentStatus } from "../../types/content";
import axios from "axios";
import { protectedApi, publicApi } from "@/utils/apiClient";

interface DashboardStats {
  users: {
    total: number;
    active: number;
    creators: number;
    reviewers: number;
    admins: number;
  };
  content: {
    stories: { total: number; pending: number; active: number };
    questions: { total: number; pending: number; active: number };
    quizzes: { total: number; pending: number; active: number };
  };
}

const AdminDashboard: React.FC = () => {
  const { user } = useAuth();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    loadDashboardStats();
  }, []);

  const loadDashboardStats = async () => {
    try {
      setLoading(true);

      // Load all data in parallel
      const [usersRes, storiesRes, questionsRes, quizzesRes] =
        await Promise.all([
          protectedApi.getUsers(),
          publicApi.getStories(),
          protectedApi.getQuestions(),
          protectedApi.getQuizs(),
        ]);

      const users = usersRes.data || [];
      const stories = storiesRes.data?.stories || [];
      const questions = questionsRes.data?.questions || [];
      const quizzes = quizzesRes.data?.quizzes || [];

      const dashboardStats: DashboardStats = {
        users: {
          total: users.length,
          active: users.filter((u: any) => u.status === "active").length,
          creators: users.filter((u: any) => u.role === "creator").length,
          reviewers: users.filter((u: any) => u.role === "reviewer").length,
          admins: users.filter((u: any) => u.role === "admin").length,
        },
        content: {
          stories: {
            total: stories.length,
            pending: stories.filter(
              (s: any) => s.status === ContentStatus.PREVIEW
            ).length,
            active: stories.filter(
              (s: any) => s.status === ContentStatus.PUBLISHED
            ).length,
          },
          questions: {
            total: questions.length,
            pending: questions.filter(
              (q: any) => q.status === ContentStatus.PREVIEW
            ).length,
            active: questions.filter(
              (q: any) => q.status === ContentStatus.PUBLISHED
            ).length,
          },
          quizzes: {
            total: quizzes.length,
            pending: quizzes.filter(
              (q: any) => q.status === ContentStatus.PREVIEW
            ).length,
            active: quizzes.filter(
              (q: any) => q.status === ContentStatus.PUBLISHED
            ).length,
          },
        },
      };

      setStats(dashboardStats);
    } catch (err) {
      console.error("Error loading dashboard stats:", err);
      setError("Failed to load dashboard statistics");
    } finally {
      setLoading(false);
    }
  };

  const totalPendingReviews = stats
    ? stats.content.stories.pending +
      stats.content.questions.pending +
      stats.content.quizzes.pending
    : 0;

  return (
    <Layout>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">
            Admin Dashboard
          </h1>
          <p className="text-muted-foreground">
            Welcome back, {user?.firstName}! Here's an overview of your
            platform.
          </p>
        </div>

        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Quick Actions */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
          <Card
            className="hover:shadow-md transition-shadow cursor-pointer"
            onClick={() => navigate("/admin/reviews")}
          >
            <CardContent className="pt-6">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold mb-2">
                    Content Reviews
                  </h3>
                  <p className="text-muted-foreground text-sm">
                    Review and approve pending content
                  </p>
                  {totalPendingReviews > 0 && (
                    <div className="mt-2">
                      <Badge variant="destructive">
                        {totalPendingReviews} pending review
                        {totalPendingReviews !== 1 ? "s" : ""}
                      </Badge>
                    </div>
                  )}
                </div>
                <div className="text-3xl">📋</div>
              </div>
            </CardContent>
          </Card>

          <Card
            className="hover:shadow-md transition-shadow cursor-pointer"
            onClick={() => navigate("/admin/users")}
          >
            <CardContent className="pt-6">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold mb-2">
                    User Management
                  </h3>
                  <p className="text-muted-foreground text-sm">
                    Manage user accounts and permissions
                  </p>
                  {stats && (
                    <div className="mt-2">
                      <Badge variant="outline">
                        {stats.users.total} total users
                      </Badge>
                    </div>
                  )}
                </div>
                <div className="text-3xl">👥</div>
              </div>
            </CardContent>
          </Card>

          <Card
            className="hover:shadow-md transition-shadow cursor-pointer border-blue-200 bg-blue-50/30"
            onClick={() => navigate("/admin/ai-generation")}
          >
            <CardContent className="pt-6">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold mb-2 flex items-center gap-2">
                    🤖 AI Generation
                  </h3>
                  <p className="text-muted-foreground text-sm">
                    Auto-generate questions and quizzes
                  </p>
                  <div className="mt-2">
                    <Badge className="bg-blue-100 text-blue-800">
                      New Feature
                    </Badge>
                  </div>
                </div>
                <div className="text-3xl">✨</div>
              </div>
            </CardContent>
          </Card>
        </div>

        {loading ? (
          <div className="text-center py-8">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            <p className="mt-2 text-muted-foreground">Loading dashboard...</p>
          </div>
        ) : (
          stats && (
            <>
              {/* User Statistics */}
              <Card className="mb-6">
                <CardHeader>
                  <CardTitle>User Statistics</CardTitle>
                  <CardDescription>Overview of platform users</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
                    <div className="text-center">
                      <div className="text-2xl font-bold text-blue-600">
                        {stats.users.total}
                      </div>
                      <div className="text-sm text-muted-foreground">
                        Total Users
                      </div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-green-600">
                        {stats.users.active}
                      </div>
                      <div className="text-sm text-muted-foreground">
                        Active
                      </div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-purple-600">
                        {stats.users.creators}
                      </div>
                      <div className="text-sm text-muted-foreground">
                        Creators
                      </div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-orange-600">
                        {stats.users.reviewers}
                      </div>
                      <div className="text-sm text-muted-foreground">
                        Reviewers
                      </div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-red-600">
                        {stats.users.admins}
                      </div>
                      <div className="text-sm text-muted-foreground">
                        Admins
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Content Statistics */}
              <Card>
                <CardHeader>
                  <CardTitle>Content Statistics</CardTitle>
                  <CardDescription>
                    Overview of platform content
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div className="text-center">
                      <h4 className="text-lg font-semibold mb-3">Stories</h4>
                      <div className="space-y-2">
                        <div>
                          <div className="text-xl font-bold text-blue-600">
                            {stats.content.stories.total}
                          </div>
                          <div className="text-sm text-muted-foreground">
                            Total
                          </div>
                        </div>
                        <div>
                          <div className="text-lg font-medium text-yellow-600">
                            {stats.content.stories.pending}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Pending Review
                          </div>
                        </div>
                        <div>
                          <div className="text-lg font-medium text-green-600">
                            {stats.content.stories.active}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Active
                          </div>
                        </div>
                      </div>
                    </div>

                    <div className="text-center">
                      <h4 className="text-lg font-semibold mb-3">Questions</h4>
                      <div className="space-y-2">
                        <div>
                          <div className="text-xl font-bold text-blue-600">
                            {stats.content.questions.total}
                          </div>
                          <div className="text-sm text-muted-foreground">
                            Total
                          </div>
                        </div>
                        <div>
                          <div className="text-lg font-medium text-yellow-600">
                            {stats.content.questions.pending}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Pending Review
                          </div>
                        </div>
                        <div>
                          <div className="text-lg font-medium text-green-600">
                            {stats.content.questions.active}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Active
                          </div>
                        </div>
                      </div>
                    </div>

                    <div className="text-center">
                      <h4 className="text-lg font-semibold mb-3">Quizzes</h4>
                      <div className="space-y-2">
                        <div>
                          <div className="text-xl font-bold text-blue-600">
                            {stats.content.quizzes.total}
                          </div>
                          <div className="text-sm text-muted-foreground">
                            Total
                          </div>
                        </div>
                        <div>
                          <div className="text-lg font-medium text-yellow-600">
                            {stats.content.quizzes.pending}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Pending Review
                          </div>
                        </div>
                        <div>
                          <div className="text-lg font-medium text-green-600">
                            {stats.content.quizzes.active}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Active
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </>
          )
        )}

        {/* Quick Links */}
        <div className="mt-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4">
          <Link to="/stories" className="block">
            <Card className="hover:shadow-md transition-shadow h-full">
              <CardContent className="pt-6 text-center">
                <div className="text-2xl mb-2">📚</div>
                <h3 className="font-semibold">View Stories</h3>
                <p className="text-sm text-muted-foreground">
                  Browse all stories
                </p>
              </CardContent>
            </Card>
          </Link>

          <Link to="/stories/create" className="block">
            <Card className="hover:shadow-md transition-shadow h-full">
              <CardContent className="pt-6 text-center">
                <div className="text-2xl mb-2">✍️</div>
                <h3 className="font-semibold">Create Story</h3>
                <p className="text-sm text-muted-foreground">Add new content</p>
              </CardContent>
            </Card>
          </Link>

          <Link to="/questions/create" className="block">
            <Card className="hover:shadow-md transition-shadow h-full">
              <CardContent className="pt-6 text-center">
                <div className="text-2xl mb-2">❓</div>
                <h3 className="font-semibold">Create Questions</h3>
                <p className="text-sm text-muted-foreground">
                  Add quiz questions
                </p>
              </CardContent>
            </Card>
          </Link>

          <Link to="/quiz/create" className="block">
            <Card className="hover:shadow-md transition-shadow h-full">
              <CardContent className="pt-6 text-center">
                <div className="text-2xl mb-2">🧩</div>
                <h3 className="font-semibold">Create Quiz</h3>
                <p className="text-sm text-muted-foreground">Build new quiz</p>
              </CardContent>
            </Card>
          </Link>

          <Link to="/admin/reviews" className="block">
            <Card className="hover:shadow-md transition-shadow h-full">
              <CardContent className="pt-6 text-center">
                <div className="text-2xl mb-2">🔍</div>
                <h3 className="font-semibold">Reviews</h3>
                <p className="text-sm text-muted-foreground">
                  Content approval
                </p>
              </CardContent>
            </Card>
          </Link>

          <Link to="/admin/users" className="block">
            <Card className="hover:shadow-md transition-shadow h-full">
              <CardContent className="pt-6 text-center">
                <div className="text-2xl mb-2">⚙️</div>
                <h3 className="font-semibold">Settings</h3>
                <p className="text-sm text-muted-foreground">User management</p>
              </CardContent>
            </Card>
          </Link>
        </div>
      </div>
    </Layout>
  );
};

const AdminDashboardPage: React.FC = () => {
  return (
    <ProtectedRoute allowedRoles={["admin", "reviewer"]}>
      <AdminDashboard />
    </ProtectedRoute>
  );
};

export default AdminDashboardPage;
