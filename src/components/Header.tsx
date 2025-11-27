import React, { useState, useEffect } from "react";
import { Link } from "gatsby";
import { useAuth } from "../contexts/AuthContext";
import { DekstopNavigation } from "./atoms/DesktopNavigation";
import { MobileNavigationDrawer } from "./atoms/MobileNavigationDrawer";
import { Menu } from "lucide-react";

const Header: React.FC = () => {
  const { logout } = useAuth();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isScrolled, setIsScrolled] = useState(false);

  const handleMenuToggle = (open: boolean) => {
    setIsMobileMenuOpen(open);
  };

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
          <DekstopNavigation />

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
          onClick={() => handleMenuToggle(false)}
        />
      )}

      {/* Mobile Navigation Drawer */}
      <MobileNavigationDrawer
        isMobileMenuOpen={isMobileMenuOpen}
        setIsMobileMenuOpen={handleMenuToggle}
      />
    </header>
  );
};

export default Header;
