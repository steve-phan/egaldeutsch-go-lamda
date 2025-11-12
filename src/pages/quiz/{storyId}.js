import React, { useState, useEffect } from "react"
import { navigate } from "gatsby"
import Layout from "../../components/layout"
import axios from "axios"

const API_URL = process.env.GATSBY_API_URL || "http://localhost:8888/.netlify/functions"

const QuizPage = ({ params }) => {
  const storyId = params.storyId
  const [quiz, setQuiz] = useState(null)
  const [answers, setAnswers] = useState([])
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [result, setResult] = useState(null)
  const [error, setError] = useState(null)

  useEffect(() => {
    fetchQuiz()
  }, [storyId])

  const fetchQuiz = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`${API_URL}/quiz/${storyId}`)
      setQuiz(response.data)
      setAnswers(new Array(response.data.questions.length).fill(null))
      setError(null)
    } catch (err) {
      console.error("Error fetching quiz:", err)
      setError("Failed to load quiz. Please try again later.")
    } finally {
      setLoading(false)
    }
  }

  const handleAnswerSelect = (questionIndex, answerIndex) => {
    const newAnswers = [...answers]
    newAnswers[questionIndex] = answerIndex
    setAnswers(newAnswers)
  }

  const handleSubmit = async () => {
    // Check if all questions are answered
    if (answers.some(answer => answer === null)) {
      alert("Please answer all questions before submitting.")
      return
    }

    try {
      setSubmitting(true)
      const response = await axios.post(`${API_URL}/quiz`, {
        storyId: storyId,
        answers: answers,
      })
      setResult(response.data)
    } catch (err) {
      console.error("Error submitting quiz:", err)
      setError("Failed to submit quiz. Please try again.")
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="text-center">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            <p className="mt-4 text-gray-600">Loading quiz...</p>
          </div>
        </div>
      </Layout>
    )
  }

  if (error || !quiz) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            {error || "Quiz not found"}
          </div>
        </div>
      </Layout>
    )
  }

  if (result) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="bg-white rounded-lg shadow-lg p-8 text-center">
            <div className="mb-6">
              <div className="inline-flex items-center justify-center w-24 h-24 rounded-full bg-primary-100 mb-4">
                <span className="text-4xl font-bold text-primary-600">
                  {result.score}/{result.totalQuestions}
                </span>
              </div>
              <h2 className="text-2xl font-bold text-gray-900 mb-2">
                Quiz Complete!
              </h2>
              <p className="text-xl text-gray-600">
                You scored {result.percentage.toFixed(1)}%
              </p>
            </div>

            <div className="space-y-4">
              <button
                onClick={() => navigate(`/story/${storyId}`)}
                className="bg-gray-200 text-gray-800 px-6 py-3 rounded-lg font-semibold hover:bg-gray-300 transition-colors duration-200 mr-4"
              >
                Read Story Again
              </button>
              <button
                onClick={() => navigate("/")}
                className="bg-primary-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors duration-200"
              >
                Browse More Stories
              </button>
            </div>
          </div>
        </div>
      </Layout>
    )
  }

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <button
          onClick={() => navigate(`/story/${storyId}`)}
          className="mb-6 text-primary-600 hover:text-primary-700 font-medium"
        >
          ← Back to Story
        </button>

        <div className="bg-white rounded-lg shadow-lg p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Quiz: {quiz.story?.title}
          </h1>
          <p className="text-gray-600 mb-8">
            Answer all {quiz.questions.length} questions to test your understanding
          </p>

          <div className="space-y-8">
            {quiz.questions.map((question, questionIndex) => (
              <div key={question.id} className="border-b pb-6 last:border-b-0">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">
                  {questionIndex + 1}. {question.question}
                </h3>
                <div className="space-y-2">
                  {question.options.map((option, optionIndex) => (
                    <label
                      key={optionIndex}
                      className={`flex items-center p-4 rounded-lg border-2 cursor-pointer transition-colors duration-200 ${
                        answers[questionIndex] === optionIndex
                          ? "border-primary-500 bg-primary-50"
                          : "border-gray-200 hover:border-gray-300"
                      }`}
                    >
                      <input
                        type="radio"
                        name={`question-${questionIndex}`}
                        value={optionIndex}
                        checked={answers[questionIndex] === optionIndex}
                        onChange={() => handleAnswerSelect(questionIndex, optionIndex)}
                        className="mr-3 h-4 w-4 text-primary-600 focus:ring-primary-500"
                      />
                      <span className="text-gray-700">{option}</span>
                    </label>
                  ))}
                </div>
              </div>
            ))}
          </div>

          <div className="mt-8 pt-6 border-t">
            <button
              onClick={handleSubmit}
              disabled={submitting || answers.some(answer => answer === null)}
              className="w-full bg-primary-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors duration-200 shadow-md hover:shadow-lg disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              {submitting ? "Submitting..." : "Submit Quiz"}
            </button>
          </div>
        </div>
      </div>
    </Layout>
  )
}

export default QuizPage
