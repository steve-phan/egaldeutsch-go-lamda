import React, { useState, useEffect } from "react";
import Layout from "../../components/layout";
import ProtectedRoute from "../../components/auth/ProtectedRoute";
import { useAuth } from "../../contexts/AuthContext";
import { ErrorAlert } from "../../components/atoms/ErrorAlert";
import { UserFilterBar } from "../../components/molecules/UserFilterBar";
import { UserManagementList } from "../../components/organisms/UserManagementList";
import { ConfirmationDialog } from "../../components/molecules/ConfirmationDialog";
import { useUserManagement, UserAction } from "../../hooks/useUserManagement";
import { useUserFilters } from "../../hooks/useUserFilters";

const UsersPage: React.FC = () => {
  const { user: currentUser } = useAuth();
  const {
    users,
    loading,
    error: apiError,
    loadUsers,
    performUserAction,
  } = useUserManagement(currentUser?.id);
  const {
    filteredUsers,
    roleFilter,
    statusFilter,
    searchQuery,
    setRoleFilter,
    setStatusFilter,
    setSearchQuery,
  } = useUserFilters(users);
  const [error, setError] = useState("");
  const [confirmationDialog, setConfirmationDialog] = useState<{
    open: boolean;
    userId: string | null;
    action: UserAction | null;
    title: string;
    description: string;
  }>({
    open: false,
    userId: null,
    action: null,
    title: "",
    description: "",
  });

  useEffect(() => {
    loadUsers();
  }, [loadUsers]);

  const handleUserAction = async (userId: string, action: UserAction) => {
    if (currentUser?.role !== "admin") {
      setError("Insufficient permissions");
      return;
    }

    // For delete action, show confirmation dialog
    if (action === "delete") {
      const user = users.find((u) => u.id === userId);
      setConfirmationDialog({
        open: true,
        userId,
        action,
        title: "Delete User",
        description: `Are you sure you want to delete ${user?.firstName} ${user?.lastName}? This action cannot be undone.`,
      });
      return;
    }

    // For other actions, execute immediately
    try {
      setError("");
      await performUserAction(userId, action);
      await loadUsers();
    } catch (err: any) {
      setError(err.message || "Failed to perform action");
    }
  };

  const handleConfirmAction = async () => {
    if (!confirmationDialog.userId || !confirmationDialog.action) return;

    try {
      setError("");
      await performUserAction(confirmationDialog.userId, confirmationDialog.action);
      await loadUsers();
    } catch (err: any) {
      setError(err.message || "Failed to perform action");
    }
  };

  const isAdmin = currentUser?.role === "admin";
  const displayError = error || apiError;

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

        <ErrorAlert message={displayError} />

        <UserFilterBar
          roleFilter={roleFilter}
          statusFilter={statusFilter}
          searchQuery={searchQuery}
          onRoleChange={setRoleFilter}
          onStatusChange={setStatusFilter}
          onSearchChange={setSearchQuery}
          onRefresh={loadUsers}
          loading={loading}
        />

        <UserManagementList
          users={filteredUsers}
          loading={loading}
          currentUserId={currentUser?.id}
          onUserAction={handleUserAction}
          isAdmin={isAdmin}
        />

        <ConfirmationDialog
          open={confirmationDialog.open}
          onOpenChange={(open) =>
            setConfirmationDialog((prev) => ({ ...prev, open }))
          }
          onConfirm={handleConfirmAction}
          title={confirmationDialog.title}
          description={confirmationDialog.description}
          confirmText="Delete"
          cancelText="Cancel"
          variant="destructive"
        />
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
