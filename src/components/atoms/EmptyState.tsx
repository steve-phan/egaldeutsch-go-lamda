import React from "react";
import { Card, CardContent } from "../ui/card";

interface EmptyStateProps {
  message?: string;
}

export const EmptyState: React.FC<EmptyStateProps> = ({ 
  message = "No items found for the selected filters." 
}) => {
  return (
    <Card>
      <CardContent className="pt-6 text-center">
        <p className="text-muted-foreground">{message}</p>
      </CardContent>
    </Card>
  );
};
