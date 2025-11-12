import React, { useState } from "react";
import { Question, QuizResult } from "../types";

interface QuizQuestionProps {
  question: Question;
  selectedAnswer?: number;
  showResult?: boolean;
  onAnswerSelect: (answerIndex: number) => void;
  questionNumber: number;
  totalQuestions: number;
}

const QuizQuestion: React.FC<QuizQuestionProps> = ({
  question,
  selectedAnswer,
  showResult = false,
  onAnswerSelect,
  questionNumber,
  totalQuestions,
}) => {
  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-6">
      <div className="flex items-center justify-between mb-4">
        <span className="text-sm text-gray-500">
          Question {questionNumber} of {totalQuestions}
        </span>
        <span
          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
            question.questionType === "comprehension"
              ? "bg-blue-100 text-blue-800"
              : question.questionType === "vocabulary"
              ? "bg-green-100 text-green-800"
              : "bg-purple-100 text-purple-800"
          }`}
        >
          {question.questionType}
        </span>
      </div>

      <h3 className="text-lg font-semibold text-gray-900 mb-4">
        {question.question}
      </h3>

      <div className="space-y-3">
        {question.options.map((option, index) => {
          const isSelected = selectedAnswer === index;
          const isCorrect = question.correctAnswer === index;
          const isIncorrect = showResult && isSelected && !isCorrect;

          let buttonClass =
            "w-full text-left p-4 rounded-lg border-2 transition-all duration-200 ";

          if (showResult) {
            if (isCorrect) {
              buttonClass += "border-green-500 bg-green-50 text-green-900";
            } else if (isIncorrect) {
              buttonClass += "border-red-500 bg-red-50 text-red-900";
            } else {
              buttonClass += "border-gray-200 bg-gray-50 text-gray-600";
            }
          } else {
            if (isSelected) {
              buttonClass += "border-blue-500 bg-blue-50 text-blue-900";
            } else {
              buttonClass +=
                "border-gray-200 hover:border-gray-300 hover:bg-gray-50";
            }
          }

          return (
            <button
              key={index}
              onClick={() => !showResult && onAnswerSelect(index)}
              disabled={showResult}
              className={buttonClass}
            >
              <div className="flex items-center">
                <span className="flex-shrink-0 w-6 h-6 rounded-full border-2 mr-3 flex items-center justify-center text-sm font-medium">
                  {String.fromCharCode(65 + index)}
                </span>
                <span className="flex-1">{option}</span>
                {showResult && isCorrect && (
                  <span className="flex-shrink-0 text-green-600 ml-2">✓</span>
                )}
                {showResult && isIncorrect && (
                  <span className="flex-shrink-0 text-red-600 ml-2">✗</span>
                )}
              </div>
            </button>
          );
        })}
      </div>

      {showResult && question.explanation && (
        <div className="mt-4 p-4 bg-blue-50 rounded-lg">
          <h4 className="text-sm font-medium text-blue-900 mb-2">
            Explanation:
          </h4>
          <p className="text-sm text-blue-800">{question.explanation}</p>
        </div>
      )}

      <div className="mt-4 text-right">
        <span className="text-sm text-gray-500">Points: {question.points}</span>
      </div>
    </div>
  );
};

interface QuizResultsProps {
  result: QuizResult;
  onRetakeQuiz: () => void;
  onReturnToStory: () => void;
}

const QuizResults: React.FC<QuizResultsProps> = ({
  result,
  onRetakeQuiz,
  onReturnToStory,
}) => {
  const getScoreColor = (percentage: number) => {
    if (percentage >= 90) return "text-green-600";
    if (percentage >= 70) return "text-blue-600";
    if (percentage >= 50) return "text-yellow-600";
    return "text-red-600";
  };

  const getScoreMessage = (percentage: number) => {
    if (percentage >= 90) return "Excellent! 🎉";
    if (percentage >= 70) return "Great job! 👏";
    if (percentage >= 50) return "Good effort! 💪";
    return "Keep practicing! 📚";
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-8 text-center">
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Quiz Results</h2>
        <div
          className={`text-4xl font-bold mb-2 ${getScoreColor(
            result.percentage
          )}`}
        >
          {result.percentage}%
        </div>
        <p className="text-lg text-gray-600">
          {result.score} out of {result.totalQuestions} correct
        </p>
        <p className="text-lg font-medium text-gray-800 mt-2">
          {getScoreMessage(result.percentage)}
        </p>
      </div>

      {result.passed ? (
        <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-6">
          <div className="flex items-center justify-center text-green-600 mb-2">
            <span className="text-2xl mr-2">🎊</span>
            <span className="font-semibold">Congratulations!</span>
          </div>
          <p className="text-green-800 text-sm">
            You passed the quiz! You now have a better understanding of this
            German story.
          </p>
        </div>
      ) : (
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6">
          <div className="flex items-center justify-center text-yellow-600 mb-2">
            <span className="text-2xl mr-2">📖</span>
            <span className="font-semibold">Keep Learning!</span>
          </div>
          <p className="text-yellow-800 text-sm">
            Don't worry! Read the story again and try the quiz once more to
            improve your understanding.
          </p>
        </div>
      )}

      <div className="flex space-x-4 justify-center">
        <button
          onClick={onReturnToStory}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors"
        >
          Read Story Again
        </button>
        <button
          onClick={onRetakeQuiz}
          className="bg-green-600 hover:bg-green-700 text-white px-6 py-3 rounded-lg font-medium transition-colors"
        >
          Retake Quiz
        </button>
      </div>
    </div>
  );
};

export { QuizQuestion, QuizResults };
