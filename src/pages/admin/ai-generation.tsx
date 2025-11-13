import React, { useState, useEffect } from "react";
import Layout from "../../components/layout";
import ProtectedRoute from "../../components/auth/ProtectedRoute";
import { useAuth } from "../../contexts/AuthContext";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "../../components/ui/card";
import { Badge } from "../../components/ui/badge";
import { AIGenerationPanel } from "../../components/organisms/AIGenerationPanel";
import { LoadingSpinner } from "../../components/atoms/LoadingSpinner";
import { ErrorAlert } from "../../components/atoms/ErrorAlert";
import axios from "axios";

const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? "/.netlify/functions"
    : "http://localhost:8888/.netlify/functions";

interface Story {
  id: string;
  title: string;
  level: string;
  status: string;
  wordCount: number;
  topics: string[];
  isAIQuestionsGenerated?: boolean;
}

const AIGenerationPage: React.FC = () => {
  const { user } = useAuth();
  const [stories, setStories] = useState<Story[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selectedStory, setSelectedStory] = useState<Story | null>(null);

  useEffect(() => {
    loadPublishedStories();
  }, []);

  const loadPublishedStories = async () => {
    try {
      setLoading(true);
      setError("");

      const response = await axios.get(`${API_BASE_URL}/stories-management`);
      const allStories = response.data?.stories || [];

      // Filter only published stories that don't have AI-generated questions yet
      const availableStories = allStories.filter(
        (story: any) =>
          story.status === "published" && !story.isAIQuestionsGenerated
      );

      setStories(availableStories);
    } catch (err: any) {
      console.error("Error loading stories:", err);
      setError("Failed to load published stories");
    } finally {
      setLoading(false);
    }
  };

  const handleGenerationComplete = () => {
    // Refresh the stories list to remove the completed story
    setSelectedStory(null);
    loadPublishedStories();
  };

  return (
    <ProtectedRoute allowedRoles={["admin"]}>
      <Layout>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-foreground mb-2">
              🤖 AI Content Generation
            </h1>
            <p className="text-muted-foreground">
              Use AI to automatically generate questions and quizzes for
              published stories
            </p>
          </div>

          <ErrorAlert message={error} />

          {loading ? (
            <LoadingSpinner message="Loading published stories..." />
          ) : stories.length === 0 ? (
            <Card>
              <CardContent className="py-12 text-center">
                <p className="text-muted-foreground mb-2">
                  No stories available for AI generation.
                </p>
                <p className="text-sm text-muted-foreground">
                  All published stories either already have AI-generated
                  questions or none are published yet.
                </p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-6">
              {stories.map((story) => (
                <Card
                  key={story.id}
                  className="hover:shadow-md transition-shadow"
                >
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div>
                        <CardTitle className="text-xl mb-2">
                          {story.title}
                        </CardTitle>
                        <CardDescription className="flex gap-2 flex-wrap">
                          <Badge variant="secondary">{story.level}</Badge>
                          <Badge variant="outline">
                            {story.wordCount} words
                          </Badge>
                          {story.topics?.map((topic, idx) => (
                            <Badge
                              key={idx}
                              variant="outline"
                              className="text-xs"
                            >
                              {topic}
                            </Badge>
                          ))}
                        </CardDescription>
                      </div>
                      <Badge className="bg-green-100 text-green-800">
                        Published
                      </Badge>
                    </div>
                  </CardHeader>
                  <CardContent>
                    {selectedStory?.id === story.id ? (
                      <>
                        <AIGenerationPanel
                          storyId={story.id}
                          storyTitle={story.title}
                          onGenerationComplete={handleGenerationComplete}
                        />
                        <button
                          onClick={() => setSelectedStory(null)}
                          className="mt-4 text-sm text-muted-foreground hover:text-foreground"
                        >
                          ← Hide AI Panel
                        </button>
                      </>
                    ) : (
                      <button
                        onClick={() => setSelectedStory(story)}
                        className="text-blue-600 hover:text-blue-700 font-medium"
                      >
                        Generate AI Content →
                      </button>
                    )}
                  </CardContent>
                </Card>
              ))}
            </div>
          )}

          <Card className="mt-8 bg-blue-50 border-blue-200">
            <CardHeader>
              <CardTitle className="text-lg">
                💡 How AI Generation Works
              </CardTitle>
            </CardHeader>
            <CardContent className="text-sm space-y-2">
              <p>
                • <strong>Questions:</strong> AI generates 8-12 diverse
                questions covering comprehension, vocabulary, and grammar based
                on story content
              </p>
              <p>
                • <strong>Quiz:</strong> AI creates quiz metadata including
                title, description, and recommended structure
              </p>
              <p>
                • <strong>Both:</strong> Generates both questions and quiz in
                one operation
              </p>
              <p className="pt-2 border-t border-blue-200">
                ⚠️ All AI-generated content starts in <strong>draft</strong>{" "}
                status and requires admin review before publishing
              </p>
            </CardContent>
          </Card>
        </div>
      </Layout>
    </ProtectedRoute>
  );
};

export default AIGenerationPage;
