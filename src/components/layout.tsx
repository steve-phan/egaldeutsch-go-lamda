import React from "react";
import { Link } from "gatsby";
import "../styles/global.css";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <Link to="/" className="text-2xl font-bold text-blue-600 hover:text-blue-700 transition-colors">
              🇩🇪 EgalDeutsch
            </Link>
            <nav className="flex space-x-6">
              <Link
                to="/"
                className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
                activeClassName="text-blue-600 bg-blue-50"
              >
                Stories
              </Link>
              <Link
                to="/about"
                className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
                activeClassName="text-blue-600 bg-blue-50"
              >
                About
              </Link>
            </nav>
          </div>
        </div>
      </header>
      <main className="flex-1">{children}</main>
      <footer className="bg-white border-t mt-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="text-center">
            <p className="text-gray-600 text-sm mb-2">
              Learn German through engaging stories and interactive quizzes
            </p>
            <p className="text-gray-500 text-xs">
              © {new Date().getFullYear()} EgalDeutsch. Practice makes perfect! 🚀
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Layout;