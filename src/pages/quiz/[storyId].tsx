import React, { useState, useEffect } from "react";
import { Link, navigate } from "gatsby";
import Layout from "../../components/layout";
import { QuizQuestion, QuizResults } from "../../components/Quiz";
import { Quiz, QuizResult } from "../../types";
import { fetchQuizByStoryId, submitQuiz } from "../../utils/api";

interface QuizPageProps {
  params: {
    storyId: string;
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
    if (params.storyId) {
      loadQuiz(params.storyId);
    }
  }, [params.storyId]);

  const loadQuiz = async (storyId: string) => {
    try {
      setLoading(true);
      setError(null);
      const fetchedQuiz = await fetchQuizByStoryId(storyId);
      setQuiz(fetchedQuiz);
      // Initialize answers array with -1 (no answer selected)
      setAnswers(new Array(fetchedQuiz.questions.length).fill(-1));
    } catch (err) {
      console.error("Error fetching quiz:", err);
      setError("Failed to load quiz. Please try again later.");
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
    if (quiz?.storyId) {
      navigate(`/story/${quiz.storyId}`);
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
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
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
              {params.storyId && (
                <>
                  <Link
                    to={`/story/${params.storyId}`}
                    className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm transition-colors mr-3"
                  >
                    Read Story First
                  </Link>
                  <button
                    onClick={() => loadQuiz(params.storyId)}
                    className="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded text-sm transition-colors"
                  >
                    Try Again
                  </button>
                </>
              )}
            </div>
          </div>
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
            to={`/story/${quiz.storyId}`}
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

          <div className="flex space-x-3">
            {/* Question indicators */}
            {quiz.questions.map((_, index) => (
              <button
                key={index}
                onClick={() => setCurrentQuestionIndex(index)}
                className={`w-8 h-8 rounded-full text-sm font-medium transition-colors ${
                  index === currentQuestionIndex
                    ? "bg-green-600 text-white"
                    : answers[index] !== -1
                    ? "bg-green-100 text-green-800"
                    : "bg-gray-200 text-gray-600"
                }`}
              >
                {index + 1}
              </button>
            ))}
          </div>

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
