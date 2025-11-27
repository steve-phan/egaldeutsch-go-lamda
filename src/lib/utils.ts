import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Utility function to merge Tailwind CSS classes with clsx
 * This is a core utility for shadcn/ui components
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Format level for display (A1, A2, B1, etc.)
 */
export function formatLevel(level: string): string {
  return level.toUpperCase();
}

/**
 * Get color classes for level badges
 */
export function getLevelColor(level: string): string {
  const levelMap: Record<string, string> = {
    A1: "bg-green-100 text-green-800",
    A2: "bg-blue-100 text-blue-800",
    B1: "bg-yellow-100 text-yellow-800",
    B2: "bg-orange-100 text-orange-800",
    C1: "bg-red-100 text-red-800",
    C2: "bg-purple-100 text-purple-800",
  };

  return levelMap[level.toUpperCase()] || "bg-gray-100 text-gray-800";
}

// Get role badge color
export const getRoleBadgeColor = (role: string) => {
  switch (role) {
    case "admin":
      return "bg-red-500/10 text-red-600 border-red-500/20";
    case "reviewer":
      return "bg-purple-500/10 text-purple-600 border-purple-500/20";
    case "creator":
      return "bg-blue-500/10 text-blue-600 border-blue-500/20";
    default:
      return "bg-gray-500/10 text-gray-600 border-gray-500/20";
  }
};
