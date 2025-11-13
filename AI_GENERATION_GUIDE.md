# AI-Powered Content Generation

## Overview

The AI Content Generation feature uses OpenAI's GPT-4o model to automatically generate high-quality German learning questions and quizzes based on published story content.

## Features

### Question Generation
- Generates 8-12 diverse questions per story
- **Question Types:**
  - 40% Comprehension (main ideas, details, inference)
  - 35% Vocabulary (definitions, usage, context)
  - 25% Grammar (structures from the text)
- **Difficulty Distribution:**
  - 30% Easy questions
  - 50% Medium questions
  - 20% Hard questions
- Points awarded: 5-10 based on difficulty
- All questions include explanations for learning

### Quiz Generation
- Creates quiz metadata including:
  - Engaging title related to story content
  - Clear description of assessment goals
  - Estimated completion time (5-30 minutes)
  - Recommended passing score (50-90%)
  - Difficulty distribution
- Associates with existing questions

### Quality Assurance
- Content starts in **draft** status for admin review
- Validation against CEFR levels (A1-C2)
- Questions must be answerable from story content
- Clear correct answers with plausible distractors

## Setup

### 1. Environment Configuration

Copy the example environment file:
```bash
cp .env.example .env
```

Add your OpenAI API key to `.env`:
```env
OPENAI_API_KEY=sk-your-api-key-here
```

Get your API key from: https://platform.openai.com/api-keys

### 2. Deploy Backend Function

The AI generator function is automatically deployed with your Netlify functions:

```bash
# Development
netlify dev

# Production
git push origin main  # Triggers Netlify deployment
```

## Usage

### Admin Interface

1. **Navigate to AI Generation Page**
   - Log in as admin
   - Go to Admin Dashboard
   - Click "🤖 AI Generation" card

2. **Select a Story**
   - Browse published stories
   - Click "Generate AI Content" for desired story

3. **Choose Generation Type**
   - **Generate Questions**: Creates 8-12 questions
   - **Generate Quiz**: Creates quiz metadata
   - **Generate Both**: Creates questions and quiz together

4. **Review Generated Content**
   - AI-generated content appears in review queue
   - Look for 🤖 AI Generated badge
   - Review and edit as needed
   - Approve for publishing

### API Endpoints

#### Generate Content

```http
POST /.netlify/functions/ai-generator?type={type}&story_id={id}
```

**Query Parameters:**
- `type`: `questions`, `quiz`, or `both`
- `story_id`: MongoDB ObjectID of published story

**Response:**
```json
{
  "success": true,
  "message": "Successfully generated 10 questions",
  "questionsCount": 10,
  "questionIds": ["id1", "id2", ...],
  "quizId": "quiz_id" 
}
```

**Error Response:**
```json
{
  "error": "Error message describing the failure"
}
```

## Architecture

### Backend Flow

1. **Validation**
   - Check story exists and is published
   - Verify OpenAI API key is configured

2. **Prompt Engineering**
   - Build context-aware prompts with:
     - Story title, content, level (CEFR)
     - Topics and existing vocabulary
     - Word count and reading time
   - Include quality guidelines and constraints

3. **OpenAI API Call**
   - Use GPT-4o model with JSON mode
   - 60-second timeout
   - Structured response format

4. **Content Validation**
   - Validate against model schemas
   - Check question count, options, answers
   - Verify difficulty and points ranges

5. **Database Storage**
   - Store in MongoDB collections
   - Set status to `draft`
   - Mark creator as AI (ID: `000000000000000000000001`)
   - Link to source story

### Frontend Flow

1. **User Interface**
   - Admin selects published story
   - Clicks generation button
   - Real-time progress feedback

2. **API Communication**
   - Axios POST request to backend
   - 90-second client timeout
   - Error handling and retry logic

3. **Result Display**
   - Success message with counts
   - Error alerts if generation fails
   - Link to review queue

## Prompt Engineering

### Question Generation Prompt Structure

