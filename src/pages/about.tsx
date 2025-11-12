import React from "react";
import Layout from "../components/layout";
import { Link } from "gatsby";

const AboutPage: React.FC = () => {
  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center mb-12">
          <h1 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            About EgalDeutsch
          </h1>
          <p className="text-xl text-gray-600">
            Learn German through engaging stories and interactive quizzes
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-12">
          <div className="bg-blue-50 rounded-lg p-6">
            <div className="text-3xl mb-4">📚</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-3">
              Story-Based Learning
            </h3>
            <p className="text-gray-600">
              Learn German naturally through engaging stories that are carefully
              crafted for different proficiency levels, from A1 (Beginner) to C2
              (Proficient).
            </p>
          </div>

          <div className="bg-green-50 rounded-lg p-6">
            <div className="text-3xl mb-4">🎯</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-3">
              Interactive Quizzes
            </h3>
            <p className="text-gray-600">
              Test your comprehension with carefully designed quizzes that cover
              vocabulary, grammar, and reading comprehension to reinforce your
              learning.
            </p>
          </div>

          <div className="bg-purple-50 rounded-lg p-6">
            <div className="text-3xl mb-4">🧠</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-3">
              Progressive Difficulty
            </h3>
            <p className="text-gray-600">
              Stories are organized by CEFR levels (A1-C2), allowing you to
              progress at your own pace and build confidence as you advance.
            </p>
          </div>

          <div className="bg-orange-50 rounded-lg p-6">
            <div className="text-3xl mb-4">📖</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-3">
              Key Vocabulary
            </h3>
            <p className="text-gray-600">
              Each story includes highlighted vocabulary with translations,
              articles (der/die/das), and word types to help expand your German
              vocabulary.
            </p>
          </div>
        </div>

        <div className="bg-gray-50 rounded-lg p-8 mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 text-center">
            How It Works
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="text-center">
              <div className="bg-blue-600 text-white w-8 h-8 rounded-full flex items-center justify-center mx-auto mb-3 font-bold">
                1
              </div>
              <h3 className="font-semibold text-gray-900 mb-2">Read a Story</h3>
              <p className="text-sm text-gray-600">
                Choose a story that matches your German level and read it
                carefully
              </p>
            </div>
            <div className="text-center">
              <div className="bg-green-600 text-white w-8 h-8 rounded-full flex items-center justify-center mx-auto mb-3 font-bold">
                2
              </div>
              <h3 className="font-semibold text-gray-900 mb-2">
                Take the Quiz
              </h3>
              <p className="text-sm text-gray-600">
                Answer 10 questions about the story to test your understanding
              </p>
            </div>
            <div className="text-center">
              <div className="bg-purple-600 text-white w-8 h-8 rounded-full flex items-center justify-center mx-auto mb-3 font-bold">
                3
              </div>
              <h3 className="font-semibold text-gray-900 mb-2">
                Learn & Improve
              </h3>
              <p className="text-sm text-gray-600">
                Review your results and vocabulary to strengthen your German
                skills
              </p>
            </div>
          </div>
        </div>

        <div className="bg-gradient-to-r from-blue-600 to-green-600 rounded-lg p-8 text-white text-center">
          <h2 className="text-2xl font-bold mb-4">Ready to Start Learning?</h2>
          <p className="text-blue-100 mb-6">
            Join thousands of learners improving their German through stories!
          </p>
          <Link
            to="/"
            className="bg-white text-blue-600 px-8 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-colors inline-flex items-center"
          >
            Start Reading Stories →
          </Link>
        </div>

        <div className="mt-12 text-center text-gray-500 text-sm">
          <p>
            EgalDeutsch uses the Common European Framework of Reference for
            Languages (CEFR) to structure learning content from beginner (A1) to
            proficient (C2) levels.
          </p>
        </div>
      </div>
    </Layout>
  );
};

export default AboutPage;
