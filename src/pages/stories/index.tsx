import React, { useState, useEffect } from "react";
import { Link } from "gatsby";
import Layout from "../../components/layout";
import StoryCard from "../../components/StoryCard";
import { Story } from "../../types";
import { fetchStories } from "../../utils/api";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { Badge } from "../../components/ui/badge";
import { Alert, AlertDescription } from "../../components/ui/alert";
import { useAuth } from "../../contexts/AuthContext";

const LEVELS = ['A1', 'A2', 'B1', 'B2', 'C1', 'C2'];

const StoriesPage: React.FC = () => {
  const { isAuthenticated, hasAnyRole } = useAuth();
  const [stories, setStories] = useState<Story[]>([]);
  const [filteredStories, setFilteredStories] = useState<Story[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  
  // Filter states
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [selectedLevel, setSelectedLevel] = useState<string>("");
  const [selectedTopic, setSelectedTopic] = useState<string>("");

  useEffect(() => {
    loadStories();
  }, []);

  useEffect(() => {
    applyFilters();
  }, [stories, searchQuery, selectedLevel, selectedTopic]);

  const loadStories = async () => {
    try {
      setLoading(true);
      setError(null);
      const fetchedStories = await fetchStories();
      setStories(fetchedStories);
    } catch (err) {
      console.error("Error fetching stories:", err);
      setError("Failed to load stories. Please try again later.");
    } finally {
      setLoading(false);
    }
  };

  const applyFilters = () => {
    let filtered = [...stories];

    // Apply search filter
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        story =>
          story.title.toLowerCase().includes(query) ||
          story.summary.toLowerCase().includes(query) ||
          story.topic.toLowerCase().includes(query)
      );
    }

    // Apply level filter
    if (selectedLevel) {
      filtered = filtered.filter(story => story.level === selectedLevel);
    }

    // Apply topic filter
    if (selectedTopic) {
      filtered = filtered.filter(story => story.topic === selectedTopic);
    }

    setFilteredStories(filtered);
  };

  const clearFilters = () => {
    setSearchQuery("");
    setSelectedLevel("");
    setSelectedTopic("");
  };

  const getUniqueTopics = (): string[] => {
    const topics = stories.map(story => story.topic);
    return Array.from(new Set(topics)).sort();
  };

  const canCreateContent = isAuthenticated && hasAnyRole(['creator', 'admin']);

  return (
    <Layout>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-2">
                German Stories
              </h1>
              <p className="text-xl text-muted-foreground">
                Learn German through engaging stories at your level
              </p>
            </div>
            {canCreateContent && (
              <Link to="/stories/create">
                <Button>Create Story</Button>
              </Link>
            )}
          </div>
        </div>

        {/* Filters */}
        <div className="bg-card rounded-lg shadow-sm border border-border p-6 mb-8">
          <h2 className="text-lg font-semibold mb-4">Filter Stories</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            {/* Search */}
            <div>
              <label className="text-sm font-medium mb-2 block">Search</label>
              <Input
                type="text"
                placeholder="Search stories..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>

            {/* Level Filter */}
            <div>
              <label className="text-sm font-medium mb-2 block">Level</label>
              <select
                value={selectedLevel}
                onChange={(e) => setSelectedLevel(e.target.value)}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              >
                <option value="">All Levels</option>
                {LEVELS.map(level => (
                  <option key={level} value={level}>{level}</option>
                ))}
              </select>
            </div>

            {/* Topic Filter */}
            <div>
              <label className="text-sm font-medium mb-2 block">Topic</label>
              <select
                value={selectedTopic}
                onChange={(e) => setSelectedTopic(e.target.value)}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              >
                <option value="">All Topics</option>
                {getUniqueTopics().map(topic => (
                  <option key={topic} value={topic}>{topic}</option>
                ))}
              </select>
            </div>
          </div>

          {/* Active Filters */}
          {(searchQuery || selectedLevel || selectedTopic) && (
            <div className="flex items-center gap-2 flex-wrap">
              <span className="text-sm text-muted-foreground">Active filters:</span>
              {searchQuery && (
                <Badge variant="secondary">
                  Search: {searchQuery}
                </Badge>
              )}
              {selectedLevel && (
                <Badge variant="secondary">
                  Level: {selectedLevel}
                </Badge>
              )}
              {selectedTopic && (
                <Badge variant="secondary">
                  Topic: {selectedTopic}
                </Badge>
              )}
              <Button 
                variant="ghost" 
                size="sm" 
                onClick={clearFilters}
              >
                Clear All
              </Button>
            </div>
          )}
        </div>

        {/* Loading State */}
        {loading && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
            <p className="mt-4 text-muted-foreground">Loading stories...</p>
          </div>
        )}

        {/* Error State */}
        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertDescription className="flex items-center justify-between">
              <span>{error}</span>
              <Button variant="outline" size="sm" onClick={loadStories}>
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
                    {filteredStories.length}
                  </div>
                  <div className="text-sm text-muted-foreground">
                    {filteredStories.length === stories.length ? 'Total Stories' : 'Filtered Stories'}
                  </div>
                </div>
                <div>
                  <div className="text-3xl font-bold text-success">
                    {new Set(filteredStories.map((s) => s.level)).size}
                  </div>
                  <div className="text-sm text-muted-foreground">
                    Difficulty Levels
                  </div>
                </div>
                <div>
                  <div className="text-3xl font-bold text-warning">
                    {new Set(filteredStories.map((s) => s.topic)).size}
                  </div>
                  <div className="text-sm text-muted-foreground">Topics</div>
                </div>
              </div>
            </div>

            {filteredStories.length === 0 ? (
              <div className="text-center py-12">
                <div className="max-w-md mx-auto">
                  <span className="text-6xl mb-4 block">📚</span>
                  <h2 className="text-2xl font-semibold text-foreground mb-4">
                    {stories.length === 0 ? 'No stories available yet' : 'No stories match your filters'}
                  </h2>
                  <p className="text-muted-foreground mb-6">
                    {stories.length === 0 
                      ? "We're working on adding amazing German learning stories. Check back soon!"
                      : "Try adjusting your filters to see more stories"}
                  </p>
                  {stories.length === 0 ? (
                    <Button onClick={loadStories}>
                      Refresh
                    </Button>
                  ) : (
                    <Button onClick={clearFilters}>
                      Clear Filters
                    </Button>
                  )}
                </div>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredStories.map((story) => (
                  <StoryCard key={story.id} story={story} />
                ))}
              </div>
            )}
          </>
        )}
      </div>
    </Layout>
  );
};

export default StoriesPage;
