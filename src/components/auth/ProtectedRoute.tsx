import React, { useEffect } from 'react';
import { navigate } from 'gatsby';
import { useAuth } from '../../contexts/AuthContext';
import { UserRole } from '../../types';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requireAuth?: boolean;
  allowedRoles?: UserRole[];
  redirectTo?: string;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  children,
  requireAuth = true,
  allowedRoles,
  redirectTo = '/auth/login',
}) => {
  const { isAuthenticated, isLoading, user, hasAnyRole } = useAuth();

  useEffect(() => {
    if (!isLoading) {
      // Redirect if authentication is required but user is not authenticated
      if (requireAuth && !isAuthenticated) {
        navigate(redirectTo);
        return;
      }

      // Redirect if user doesn't have required role
      if (isAuthenticated && allowedRoles && !hasAnyRole(allowedRoles)) {
        navigate('/'); // Redirect to home if doesn't have permission
        return;
      }
    }
  }, [isAuthenticated, isLoading, requireAuth, allowedRoles, redirectTo, user, hasAnyRole]);

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
          <p className="mt-4 text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  // Don't render protected content if not authenticated and auth is required
  if (requireAuth && !isAuthenticated) {
    return null;
  }

  // Don't render if user doesn't have required role
  if (isAuthenticated && allowedRoles && !hasAnyRole(allowedRoles)) {
    return null;
  }

  return <>{children}</>;
};

export default ProtectedRoute;