```
You are an expert German language teacher creating assessment materials for CEFR level {level}.

STORY CONTENT:
Title: {title}
Level: {level}
Content: {content}
Topics: {topics}
Existing Vocabulary: {vocabulary}

GENERATION REQUIREMENTS:
- Create 8-12 diverse questions
- Test genuine comprehension, not just recall
- Vocabulary questions use words from story
- Grammar questions focus on text structures
- All answerable from story content
- Difficulty matches CEFR level {level}

QUALITY STANDARDS:
- Unambiguous questions with clear correct answers
- Plausible but clearly incorrect distractors
- Accurate German grammar terminology
- Explanations aid learning
```

### Quiz Generation Prompt Structure

```
You are an expert German language teacher creating a quiz for CEFR level {level} learners.

STORY CONTENT:
Title: {title}
Level: {level}
Topics: {topics}
Word Count: {wordCount}

QUIZ GENERATION REQUIREMENTS:
- Engaging title related to story
- Clear description of assessment goals
- Appropriate time estimate (5-30 minutes)
- Recommended question type distribution
- Balanced difficulty progression
- Passing score between 50-90%
```

## Cost Management

### OpenAI API Costs

**Typical Generation Costs (GPT-4o):**
- Question Generation: ~$0.02-0.05 per story
- Quiz Generation: ~$0.01-0.02 per story
- Both: ~$0.03-0.07 per story

**Token Usage:**
- Input tokens: 500-2000 (depends on story length)
- Output tokens: 1500-3000 (for 10 questions)

### Rate Limiting

Current implementation:
- No built-in rate limiting
- Recommend implementing per-user limits
- Consider daily/monthly caps

**Suggested Rate Limits:**
```go
// Example rate limiting (not implemented)
type RateLimiter struct {
    MaxRequestsPerHour int // e.g., 10
    MaxRequestsPerDay  int // e.g., 50
}
```

## Monitoring

### Logging

The function logs:
- Generation requests (story ID, type)
- OpenAI API calls and responses
- Validation errors
- Database operations

### Metrics to Track

1. **Generation Success Rate**
   - Successful generations / Total attempts
   - Target: > 95%

2. **Content Approval Rate**
   - Approved AI content / Total AI content
   - Target: > 80%

3. **API Response Times**
   - Average generation time
   - Target: < 30 seconds

4. **Cost per Generation**
   - OpenAI API costs
   - Track monthly spending

## Troubleshooting

### Common Issues

#### "OpenAI API key not configured"
- Ensure `OPENAI_API_KEY` is set in environment
- Verify key is valid on OpenAI platform
- Check Netlify environment variables

#### "Story must be published to generate content"
- Only published stories can generate content
- Check story status in database
- Publish story through admin interface

#### "Generation timeout"
- Story may be too long (>15,000 characters)
- Try again (temporary API issues)
- Check OpenAI API status page

#### "Failed to parse OpenAI response"
- Response didn't match expected format
- Check OpenAI model version
- Review prompt structure

### Debug Mode

Enable detailed logging (not currently implemented):
```go
// Future enhancement
const DEBUG_MODE = os.Getenv("AI_GENERATOR_DEBUG") == "true"
```

## Future Enhancements

### Planned Features

1. **Customization Options**
   - Adjust question count (5-20)
   - Specify difficulty distribution
   - Choose question types

2. **Batch Generation**
   - Generate for multiple stories
   - Schedule bulk operations
   - Progress tracking

3. **Quality Scoring**
   - AI confidence scores
   - Automatic quality checks
   - Flag low-quality content

4. **Learning from Edits**
   - Track admin edits
   - Improve prompts over time
   - A/B test different approaches

5. **Multi-language Support**
   - Support for other target languages
   - Language-specific prompts
   - Cultural adaptations

## Security Considerations

### API Key Protection
- Never commit API keys to repository
- Use environment variables
- Rotate keys periodically
- Monitor usage for anomalies

### Input Validation
- Sanitize story content before sending to OpenAI
- Validate MongoDB ObjectIDs
- Check user permissions
- Rate limit requests

### Content Moderation
- Review AI-generated content before publishing
- Flag inappropriate content
- Maintain edit history
- Human oversight required

## Support

### Resources
- OpenAI Documentation: https://platform.openai.com/docs
- OpenAI Go SDK: https://github.com/openai/openai-go
- CEFR Levels: https://www.coe.int/en/web/common-european-framework-reference-languages

### Contact
For issues or questions:
1. Check documentation
2. Review error logs
3. Contact admin team
