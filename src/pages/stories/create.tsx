import React, { useState } from "react";
import { navigate } from "gatsby";
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
import { Label } from "../../components/ui/label";
import { Textarea } from "../../components/ui/textarea";
import { Alert, AlertDescription } from "../../components/ui/alert";
import { Badge } from "../../components/ui/badge";
import axios from "axios";

const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? "/.netlify/functions"
    : "http://localhost:8888/.netlify/functions";

const LEVELS = ["A1", "A2", "B1", "B2", "C1", "C2"];
const TOPICS = [
  "Family",
  "Travel",
  "Food & Cooking",
  "Work",
  "Shopping",
  "Health",
  "Education",
  "Hobbies",
  "Environment",
  "Technology",
  "Culture",
  "Daily Life",
];

interface VocabularyItem {
  german: string;
  english: string;
  wordType: string;
  article?: string;
}

const CreateStoryPage: React.FC = () => {
  const { user } = useAuth();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const [formData, setFormData] = useState({
    title: "",
    content: "",
    summary: "",
    level: "A1",
    topics: ["Daily Life"], // Changed to array
  });

  const [vocabularyItems, setVocabularyItems] = useState<VocabularyItem[]>([]);
  const [newVocabItem, setNewVocabItem] = useState<VocabularyItem>({
    german: "",
    english: "",
    wordType: "noun",
    article: "",
  });

  const handleInputChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    if (error) setError("");
  };

  const handleVocabInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setNewVocabItem((prev) => ({ ...prev, [name]: value }));
  };

  const addVocabularyItem = () => {
    if (!newVocabItem.german || !newVocabItem.english) {
      setError("Please fill in both German and English words");
      return;
    }

    setVocabularyItems((prev) => [...prev, { ...newVocabItem }]);
    setNewVocabItem({
      german: "",
      english: "",
      wordType: "noun",
      article: "",
    });
    setError("");
  };

  const removeVocabularyItem = (index: number) => {
    setVocabularyItems((prev) => prev.filter((_, i) => i !== index));
  };

  const addTopicToForm = (topic: string) => {
    if (
      topic &&
      !formData.topics.includes(topic) &&
      formData.topics.length < 5
    ) {
      setFormData((prev) => ({
        ...prev,
        topics: [...prev.topics, topic],
      }));
    }
  };

  const removeTopicFromForm = (topicToRemove: string) => {
    if (formData.topics.length > 1) {
      // Keep at least one topic
      setFormData((prev) => ({
        ...prev,
        topics: prev.topics.filter((topic) => topic !== topicToRemove),
      }));
    }
  };

  const calculateReadingTime = (text: string): number => {
    const wordsPerMinute = 200;
    const wordCount = text.trim().split(/\s+/).length;
    return Math.ceil(wordCount / wordsPerMinute);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    // Validation
    if (!formData.title || !formData.content || !formData.summary) {
      setError("Please fill in all required fields");
      return;
    }

    if (formData.content.length < 100) {
      setError("Story content must be at least 100 characters");
      return;
    }

    if (formData.topics.length === 0 || formData.topics.length > 5) {
      setError("Please select between 1 and 5 topics");
      return;
    }

    try {
      setIsSubmitting(true);

      const wordCount = formData.content.trim().split(/\s+/).length;
      const readingTime = calculateReadingTime(formData.content);

      const storyData = {
        title: formData.title,
        content: formData.content,
        summary: formData.summary,
        level: formData.level,
        topics: formData.topics,
        wordCount,
        readingTime,
        vocabulary: vocabularyItems,
        isActive: true,
      };

      await axios.post(`${API_BASE_URL}/stories-management`, storyData);

      setSuccess(true);
      setTimeout(() => {
        navigate("/stories");
      }, 2000);
    } catch (err: any) {
      console.error("Error creating story:", err);
      if (err.response) {
        setError(err.response.data?.error || "Failed to create story");
      } else {
        setError("Network error. Please check your connection.");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const wordCount = formData.content
    .trim()
    .split(/\s+/)
    .filter((w) => w).length;
  const readingTime = calculateReadingTime(formData.content);

  return (
    <Layout>
      <ProtectedRoute allowedRoles={["creator", "admin"]}>
        <div className="max-w-4xl mx-auto px-4 py-12">
          <div className="mb-6">
            <h1 className="text-3xl font-bold mb-2">Create New Story</h1>
            <p className="text-muted-foreground">
              Write an engaging German learning story for students
            </p>
          </div>

          {success ? (
            <Alert variant="success" className="mb-6">
              <AlertDescription>
                Story created successfully! Redirecting to stories page...
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
                  <CardTitle>Story Details</CardTitle>
                  <CardDescription>
                    Basic information about the story
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="title">Title *</Label>
                    <Input
                      id="title"
                      name="title"
                      placeholder="e.g., Ein Tag in Berlin"
                      value={formData.title}
                      onChange={handleInputChange}
                      disabled={isSubmitting}
                      required
                    />
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="level">Level *</Label>
                      <select
                        id="level"
                        name="level"
                        value={formData.level}
                        onChange={handleInputChange}
                        disabled={isSubmitting}
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                        required
                      >
                        {LEVELS.map((level) => (
                          <option key={level} value={level}>
                            {level}
                          </option>
                        ))}
                      </select>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="topics">Topics * (1-5 topics)</Label>
                      <div className="space-y-2">
                        <div className="flex flex-wrap gap-2">
                          {formData.topics.map((topic, index) => (
                            <Badge
                              key={index}
                              variant="secondary"
                              className="flex items-center gap-1 px-3 py-1"
                            >
                              {topic}
                              <button
                                type="button"
                                onClick={() => removeTopicFromForm(topic)}
                                className="ml-1 text-xs hover:text-red-500"
                                disabled={isSubmitting}
                              >
                                ×
                              </button>
                            </Badge>
                          ))}
                        </div>
                        <select
                          value=""
                          onChange={(e) => addTopicToForm(e.target.value)}
                          disabled={isSubmitting || formData.topics.length >= 5}
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                        >
                          <option value="">Add a topic...</option>
                          {TOPICS.filter(
                            (topic) => !formData.topics.includes(topic)
                          ).map((topic) => (
                            <option key={topic} value={topic}>
                              {topic}
                            </option>
                          ))}
                        </select>
                      </div>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="summary">Summary *</Label>
                    <Textarea
                      id="summary"
                      name="summary"
                      placeholder="Brief summary of the story (2-3 sentences)"
                      value={formData.summary}
                      onChange={handleInputChange}
                      disabled={isSubmitting}
                      rows={3}
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="content">Story Content *</Label>
                    <Textarea
                      id="content"
                      name="content"
                      placeholder="Write your German story here..."
                      value={formData.content}
                      onChange={handleInputChange}
                      disabled={isSubmitting}
                      rows={12}
                      required
                    />
                    <div className="flex gap-4 text-sm text-muted-foreground">
                      <span>{wordCount} words</span>
                      <span>~{readingTime} min read</span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card className="mb-6">
                <CardHeader>
                  <CardTitle>Vocabulary</CardTitle>
                  <CardDescription>
                    Add key vocabulary words from the story
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  {vocabularyItems.length > 0 && (
                    <div className="space-y-2 mb-4">
                      {vocabularyItems.map((item, index) => (
                        <div
                          key={index}
                          className="flex items-center justify-between p-3 bg-muted rounded-lg"
                        >
                          <div className="flex-1">
                            <span className="font-medium">
                              {item.article && `${item.article} `}
                              {item.german}
                            </span>
                            <span className="text-muted-foreground mx-2">
                              →
                            </span>
                            <span>{item.english}</span>
                            <Badge variant="outline" className="ml-2 text-xs">
                              {item.wordType}
                            </Badge>
                          </div>
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => removeVocabularyItem(index)}
                            disabled={isSubmitting}
                          >
                            Remove
                          </Button>
                        </div>
                      ))}
                    </div>
                  )}

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="german">German Word</Label>
                      <Input
                        id="german"
                        name="german"
                        placeholder="e.g., Haus"
                        value={newVocabItem.german}
                        onChange={handleVocabInputChange}
                        disabled={isSubmitting}
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="english">English Translation</Label>
                      <Input
                        id="english"
                        name="english"
                        placeholder="e.g., house"
                        value={newVocabItem.english}
                        onChange={handleVocabInputChange}
                        disabled={isSubmitting}
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="wordType">Word Type</Label>
                      <select
                        id="wordType"
                        name="wordType"
                        value={newVocabItem.wordType}
                        onChange={handleVocabInputChange}
                        disabled={isSubmitting}
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                      >
                        <option value="noun">Noun</option>
                        <option value="verb">Verb</option>
                        <option value="adjective">Adjective</option>
                        <option value="adverb">Adverb</option>
                        <option value="phrase">Phrase</option>
                      </select>
                    </div>

                    {newVocabItem.wordType === "noun" && (
                      <div className="space-y-2">
                        <Label htmlFor="article">Article (for nouns)</Label>
                        <select
                          id="article"
                          name="article"
                          value={newVocabItem.article}
                          onChange={handleVocabInputChange}
                          disabled={isSubmitting}
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                        >
                          <option value="">None</option>
                          <option value="der">der</option>
                          <option value="die">die</option>
                          <option value="das">das</option>
                        </select>
                      </div>
                    )}
                  </div>

                  <Button
                    type="button"
                    variant="outline"
                    onClick={addVocabularyItem}
                    disabled={isSubmitting}
                  >
                    Add Vocabulary Item
                  </Button>
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
                    "Create Story"
                  )}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigate("/stories")}
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

export default CreateStoryPage;
