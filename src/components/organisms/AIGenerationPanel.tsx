import React, { useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { Button } from "../ui/button";
import { Alert, AlertDescription } from "../ui/alert";
import { LoadingSpinner } from "../atoms/LoadingSpinner";
import { protectedApi } from "@/utils/apiClient";

interface AIGenerationPanelProps {
  storyId: string;
  storyTitle: string;
  onGenerationComplete?: () => void;
}

type GenerationType = "questions";

interface GenerationResult {
  success: boolean;
  message: string;
  questionsCount?: number;
  quizId?: string;
  questionIds?: string[];
}

export const AIGenerationPanel: React.FC<AIGenerationPanelProps> = ({
  storyId,
  storyTitle,
  onGenerationComplete,
}) => {
  const [generating, setGenerating] = useState(false);
  const [progress, setProgress] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [result, setResult] = useState<GenerationResult | null>(null);

  const handleGenerate = async () => {
    setGenerating(true);
    setProgress(`🤖 Generating questions with AI...`);
    setError("");
    setResult(null);

    try {
      const response = await protectedApi.generateQuestions(storyId);

      setResult(response.data);
      setProgress("");

      // Auto-hide success message after 5 seconds
      setTimeout(() => {
        if (onGenerationComplete) {
          onGenerationComplete();
        }
      }, 5000);
    } catch (err: any) {
      const errorMsg =
        err.response?.data?.error ||
        err.message ||
        "Failed to generate content";
      setError(errorMsg);
      setProgress("");
    } finally {
      setGenerating(false);
    }
  };

  const handleClose = () => {
    if (onGenerationComplete) {
      onGenerationComplete();
    }
  };

  return (
    <Card className="mt-4 border-blue-200">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          🤖 AI Question Generation
        </CardTitle>
        <CardDescription>
          Generate questions using OpenAI for "{storyTitle}"
        </CardDescription>
      </CardHeader>
      <CardContent>
        {!result && !generating && (
          <Button
            onClick={handleGenerate}
            disabled={generating}
            className="w-full"
            size="lg"
          >
            {generating ? "Generating Questions..." : "🤖 Generate Questions"}
          </Button>
        )}

        {progress && (
          <Alert className="mb-4 bg-blue-50 border-blue-200">
            <AlertDescription className="flex items-center gap-2">
              <LoadingSpinner />
              <span>{progress}</span>
            </AlertDescription>
          </Alert>
        )}

        {error && (
          <Alert className="mb-4 bg-red-50 border-red-200">
            <AlertDescription className="text-red-800">
              ⚠️ {error}
            </AlertDescription>
          </Alert>
        )}

        {result && result.success && (
          <div className="space-y-4">
            <Alert className="bg-green-50 border-green-200">
              <AlertDescription>
                <div className="flex items-start gap-3">
                  {/* Animated checkmark */}
                  <div className="flex-shrink-0">
                    <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center animate-bounce">
                      <svg
                        className="w-6 h-6 text-green-600"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M5 13l4 4L19 7"
                        />
                      </svg>
                    </div>
                  </div>
                  <div className="flex-1">
                    <h4 className="font-semibold text-green-800 mb-1">
                      Success! Questions Generated
                    </h4>
                    <p className="text-green-700 text-sm">{result.message}</p>
                    {result.questionsCount && (
                      <div className="mt-2 flex items-center gap-2">
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                          {result.questionsCount} Questions
                        </span>
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                          Draft Status
                        </span>
                      </div>
                    )}
                  </div>
                </div>
              </AlertDescription>
            </Alert>

            {/* Action buttons */}
            <div className="flex gap-2">
              <Button
                onClick={() => window.location.href = "/admin/reviews"}
                variant="default"
                className="flex-1"
              >
                Review Questions →
              </Button>
              <Button
                onClick={handleClose}
                variant="outline"
                className="flex-1"
              >
                Generate Another
              </Button>
            </div>
          </div>
        )}

        {!result && !generating && (
          <div className="mt-4 text-sm text-muted-foreground">
            <p className="mb-2">
              💡 <strong>How it works:</strong>
            </p>
            <ul className="list-disc list-inside space-y-1 text-xs">
              <li>
                AI generates 8-12 diverse questions (comprehension, vocabulary,
                grammar)
              </li>
              <li>
                Content is saved in <strong>draft</strong> status for admin review
              </li>
              <li>Review and approve content before publishing</li>
              <li>Generation typically takes 30-60 seconds</li>
            </ul>
          </div>
        )}
      </CardContent>
    </Card>
  );
};
