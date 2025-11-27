import React, { useEffect, useState } from "react";
import { Link, navigate } from "gatsby";
import Layout from "../../components/layout";
import { Story } from "../../types";
import { formatLevel, getLevelColor } from "../../utils/api";
import { publicApi } from "../../utils/apiClient";

interface StoryPageProps {
  params: {
    slug: string;
  };
}

const StoryPage: React.FC<StoryPageProps> = ({ params }) => {
  const [story, setStory] = useState<Story | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (params.slug) {
      loadStory(params.slug);
    }
  }, [params.slug]);

  const loadStory = async (slug: string) => {
    try {
      setLoading(true);
      setError(null);
      const response = await publicApi.getStoryBySlug(slug);
       
      if (response.data.success && response.data.data) {
        setStory(response.data.data);
      } else {
        setError(response.data.error || "Story not found");
      }
    } catch (err: any) {
      console.error("Error fetching story:", err);
      setError(err.response?.data?.error || "Failed to load story");
    } finally {
      setLoading(false);
    }
  };

  const handleTakeQuiz = () => {
    if (story && story.slug) {
      navigate(`/quiz/${story.slug}`);
    }
  };

  if (loading) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">Loading story...</p>
          </div>
        </div>
      </Layout>
    );
  }

  if (error || !story) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
            <div className="flex items-center">
              <span className="text-red-500 mr-2">⚠️</span>
              {error || "Story not found"}
            </div>
            <div className="mt-4">
              <Link
                to="/"
                className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded text-sm transition-colors"
              >
                Back to Stories
              </Link>
            </div>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Header */}
        <div className="mb-8">
          <Link
            to="/"
            className="inline-flex items-center text-blue-600 hover:text-blue-700 mb-4"
          >
            ← Back to Stories
          </Link>

          <div className="flex items-start justify-between mb-4">
            <h1 className="text-3xl md:text-4xl font-bold text-gray-900 flex-1 mr-4">
              {story.title}
            </h1>
            <span
              className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${getLevelColor(
                story.level
              )} flex-shrink-0`}
            >
              {story.level} - {formatLevel(story.level)}
            </span>
          </div>

          <div className="flex flex-wrap items-center gap-4 text-sm text-gray-600 mb-6">
            <span className="flex items-center">
              📚 {story.wordCount} words
            </span>
            <span className="flex items-center">
              ⏱️ {story.readingTime} min read
            </span>
            <span className="flex items-center">
              🎯 {story.topics?.join(", ")}
            </span>
          </div>

          {story.summary && (
            <div className="bg-blue-50 border-l-4 border-blue-400 p-4 mb-6">
              <h3 className="text-sm font-medium text-blue-900 mb-2">
                Story Summary
              </h3>
              <p className="text-blue-800 text-sm">{story.summary}</p>
            </div>
          )}
        </div>

        {/* Story Content */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-8 mb-8">
          <div className="prose prose-lg max-w-none">
            {story.content?.split("\n").map((paragraph, index) => (
              <p key={index} className="mb-4 text-gray-800 leading-relaxed">
                {paragraph}
              </p>
            ))}
          </div>
        </div>

        {/* Vocabulary Section */}
        {story.vocabulary && story.vocabulary.length > 0 && (
          <div className="bg-gray-50 rounded-lg p-6 mb-8">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Key Vocabulary 📝
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {story.vocabulary.map((word, index) => (
                <div
                  key={index}
                  className="bg-white rounded-lg p-4 border border-gray-200"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center mb-1">
                        <span className="font-semibold text-gray-900 mr-2">
                          {word.article && `${word.article} `}
                          {word.german}
                        </span>
                        <span className="text-xs px-2 py-1 bg-gray-100 text-gray-600 rounded">
                          {word.wordType}
                        </span>
                      </div>
                      <p className="text-gray-600 text-sm">{word.english}</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Action Buttons */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <button
            onClick={handleTakeQuiz}
            className="bg-green-600 hover:bg-green-700 text-white px-8 py-3 rounded-lg font-medium text-lg transition-colors"
          >
            Take Quiz 🎯
          </button>
          <Link
            to="/"
            className="bg-gray-600 hover:bg-gray-700 text-white px-8 py-3 rounded-lg font-medium text-lg transition-colors text-center"
          >
            More Stories 📚
          </Link>
        </div>
      </div>
    </Layout>
  );
};

export default StoryPage;
