# Learn German Platform

A serverless learning platform for German language learners. Read stories and test your understanding with interactive quizzes.

## Tech Stack

- **Frontend**: Gatsby.js with React and TailwindCSS
- **Backend**: Go Lambda Functions (Netlify Functions)
- **Database**: MongoDB
- **Deployment**: Netlify

## Features

- 📖 Read German stories at different difficulty levels
- ❓ Take quizzes to test comprehension
- 📊 View quiz results and scores
- 🎨 Modern, responsive UI with TailwindCSS
- 🔧 Health monitoring with `/api/health` endpoint
- ⚡ Optimized MongoDB connection pooling for serverless

## Serverless MongoDB Optimization

This project implements several optimizations for MongoDB in serverless environments:

- **Connection Reuse**: Uses singleton pattern with `sync.Once` to reuse connections across lambda invocations
- **Connection Pooling**: Configured with optimal pool settings for serverless:
  - MaxPoolSize: 10 (limit concurrent connections)
  - MinPoolSize: 0 (no minimum, serverless-friendly)
  - MaxConnIdleTime: 30s (close idle connections quickly)
- **Health Checks**: Automatic connection health monitoring and reconnection
- **Error Handling**: Graceful degradation and retry mechanisms

### Health Endpoint

The `/api/health` endpoint provides comprehensive health monitoring:

```bash
GET /api/health
```

**Response Example:**

```json
{
  "status": "healthy",
  "timestamp": "2025-11-12T10:30:00Z",
  "services": {
    "mongodb": "healthy"
  },
  "version": "1.0.0"
}
```

**MongoDB Status Values:**

- `healthy`: All systems operational
- `connection_failed`: Cannot establish connection
- `ping_failed`: Connection exists but ping failed
- `database_access_failed`: Cannot access database collections

## Project Structure

```
.
├── src/
│   ├── components/       # React components
│   ├── pages/           # Gatsby pages
│   └── styles/          # Global styles
├── netlify/
│   └── functions/       # Go Lambda functions
│       ├── stories/     # Story CRUD operations
│       ├── quiz/        # Quiz retrieval and submission
│       └── questions/   # Question management
├── models/              # Go data models
├── db/                  # MongoDB connection
└── netlify.toml         # Netlify configuration
```

## Setup

### Prerequisites

- Node.js 18+
- Go 1.21+
- MongoDB (local or cloud instance)

### Installation

1. Install Node.js dependencies:

```bash
npm install
```

2. Install Go dependencies:

```bash
go mod download
```

3. Set up environment variables:
   Create a `.env` file in the root directory:

```
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=egaldeutsch
GATSBY_API_URL=http://localhost:8888/.netlify/functions
```

For Netlify deployment, add these as environment variables in the Netlify dashboard.

### Development

1. Start MongoDB (if running locally)

2. Start Gatsby development server:

```bash
npm run develop
```

3. For local Lambda function testing, use Netlify Dev:

```bash
netlify dev
```

## Data Models

### Story

- `id`: ObjectID
- `title`: string
- `content`: string
- `level`: string (A1, A2, B1, B2, etc.)
- `wordCount`: number
- `createdAt`: timestamp
- `updatedAt`: timestamp

### Question

- `id`: ObjectID
- `storyId`: ObjectID
- `question`: string
- `options`: array of strings
- `correctAnswer`: number (index of correct option)
- `explanation`: string
- `createdAt`: timestamp

### Quiz Submission

- `id`: ObjectID
- `storyId`: ObjectID
- `answers`: array of numbers
- `score`: number
- `totalQuestions`: number
- `submittedAt`: timestamp

## API Endpoints

### Stories

- `GET /.netlify/functions/stories` - Get all stories
- `GET /.netlify/functions/stories/:id` - Get story by ID
- `POST /.netlify/functions/stories` - Create new story

### Quiz

- `GET /.netlify/functions/quiz/:storyId` - Get quiz for a story
- `POST /.netlify/functions/quiz` - Submit quiz answers

### Questions

- `GET /.netlify/functions/questions/:storyId` - Get questions for a story
- `POST /.netlify/functions/questions` - Create new question

## Deployment

1. Build the project:

```bash
npm run build
```

2. Deploy to Netlify:

- Connect your repository to Netlify
- Set environment variables in Netlify dashboard
- Netlify will automatically build and deploy

## Adding Sample Data

You can add sample stories and questions through the API endpoints or directly to MongoDB. Each story should have 10 questions for the quiz.

## License

MIT
