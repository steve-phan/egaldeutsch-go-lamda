# Atom UI Components Documentation

This document describes the reusable, production-grade UI component system built with Tailwind CSS and shadcn/ui principles.

## Overview

The component system follows atomic design patterns with a global theme system built on CSS variables, enabling easy theme customization and maintainability at a production-grade level.

## Features

- 🎨 **Global Theme System**: CSS variables for easy theme switching (light/dark mode ready)
- 🔧 **Type-Safe**: Full TypeScript support with proper typing
- ♿ **Accessible**: Built with accessibility best practices using Radix UI primitives
- 🎯 **Reusable**: Atomic components that can be composed into larger components
- 📦 **Production-Ready**: Battle-tested patterns from shadcn/ui
- 🚀 **Performance**: Optimized with minimal re-renders and efficient bundling

## Theme System

### CSS Variables

All components use CSS variables defined in `src/styles/global.css`. This allows for easy theme customization without modifying component code.

**Available Theme Colors:**
- `--primary` / `--primary-foreground`
- `--secondary` / `--secondary-foreground`
- `--success` / `--success-foreground`
- `--warning` / `--warning-foreground`
- `--destructive` / `--destructive-foreground`
- `--muted` / `--muted-foreground`
- `--accent` / `--accent-foreground`
- `--card` / `--card-foreground`
- `--popover` / `--popover-foreground`
- `--background` / `--foreground`
- `--border` / `--input` / `--ring`

### Dark Mode Support

Dark mode is supported through the `.dark` class. Toggle dark mode by adding/removing this class from the root element.

```typescript
// Example: Toggle dark mode
document.documentElement.classList.toggle('dark');
```

## Core Utilities

### `cn()` Function

The `cn()` utility merges Tailwind CSS classes intelligently, handling conflicts and conditional classes.

```typescript
import { cn } from "@/lib/utils";

// Basic usage
cn("text-red-500", "text-blue-500") // => "text-blue-500"

// Conditional classes
cn("base-class", condition && "conditional-class")

// Object syntax
cn("base", { "active": isActive, "disabled": isDisabled })
```

## Component Library

### Button

A versatile button component with multiple variants and sizes.

**Variants:** `default` | `destructive` | `success` | `outline` | `secondary` | `ghost` | `link`

**Sizes:** `default` | `sm` | `lg` | `icon`

```tsx
import { Button } from "@/components/ui";

// Basic usage
<Button>Click me</Button>

// With variants
<Button variant="destructive">Delete</Button>
<Button variant="success">Save</Button>
<Button variant="outline">Cancel</Button>

// Different sizes
<Button size="sm">Small</Button>
<Button size="lg">Large</Button>
<Button size="icon">🔍</Button>

// Disabled state
<Button disabled>Disabled</Button>
```

### Card

A flexible container component for grouping related content.

**Sub-components:** `Card`, `CardHeader`, `CardTitle`, `CardDescription`, `CardContent`, `CardFooter`

```tsx
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "@/components/ui";

<Card>
  <CardHeader>
    <CardTitle>Card Title</CardTitle>
    <CardDescription>Card description goes here</CardDescription>
  </CardHeader>
  <CardContent>
    <p>Main content of the card</p>
  </CardContent>
  <CardFooter>
    <Button>Action</Button>
  </CardFooter>
</Card>
```

### Badge

Small status indicators or labels.

**Variants:** `default` | `secondary` | `destructive` | `success` | `warning` | `outline`

```tsx
import { Badge } from "@/components/ui";

<Badge>Default</Badge>
<Badge variant="success">Success</Badge>
<Badge variant="destructive">Error</Badge>
<Badge variant="warning">Warning</Badge>
<Badge variant="outline">Outline</Badge>
```

### Input

A styled text input component.

```tsx
import { Input, Label } from "@/components/ui";

<div>
  <Label htmlFor="email">Email</Label>
  <Input 
    id="email" 
    type="email" 
    placeholder="Enter your email"
  />
</div>
```

### Label

Accessible form labels using Radix UI primitives.

```tsx
import { Label } from "@/components/ui";

<Label htmlFor="username">Username</Label>
<Input id="username" />
```

### Alert

Display important messages to users.

**Variants:** `default` | `destructive` | `success` | `warning`

```tsx
import { Alert, AlertTitle, AlertDescription } from "@/components/ui";

<Alert variant="success">
  <AlertTitle>Success!</AlertTitle>
  <AlertDescription>
    Your changes have been saved successfully.
  </AlertDescription>
</Alert>

<Alert variant="destructive">
  <AlertTitle>Error</AlertTitle>
  <AlertDescription>
    Something went wrong. Please try again.
  </AlertDescription>
</Alert>
```

### Separator

Visual divider for content sections.

