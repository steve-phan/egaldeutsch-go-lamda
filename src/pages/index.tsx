import React, { useState, useEffect } from "react";
import { Link } from "gatsby";
import Layout from "../components/layout";
import StoryCard from "../components/StoryCard";
import { Story } from "../types";
import { fetchStories } from "../utils/api";

const IndexPage: React.FC = () => {
  const [stories, setStories] = useState<Story[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadStories();
  }, []);

  const loadStories = async () => {
    try {
      setLoading(true);
      setError(null);
      const fetchedStories = await fetchStories();
      setStories(fetchedStories);
    } catch (err) {
      console.error("Error fetching stories:", err);
      setError("Failed to load stories. Please try again later.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Hero Section */}
        <div className="text-center mb-12">
          <h1 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            Learn German Through Stories
          </h1>
          <p className="text-xl text-gray-600 mb-2">
            Read engaging stories and test your understanding with interactive
            quizzes
          </p>
          <p className="text-sm text-gray-500">
            Choose your level: A1 (Beginner) to C2 (Proficient)
          </p>
        </div>

        {/* Loading State */}
        {loading && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">Loading stories...</p>
          </div>
        )}

        {/* Error State */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
            <div className="flex items-center">
              <span className="text-red-500 mr-2">⚠️</span>
              {error}
            </div>
            <button
              onClick={loadStories}
              className="mt-3 bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded text-sm transition-colors"
            >
              Try Again
            </button>
          </div>
        )}

        {/* Content */}
        {!loading && !error && (
          <>
            {stories.length === 0 ? (
              <div className="text-center py-12">
                <div className="max-w-md mx-auto">
                  <span className="text-6xl mb-4 block">📚</span>
                  <h2 className="text-2xl font-semibold text-gray-900 mb-4">
                    No stories available yet
                  </h2>
                  <p className="text-gray-600 mb-6">
                    We're working on adding amazing German learning stories.
                    Check back soon!
                  </p>
                  <button
                    onClick={loadStories}
                    className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors"
                  >
                    Refresh
                  </button>
                </div>
              </div>
            ) : (
              <>
                {/* Stats */}
                <div className="bg-blue-50 rounded-lg p-6 mb-8">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-center">
                    <div>
                      <div className="text-2xl font-bold text-blue-600">
                        {stories.length}
                      </div>
                      <div className="text-sm text-gray-600">
                        Stories Available
                      </div>
                    </div>
                    <div>
                      <div className="text-2xl font-bold text-green-600">
                        {new Set(stories.map((s) => s.level)).size}
                      </div>
                      <div className="text-sm text-gray-600">
                        Difficulty Levels
                      </div>
                    </div>
                    <div>
                      <div className="text-2xl font-bold text-purple-600">
                        {new Set(stories.map((s) => s.topic)).size}
                      </div>
                      <div className="text-sm text-gray-600">Topics</div>
                    </div>
                  </div>
                </div>

                {/* Stories Grid */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {stories.map((story) => (
                    <StoryCard key={story.id} story={story} />
                  ))}
                </div>
              </>
            )}
          </>
        )}
      </div>
    </Layout>
  );
};

export default IndexPage;
