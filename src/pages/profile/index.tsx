import React from "react";
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
import { Badge } from "../../components/ui/badge";
import { Separator } from "../../components/ui/separator";
import { LoadingSpinner } from "@/components/atoms/LoadingSpinner";

const ProfilePage: React.FC = () => {
  const { user, isLoading } = useAuth();

  console.log({ user, isLoading });

  if (isLoading) return <LoadingSpinner message="Loading..." />;

  if (!user) return null;

  const formatDate = (dateString?: string) => {
    if (!dateString) return "Never";
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const getRoleBadgeVariant = (role: string) => {
    switch (role) {
      case "admin":
        return "destructive";
      case "reviewer":
        return "secondary";
      case "creator":
      default:
        return "default";
    }
  };

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case "active":
        return "success";
      case "suspended":
        return "destructive";
      case "pending":
      default:
        return "warning";
    }
  };

  return (
    <Layout>
      <ProtectedRoute>
        <div className="max-w-4xl mx-auto px-4 py-12">
          <div className="mb-6">
            <h1 className="text-3xl font-bold mb-2">Profile</h1>
            <p className="text-muted-foreground">
              Manage your account settings and preferences
            </p>
          </div>

          {/* User Information Card */}
          <Card className="mb-6">
            <CardHeader>
              <div className="flex items-start justify-between">
                <div>
                  <CardTitle>Personal Information</CardTitle>
                  <CardDescription>
                    Your account details and role information
                  </CardDescription>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => navigate("/profile/settings")}
                >
                  Edit Profile
                </Button>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-muted-foreground">
                    Full Name
                  </label>
                  <p className="text-lg font-medium">
                    {user.firstName} {user.lastName}
                  </p>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">
                    Username
                  </label>
                  <p className="text-lg font-medium">@{user.username}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">
                    Email
                  </label>
                  <p className="text-lg font-medium">{user.email}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">
                    User ID
                  </label>
                  <p className="text-sm font-mono text-muted-foreground">
                    {user.id}
                  </p>
                </div>
              </div>

              <Separator />

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-muted-foreground">
                    Role
                  </label>
                  <div className="mt-1">
                    <Badge variant={getRoleBadgeVariant(user.role)}>
                      {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                    </Badge>
                  </div>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">
                    Account Status
                  </label>
                  <div className="mt-1">
                    <Badge variant={getStatusBadgeVariant(user.status)}>
                      {user.status.charAt(0).toUpperCase() +
                        user.status.slice(1)}
                    </Badge>
                  </div>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">
                    Member Since
                  </label>
                  <p className="text-sm">{formatDate(user.createdAt)}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">
                    Last Login
                  </label>
                  <p className="text-sm">{formatDate(user.lastLoginAt)}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Role-Specific Information */}
          <Card>
            <CardHeader>
              <CardTitle>Role Capabilities</CardTitle>
              <CardDescription>
                What you can do with your {user.role} role
              </CardDescription>
            </CardHeader>
            <CardContent>
              {user.role === "admin" && (
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">
                    As an <strong>Admin</strong>, you have full access to:
                  </p>
                  <ul className="list-disc list-inside space-y-1 text-sm ml-4">
                    <li>Manage all users and their roles</li>
                    <li>Review and approve all content</li>
                    <li>Create stories and quizzes</li>
                    <li>Access system analytics and reports</li>
                    <li>Configure platform settings</li>
                  </ul>
                </div>
              )}
              {user.role === "reviewer" && (
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">
                    As a <strong>Reviewer</strong>, you can:
                  </p>
                  <ul className="list-disc list-inside space-y-1 text-sm ml-4">
                    <li>Review content submitted by creators</li>
                    <li>Approve or request changes to stories and quizzes</li>
                    <li>Provide feedback and quality assurance</li>
                    <li>View content analytics</li>
                  </ul>
                </div>
              )}
              {user.role === "creator" && (
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">
                    As a <strong>Creator</strong>, you can:
                  </p>
                  <ul className="list-disc list-inside space-y-1 text-sm ml-4">
                    <li>Create new German learning stories</li>
                    <li>Design quizzes and questions</li>
                    <li>Submit content for review</li>
                    <li>View your content statistics</li>
                  </ul>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Quick Actions */}
          <div className="mt-6 flex flex-wrap gap-4">
            <Button onClick={() => navigate("/profile/settings")}>
              Edit Profile
            </Button>
            {(user.role === "creator" || user.role === "admin") && (
              <Button
                variant="outline"
                onClick={() => navigate("/stories/create")}
              >
                Create Story
              </Button>
            )}
            {(user.role === "reviewer" || user.role === "admin") && (
              <Button
                variant="outline"
                onClick={() => navigate("/admin/reviews")}
              >
                Review Content
              </Button>
            )}
            {user.role === "admin" && (
              <Button
                variant="outline"
                onClick={() => navigate("/admin/users")}
              >
                Manage Users
              </Button>
            )}
          </div>
        </div>
      </ProtectedRoute>
    </Layout>
  );
};

export default ProfilePage;
