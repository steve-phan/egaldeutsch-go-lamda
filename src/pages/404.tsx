import React from "react";
import Layout from "../components/layout";
import { Link } from "gatsby";

const NotFoundPage: React.FC = () => {
  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 text-center">
        <div className="mb-8">
          <span className="text-8xl mb-4 block">🤔</span>
          <h1 className="text-6xl font-bold text-gray-900 mb-4">404</h1>
          <p className="text-xl text-gray-600 mb-4">Seite nicht gefunden</p>
          <p className="text-gray-500 mb-8">
            The page you're looking for doesn't exist or has been moved.
          </p>
        </div>

        <div className="space-y-4">
          <Link
            to="/"
            className="inline-block bg-blue-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-blue-700 transition-colors duration-200 mr-4"
          >
            Back to Stories
          </Link>
          <Link
            to="/about"
            className="inline-block bg-gray-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-gray-700 transition-colors duration-200"
          >
            About EgalDeutsch
          </Link>
        </div>
      </div>
    </Layout>
  );
};

export default NotFoundPage;
