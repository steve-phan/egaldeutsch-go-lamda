import React, { useState } from "react";
import { Question, QuizResult } from "../types";
import { Card, CardContent, Badge, Button, Alert, AlertTitle, AlertDescription } from "./ui";
import { cn } from "@/lib/utils";

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
    <Card className="mb-6">
      <CardContent className="pt-6">
        <div className="flex items-center justify-between mb-4">
          <span className="text-sm text-muted-foreground">
            Question {questionNumber} of {totalQuestions}
          </span>
          <Badge
            variant={
              question.questionType === "comprehension"
                ? "default"
                : question.questionType === "vocabulary"
                ? "success"
                : "secondary"
            }
          >
            {question.questionType}
          </Badge>
        </div>

        <h3 className="text-lg font-semibold text-foreground mb-4">
          {question.question}
        </h3>

        <div className="space-y-3">
          {question.options.map((option, index) => {
            const isSelected = selectedAnswer === index;
            const isCorrect = question.correctAnswer === index;
            const isIncorrect = showResult && isSelected && !isCorrect;

            let buttonClass = cn(
              "w-full text-left p-4 rounded-lg border-2 transition-all duration-200",
              {
                "border-success bg-success/10 text-success": showResult && isCorrect,
                "border-destructive bg-destructive/10 text-destructive": isIncorrect,
                "border-muted bg-muted/50 text-muted-foreground": showResult && !isCorrect && !isSelected,
                "border-primary bg-primary/10 text-primary": !showResult && isSelected,
                "border-input hover:border-muted-foreground hover:bg-accent": !showResult && !isSelected,
              }
            );

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
                    <span className="flex-shrink-0 text-success ml-2">✓</span>
                  )}
                  {showResult && isIncorrect && (
                    <span className="flex-shrink-0 text-destructive ml-2">✗</span>
                  )}
                </div>
              </button>
            );
          })}
        </div>

        {showResult && question.explanation && (
          <Alert variant="default" className="mt-4 bg-primary/10 border-primary/20">
            <AlertTitle className="text-primary">Explanation:</AlertTitle>
            <AlertDescription className="text-primary/90">
              {question.explanation}
            </AlertDescription>
          </Alert>
        )}

        <div className="mt-4 text-right">
          <span className="text-sm text-muted-foreground">Points: {question.points}</span>
        </div>
      </CardContent>
    </Card>
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
    if (percentage >= 90) return "text-success";
    if (percentage >= 70) return "text-primary";
    if (percentage >= 50) return "text-warning";
    return "text-destructive";
  };

  const getScoreMessage = (percentage: number) => {
    if (percentage >= 90) return "Excellent! 🎉";
    if (percentage >= 70) return "Great job! 👏";
    if (percentage >= 50) return "Good effort! 💪";
    return "Keep practicing! 📚";
  };

  return (
    <Card className="text-center">
      <CardContent className="pt-8 pb-8">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-foreground mb-2">Quiz Results</h2>
          <div
            className={cn(
              "text-4xl font-bold mb-2",
              getScoreColor(result.percentage)
            )}
          >
            {result.percentage}%
          </div>
          <p className="text-lg text-muted-foreground">
            {result.score} out of {result.totalQuestions} correct
          </p>
          <p className="text-lg font-medium text-foreground mt-2">
            {getScoreMessage(result.percentage)}
          </p>
        </div>

        {result.passed ? (
          <Alert variant="success" className="mb-6">
            <AlertTitle className="flex items-center justify-center">
              <span className="text-2xl mr-2">🎊</span>
              <span className="font-semibold">Congratulations!</span>
            </AlertTitle>
            <AlertDescription>
              You passed the quiz! You now have a better understanding of this
              German story.
            </AlertDescription>
          </Alert>
        ) : (
          <Alert variant="warning" className="mb-6">
            <AlertTitle className="flex items-center justify-center">
              <span className="text-2xl mr-2">📖</span>
              <span className="font-semibold">Keep Learning!</span>
            </AlertTitle>
            <AlertDescription>
              Don't worry! Read the story again and try the quiz once more to
              improve your understanding.
            </AlertDescription>
          </Alert>
        )}

        <div className="flex space-x-4 justify-center">
          <Button
            onClick={onReturnToStory}
            variant="default"
            size="lg"
          >
            Read Story Again
          </Button>
          <Button
            onClick={onRetakeQuiz}
            variant="success"
            size="lg"
          >
            Retake Quiz
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

export { QuizQuestion, QuizResults };
