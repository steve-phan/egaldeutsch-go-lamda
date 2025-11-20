import React, { useState, useEffect } from "react";
import { Link, navigate } from "gatsby";
import Layout from "../../components/layout";
import { QuizQuestion, QuizResults } from "../../components/Quiz";
import { Quiz, QuizResult } from "../../types";
import { fetchQuizByStorySlug, submitQuiz } from "../../utils/api";

interface QuizPageProps {
  params: {
    slug: string;
  };
}

const QuizPage: React.FC<QuizPageProps> = ({ params }) => {
  const [quiz, setQuiz] = useState<Quiz | null>(null);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState<number>(0);
  const [answers, setAnswers] = useState<number[]>([]);
  const [quizResult, setQuizResult] = useState<QuizResult | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [showResults, setShowResults] = useState<boolean>(false);

  useEffect(() => {
    if (params.slug) {
      loadQuiz(params.slug);
    }
  }, [params.slug]);

  const loadQuiz = async (slug: string) => {
    try {
      setLoading(true);
      setError(null);
      const fetchedQuiz = await fetchQuizByStorySlug(slug);
      setQuiz(fetchedQuiz);
      // Initialize answers array with -1 (no answer selected)
      setAnswers(new Array(fetchedQuiz.questions.length).fill(-1));
    } catch (err: any) {
      console.error("Error fetching quiz:", err);
      // Check if it's a "no questions" error
      const errorMessage =
        err?.response?.data?.error ||
        err?.message ||
        "Failed to load quiz. Please try again later.";
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleAnswerSelect = (answerIndex: number) => {
    const newAnswers = [...answers];
    newAnswers[currentQuestionIndex] = answerIndex;
    setAnswers(newAnswers);
  };

  const handleNextQuestion = () => {
    if (quiz && currentQuestionIndex < quiz.questions.length - 1) {
      setCurrentQuestionIndex(currentQuestionIndex + 1);
    }
  };

  const handlePreviousQuestion = () => {
    if (currentQuestionIndex > 0) {
      setCurrentQuestionIndex(currentQuestionIndex - 1);
    }
  };

  const handleSubmitQuiz = async () => {
    if (!quiz || answers.includes(-1)) {
      alert("Please answer all questions before submitting.");
      return;
    }

    try {
      setSubmitting(true);
      const result = await submitQuiz(quiz.storyId, answers);
      setQuizResult(result);
      setShowResults(true);
    } catch (err) {
      console.error("Error submitting quiz:", err);
      alert("Failed to submit quiz. Please try again.");
    } finally {
      setSubmitting(false);
    }
  };

  const handleRetakeQuiz = () => {
    setAnswers(new Array(quiz?.questions.length || 0).fill(-1));
    setCurrentQuestionIndex(0);
    setQuizResult(null);
    setShowResults(false);
  };

  const handleReturnToStory = () => {
    if (quiz?.story?.slug) {
      navigate(`/story/${quiz.story.slug}`);
    }
  };

  if (loading) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-green-600"></div>
            <p className="mt-4 text-gray-600">Loading quiz...</p>
          </div>
        </div>
      </Layout>
    );
  }

  if (error || !quiz) {
    const isNoQuestionsError =
      error && error.includes("No questions available");

    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          {isNoQuestionsError ? (
            // Friendly "no questions" message
            <div className="text-center py-12">
              <div className="max-w-md mx-auto">
                <span className="text-6xl mb-4 block">📝</span>
                <h2 className="text-2xl font-semibold text-gray-900 mb-4">
                  No Quiz Questions Yet
                </h2>
                <p className="text-gray-600 mb-6">
                  This story doesn't have quiz questions available at the
                  moment. You can still read the story and check back later for
                  the quiz!
                </p>
                <div className="space-y-3">
                  <Link
                    to={`/story/${params.slug}`}
                    className="block bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors"
                  >
                    Read Story Instead
                  </Link>
                  <Link
                    to="/"
                    className="block bg-gray-100 hover:bg-gray-200 text-gray-700 px-6 py-3 rounded-lg font-medium transition-colors"
                  >
                    Browse Other Stories
                  </Link>
                </div>
              </div>
            </div>
          ) : (
            // Regular error message
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
              <div className="flex items-center">
                <span className="text-red-500 mr-2">⚠️</span>
                {error || "Quiz not found"}
              </div>
              <div className="mt-4">
                <Link
                  to="/"
                  className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded text-sm transition-colors mr-3"
                >
                  Back to Stories
                </Link>
                {params.slug && (
                  <>
                    <Link
                      to={`/story/${params.slug}`}
                      className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm transition-colors mr-3"
                    >
                      Read Story First
                    </Link>
                    <button
                      onClick={() => loadQuiz(params.slug)}
                      className="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded text-sm transition-colors"
                    >
                      Try Again
                    </button>
                  </>
                )}
              </div>
            </div>
          )}
        </div>
      </Layout>
    );
  }

  if (showResults && quizResult) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <QuizResults
            result={quizResult}
            onRetakeQuiz={handleRetakeQuiz}
            onReturnToStory={handleReturnToStory}
          />
        </div>
      </Layout>
    );
  }

  const currentQuestion = quiz.questions[currentQuestionIndex];
  const progress = ((currentQuestionIndex + 1) / quiz.questions.length) * 100;
  const isLastQuestion = currentQuestionIndex === quiz.questions.length - 1;
  const allQuestionsAnswered = !answers.includes(-1);

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Header */}
        <div className="mb-8">
          <Link
            to={`/story/${quiz.story?.slug || params.slug}`}
            className="inline-flex items-center text-blue-600 hover:text-blue-700 mb-4"
          >
            ← Back to Story
          </Link>

          <div className="flex items-center justify-between mb-4">
            <h1 className="text-2xl md:text-3xl font-bold text-gray-900">
              Quiz: {quiz.story?.title}
            </h1>
            <span className="text-sm text-gray-500">
              {currentQuestionIndex + 1} of {quiz.questions.length}
            </span>
          </div>

          {/* Progress Bar */}
          <div className="w-full bg-gray-200 rounded-full h-2 mb-6">
            <div
              className="bg-green-600 h-2 rounded-full transition-all duration-300"
              style={{ width: `${progress}%` }}
            ></div>
          </div>
        </div>

        {/* Question */}
        <QuizQuestion
          question={currentQuestion}
          selectedAnswer={answers[currentQuestionIndex]}
          onAnswerSelect={handleAnswerSelect}
          questionNumber={currentQuestionIndex + 1}
          totalQuestions={quiz.questions.length}
        />

        {/* Navigation */}
        <div className="flex justify-between items-center mt-8">
          <button
            onClick={handlePreviousQuestion}
            disabled={currentQuestionIndex === 0}
            className="bg-gray-600 hover:bg-gray-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white px-6 py-3 rounded-lg font-medium transition-colors"
          >
            Previous
          </button>

          {isLastQuestion ? (
            <button
              onClick={handleSubmitQuiz}
              disabled={!allQuestionsAnswered || submitting}
              className="bg-green-600 hover:bg-green-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white px-6 py-3 rounded-lg font-medium transition-colors"
            >
              {submitting ? "Submitting..." : "Submit Quiz"}
            </button>
          ) : (
            <button
              onClick={handleNextQuestion}
              className="bg-green-600 hover:bg-green-700 text-white px-6 py-3 rounded-lg font-medium transition-colors"
            >
              Next
            </button>
          )}
        </div>

        {/* Answer status */}
        <div className="mt-6 text-center">
          <p className="text-sm text-gray-600">
            {answers.filter((a) => a !== -1).length} of {quiz.questions.length}{" "}
            questions answered
          </p>
          {!allQuestionsAnswered && (
            <p className="text-xs text-orange-600 mt-1">
              Please answer all questions to submit the quiz
            </p>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default QuizPage;
