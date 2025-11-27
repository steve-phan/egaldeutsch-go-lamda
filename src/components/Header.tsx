import React, { useState, useEffect } from "react";
import { Link } from "gatsby";
import {
  Menu,
  User,
  BookOpen,
  HelpCircle,
  LayoutDashboard,
  FileQuestion,
  Plus,
  X,
  LogOut,
  Bell,
  ChevronDown,
  Sparkles,
} from "lucide-react";
import { Button } from "./ui/button";
import { Badge } from "./ui/badge";
import { DropdownMenu, DropdownMenuItem } from "./ui/dropdown-menu";
import { useAuth } from "../contexts/AuthContext";
import NotificationBell from "./NotificationBell";
import { DekstopNavigation } from "./atoms/DesktopNavigation";

const Header: React.FC = () => {
  const { isAuthenticated, user, logout } = useAuth();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isScrolled, setIsScrolled] = useState(false);

  const isCreatorOrAdmin = user?.role === "creator" || user?.role === "admin";
  const isReviewerOrAdmin = user?.role === "reviewer" || user?.role === "admin";

  // Add scroll listener for header backdrop effect
  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 10);
    };
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  // Get role badge color
  const getRoleBadgeColor = (role: string) => {
    switch (role) {
      case "admin":
        return "bg-red-500/10 text-red-600 border-red-500/20";
      case "reviewer":
        return "bg-purple-500/10 text-purple-600 border-purple-500/20";
      case "creator":
        return "bg-blue-500/10 text-blue-600 border-blue-500/20";
      default:
        return "bg-gray-500/10 text-gray-600 border-gray-500/20";
    }
  };

  return (
    <header
      className={`sticky top-0 z-50 transition-all duration-300 ${
        isScrolled
          ? "bg-background/80 backdrop-blur-lg shadow-lg border-b border-border/50"
          : "bg-background border-b border-border"
      }`}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link
            to="/"
            className="flex items-center gap-2 text-2xl font-bold bg-gradient-to-r from-primary via-purple-600 to-pink-600 bg-clip-text text-transparent hover:opacity-80 transition-opacity flex-shrink-0 group"
          >
            <span className="group-hover:scale-110 transition-transform duration-200">
              🇩🇪
            </span>
            <span>EgalDeutsch</span>
          </Link>
          <DekstopNavigation
            isAuthenticated={isAuthenticated}
            user={user}
            isCreatorOrAdmin={isCreatorOrAdmin}
            isReviewerOrAdmin={isReviewerOrAdmin}
            logout={logout}
          />

          {/* Mobile Menu Button */}
          <button
            className="lg:hidden p-2 rounded-lg hover:bg-accent/50 transition-colors"
            onClick={() => setIsMobileMenuOpen(true)}
            aria-label="Open menu"
          >
            <Menu size={24} />
          </button>
        </div>
      </div>

      {/* Mobile Navigation Overlay */}
      {isMobileMenuOpen && (
        <div
          className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 lg:hidden"
          onClick={() => setIsMobileMenuOpen(false)}
        />
      )}

      {/* Mobile Navigation Drawer */}
      <div
        className={`fixed top-0 right-0 h-full w-80 max-w-[85vw] bg-background shadow-2xl z-50 transform transition-transform duration-300 ease-in-out lg:hidden ${
          isMobileMenuOpen ? "translate-x-0" : "translate-x-full"
        }`}
      >
        <div className="flex flex-col h-full">
          {/* Mobile Menu Header */}
          <div className="flex items-center justify-between p-4 border-b border-border bg-gradient-to-r from-primary/5 to-purple-600/5">
            <h2 className="text-lg font-semibold">Menu</h2>
            <button
              onClick={() => setIsMobileMenuOpen(false)}
              className="p-2 rounded-lg hover:bg-accent/50 transition-colors"
              aria-label="Close menu"
            >
              <X size={20} />
            </button>
          </div>

          {/* Mobile Menu Content */}
          <div className="flex-1 overflow-y-auto p-4">
            {/* User Info (if authenticated) */}
            {isAuthenticated && user && (
              <div className="mb-6 p-4 rounded-xl bg-gradient-to-br from-primary/10 to-purple-600/10 border border-primary/20">
                <Link
                  to="/profile"
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="flex items-center gap-3"
                >
                  <div className="w-12 h-12 rounded-full bg-gradient-to-br from-primary to-purple-600 flex items-center justify-center text-white text-lg font-semibold shadow-lg">
                    {user.firstName?.charAt(0).toUpperCase()}
                  </div>
                  <div className="flex-1">
                    <div className="font-semibold text-base">
                      {user.firstName} {user.lastName}
                    </div>
                    <Badge
                      className={`text-xs mt-1 ${getRoleBadgeColor(user.role)}`}
                    >
                      {user.role}
                    </Badge>
                  </div>
                </Link>
              </div>
            )}

            {/* Main Navigation */}
            <div className="mb-6">
              <h3 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-3 px-2">
                Navigation
              </h3>
              <div className="space-y-1">
                <Link
                  to="/"
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-accent/50 transition-all text-foreground font-medium"
                >
                  <BookOpen size={20} />
                  Stories
                </Link>
                <Link
                  to="/leaderboard"
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-accent/50 transition-all text-foreground font-medium"
                >
                  <LayoutDashboard size={20} />
                  Leaderboard
                </Link>
                <Link
                  to="/about"
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-accent/50 transition-all text-foreground font-medium"
                >
                  <HelpCircle size={20} />
                  About
                </Link>
                {isAuthenticated && user && (
                  <Link
                    to="/notifications"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-accent/50 transition-all text-foreground font-medium"
                  >
                    <Bell size={20} />
                    Notifications
                  </Link>
                )}
                {isReviewerOrAdmin && (
                  <Link
                    to="/admin"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-accent/50 transition-all text-foreground font-medium"
                  >
                    <LayoutDashboard size={20} />
                    Admin Dashboard
                  </Link>
                )}
              </div>
            </div>

            {/* Content Management (for creators/admins) */}
            {isCreatorOrAdmin && (
              <div className="mb-6">
                <h3 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-3 px-2">
                  Content Management
                </h3>
                <div className="space-y-1">
                  <Link
                    to="/stories/create"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-accent/50 transition-all text-foreground font-medium"
                  >
                    <Plus size={20} />
                    Create Story
                  </Link>
                  <Link
                    to="/stories"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-accent/50 transition-all text-foreground font-medium"
                  >
                    <BookOpen size={20} />
                    Manage Stories
                  </Link>
                  <Link
                    to="/questions/create"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-accent/50 transition-all text-foreground font-medium"
                  >
                    <Plus size={20} />
                    Create Question
                  </Link>
                  <Link
                    to="/questions"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-accent/50 transition-all text-foreground font-medium"
                  >
                    <FileQuestion size={20} />
                    Manage Questions
                  </Link>
                  <Link
                    to="/quiz/create"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-accent/50 transition-all text-foreground font-medium"
                  >
                    <Plus size={20} />
                    Create Quiz
                  </Link>
                </div>
              </div>
            )}
          </div>

          {/* Mobile Menu Footer */}
          <div className="p-4 border-t border-border bg-muted/30">
            {isAuthenticated && user ? (
              <Button
                variant="ghost"
                className="w-full justify-start text-destructive hover:text-destructive hover:bg-destructive/10 font-medium"
                onClick={() => {
                  logout();
                  setIsMobileMenuOpen(false);
                }}
              >
                <LogOut size={18} className="mr-2" />
                Logout
              </Button>
            ) : (
              <div className="space-y-2">
                <Link
                  to="/auth/login"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  <Button variant="outline" className="w-full font-medium">
                    Login
                  </Button>
                </Link>
                <Link
                  to="/auth/register"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  <Button className="w-full bg-gradient-to-r from-primary to-purple-600 hover:from-primary/90 hover:to-purple-600/90 shadow-md font-medium">
                    Sign Up
                  </Button>
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
