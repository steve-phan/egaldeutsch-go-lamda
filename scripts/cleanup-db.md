# Database Cleanup Script

## Remove redundant isActive field

The `isActive` field is now redundant since we use the `status` field with the 4-state workflow.

### MongoDB Commands to clean up:

```javascript
// Remove isActive field from all stories
db.stories.updateMany({}, { $unset: { isActive: "" } });

// Remove isActive field from all questions
db.questions.updateMany({}, { $unset: { isActive: "" } });

// Remove isActive field from all quizzes
db.quizzes.updateMany({}, { $unset: { isActive: "" } });
```

### Verify the cleanup:

```javascript
// Check that no documents have isActive field
db.stories.find({ isActive: { $exists: true } }).count();
db.questions.find({ isActive: { $exists: true } }).count();
db.quizzes.find({ isActive: { $exists: true } }).count();
```

All counts should return 0 after cleanup.

## Status Migration (if needed)

If you have old documents with different status values, you can migrate them:

```javascript
// Example: Convert old "active" status to "published"
db.stories.updateMany({ status: "active" }, { $set: { status: "published" } });

// Example: Convert old "pending_review" to "preview"
db.stories.updateMany(
  { status: "pending_review" },
  { $set: { status: "preview" } }
);
```
