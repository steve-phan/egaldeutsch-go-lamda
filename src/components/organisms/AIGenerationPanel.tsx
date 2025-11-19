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

      if (onGenerationComplete) {
        onGenerationComplete();
      }
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
        <Button
          onClick={handleGenerate}
          disabled={generating}
          className="w-full"
          size="lg"
        >
          {generating ? "Generating Questions..." : "🤖 Generate Questions"}
        </Button>

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
          <Alert className="bg-green-50 border-green-200">
            <AlertDescription className="text-green-800">
              ✅ {result.message}
              {result.questionsCount && (
                <div className="mt-2 text-sm">
                  Generated {result.questionsCount} questions. Content saved to
                  draft status for review.
                </div>
              )}
            </AlertDescription>
          </Alert>
        )}

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
      </CardContent>
    </Card>
  );
};
