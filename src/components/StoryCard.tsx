import React from "react";
import { Link } from "gatsby";
import { Story } from "../types";
import { formatLevel, getLevelColor } from "../utils/api";
import { Card, CardContent, CardFooter, Badge, Button } from "./ui";

interface StoryCardProps {
  story: Story;
  className?: string;
}

const StoryCard: React.FC<StoryCardProps> = ({ story, className = "" }) => {
  return (
    <Card className={`flex flex-col ${className}`}>
      <CardContent className="pt-6 flex-1 flex flex-col">
        <div className="flex items-start justify-between mb-3">
          <h3 className="text-xl font-semibold text-foreground line-clamp-2 flex-1 mr-3">
            {story.title}
          </h3>
          <Badge className={getLevelColor(story.level)} variant="outline">
            {story.level}
          </Badge>
        </div>

        <p className="text-muted-foreground text-sm mb-4 line-clamp-3">
          {story.summary}
        </p>

        <div className="flex flex-wrap gap-2 text-sm text-muted-foreground mb-4">
          <span className="flex items-center bg-muted/50 px-2 py-1 rounded-md">
            📚 {story.wordCount} words
          </span>
          <span className="flex items-center bg-muted/50 px-2 py-1 rounded-md">
            ⏱️ {story.readingTime} min read
          </span>
        </div>

        <div className="flex flex-wrap gap-1 mb-4">
          {story.topics.map((topic, index) => (
            <Badge key={index} variant="secondary" className="text-xs">
              🎯 {topic}
            </Badge>
          ))}
        </div>

        {story.vocabulary && story.vocabulary.length > 0 && (
          <div className="mb-4 mt-auto">
            <p className="text-xs text-muted-foreground mb-2">
              Key vocabulary:
            </p>
            <div className="flex flex-wrap gap-1">
              {story.vocabulary.slice(0, 3).map((word, index) => (
                <Badge
                  key={index}
                  variant="secondary"
                  className="bg-primary/10 text-primary text-xs"
                >
                  {word.german}
                </Badge>
              ))}
              {story.vocabulary.length > 3 && (
                <Badge
                  variant="secondary"
                  className="bg-muted text-muted-foreground text-xs"
                >
                  +{story.vocabulary.length - 3} more
                </Badge>
              )}
            </div>
          </div>
        )}
      </CardContent>

      <CardFooter className="flex space-x-3 mt-auto">
        <Link to={`/story/${story.slug}`} className="flex-1">
          <Button variant="default" className="w-full">
            Read Story
          </Button>
        </Link>
        <Link to={`/quiz/${story.slug}`} className="flex-1">
          <Button variant="success" className="w-full">
            Take Quiz
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
};

export default StoryCard;
