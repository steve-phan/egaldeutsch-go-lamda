import React, { useState, useEffect } from "react";
import Layout from "../components/layout";
import axios from "axios";

interface LeaderboardEntry {
  userId: string;
  username: string;
  totalPoints: number;
  quizzesTaken: number;
  quizzesPassed: number;
  averageScore: number;
  rank: number;
  lastSubmittedAt: string;
}

interface LeaderboardData {
  entries: LeaderboardEntry[];
  totalParticipants: number;
  totalQuizzes: number;
  generatedAt: string;
}

const LeaderboardPage: React.FC = () => {
  const [data, setData] = useState<LeaderboardData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadLeaderboard();
  }, []);

  const loadLeaderboard = async () => {
    try {
      setLoading(true);
      setError(null);

      const apiUrl = process.env.GATSBY_API_URL || "/.netlify/functions";
      const response = await axios.get(`${apiUrl}/leaderboard?limit=100`);

      if (response.data.success) {
        setData(response.data.data);
      } else {
        setError(response.data.error || "Failed to load leaderboard");
      }
    } catch (err) {
      console.error("Error fetching leaderboard:", err);
      setError("Failed to load leaderboard. Please try again later.");
    } finally {
      setLoading(false);
    }
  };

  const getRankBadgeColor = (rank: number) => {
    if (rank === 1) return "bg-yellow-500 text-white";
    if (rank === 2) return "bg-gray-400 text-white";
    if (rank === 3) return "bg-orange-600 text-white";
    return "bg-blue-100 text-blue-800";
  };

  const getRankEmoji = (rank: number) => {
    if (rank === 1) return "🥇";
    if (rank === 2) return "🥈";
    if (rank === 3) return "🥉";
    return "";
  };

  return (
    <Layout>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Hero Section */}
        <div className="text-center mb-12">
          <h1 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            🏆 Leaderboard
          </h1>
          <p className="text-xl text-gray-600 mb-2">
            Top German learners on EgalDeutsch
          </p>
          <p className="text-sm text-gray-500">
            Rankings based on total points earned from quizzes
          </p>
        </div>

        {/* Loading State */}
        {loading && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">Loading leaderboard...</p>
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
              onClick={loadLeaderboard}
              className="mt-3 bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded text-sm transition-colors"
            >
              Try Again
            </button>
          </div>
        )}

        {/* Content */}
        {!loading && !error && data && (
          <>
            {/* Statistics Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
              <div className="bg-blue-50 rounded-lg p-6 text-center">
                <div className="text-3xl font-bold text-blue-600 mb-2">
                  {data.totalParticipants}
                </div>
                <div className="text-sm text-gray-600">Total Participants</div>
              </div>
              <div className="bg-green-50 rounded-lg p-6 text-center">
                <div className="text-3xl font-bold text-green-600 mb-2">
                  {data.totalQuizzes}
                </div>
                <div className="text-sm text-gray-600">Quizzes Completed</div>
              </div>
              <div className="bg-purple-50 rounded-lg p-6 text-center">
                <div className="text-3xl font-bold text-purple-600 mb-2">
                  {data.entries?.length}
                </div>
                <div className="text-sm text-gray-600">Ranked Players</div>
              </div>
            </div>

            {/* Leaderboard Table */}
            {data.entries?.length === 0 ? (
              <div className="text-center py-12">
                <div className="max-w-md mx-auto">
                  <span className="text-6xl mb-4 block">🏆</span>
                  <h2 className="text-2xl font-semibold text-gray-900 mb-4">
                    No rankings yet
                  </h2>
                  <p className="text-gray-600 mb-6">
                    Be the first to complete a quiz and make it to the
                    leaderboard!
                  </p>
                </div>
              </div>
            ) : (
              <div className="bg-white rounded-lg shadow-md overflow-hidden">
                {/* Desktop View */}
                <div className="hidden md:block overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Rank
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Username
                        </th>
                        <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Total Points
                        </th>
                        <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Quizzes Taken
                        </th>
                        <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Pass Rate
                        </th>
                        <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Avg Score
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {data.entries?.map((entry) => (
                        <tr
                          key={entry.userId}
                          className="hover:bg-gray-50 transition-colors"
                        >
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center">
                              <span
                                className={`inline-flex items-center justify-center w-10 h-10 rounded-full font-bold text-sm ${getRankBadgeColor(
                                  entry.rank
                                )}`}
                              >
                                {getRankEmoji(entry.rank) || `#${entry.rank}`}
                              </span>
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="text-sm font-medium text-gray-900">
                              {entry.username || "Anonymous"}
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-center">
                            <div className="text-sm font-semibold text-blue-600">
                              {entry.totalPoints}
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-center">
                            <div className="text-sm text-gray-900">
                              {entry.quizzesTaken}
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-center">
                            <div className="text-sm text-gray-900">
                              {entry.quizzesTaken > 0
                                ? `${Math.round(
                                    (entry.quizzesPassed / entry.quizzesTaken) *
                                      100
                                  )}%`
                                : "N/A"}
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-center">
                            <div className="text-sm text-gray-900">
                              {entry.averageScore.toFixed(1)}%
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {/* Mobile View */}
                <div className="md:hidden divide-y divide-gray-200">
                  {data.entries?.map((entry) => (
                    <div
                      key={entry.userId}
                      className="p-4 hover:bg-gray-50 transition-colors"
                    >
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center space-x-3">
                          <span
                            className={`inline-flex items-center justify-center w-12 h-12 rounded-full font-bold ${getRankBadgeColor(
                              entry.rank
                            )}`}
                          >
                            {getRankEmoji(entry.rank) || `#${entry.rank}`}
                          </span>
                          <div>
                            <div className="font-medium text-gray-900">
                              {entry.username || "Anonymous"}
                            </div>
                            <div className="text-sm text-gray-500">
                              {entry.quizzesTaken} quizzes
                            </div>
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="text-lg font-bold text-blue-600">
                            {entry.totalPoints}
                          </div>
                          <div className="text-xs text-gray-500">points</div>
                        </div>
                      </div>
                      <div className="grid grid-cols-2 gap-2 text-sm">
                        <div>
                          <span className="text-gray-500">Pass Rate:</span>{" "}
                          <span className="font-medium">
                            {entry.quizzesTaken > 0
                              ? `${Math.round(
                                  (entry.quizzesPassed / entry.quizzesTaken) *
                                    100
                                )}%`
                              : "N/A"}
                          </span>
                        </div>
                        <div>
                          <span className="text-gray-500">Avg Score:</span>{" "}
                          <span className="font-medium">
                            {entry.averageScore.toFixed(1)}%
                          </span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Footer Info */}
            <div className="mt-8 text-center text-sm text-gray-500">
              <p>
                Leaderboard updates in real-time. Keep learning to climb the
                ranks! 🚀
              </p>
              <p className="mt-2">
                Last updated: {new Date(data.generatedAt).toLocaleString()}
              </p>
            </div>
          </>
        )}
      </div>
    </Layout>
  );
};

export default LeaderboardPage;
