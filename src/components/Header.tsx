import React, { useState } from "react";
import { Link } from "gatsby";
import { Menu, User, BookOpen, HelpCircle, LayoutDashboard, FileQuestion, Plus } from "lucide-react";
import { Button } from "./ui/button";
import { Badge } from "./ui/badge";
import { DropdownMenu, DropdownMenuItem } from "./ui/dropdown-menu";
import { MobileMenu, MobileMenuItem, MobileMenuGroup } from "./ui/mobile-menu";
import { useAuth } from "../contexts/AuthContext";
import NotificationBell from "./NotificationBell";

const Header: React.FC = () => {
  const { isAuthenticated, user, logout } = useAuth();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const isCreatorOrAdmin = user?.role === "creator" || user?.role === "admin";
  const isReviewerOrAdmin = user?.role === "reviewer" || user?.role === "admin";

  return (
    <header className="bg-card shadow-sm border-b border-border sticky top-0 z-30">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center py-4">
          {/* Logo */}
          <Link
            to="/"
            className="text-2xl font-bold text-primary hover:text-primary/90 transition-colors flex-shrink-0"
          >
            🇩🇪 EgalDeutsch
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden lg:flex items-center space-x-2">
            {/* Public Links */}
            <Link
              to="/"
              className="text-foreground hover:text-primary px-3 py-2 rounded-md text-sm font-medium transition-colors"
              activeClassName="text-primary bg-primary/10"
            >
              Stories
            </Link>
            <Link
              to="/leaderboard"
              className="text-foreground hover:text-primary px-3 py-2 rounded-md text-sm font-medium transition-colors"
              activeClassName="text-primary bg-primary/10"
            >
              Leaderboard
            </Link>
            <Link
              to="/about"
              className="text-foreground hover:text-primary px-3 py-2 rounded-md text-sm font-medium transition-colors"
              activeClassName="text-primary bg-primary/10"
            >
              About
            </Link>

            {/* Content Management - For Creators/Admins */}
            {isCreatorOrAdmin && (
              <DropdownMenu trigger="Content" align="left">
                <DropdownMenuItem>
                  <Link to="/stories/create" className="flex items-center gap-2">
                    <Plus size={16} />
                    Create Story
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Link to="/stories" className="flex items-center gap-2">
                    <BookOpen size={16} />
                    Manage Stories
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Link to="/questions/create" className="flex items-center gap-2">
                    <Plus size={16} />
                    Create Question
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Link to="/questions" className="flex items-center gap-2">
                    <HelpCircle size={16} />
                    Manage Questions
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Link to="/quiz/create" className="flex items-center gap-2">
                    <Plus size={16} />
                    Create Quiz
                  </Link>
                </DropdownMenuItem>
              </DropdownMenu>
            )}

            {/* Admin Panel - For Admins/Reviewers */}
            {isReviewerOrAdmin && (
              <Link
                to="/admin"
                className="text-foreground hover:text-primary px-3 py-2 rounded-md text-sm font-medium transition-colors"
                activeClassName="text-primary bg-primary/10"
              >
                Admin
              </Link>
            )}

            {/* User Menu */}
            {isAuthenticated && user ? (
              <>
                <NotificationBell className="mx-1" />
                
                <DropdownMenu
                  trigger={
                    <div className="flex items-center gap-2">
                      <User size={16} />
                      <span>{user.firstName}</span>
                      <Badge variant="outline" className="text-xs">
                        {user.role}
                      </Badge>
                    </div>
                  }
                  align="right"
                >
                  <DropdownMenuItem>
                    <Link to="/profile" className="flex items-center gap-2">
                      <User size={16} />
                      Profile
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <Link to="/notifications" className="flex items-center gap-2">
                      <span>Notifications</span>
                    </Link>
                  </DropdownMenuItem>
                  <div className="border-t border-border my-1" />
                  <DropdownMenuItem onClick={logout}>
                    <span className="text-destructive">Logout</span>
                  </DropdownMenuItem>
                </DropdownMenu>
              </>
            ) : (
              <div className="flex items-center space-x-2">
                <Link to="/auth/login">
                  <Button variant="ghost" size="sm">
                    Login
                  </Button>
                </Link>
                <Link to="/auth/register">
                  <Button size="sm">Sign Up</Button>
                </Link>
              </div>
            )}
          </nav>

          {/* Mobile Menu Button */}
          <button
            className="lg:hidden p-2 rounded-md hover:bg-accent transition-colors"
            onClick={() => setIsMobileMenuOpen(true)}
            aria-label="Open menu"
          >
            <Menu size={24} />
          </button>
        </div>
      </div>

      {/* Mobile Navigation */}
      <MobileMenu
        isOpen={isMobileMenuOpen}
        onClose={() => setIsMobileMenuOpen(false)}
      >
        {/* User Info (if authenticated) */}
        {isAuthenticated && user && (
          <MobileMenuGroup title="Account">
            <MobileMenuItem>
              <Link
                to="/profile"
                className="flex items-center gap-3"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                <User size={18} />
                <div>
                  <div className="font-medium">{user.firstName}</div>
                  <Badge variant="outline" className="text-xs mt-1">
                    {user.role}
                  </Badge>
                </div>
              </Link>
            </MobileMenuItem>
            <MobileMenuItem>
              <Link
                to="/notifications"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Notifications
              </Link>
            </MobileMenuItem>
          </MobileMenuGroup>
        )}

        {/* Main Navigation */}
        <MobileMenuGroup title="Navigation">
          <MobileMenuItem>
            <Link to="/" onClick={() => setIsMobileMenuOpen(false)}>
              Stories
            </Link>
          </MobileMenuItem>
          <MobileMenuItem>
            <Link to="/leaderboard" onClick={() => setIsMobileMenuOpen(false)}>
              Leaderboard
            </Link>
          </MobileMenuItem>
          <MobileMenuItem>
            <Link to="/about" onClick={() => setIsMobileMenuOpen(false)}>
              About
            </Link>
          </MobileMenuItem>
          {isReviewerOrAdmin && (
            <MobileMenuItem>
              <Link to="/admin" onClick={() => setIsMobileMenuOpen(false)}>
                Admin Dashboard
              </Link>
            </MobileMenuItem>
          )}
        </MobileMenuGroup>

        {/* Content Management (for creators/admins) */}
        {isCreatorOrAdmin && (
          <MobileMenuGroup title="Content Management">
            <MobileMenuItem>
              <Link
                to="/stories/create"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Create Story
              </Link>
            </MobileMenuItem>
            <MobileMenuItem>
              <Link to="/stories" onClick={() => setIsMobileMenuOpen(false)}>
                Manage Stories
              </Link>
            </MobileMenuItem>
            <MobileMenuItem>
              <Link
                to="/questions/create"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Create Question
              </Link>
            </MobileMenuItem>
            <MobileMenuItem>
              <Link to="/questions" onClick={() => setIsMobileMenuOpen(false)}>
                Manage Questions
              </Link>
            </MobileMenuItem>
            <MobileMenuItem>
              <Link
                to="/quiz/create"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Create Quiz
              </Link>
            </MobileMenuItem>
          </MobileMenuGroup>
        )}

        {/* Auth Actions */}
        {isAuthenticated && user ? (
          <div className="mt-6 pt-6 border-t border-border">
            <Button
              variant="ghost"
              className="w-full justify-start text-destructive hover:text-destructive hover:bg-destructive/10"
              onClick={() => {
                logout();
                setIsMobileMenuOpen(false);
              }}
            >
              Logout
            </Button>
          </div>
        ) : (
          <div className="mt-6 pt-6 border-t border-border space-y-2">
            <Link to="/auth/login" onClick={() => setIsMobileMenuOpen(false)}>
              <Button variant="ghost" className="w-full">
                Login
              </Button>
            </Link>
            <Link
              to="/auth/register"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              <Button className="w-full">Sign Up</Button>
            </Link>
          </div>
        )}
      </MobileMenu>
    </header>
  );
};

export default Header;
