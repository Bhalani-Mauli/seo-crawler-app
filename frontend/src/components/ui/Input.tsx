import React from "react";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  error?: string;
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ error, className = "", ...props }, ref) => {
    const baseClasses =
      "w-full px-4 py-2 border text-inherit rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors duration-200";
    const errorClasses = error
      ? "border-red-300 focus:ring-red-500"
      : "border-gray-300";
    const disabledClasses = props.disabled
      ? "bg-gray-100 cursor-not-allowed"
      : "bg-white";

    return (
      <div className="w-full">
        <input
          ref={ref}
          className={`${baseClasses} ${errorClasses} ${disabledClasses} ${className}`}
          {...props}
        />
        {error && <p className="mt-1 text-sm text-red-600">{error}</p>}
      </div>
    );
  }
);

Input.displayName = "Input";
