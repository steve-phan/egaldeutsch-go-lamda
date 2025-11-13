import React, { useState, useEffect } from "react";
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
import { Badge } from "../../components/ui/badge";
import { Alert, AlertDescription } from "../../components/ui/alert";
import { Separator } from "../../components/ui/separator";
import axios from "axios";

const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? "/.netlify/functions"
    : "http://localhost:8888/.netlify/functions";

interface User {
  id: string;
  username: string;
  email: string;
  firstName: string;
  lastName: string;
  role: "creator" | "reviewer" | "admin";
  status: "active" | "suspended" | "pending";
  createdAt: string;
  updatedAt: string;
  lastLoginAt?: string;
}

const USER_ROLES = [
  { value: "creator", label: "Creator", color: "bg-blue-100 text-blue-800" },
  {
    value: "reviewer",
    label: "Reviewer",
    color: "bg-green-100 text-green-800",
  },
  { value: "admin", label: "Admin", color: "bg-purple-100 text-purple-800" },
];

const USER_STATUSES = [
  { value: "active", label: "Active", color: "bg-green-100 text-green-800" },
  { value: "suspended", label: "Suspended", color: "bg-red-100 text-red-800" },
  {
    value: "pending",
    label: "Pending",
    color: "bg-yellow-100 text-yellow-800",
  },
];

