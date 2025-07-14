import React from "react";

interface MessageDisplayProps {
  message: string | null;
}

export const MessageDisplay: React.FC<MessageDisplayProps> = ({ message }) => {
  if (!message) return null;

  const isError = message.includes("Error");
  const messageClasses = isError
    ? "bg-red-100 text-red-700 border border-red-200"
    : "bg-green-100 text-green-700 border border-green-200";

  return <div className={`p-3 rounded-md ${messageClasses}`}>{message}</div>;
};
