import React from "react";
import { Link } from "gatsby";
import "../styles/global.css";
import { Separator } from "./ui";
import { Button } from "./ui/button";
import { Badge } from "./ui/badge";
import { cn } from "@/lib/utils";
import { useAuth } from "../contexts/AuthContext";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { isAuthenticated, user, logout } = useAuth();

  return (
    <div className="min-h-screen bg-background">
      <header className="bg-card shadow-sm border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <Link
              to="/"
              className="text-2xl font-bold text-primary hover:text-primary/90 transition-colors"
            >
              🇩🇪 EgalDeutsch
            </Link>
            <nav className="flex items-center space-x-6">
              <Link
                to="/"
                className="text-foreground hover:text-primary px-3 py-2 rounded-md text-sm font-medium transition-colors"
                activeClassName="text-primary bg-primary/10"
              >
                Stories
              </Link>
              <Link
                to="/about"
                className="text-foreground hover:text-primary px-3 py-2 rounded-md text-sm font-medium transition-colors"
                activeClassName="text-primary bg-primary/10"
              >
                About
              </Link>
              {isAuthenticated && user ? (
                <>
                  <Link
                    to="/profile"
                    className="text-foreground hover:text-primary px-3 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
                    activeClassName="text-primary bg-primary/10"
                  >
                    <span>{user.firstName}</span>
                    <Badge variant="outline" className="text-xs">
                      {user.role}
                    </Badge>
                  </Link>
                  <Button 
                    variant="ghost" 
                    size="sm" 
                    onClick={logout}
                    className="text-sm"
                  >
                    Logout
                  </Button>
                </>
              ) : (
                <>
                  <Link to="/auth/login">
                    <Button variant="ghost" size="sm">
                      Login
                    </Button>
                  </Link>
                  <Link to="/auth/register">
                    <Button size="sm">
                      Sign Up
                    </Button>
                  </Link>
                </>
              )}
            </nav>
          </div>
        </div>
      </header>
      <main className="flex-1">{children}</main>
      <footer className="bg-card border-t border-border mt-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="text-center">
            <p className="text-muted-foreground text-sm mb-2">
              Learn German through engaging stories and interactive quizzes
            </p>
            <p className="text-muted-foreground text-xs">
              © {new Date().getFullYear()} EgalDeutsch. Practice makes perfect!
              🚀
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Layout;
