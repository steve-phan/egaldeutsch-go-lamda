import React from "react";
import { Link } from "gatsby";
import { Separator } from "./ui/separator";
import { useAuth } from "../contexts/AuthContext";

const Footer: React.FC = () => {
  const { isAuthenticated, user } = useAuth();
  const currentYear = new Date().getFullYear();

  const isCreatorOrAdmin = user?.role === "creator" || user?.role === "admin";
  const isReviewerOrAdmin = user?.role === "reviewer" || user?.role === "admin";

  return (
    <footer className="bg-card border-t border-border mt-12">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Main Footer Content */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-8">
          {/* About Section */}
          <div>
            <h3 className="text-lg font-semibold mb-3 text-foreground">
              EgalDeutsch
            </h3>
            <p className="text-muted-foreground text-sm mb-3">
              Learn German through engaging stories and interactive quizzes.
              Practice makes perfect! 🚀
            </p>
          </div>

          {/* Quick Links */}
          <div>
            <h3 className="text-sm font-semibold mb-3 text-foreground uppercase tracking-wide">
              Quick Links
            </h3>
            <ul className="space-y-2">
              <li>
                <Link
                  to="/"
                  className="text-muted-foreground text-sm hover:text-primary transition-colors"
                >
                  Home
                </Link>
              </li>
              <li>
                <Link
                  to="/about"
                  className="text-muted-foreground text-sm hover:text-primary transition-colors"
                >
                  About Us
                </Link>
              </li>
              <li>
                <Link
                  to="/leaderboard"
                  className="text-muted-foreground text-sm hover:text-primary transition-colors"
                >
                  Leaderboard
                </Link>
              </li>
              {isReviewerOrAdmin && (
                <li>
                  <Link
                    to="/admin"
                    className="text-muted-foreground text-sm hover:text-primary transition-colors"
                  >
                    Admin Dashboard
                  </Link>
                </li>
              )}
            </ul>
          </div>

          {/* Resources & Tools */}
          <div>
            <h3 className="text-sm font-semibold mb-3 text-foreground uppercase tracking-wide">
              {isAuthenticated ? "Resources" : "Get Started"}
            </h3>
            <ul className="space-y-2">
              {isAuthenticated ? (
                <>
                  <li>
                    <Link
                      to="/profile"
                      className="text-muted-foreground text-sm hover:text-primary transition-colors"
                    >
                      My Profile
                    </Link>
                  </li>
                  {isCreatorOrAdmin && (
                    <>
                      <li>
                        <Link
                          to="/stories/create"
                          className="text-muted-foreground text-sm hover:text-primary transition-colors"
                        >
                          Create Story
                        </Link>
                      </li>
                      <li>
                        <Link
                          to="/questions/create"
                          className="text-muted-foreground text-sm hover:text-primary transition-colors"
                        >
                          Create Question
                        </Link>
                      </li>
                    </>
                  )}
                  <li>
                    <Link
                      to="/notifications"
                      className="text-muted-foreground text-sm hover:text-primary transition-colors"
                    >
                      Notifications
                    </Link>
                  </li>
                </>
              ) : (
                <>
                  <li>
                    <Link
                      to="/auth/register"
                      className="text-muted-foreground text-sm hover:text-primary transition-colors"
                    >
                      Sign Up Free
                    </Link>
                  </li>
                  <li>
                    <Link
                      to="/auth/login"
                      className="text-muted-foreground text-sm hover:text-primary transition-colors"
                    >
                      Login
                    </Link>
                  </li>
                  <li>
                    <Link
                      to="/about"
                      className="text-muted-foreground text-sm hover:text-primary transition-colors"
                    >
                      How It Works
                    </Link>
                  </li>
                </>
              )}
            </ul>
          </div>

          {/* Legal & Support */}
          <div>
            <h3 className="text-sm font-semibold mb-3 text-foreground uppercase tracking-wide">
              Legal & Support
            </h3>
            <ul className="space-y-2">
              <li>
                <Link
                  to="/privacy-policy"
                  className="text-muted-foreground text-sm hover:text-primary transition-colors"
                >
                  Privacy Policy
                </Link>
              </li>
              <li>
                <Link
                  to="/terms-of-service"
                  className="text-muted-foreground text-sm hover:text-primary transition-colors"
                >
                  Terms of Service
                </Link>
              </li>
              <li>
                <a
                  href="mailto:support@egaldeutsch.com"
                  className="text-muted-foreground text-sm hover:text-primary transition-colors"
                >
                  Contact Support
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/steve-phan/egaldeutsch-go-lamda"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-muted-foreground text-sm hover:text-primary transition-colors"
                >
                  GitHub
                </a>
              </li>
            </ul>
          </div>
        </div>

        <Separator className="mb-6" />

        {/* Bottom Footer */}
        <div className="flex flex-col sm:flex-row justify-between items-center gap-4">
          <p className="text-muted-foreground text-xs text-center sm:text-left">
            © {currentYear} EgalDeutsch. All rights reserved.
          </p>
          <div className="flex items-center gap-4">
            <span className="text-muted-foreground text-xs">
              Made with ❤️ for German learners
            </span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
