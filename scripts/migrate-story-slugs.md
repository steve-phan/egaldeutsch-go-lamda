# Story Slug Migration

This script generates SEO-friendly slugs for all existing stories in the database.

## What it does

1. Connects to MongoDB
2. Finds all stories without a `slug` field (or with empty slugs)
3. Generates a URL-friendly slug from each story's title
4. Handles duplicates by appending numbers (e.g., `ein-tag-1`, `ein-tag-2`)
5. Updates each story with its unique slug

## Prerequisites

- Go 1.18 or later
- MongoDB connection configured in `.env`
- Stories collection exists with stories

## How to run

### From project root:

```bash
# Run the migration
go run scripts/migrate-story-slugs.go
```

### Expected output:

```
Starting slug migration for existing stories...
Connected to database
Found 15 stories without slugs
Processing story 1/15: Ein schöner Tag
  ✅ Updated with slug: ein-schoener-tag
Processing story 2/15: Der Familienausflug
  ✅ Updated with slug: der-familienausflug
...
==================================================
Migration complete!
Total stories: 15
Successfully updated: 15
Failed: 0
==================================================
```

## Slug Format

- Lowercase
- Spaces replaced with hyphens
- German characters transliterated (ä→ae, ö→oe, ü→ue, ß→ss)
- Special characters removed
- Maximum 60 characters
- No consecutive hyphens

## Examples

| Title | Generated Slug |
|-------|---------------|
| Ein schöner Tag | `ein-schoener-tag` |
| Übung macht den Meister | `uebung-macht-den-meister` |
| Das große Abenteuer! | `das-grosse-abenteuer` |

## Safety

- ✅ **Idempotent**: Safe to run multiple times
- ✅ **Non-destructive**: Only adds/updates the `slug` field
- ✅ **Uniqueness**: Guarantees unique slugs by appending numbers
- ✅ **Backwards compatible**: ID-based URLs continue to work

## Verification

After running, you can verify:

```bash
# Check that all stories have slugs
mongosh <your-db-url>
> use <your-database>
> db.stories.find({ $or: [{ slug: { $exists: false } }, { slug: "" }] }).count()
# Should return 0
```

## Rollback

To remove all slugs (if needed):

```bash
mongosh <your-db-url>
> use <your-database>
> db.stories.updateMany({}, { $unset: { slug: "" } })
```