const UsersPage: React.FC = () => {
  const { user: currentUser } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [filteredUsers, setFilteredUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selectedRole, setSelectedRole] = useState("all");
  const [selectedStatus, setSelectedStatus] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    loadUsers();
  }, []);

  useEffect(() => {
    filterUsers();
  }, [users, selectedRole, selectedStatus, searchQuery]);

  const loadUsers = async () => {
    try {
      setLoading(true);
      const response = await axios.get(`${API_BASE_URL}/user-management`);
      setUsers(response.data);
    } catch (err: any) {
      console.error("Error loading users:", err);
      setError(err.response?.data?.error || "Failed to load users");
    } finally {
      setLoading(false);
    }
  };

  const filterUsers = () => {
    let filtered = users;

    if (selectedRole && selectedRole !== "all") {
      filtered = filtered.filter((user) => user.role === selectedRole);
    }

    if (selectedStatus && selectedStatus !== "all") {
      filtered = filtered.filter((user) => user.status === selectedStatus);
    }

    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (user) =>
          user.username.toLowerCase().includes(query) ||
          user.email.toLowerCase().includes(query) ||
          user.firstName.toLowerCase().includes(query) ||
          user.lastName.toLowerCase().includes(query)
      );
    }

    setFilteredUsers(filtered);
  };

  const handleUserAction = async (
    userId: string,
    action: "suspend" | "activate" | "delete" | "promote" | "demote"
  ) => {
    if (!currentUser || currentUser.role !== "admin") {
      setError("Insufficient permissions");
      return;
    }

    if (
      userId === currentUser.id &&
      (action === "suspend" || action === "delete")
    ) {
      setError("Cannot perform this action on your own account");
      return;
    }

    try {
      setIsSubmitting(true);
      setError("");

      switch (action) {
        case "suspend":
          await axios.put(`${API_BASE_URL}/user-management/${userId}`, {
            status: "suspended",
          });
          break;
        case "activate":
          await axios.put(`${API_BASE_URL}/user-management/${userId}`, {
            status: "active",
          });
          break;
        case "delete":
          if (
            window.confirm(
              "Are you sure you want to delete this user? This action cannot be undone."
            )
          ) {
            await axios.delete(`${API_BASE_URL}/user-management/${userId}`);
          } else {
            return;
          }
          break;
        case "promote":
          const targetUser = users.find((u) => u.id === userId);
          if (targetUser) {
            const newRole =
              targetUser.role === "creator"
                ? "reviewer"
                : targetUser.role === "reviewer"
                ? "admin"
                : "admin";
            await axios.put(`${API_BASE_URL}/user-management/${userId}`, {
              role: newRole,
            });
          }
          break;
        case "demote":
          const targetUser2 = users.find((u) => u.id === userId);
          if (targetUser2) {
            const newRole =
              targetUser2.role === "admin"
                ? "reviewer"
                : targetUser2.role === "reviewer"
                ? "creator"
                : "creator";
            await axios.put(`${API_BASE_URL}/user-management/${userId}`, {
              role: newRole,
            });
          }
          break;
      }

      // Reload users
      await loadUsers();
      setSelectedUser(null);
    } catch (err: any) {
      console.error("Error performing user action:", err);
      setError(err.response?.data?.error || "Failed to perform action");
    } finally {
      setIsSubmitting(false);
    }
  };

  const getRoleBadge = (role: string) => {
    const roleInfo = USER_ROLES.find((r) => r.value === role);
    return (
      <Badge className={roleInfo?.color || "bg-gray-100 text-gray-800"}>
        {roleInfo?.label || role}
      </Badge>
    );
  };

  const getStatusBadge = (status: string) => {
    const statusInfo = USER_STATUSES.find((s) => s.value === status);
    return (
      <Badge className={statusInfo?.color || "bg-gray-100 text-gray-800"}>
        {statusInfo?.label || status}
      </Badge>
    );
  };

  const getAvailableActions = (user: User) => {
    if (user.id === currentUser?.id) return [];

    const actions = [];

    if (user.status === "active") {
      actions.push({
        action: "suspend",
        label: "Suspend",
        variant: "destructive" as const,
      });
    } else if (user.status === "suspended") {
      actions.push({
        action: "activate",
        label: "Activate",
        variant: "default" as const,
      });
    }

    if (user.role === "creator") {
      actions.push({
        action: "promote",
        label: "Promote to Reviewer",
        variant: "outline" as const,
      });
    } else if (user.role === "reviewer") {
      actions.push({
        action: "promote",
        label: "Promote to Admin",
        variant: "outline" as const,
      });
      actions.push({
        action: "demote",
        label: "Demote to Creator",
        variant: "outline" as const,
      });
    } else if (user.role === "admin") {
      actions.push({
        action: "demote",
        label: "Demote to Reviewer",
        variant: "outline" as const,
      });
    }

    actions.push({
      action: "delete",
      label: "Delete User",
      variant: "destructive" as const,
    });

    return actions;
  };

  return (
    <Layout>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">
            User Management
          </h1>
          <p className="text-muted-foreground">
            Manage user accounts, roles, and permissions
          </p>
        </div>

        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <Card>
            <CardContent className="pt-6">
              <div className="text-2xl font-bold text-blue-600">
                {users.filter((u) => u.status === "active").length}
              </div>
              <div className="text-sm text-muted-foreground">Active Users</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="pt-6">
              <div className="text-2xl font-bold text-green-600">
                {users.filter((u) => u.role === "creator").length}
              </div>
              <div className="text-sm text-muted-foreground">Creators</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="pt-6">
              <div className="text-2xl font-bold text-purple-600">
                {users.filter((u) => u.role === "reviewer").length}
              </div>
              <div className="text-sm text-muted-foreground">Reviewers</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="pt-6">
              <div className="text-2xl font-bold text-orange-600">
                {users.filter((u) => u.role === "admin").length}
              </div>
              <div className="text-sm text-muted-foreground">Admins</div>
            </CardContent>
          </Card>
        </div>

        {/* Filters */}
        <Card className="mb-6">
          <CardContent className="pt-6">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div>
                <Label htmlFor="search">Search Users</Label>
                <Input
                  id="search"
                  type="text"
                  placeholder="Search by name, username, or email..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
              </div>
              <div>
                <Label htmlFor="role-filter">Filter by Role</Label>
                <select
                  id="role-filter"
                  value={selectedRole}
                  onChange={(e) => setSelectedRole(e.target.value)}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                >
                  <option value="all">All Roles</option>
                  {USER_ROLES.map((role) => (
                    <option key={role.value} value={role.value}>
                      {role.label}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <Label htmlFor="status-filter">Filter by Status</Label>
                <select
                  id="status-filter"
                  value={selectedStatus}
                  onChange={(e) => setSelectedStatus(e.target.value)}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                >
                  <option value="all">All Statuses</option>
                  {USER_STATUSES.map((status) => (
                    <option key={status.value} value={status.value}>
                      {status.label}
                    </option>
                  ))}
                </select>
              </div>
              <div className="flex items-end">
                <Button onClick={loadUsers} disabled={loading}>
                  {loading ? "Loading..." : "Refresh"}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Users List */}
        {loading ? (
          <div className="text-center py-8">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            <p className="mt-2 text-muted-foreground">Loading users...</p>
          </div>
        ) : (
          <div className="grid gap-4">
            {filteredUsers.length === 0 ? (
              <Card>
                <CardContent className="pt-6 text-center">
                  <p className="text-muted-foreground">
                    No users found for the selected filters.
                  </p>
                </CardContent>
              </Card>
            ) : (
              filteredUsers.map((user) => (
                <Card
                  key={user.id}
                  className="hover:shadow-md transition-shadow"
                >
                  <CardContent className="pt-6">
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          <h3 className="text-lg font-semibold">
                            {user.firstName} {user.lastName}
                          </h3>
                          {user.id === currentUser?.id && (
                            <Badge variant="outline" className="text-xs">
                              You
                            </Badge>
                          )}
                          {getRoleBadge(user.role)}
                          {getStatusBadge(user.status)}
                        </div>
                        <div className="text-sm text-muted-foreground space-y-1">
                          <div>
                            @{user.username} • {user.email}
                          </div>
                          <div>
                            Joined:{" "}
                            {new Date(user.createdAt).toLocaleDateString()}
                            {user.lastLoginAt && (
                              <>
                                {" "}
                                • Last login:{" "}
                                {new Date(
                                  user.lastLoginAt
                                ).toLocaleDateString()}
                              </>
                            )}
                          </div>
                        </div>
                      </div>
                      <div className="flex gap-2">
                        {currentUser?.role === "admin" &&
                          getAvailableActions(user).length > 0 && (
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => setSelectedUser(user)}
                              disabled={isSubmitting}
                            >
                              Manage
                            </Button>
                          )}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))
            )}
          </div>
        )}

        {/* User Management Modal */}
        {selectedUser && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <Card className="w-full max-w-md">
              <CardHeader>
                <CardTitle>Manage User</CardTitle>
                <CardDescription>
                  {selectedUser.firstName} {selectedUser.lastName} (@
                  {selectedUser.username})
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">Current Role:</span>
                  {getRoleBadge(selectedUser.role)}
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">Current Status:</span>
                  {getStatusBadge(selectedUser.status)}
                </div>

                <Separator />

                <div className="space-y-2">
                  <Label>Available Actions:</Label>
                  <div className="space-y-2">
                    {getAvailableActions(selectedUser).map((actionInfo) => (
                      <Button
                        key={actionInfo.action}
                        variant={actionInfo.variant}
                        size="sm"
                        className="w-full justify-start"
                        onClick={() =>
                          handleUserAction(
                            selectedUser.id,
                            actionInfo.action as any
                          )
                        }
                        disabled={isSubmitting}
                      >
                        {isSubmitting ? "Processing..." : actionInfo.label}
                      </Button>
                    ))}
                  </div>
                </div>

                <div className="flex gap-2 pt-4">
                  <Button
                    variant="outline"
                    onClick={() => setSelectedUser(null)}
                    disabled={isSubmitting}
                    className="flex-1"
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

const AdminUsersPage: React.FC = () => {
  return (
    <ProtectedRoute allowedRoles={["admin"]}>
      <UsersPage />
    </ProtectedRoute>
  );
};

export default AdminUsersPage;
