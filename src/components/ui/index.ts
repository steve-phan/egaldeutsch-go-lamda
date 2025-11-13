/**
 * Atom UI Components Library
 * 
 * This file exports all reusable, production-grade UI components
 * built with Tailwind CSS and shadcn/ui principles.
 * 
 * Components follow:
 * - Atomic design patterns
 * - Global theme system with CSS variables
 * - Full TypeScript support
 * - Accessibility best practices
 */

export { Button, buttonVariants, type ButtonProps } from "./button";
export { Badge, badgeVariants, type BadgeProps } from "./badge";
export {
  Card,
  CardHeader,
  CardFooter,
  CardTitle,
  CardDescription,
  CardContent,
} from "./card";
export { Input, type InputProps } from "./input";
export { Label } from "./label";
export { Alert, AlertTitle, AlertDescription } from "./alert";
export { Separator } from "./separator";
export { Textarea, type TextareaProps } from "./textarea";
export { DropdownMenu, DropdownMenuItem } from "./dropdown-menu";
export { MobileMenu, MobileMenuItem, MobileMenuGroup } from "./mobile-menu";
