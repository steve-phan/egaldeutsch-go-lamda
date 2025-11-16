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
import { protectedApi, publicApi } from "@/utils/apiClient";

interface Story {
  id: string;
  title: string;
  level: string;
}

const QUESTION_TYPES = [
  { value: "comprehension", label: "Comprehension" },
  { value: "vocabulary", label: "Vocabulary" },
  { value: "grammar", label: "Grammar" },
];

const DIFFICULTY_LEVELS = ["easy", "medium", "hard"];

const CreateQuestionPage: React.FC = () => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [stories, setStories] = useState<Story[]>([]);
  const [loadingStories, setLoadingStories] = useState(true);

  const [formData, setFormData] = useState({
    storyId: "",
    question: "",
    questionType: "comprehension",
    difficulty: "medium",
    options: ["", "", "", ""],
    correctAnswer: 0,
    explanation: "",
    points: 10,
    order: 1,
  });

  useEffect(() => {
    loadStories();
  }, []);

  const loadStories = async () => {
    try {
      setLoadingStories(true);
      const response = await publicApi.getStories();
      setStories(response.data.data || []);
    } catch (err) {
      console.error("Error loading stories:", err);
    } finally {
      setLoadingStories(false);
    }
  };

  const handleInputChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value, type } = e.target;

    // Convert numeric fields to numbers
    let processedValue: any = value;
    if (type === "number" && value !== "") {
      processedValue = parseInt(value, 10);
    }

    setFormData((prev) => ({ ...prev, [name]: processedValue }));
    if (error) setError("");
  };

  const handleOptionChange = (index: number, value: string) => {
    const newOptions = [...formData.options];
    newOptions[index] = value;
    setFormData((prev) => ({ ...prev, options: newOptions }));
  };

  const handleCorrectAnswerChange = (index: number) => {
    setFormData((prev) => ({ ...prev, correctAnswer: index }));
  };

  const validateForm = (): boolean => {
    if (!formData.storyId) {
      setError("Please select a story");
      return false;
    }

    if (!formData.question.trim()) {
      setError("Please enter a question");
      return false;
    }

    const filledOptions = formData.options.filter((opt) => opt.trim());
    if (filledOptions.length < 4) {
      setError("Please fill in all 4 options");
      return false;
    }

    if (!formData.options[formData.correctAnswer].trim()) {
      setError("The correct answer option cannot be empty");
      return false;
    }

    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!validateForm()) {
      return;
    }

    try {
      setIsSubmitting(true);

      const questionData = {
        storyId: formData.storyId,
        question: formData.question,
        questionType: formData.questionType,
        options: formData.options,
        correctAnswer: formData.correctAnswer,
        explanation: formData.explanation,
        points: formData.points,
        order: formData.order,
      };
      await protectedApi.createQuestion(questionData);

      setSuccess(true);
      setTimeout(() => {
        navigate("/questions");
      }, 2000);
    } catch (err: any) {
      console.error("Error creating question:", err);
      if (err.response) {
        setError(err.response.data?.error || "Failed to create question");
      } else {
        setError("Network error. Please check your connection.");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Layout>
      <ProtectedRoute allowedRoles={["creator", "admin"]}>
        <div className="max-w-4xl mx-auto px-4 py-12">
          <div className="mb-6">
            <h1 className="text-3xl font-bold mb-2">Create New Question</h1>
            <p className="text-muted-foreground">
              Add a quiz question for a German learning story
            </p>
          </div>

          {success ? (
            <Alert variant="success" className="mb-6">
              <AlertDescription>
                Question created successfully! Redirecting to questions page...
              </AlertDescription>
            </Alert>
          ) : (
            <form onSubmit={handleSubmit}>
              {error && (
                <Alert variant="destructive" className="mb-6">
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              )}

              <Card className="mb-6">
                <CardHeader>
                  <CardTitle>Question Details</CardTitle>
                  <CardDescription>
                    Basic information about the question
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="storyId">Story *</Label>
                    {loadingStories ? (
                      <p className="text-sm text-muted-foreground">
                        Loading stories...
                      </p>
                    ) : (
                      <select
                        id="storyId"
                        name="storyId"
                        value={formData.storyId}
                        onChange={handleInputChange}
                        disabled={isSubmitting}
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                        required
                      >
                        <option value="">Select a story</option>
                        {stories.map((story) => (
                          <option key={story.id} value={story.id}>
                            {story.title} ({story.level})
                          </option>
                        ))}
                      </select>
                    )}
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="questionType">Question Type *</Label>
                      <select
                        id="questionType"
                        name="questionType"
                        value={formData.questionType}
                        onChange={handleInputChange}
                        disabled={isSubmitting}
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                        required
                      >
                        {QUESTION_TYPES.map((type) => (
                          <option key={type.value} value={type.value}>
                            {type.label}
                          </option>
                        ))}
                      </select>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="difficulty">Difficulty</Label>
                      <select
                        id="difficulty"
                        name="difficulty"
                        value={formData.difficulty}
                        onChange={handleInputChange}
                        disabled={isSubmitting}
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                      >
                        {DIFFICULTY_LEVELS.map((level) => (
                          <option key={level} value={level}>
                            {level.charAt(0).toUpperCase() + level.slice(1)}
                          </option>
                        ))}
                      </select>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="question">Question Text *</Label>
                    <Textarea
                      id="question"
                      name="question"
                      placeholder="Enter your question here..."
                      value={formData.question}
                      onChange={handleInputChange}
                      disabled={isSubmitting}
                      rows={3}
                      required
                    />
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="points">Points</Label>
                      <Input
                        id="points"
                        name="points"
                        type="number"
                        min="1"
                        max="100"
                        value={formData.points}
                        onChange={handleInputChange}
                        disabled={isSubmitting}
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="order">Question Order</Label>
                      <Input
                        id="order"
                        name="order"
                        type="number"
                        min="1"
                        value={formData.order}
                        onChange={handleInputChange}
                        disabled={isSubmitting}
                        required
                      />
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card className="mb-6">
                <CardHeader>
                  <CardTitle>Answer Options</CardTitle>
                  <CardDescription>
                    Provide 4 options and select the correct answer
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  {formData.options.map((option, index) => (
                    <div key={index} className="space-y-2">
                      <div className="flex items-center gap-2">
                        <input
                          type="radio"
                          id={`correct-${index}`}
                          name="correctAnswer"
                          checked={formData.correctAnswer === index}
                          onChange={() => handleCorrectAnswerChange(index)}
                          disabled={isSubmitting}
                          className="h-4 w-4"
                        />
                        <Label htmlFor={`option-${index}`} className="flex-1">
                          Option {String.fromCharCode(65 + index)} *
                        </Label>
                        {formData.correctAnswer === index && (
                          <Badge variant="success">Correct Answer</Badge>
                        )}
                      </div>
                      <Input
                        id={`option-${index}`}
                        placeholder={`Enter option ${String.fromCharCode(
                          65 + index
                        )}`}
                        value={option}
                        onChange={(e) =>
                          handleOptionChange(index, e.target.value)
                        }
                        disabled={isSubmitting}
                        required
                      />
                    </div>
                  ))}
                </CardContent>
              </Card>

              <Card className="mb-6">
                <CardHeader>
                  <CardTitle>Explanation</CardTitle>
                  <CardDescription>
                    Explain why the correct answer is right (optional but
                    recommended)
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <Textarea
                    id="explanation"
                    name="explanation"
                    placeholder="Provide an explanation for the correct answer..."
                    value={formData.explanation}
                    onChange={handleInputChange}
                    disabled={isSubmitting}
                    rows={4}
                  />
                </CardContent>
              </Card>

              <div className="flex gap-4">
                <Button
                  type="submit"
                  disabled={isSubmitting}
                  className="flex-1"
                >
                  {isSubmitting ? (
                    <>
                      <span className="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></span>
                      Creating...
                    </>
                  ) : (
                    "Create Question"
                  )}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigate("/questions")}
                  disabled={isSubmitting}
                >
                  Cancel
                </Button>
              </div>
            </form>
          )}
        </div>
      </ProtectedRoute>
    </Layout>
  );
};

export default CreateQuestionPage;
