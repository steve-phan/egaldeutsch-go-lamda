import React, { useState, useEffect } from "react"
import { Link } from "gatsby"
import Layout from "../components/layout"
import axios from "axios"

const API_URL = process.env.GATSBY_API_URL || "http://localhost:8888/.netlify/functions"

const IndexPage = () => {
  const [stories, setStories] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    fetchStories()
  }, [])

  const fetchStories = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`${API_URL}/stories`)
      setStories(response.data)
      setError(null)
    } catch (err) {
      console.error("Error fetching stories:", err)
      setError("Failed to load stories. Please try again later.")
    } finally {
      setLoading(false)
    }
  }

  return (
    <Layout>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Learn German Through Stories
          </h1>
          <p className="text-xl text-gray-600">
            Read engaging stories and test your understanding with quizzes
          </p>
        </div>

        {loading && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            <p className="mt-4 text-gray-600">Loading stories...</p>
          </div>
        )}

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        {!loading && !error && (
          <>
            {stories.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-600 text-lg mb-4">
                  No stories available yet.
                </p>
                <p className="text-gray-500">
                  Check back later for new German learning stories!
                </p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {stories.map((story) => (
                  <Link
                    key={story.id}
                    to={`/story/${story.id}`}
                    className="block bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 overflow-hidden"
                  >
                    <div className="p-6">
                      <div className="flex items-center justify-between mb-3">
                        <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-primary-100 text-primary-800">
                          {story.level || "A1"}
                        </span>
                        <span className="text-sm text-gray-500">
                          {story.wordCount || 0} words
                        </span>
                      </div>
                      <h2 className="text-xl font-semibold text-gray-900 mb-2">
                        {story.title}
                      </h2>
                      <p className="text-gray-600 text-sm line-clamp-3">
                        {story.content?.substring(0, 150)}...
                      </p>
                      <div className="mt-4 flex items-center text-primary-600 font-medium">
                        Read Story →
                      </div>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </>
        )}
      </div>
    </Layout>
  )
}

export default IndexPage
