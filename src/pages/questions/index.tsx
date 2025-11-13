import React, { useState, useEffect } from "react";
import { Link, navigate } from "gatsby";
import Layout from "../../components/layout";
import ProtectedRoute from "../../components/auth/ProtectedRoute";
import { useAuth } from "../../contexts/AuthContext";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { Badge } from "../../components/ui/badge";
import { Alert, AlertDescription } from "../../components/ui/alert";
import axios from "axios";

const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? "/.netlify/functions"
    : "http://localhost:8888/.netlify/functions";

interface Question {
  id: string;
  storyId: string;
  storyTitle?: string;
  question: string;
  questionType: string;
  options: string[];
  correctAnswer: number;
  explanation: string;
  points: number;
  order: number;
  createdAt: string;
}

const QUESTION_TYPES = [
  { value: "comprehension", label: "Comprehension" },
  { value: "vocabulary", label: "Vocabulary" },
  { value: "grammar", label: "Grammar" },
];

const QuestionsPage: React.FC = () => {
  const { hasAnyRole } = useAuth();
  const [questions, setQuestions] = useState<Question[]>([]);
  const [filteredQuestions, setFilteredQuestions] = useState<Question[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Filters
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedType, setSelectedType] = useState("");
  const [selectedStory, setSelectedStory] = useState("");

  useEffect(() => {
    loadQuestions();
  }, []);

  useEffect(() => {
    applyFilters();
  }, [questions, searchQuery, selectedType, selectedStory]);

  const loadQuestions = async () => {
    try {
      setLoading(true);
      setError("");
      const response = await axios.get(`${API_BASE_URL}/questions-management`);
      setQuestions(response.data?.questions || []);
    } catch (err: any) {
      console.error("Error fetching questions:", err);
      setError("Failed to load questions. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const applyFilters = () => {
    if (!Array.isArray(questions)) {
      setFilteredQuestions([]);
      return;
    }

    let filtered = [...questions];

    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (q) =>
          q.question.toLowerCase().includes(query) ||
          q.storyTitle?.toLowerCase().includes(query)
      );
    }

    if (selectedType) {
      filtered = filtered.filter((q) => q.questionType === selectedType);
    }

    if (selectedStory) {
      filtered = filtered.filter((q) => q.storyId === selectedStory);
    }

    setFilteredQuestions(filtered);
  };

  const clearFilters = () => {
    setSearchQuery("");
    setSelectedType("");
    setSelectedStory("");
  };

  const getUniqueStories = () => {
    if (!Array.isArray(questions)) {
      return [];
    }

    const stories = questions
      .filter((q): q is Question & { storyTitle: string } =>
        Boolean(q.storyTitle)
      )
      .map((q) => ({ id: q.storyId, title: q.storyTitle }));

    const uniqueMap = new Map(stories.map((s) => [s.id, s]));
    return Array.from(uniqueMap.values());
  };

  const getTypeBadgeVariant = (type: string) => {
    switch (type) {
      case "comprehension":
        return "default";
      case "vocabulary":
        return "secondary";
      case "grammar":
        return "outline";
      default:
        return "default";
    }
  };

  const canCreateContent = hasAnyRole(["creator", "admin"]);

  return (
    <Layout>
      <ProtectedRoute allowedRoles={["creator", "reviewer", "admin"]}>
        <div className="max-w-7xl mx-auto px-4 py-12">
          {/* Header */}
          <div className="mb-8">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h1 className="text-4xl font-bold mb-2">Questions</h1>
                <p className="text-muted-foreground">
                  Manage quiz questions for German learning stories
                </p>
              </div>
              {canCreateContent && (
                <Link to="/questions/create">
                  <Button>Create Question</Button>
                </Link>
              )}
            </div>
          </div>

          {/* Filters */}
          <Card className="mb-8">
            <CardContent className="pt-6">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                <div>
                  <label className="text-sm font-medium mb-2 block">
                    Search
                  </label>
                  <Input
                    type="text"
                    placeholder="Search questions..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>

                <div>
                  <label className="text-sm font-medium mb-2 block">Type</label>
                  <select
                    value={selectedType}
                    onChange={(e) => setSelectedType(e.target.value)}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  >
                    <option value="">All Types</option>
                    {QUESTION_TYPES.map((type) => (
                      <option key={type.value} value={type.value}>
                        {type.label}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="text-sm font-medium mb-2 block">
                    Story
                  </label>
                  <select
                    value={selectedStory}
                    onChange={(e) => setSelectedStory(e.target.value)}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  >
                    <option value="">All Stories</option>
                    {getUniqueStories().map((story) => (
                      <option key={story.id} value={story.id}>
                        {story.title}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              {(searchQuery || selectedType || selectedStory) && (
                <div className="flex items-center gap-2 flex-wrap">
                  <span className="text-sm text-muted-foreground">
                    Active filters:
                  </span>
                  {searchQuery && (
                    <Badge variant="secondary">Search: {searchQuery}</Badge>
                  )}
                  {selectedType && (
                    <Badge variant="secondary">Type: {selectedType}</Badge>
                  )}
                  {selectedStory && (
                    <Badge variant="secondary">Story selected</Badge>
                  )}
                  <Button variant="ghost" size="sm" onClick={clearFilters}>
                    Clear All
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Loading State */}
          {loading && (
            <div className="text-center py-12">
              <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
              <p className="mt-4 text-muted-foreground">Loading questions...</p>
            </div>
          )}

          {/* Error State */}
          {error && (
            <Alert variant="destructive" className="mb-6">
              <AlertDescription className="flex items-center justify-between">
                <span>{error}</span>
                <Button variant="outline" size="sm" onClick={loadQuestions}>
                  Try Again
                </Button>
              </AlertDescription>
            </Alert>
          )}

          {/* Content */}
          {!loading && !error && (
            <>
              {/* Stats */}
              <div className="bg-muted/50 rounded-lg p-6 mb-8">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-center">
                  <div>
                    <div className="text-3xl font-bold text-primary">
                      {filteredQuestions.length}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {filteredQuestions.length === questions.length
                        ? "Total Questions"
                        : "Filtered Questions"}
                    </div>
                  </div>
                  <div>
                    <div className="text-3xl font-bold text-success">
                      {
                        new Set(filteredQuestions.map((q) => q.questionType))
                          .size
                      }
                    </div>
                    <div className="text-sm text-muted-foreground">
                      Question Types
                    </div>
                  </div>
                  <div>
                    <div className="text-3xl font-bold text-warning">
                      {new Set(filteredQuestions.map((q) => q.storyId)).size}
                    </div>
                    <div className="text-sm text-muted-foreground">Stories</div>
                  </div>
                </div>
              </div>

              {filteredQuestions.length === 0 ? (
                <div className="text-center py-12">
                  <span className="text-6xl mb-4 block">❓</span>
                  <h2 className="text-2xl font-semibold mb-4">
                    {questions.length === 0
                      ? "No questions yet"
                      : "No questions match your filters"}
                  </h2>
                  <p className="text-muted-foreground mb-6">
                    {questions.length === 0
                      ? "Create your first question to get started"
                      : "Try adjusting your filters to see more questions"}
                  </p>
                  {canCreateContent && questions.length === 0 && (
                    <Button onClick={() => navigate("/questions/create")}>
                      Create First Question
                    </Button>
                  )}
                  {questions.length > 0 && (
                    <Button onClick={clearFilters}>Clear Filters</Button>
                  )}
                </div>
              ) : (
                <div className="space-y-4">
                  {filteredQuestions.map((question) => (
                    <Card key={question.id}>
                      <CardHeader>
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <div className="flex items-center gap-2 mb-2">
                              <Badge
                                variant={getTypeBadgeVariant(
                                  question.questionType
                                )}
                              >
                                {question.questionType}
                              </Badge>
                              <span className="text-sm text-muted-foreground">
                                Question #{question.order}
                              </span>
                            </div>
                            <CardTitle className="mb-2">
                              {question.question}
                            </CardTitle>
                            {question.storyTitle && (
                              <CardDescription>
                                Story: {question.storyTitle}
                              </CardDescription>
                            )}
                          </div>
                          {canCreateContent && (
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() =>
                                navigate(`/questions/edit/${question.id}`)
                              }
                            >
                              Edit
                            </Button>
                          )}
                        </div>
                      </CardHeader>
                      <CardContent>
                        <div className="space-y-2">
                          <p className="text-sm font-medium">Options:</p>
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                            {question.options.map((option, index) => (
                              <div
                                key={index}
                                className={`p-2 rounded border ${
                                  index === question.correctAnswer
                                    ? "border-success bg-success/10"
                                    : "border-border"
                                }`}
                              >
                                <span className="font-medium mr-2">
                                  {String.fromCharCode(65 + index)}.
                                </span>
                                {option}
                                {index === question.correctAnswer && (
                                  <Badge
                                    variant="success"
                                    className="ml-2 text-xs"
                                  >
                                    Correct
                                  </Badge>
                                )}
                              </div>
                            ))}
                          </div>
                          {question.explanation && (
                            <div className="mt-4 p-3 bg-muted rounded">
                              <p className="text-sm font-medium mb-1">
                                Explanation:
                              </p>
                              <p className="text-sm text-muted-foreground">
                                {question.explanation}
                              </p>
                            </div>
                          )}
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </>
          )}
        </div>
      </ProtectedRoute>
    </Layout>
  );
};

export default QuestionsPage;
