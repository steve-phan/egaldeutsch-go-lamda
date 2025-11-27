import { Link } from "gatsby";
import { Badge, Button, DropdownMenu, DropdownMenuItem } from "../ui";
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
import NotificationBell from "../NotificationBell";

interface DesktopNavigationProps {
  isAuthenticated: boolean;
  user: {
    firstName: string;
    lastName: string;
    email: string;
    role: string;
  } | null;
  isCreatorOrAdmin: boolean;
  isReviewerOrAdmin: boolean;
  logout: () => void;
}

export const DekstopNavigation = ({
  isAuthenticated,
  user,
  isCreatorOrAdmin,
  isReviewerOrAdmin,
  logout,
}: DesktopNavigationProps) => {
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
    <nav className="hidden lg:flex items-center space-x-1">
      {/* Public Links */}
      <Link
        to="/"
        className="relative text-foreground hover:text-primary px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-accent/50 group"
        activeClassName="text-primary bg-primary/10"
      >
        <span className="relative z-10">Stories</span>
        <span className="absolute inset-0 bg-gradient-to-r from-primary/0 via-primary/5 to-primary/0 opacity-0 group-hover:opacity-100 transition-opacity rounded-lg" />
      </Link>
      <Link
        to="/leaderboard"
        className="relative text-foreground hover:text-primary px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-accent/50 group"
        activeClassName="text-primary bg-primary/10"
      >
        <span className="relative z-10">Leaderboard</span>
        <span className="absolute inset-0 bg-gradient-to-r from-primary/0 via-primary/5 to-primary/0 opacity-0 group-hover:opacity-100 transition-opacity rounded-lg" />
      </Link>
      <Link
        to="/about"
        className="relative text-foreground hover:text-primary px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-accent/50 group"
        activeClassName="text-primary bg-primary/10"
      >
        <span className="relative z-10">About</span>
        <span className="absolute inset-0 bg-gradient-to-r from-primary/0 via-primary/5 to-primary/0 opacity-0 group-hover:opacity-100 transition-opacity rounded-lg" />
      </Link>

      {/* Content Management - For Creators/Admins */}
      {isCreatorOrAdmin && (
        <DropdownMenu
          trigger={
            <span className="flex items-center gap-1">
              <Sparkles size={16} />
              Content
            </span>
          }
          align="left"
        >
          <DropdownMenuItem>
            <Link
              to="/stories/create"
              className="flex items-center gap-2 w-full"
            >
              <Plus size={16} />
              Create Story
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Link to="/stories" className="flex items-center gap-2 w-full">
              <BookOpen size={16} />
              Manage Stories
            </Link>
          </DropdownMenuItem>
          <div className="h-px bg-border my-1" />
          <DropdownMenuItem>
            <Link
              to="/questions/create"
              className="flex items-center gap-2 w-full"
            >
              <Plus size={16} />
              Create Question
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Link to="/questions" className="flex items-center gap-2 w-full">
              <HelpCircle size={16} />
              Manage Questions
            </Link>
          </DropdownMenuItem>
          <div className="h-px bg-border my-1" />
          <DropdownMenuItem>
            <Link to="/quiz/create" className="flex items-center gap-2 w-full">
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
          className="relative text-foreground hover:text-primary px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-accent/50 group flex items-center gap-2"
          activeClassName="text-primary bg-primary/10"
        >
          <LayoutDashboard size={16} />
          <span className="relative z-10">Admin</span>
          <span className="absolute inset-0 bg-gradient-to-r from-primary/0 via-primary/5 to-primary/0 opacity-0 group-hover:opacity-100 transition-opacity rounded-lg" />
        </Link>
      )}

      {/* User Menu */}
      {isAuthenticated && user ? (
        <div className="flex items-center gap-2 ml-2">
          <NotificationBell className="mx-1" />

          <DropdownMenu
            trigger={
              <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg hover:bg-accent/50 transition-all cursor-pointer">
                <div className="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-purple-600 flex items-center justify-center text-white text-sm font-semibold shadow-md">
                  {user.firstName?.charAt(0).toUpperCase()}
                </div>
                <div className="hidden xl:block text-left">
                  <div className="text-sm font-medium">{user.firstName}</div>
                  <div className="text-xs text-muted-foreground capitalize">
                    {user.role}
                  </div>
                </div>
              </div>
            }
            align="right"
          >
            <div className="px-3 py-2 border-b border-border">
              <div className="font-medium">
                {user.firstName} {user.lastName}
              </div>
              <div className="text-xs text-muted-foreground">{user.email}</div>
              <Badge className={`text-xs mt-2 ${getRoleBadgeColor(user.role)}`}>
                {user.role}
              </Badge>
            </div>
            <DropdownMenuItem>
              <Link to="/profile" className="flex items-center gap-2 w-full">
                <User size={16} />
                Profile
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem>
              <Link
                to="/notifications"
                className="flex items-center gap-2 w-full"
              >
                <Bell size={16} />
                Notifications
              </Link>
            </DropdownMenuItem>
            <div className="h-px bg-border my-1" />
            <DropdownMenuItem onClick={logout}>
              <span className="flex items-center gap-2 w-full text-destructive">
                <LogOut size={16} />
                Logout
              </span>
            </DropdownMenuItem>
          </DropdownMenu>
        </div>
      ) : (
        <div className="flex items-center gap-2 ml-2">
          <Link to="/auth/login">
            <Button variant="ghost" size="sm" className="hover:bg-accent/50">
              Login
            </Button>
          </Link>
          <Link to="/auth/register">
            <Button
              size="sm"
              className="bg-gradient-to-r from-primary to-purple-600 hover:from-primary/90 hover:to-purple-600/90 shadow-md"
            >
              Sign Up
            </Button>
          </Link>
        </div>
      )}
    </nav>
  );
};
