import React, { useState, useEffect } from "react"
import { navigate } from "gatsby"
import Layout from "../../components/layout"
import axios from "axios"

const API_URL = process.env.GATSBY_API_URL || "http://localhost:8888/.netlify/functions"

const StoryPage = ({ params }) => {
  const storyId = params.storyId
  const [story, setStory] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [showQuiz, setShowQuiz] = useState(false)

  useEffect(() => {
    fetchStory()
  }, [storyId])

  const fetchStory = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`${API_URL}/stories/${storyId}`)
      setStory(response.data)
      setError(null)
    } catch (err) {
      console.error("Error fetching story:", err)
      setError("Failed to load story. Please try again later.")
    } finally {
      setLoading(false)
    }
  }

  const handleStartQuiz = () => {
    setShowQuiz(true)
    navigate(`/quiz/${storyId}`)
  }

  if (loading) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="text-center">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            <p className="mt-4 text-gray-600">Loading story...</p>
          </div>
        </div>
      </Layout>
    )
  }

  if (error || !story) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            {error || "Story not found"}
          </div>
        </div>
      </Layout>
    )
  }

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <button
          onClick={() => navigate("/")}
          className="mb-6 text-primary-600 hover:text-primary-700 font-medium"
        >
          ← Back to Stories
        </button>

        <div className="bg-white rounded-lg shadow-lg p-8">
          <div className="flex items-center justify-between mb-6">
            <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-primary-100 text-primary-800">
              {story.level || "A1"}
            </span>
            <span className="text-sm text-gray-500">
              {story.wordCount || 0} words
            </span>
          </div>

          <h1 className="text-3xl font-bold text-gray-900 mb-6">
            {story.title}
          </h1>

          <div className="prose prose-lg max-w-none mb-8">
            <div className="text-gray-700 leading-relaxed whitespace-pre-line">
              {story.content}
            </div>
          </div>

          <div className="border-t pt-6 mt-8">
            <button
              onClick={handleStartQuiz}
              className="w-full bg-primary-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors duration-200 shadow-md hover:shadow-lg"
            >
              Start Quiz →
            </button>
          </div>
        </div>
      </div>
    </Layout>
  )
}

export default StoryPage