```tsx
import { Separator } from "@/components/ui";

<div>
  <p>Content above</p>
  <Separator />
  <p>Content below</p>
</div>

{/* Vertical separator */}
<div className="flex h-5 items-center space-x-4">
  <div>Item 1</div>
  <Separator orientation="vertical" />
  <div>Item 2</div>
</div>
```

### Textarea

Multi-line text input component.

```tsx
import { Textarea, Label } from "@/components/ui";

<div>
  <Label htmlFor="message">Message</Label>
  <Textarea 
    id="message" 
    placeholder="Type your message here"
    rows={4}
  />
</div>
```

## Migration Guide

### Migrating Existing Components

When refactoring existing components to use the new atom UI components:

1. **Import atom components:**
   ```tsx
   import { Button, Card, Badge } from "@/components/ui";
   ```

2. **Replace hardcoded classes with semantic components:**
   ```tsx
   // Before
   <button className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded">
     Click me
   </button>

   // After
   <Button>Click me</Button>
   ```

3. **Use theme colors instead of hardcoded Tailwind colors:**
   ```tsx
   // Before
   <div className="text-gray-900 bg-white">Content</div>

   // After
   <div className="text-foreground bg-background">Content</div>
   ```

4. **Leverage component variants:**
   ```tsx
   // Before - multiple conditional classes
   <button className={`px-4 py-2 rounded ${
     type === 'danger' ? 'bg-red-600' : 'bg-blue-600'
   }`}>
     Button
   </button>

   // After - clean variant prop
   <Button variant={type === 'danger' ? 'destructive' : 'default'}>
     Button
   </Button>
   ```

## Best Practices

1. **Use semantic component names**: Button, Card, Alert are clearer than div with classes
2. **Leverage the theme system**: Use CSS variables through Tailwind classes (e.g., `bg-primary` instead of `bg-blue-600`)
3. **Compose components**: Build complex UIs by composing simple atom components
4. **Maintain consistency**: Use the same variants across your application
5. **Extend wisely**: When extending components, use the `cn()` utility to merge classes properly

## Examples

### Example 1: Story Card (Refactored)

```tsx
import { Card, CardContent, CardFooter, Badge, Button } from "@/components/ui";

const StoryCard = ({ story }) => (
  <Card>
    <CardContent className="pt-6">
      <div className="flex items-start justify-between mb-3">
        <h3 className="text-xl font-semibold">{story.title}</h3>
        <Badge variant="outline">{story.level}</Badge>
      </div>
      <p className="text-muted-foreground">{story.summary}</p>
    </CardContent>
    <CardFooter className="flex gap-3">
      <Button className="flex-1">Read Story</Button>
      <Button variant="success" className="flex-1">Take Quiz</Button>
    </CardFooter>
  </Card>
);
```

### Example 2: Quiz Results (Refactored)

```tsx
import { Card, CardContent, Alert, AlertTitle, AlertDescription, Button } from "@/components/ui";

const QuizResults = ({ result, onRetake, onReturn }) => (
  <Card>
    <CardContent className="pt-8 pb-8 text-center">
      <h2 className="text-2xl font-bold mb-2">Quiz Results</h2>
      <div className="text-4xl font-bold text-primary mb-2">
        {result.percentage}%
      </div>
      
      {result.passed ? (
        <Alert variant="success">
          <AlertTitle>🎊 Congratulations!</AlertTitle>
          <AlertDescription>You passed the quiz!</AlertDescription>
        </Alert>
      ) : (
        <Alert variant="warning">
          <AlertTitle>📖 Keep Learning!</AlertTitle>
          <AlertDescription>Try again to improve.</AlertDescription>
        </Alert>
      )}
      
      <div className="flex gap-4 mt-6">
        <Button onClick={onReturn}>Read Story Again</Button>
        <Button variant="success" onClick={onRetake}>Retake Quiz</Button>
      </div>
    </CardContent>
  </Card>
);
```

## Customization

### Adjusting Theme Colors

To customize the theme, modify the CSS variables in `src/styles/global.css`:

```css
:root {
  --primary: 199.5 89.1% 48.2%; /* Change this */
  --success: 142.1 76.2% 36.3%; /* And this */
  /* ... other variables */
}
```

### Creating New Variants

Add new variants to existing components using `class-variance-authority`:

```tsx
// In src/components/ui/button.tsx
const buttonVariants = cva(
  "base-classes",
  {
    variants: {
      variant: {
        // ... existing variants
        info: "bg-blue-500 text-white hover:bg-blue-600", // New variant
      },
    },
  }
);
```

## Future Enhancements

Potential additions to the component library:

- Dialog/Modal component
- Dropdown Menu
- Tabs component
- Toast notifications
- Form components (Checkbox, Radio, Select)
- Progress indicators
- Skeleton loaders
- Tooltip component

## Support

For issues or questions about the component system:
- Check the [shadcn/ui documentation](https://ui.shadcn.com)
- Review component source code in `src/components/ui/`
- Consult the Tailwind CSS documentation for styling options
