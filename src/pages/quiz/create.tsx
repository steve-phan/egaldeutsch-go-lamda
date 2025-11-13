import React, { useState, useEffect } from "react";
import { navigate } from "gatsby";
import Layout from "../../components/layout";
import ProtectedRoute from "../../components/auth/ProtectedRoute";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { Label } from "../../components/ui/label";
import { Textarea } from "../../components/ui/textarea";
import { Alert, AlertDescription } from "../../components/ui/alert";
import { Badge } from "../../components/ui/badge";
// import { Checkbox } from '../../components/ui/checkbox'; // Will implement checkbox inline
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
}

interface Question {
  id: string;
  question: string;
  questionType: string;
  difficulty: string;
  storyId: string;
  status: string;
}

const QUIZ_TYPES = [
  { value: "comprehension", label: "Reading Comprehension" },
  { value: "vocabulary", label: "Vocabulary Quiz" },
  { value: "mixed", label: "Mixed Questions" },
];

const CreateQuizPage: React.FC = () => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [stories, setStories] = useState<Story[]>([]);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [loadingStories, setLoadingStories] = useState(true);
  const [loadingQuestions, setLoadingQuestions] = useState(false);

  const [formData, setFormData] = useState({
    storyId: "",
    title: "",
    description: "",
    quizType: "mixed",
    totalQuestions: 5,
    passingScore: 70,
    questionIds: [] as string[],
  });

  useEffect(() => {
    loadStories();
  }, []);

  useEffect(() => {
    if (formData.storyId) {
      loadQuestions(formData.storyId);
    }
  }, [formData.storyId]);

  const loadStories = async () => {
    try {
      setLoadingStories(true);
      // Load published stories for quiz creation
      const response = await axios.get(`${API_BASE_URL}/stories-management`);
      if (response.data?.stories) {
        setStories(
          response.data.stories.filter((s: Story) => s.status === "published")
        );
      }
    } catch (err) {
      console.error("Error loading stories:", err);
      setError("Failed to load stories");
    } finally {
      setLoadingStories(false);
    }
  };

  const loadQuestions = async (storyId: string) => {
    try {
      setLoadingQuestions(true);
      const response = await axios.get(
        `${API_BASE_URL}/questions-management?storyId=${storyId}`
      );
      if (response.data?.questions) {
        // Only show published questions
        setQuestions(
          response.data.questions.filter(
            (q: Question) => q.status === "published"
          )
        );
      }
    } catch (err) {
      console.error("Error loading questions:", err);
      setError("Failed to load questions for this story");
    } finally {
      setLoadingQuestions(false);
    }
  };

  const handleInputChange = (field: string, value: any) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }));
    setError("");
  };

  const handleQuestionToggle = (questionId: string) => {
    setFormData((prev) => ({
      ...prev,
      questionIds: prev.questionIds.includes(questionId)
        ? prev.questionIds.filter((id) => id !== questionId)
        : [...prev.questionIds, questionId],
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (
      !formData.storyId ||
      !formData.title ||
      formData.questionIds.length === 0
    ) {
      setError(
        "Please fill in all required fields and select at least one question"
      );
      return;
    }

    if (formData.questionIds.length < formData.totalQuestions) {
      setError(
        `You need to select at least ${formData.totalQuestions} questions`
      );
      return;
    }

    try {
      setIsSubmitting(true);
      setError("");

      const quizData = {
        storyId: formData.storyId,
        title: formData.title,
        description: formData.description,
        quizType: formData.quizType,
        totalQuestions: formData.totalQuestions,
        passingScore: formData.passingScore,
        questionIds: formData.questionIds,
      };

      const response = await axios.post(
        `${API_BASE_URL}/quiz-management`,
        quizData
      );

      if (response.data) {
        setSuccess(true);
        setTimeout(() => {
          navigate("/admin/reviews");
        }, 2000);
      }
    } catch (err: any) {
      console.error("Error creating quiz:", err);
      setError(err.response?.data?.error || "Failed to create quiz");
    } finally {
      setIsSubmitting(false);
    }
  };

  const selectedStory = stories.find((s) => s.id === formData.storyId);

  if (success) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <Alert className="mb-6">
            <AlertDescription>
              Quiz created successfully! You'll be redirected to the admin
              panel.
            </AlertDescription>
          </Alert>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Create Quiz</h1>
          <p className="text-gray-600">
            Create a new quiz for published stories with existing questions.
          </p>
        </div>

        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Quiz Information</CardTitle>
              <CardDescription>
                Basic information about the quiz
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label htmlFor="storyId">Story *</Label>
                <select
                  id="storyId"
                  value={formData.storyId}
                  onChange={(e) => handleInputChange("storyId", e.target.value)}
                  className="w-full p-2 border border-gray-300 rounded-md"
                  disabled={loadingStories}
                >
                  <option value="">
                    {loadingStories ? "Loading stories..." : "Select a story"}
                  </option>
                  {stories.map((story) => (
                    <option key={story.id} value={story.id}>
                      {story.title} ({story.level})
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <Label htmlFor="title">Quiz Title *</Label>
                <Input
                  id="title"
                  value={formData.title}
                  onChange={(e) => handleInputChange("title", e.target.value)}
                  placeholder="Enter quiz title"
                  required
                />
              </div>

              <div>
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) =>
                    handleInputChange("description", e.target.value)
                  }
                  placeholder="Brief description of the quiz"
                  rows={3}
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <Label htmlFor="quizType">Quiz Type</Label>
                  <select
                    id="quizType"
                    value={formData.quizType}
                    onChange={(e) =>
                      handleInputChange("quizType", e.target.value)
                    }
                    className="w-full p-2 border border-gray-300 rounded-md"
                  >
                    {QUIZ_TYPES.map((type) => (
                      <option key={type.value} value={type.value}>
                        {type.label}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <Label htmlFor="totalQuestions">Total Questions</Label>
                  <Input
                    id="totalQuestions"
                    type="number"
                    min="1"
                    max="20"
                    value={formData.totalQuestions}
                    onChange={(e) =>
                      handleInputChange(
                        "totalQuestions",
                        parseInt(e.target.value)
                      )
                    }
                  />
                </div>

                <div>
                  <Label htmlFor="passingScore">Passing Score (%)</Label>
                  <Input
                    id="passingScore"
                    type="number"
                    min="0"
                    max="100"
                    value={formData.passingScore}
                    onChange={(e) =>
                      handleInputChange(
                        "passingScore",
                        parseInt(e.target.value)
                      )
                    }
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {formData.storyId && (
            <Card>
              <CardHeader>
                <CardTitle>Select Questions</CardTitle>
                <CardDescription>
                  Choose questions for this quiz from the selected story
                  {selectedStory && (
                    <span className="ml-2">
                      <Badge variant="outline">{selectedStory.title}</Badge>
                    </span>
                  )}
                </CardDescription>
              </CardHeader>
              <CardContent>
                {loadingQuestions ? (
                  <div className="text-center py-4">
                    <div className="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
                    <p className="mt-2 text-gray-600">Loading questions...</p>
                  </div>
                ) : questions.length === 0 ? (
                  <Alert>
                    <AlertDescription>
                      No published questions found for this story. You need to
                      create and publish questions first.
                    </AlertDescription>
                  </Alert>
                ) : (
                  <div className="space-y-3">
                    <p className="text-sm text-gray-600 mb-4">
                      Selected: {formData.questionIds.length} /{" "}
                      {formData.totalQuestions} required
                    </p>
                    {questions.map((question) => (
                      <div
                        key={question.id}
                        className="flex items-start space-x-3 p-3 border rounded-lg"
                      >
                        <input
                          type="checkbox"
                          id={question.id}
                          checked={formData.questionIds.includes(question.id)}
                          onChange={() => handleQuestionToggle(question.id)}
                          className="mt-1"
                        />
                        <div className="flex-1">
                          <label
                            htmlFor={question.id}
                            className="cursor-pointer"
                          >
                            <p className="font-medium">{question.question}</p>
                            <div className="flex space-x-2 mt-1">
                              <Badge variant="secondary" className="text-xs">
                                {question.questionType}
                              </Badge>
                              <Badge variant="outline" className="text-xs">
                                {question.difficulty}
                              </Badge>
                            </div>
                          </label>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          )}

          <div className="flex justify-end space-x-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => navigate("/admin")}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={
                isSubmitting ||
                formData.questionIds.length < formData.totalQuestions
              }
            >
              {isSubmitting ? "Creating..." : "Create Quiz"}
            </Button>
          </div>
        </form>
      </div>
    </Layout>
  );
};

const CreateQuizPageWrapper: React.FC = () => {
  return (
    <ProtectedRoute allowedRoles={["admin", "creator"]}>
      <CreateQuizPage />
    </ProtectedRoute>
  );
};

export default CreateQuizPageWrapper;
