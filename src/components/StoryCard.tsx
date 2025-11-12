import React from "react";
import { Link } from "gatsby";
import { Story } from "@/types";
import { formatLevel, getLevelColor } from "@/utils/api";

interface StoryCardProps {
  story: Story;
  className?: string;
}

const StoryCard: React.FC<StoryCardProps> = ({ story, className = "" }) => {
  return (
    <div className={`bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300 overflow-hidden ${className}`}>
      <div className="p-6">
        <div className="flex items-start justify-between mb-3">
          <h3 className="text-xl font-semibold text-gray-900 line-clamp-2 flex-1 mr-3">
            {story.title}
          </h3>
          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getLevelColor(story.level)} flex-shrink-0`}>
            {story.level}
          </span>
        </div>
        
        <p className="text-gray-600 text-sm mb-4 line-clamp-3">
          {story.summary}
        </p>
        
        <div className="flex items-center justify-between text-sm text-gray-500 mb-4">
          <div className="flex items-center space-x-4">
            <span className="flex items-center">
              📚 {story.wordCount} words
            </span>
            <span className="flex items-center">
              ⏱️ {story.readingTime} min read
            </span>
            <span className="flex items-center">
              🎯 {story.topic}
            </span>
          </div>
        </div>
        
        {story.vocabulary && story.vocabulary.length > 0 && (
          <div className="mb-4">
            <p className="text-xs text-gray-500 mb-2">Key vocabulary:</p>
            <div className="flex flex-wrap gap-1">
              {story.vocabulary.slice(0, 3).map((word, index) => (
                <span 
                  key={index}
                  className="inline-flex items-center px-2 py-1 rounded text-xs bg-blue-50 text-blue-700"
                >
                  {word.german}
                </span>
              ))}
              {story.vocabulary.length > 3 && (
                <span className="inline-flex items-center px-2 py-1 rounded text-xs bg-gray-50 text-gray-500">
                  +{story.vocabulary.length - 3} more
                </span>
              )}
            </div>
          </div>
        )}
        
        <div className="flex space-x-3">
          <Link
            to={`/story/${story.id}`}
            className="flex-1 bg-blue-600 hover:bg-blue-700 text-white text-center py-2 px-4 rounded-md text-sm font-medium transition-colors"
          >
            Read Story
          </Link>
          <Link
            to={`/quiz/${story.id}`}
            className="flex-1 bg-green-600 hover:bg-green-700 text-white text-center py-2 px-4 rounded-md text-sm font-medium transition-colors"
          >
            Take Quiz
          </Link>
        </div>
      </div>
    </div>
  );
};

export default StoryCard;