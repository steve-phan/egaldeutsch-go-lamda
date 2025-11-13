import React, { useState } from "react";
import axios from "axios";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Button } from "../ui/button";
import { Alert, AlertDescription } from "../ui/alert";
import { LoadingSpinner } from "../atoms/LoadingSpinner";

interface AIGenerationPanelProps {
  storyId: string;
  storyTitle: string;
  onGenerationComplete?: () => void;
}

type GenerationType = "questions" | "quiz" | "both";

interface GenerationResult {
  success: boolean;
  message: string;
  questionsCount?: number;
  quizId?: string;
  questionIds?: string[];
}

const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? "/.netlify/functions"
    : "http://localhost:8888/.netlify/functions";

export const AIGenerationPanel: React.FC<AIGenerationPanelProps> = ({
  storyId,
  storyTitle,
  onGenerationComplete,
}) => {
  const [generating, setGenerating] = useState(false);
  const [progress, setProgress] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [result, setResult] = useState<GenerationResult | null>(null);

  const handleGenerate = async (type: GenerationType) => {
    setGenerating(true);
    setProgress(`🤖 Generating ${type === "both" ? "questions and quiz" : type} with AI...`);
    setError("");
    setResult(null);

    try {
      const response = await axios.post<GenerationResult>(
        `${API_BASE_URL}/ai-generator?type=${type}&story_id=${storyId}`,
        {},
        {
          timeout: 90000, // 90 second timeout for AI generation
        }
      );

      setResult(response.data);
      setProgress("");
      
      if (onGenerationComplete) {
        onGenerationComplete();
      }
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || err.message || "Failed to generate content";
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
          🤖 AI Content Generation
        </CardTitle>
        <CardDescription>
          Generate questions and quizzes using OpenAI for "{storyTitle}"
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex flex-wrap gap-2 mb-4">
          <Button
            onClick={() => handleGenerate("questions")}
            disabled={generating}
            variant="outline"
            className="flex-1 min-w-[150px]"
          >
            {generating ? "..." : "Generate Questions"}
          </Button>
          <Button
            onClick={() => handleGenerate("quiz")}
            disabled={generating}
            variant="outline"
            className="flex-1 min-w-[150px]"
          >
            {generating ? "..." : "Generate Quiz"}
          </Button>
          <Button
            onClick={() => handleGenerate("both")}
            disabled={generating}
            className="flex-1 min-w-[150px]"
          >
            {generating ? "..." : "Generate Both"}
          </Button>
        </div>

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
                  Generated {result.questionsCount} questions. Content saved to draft status for review.
                </div>
              )}
            </AlertDescription>
          </Alert>
        )}

        <div className="mt-4 text-sm text-muted-foreground">
          <p className="mb-2">💡 <strong>How it works:</strong></p>
          <ul className="list-disc list-inside space-y-1 text-xs">
            <li>AI generates 8-12 diverse questions (comprehension, vocabulary, grammar)</li>
            <li>Content is saved in <strong>draft</strong> status for admin review</li>
            <li>Review and approve content before publishing</li>
            <li>Generation typically takes 30-60 seconds</li>
          </ul>
        </div>
      </CardContent>
    </Card>
  );
};
