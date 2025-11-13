import React from "react";
import { Alert, AlertDescription } from "../ui/alert";

interface ErrorAlertProps {
  message: string;
}

export const ErrorAlert: React.FC<ErrorAlertProps> = ({ message }) => {
  if (!message) return null;
  
  return (
    <Alert variant="destructive" className="mb-6">
      <AlertDescription>{message}</AlertDescription>
    </Alert>
  );
};
