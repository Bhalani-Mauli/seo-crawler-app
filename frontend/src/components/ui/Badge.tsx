import React from "react";

interface BadgeProps {
  children: React.ReactNode;
  variant?: "default" | "success" | "warning" | "error" | "info";
  size?: "sm" | "md";
  className?: string;
}

const badgeVariants = {
  default: "bg-gray-100 text-gray-800",
  success: "bg-green-100 text-green-800",
  warning: "bg-yellow-100 text-yellow-800",
  error: "bg-red-100 text-red-800",
  info: "bg-blue-100 text-blue-800",
} as const;

const badgeSizes = {
  sm: "px-2 py-1 text-xs",
  md: "px-3 py-1 text-sm",
} as const;

export const Badge: React.FC<BadgeProps> = ({
  children,
  variant = "default",
  size = "sm",
  className = "",
}) => {
  const baseClasses = "inline-flex items-center font-medium rounded-full";
  const variantClasses = badgeVariants[variant];
  const sizeClasses = badgeSizes[size];

  return (
    <span
      className={`${baseClasses} ${variantClasses} ${sizeClasses} ${className}`}
    >
      {children}
    </span>
  );
};
