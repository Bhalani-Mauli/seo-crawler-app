import React, { useEffect } from "react";

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: React.ReactNode;
  size?: "sm" | "md" | "lg" | "xl";
}

const modalSizes = {
  sm: "max-w-md",
  md: "max-w-lg",
  lg: "max-w-2xl",
  xl: "max-w-6xl",
} as const;

export const Modal: React.FC<ModalProps> = ({
  isOpen,
  onClose,
  title,
  children,
  size = "md",
}) => {
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener("keydown", handleEscape);
      document.body.style.overflow = "hidden";
    }

    return () => {
      document.removeEventListener("keydown", handleEscape);
      document.body.style.overflow = "unset";
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50 p-4">
      <div
        className={`bg-white rounded-lg shadow-lg w-full ${modalSizes[size]} max-h-[90vh] flex flex-col`}
      >
        {title && (
          <div className="p-6 border-b border-gray-200 flex-shrink-0">
            <div className="flex justify-between items-center">
              <h3 className="text-2xl font-bold text-gray-900">{title}</h3>
              <button
                onClick={onClose}
                className="text-gray-500 hover:text-gray-700 text-2xl font-bold transition-colors"
                aria-label="Close modal"
              >
                &times;
              </button>
            </div>
          </div>
        )}
        <div className="flex-1 overflow-scroll p-6">{children}</div>
      </div>
    </div>
  );
};
