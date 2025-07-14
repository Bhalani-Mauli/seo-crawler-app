import React from "react";
import type { HeadingData } from "../../types";
import { Badge } from "../ui/Badge";

interface HeadingsTableProps {
  headings: HeadingData[];
  maxHeight?: string;
}

export const HeadingsTable: React.FC<HeadingsTableProps> = ({
  headings,
  maxHeight = "300px",
}) => {
  return (
    <div className="flex flex-col">
      <h4 className="font-semibold mb-3 text-lg flex-shrink-0">
        Headings ({headings.length})
      </h4>
      <div className="flex flex-col border border-gray-200 rounded-lg overflow-hidden">
        {/* Fixed Header */}
        <div className="flex-shrink-0 bg-gray-50 border-b border-gray-200">
          <div className="grid grid-cols-12 gap-0 text-xs font-medium text-gray-700">
            <div className="col-span-2 px-3 py-2">Level</div>
            <div className="col-span-8 px-3 py-2">Text</div>
            <div className="col-span-2 px-3 py-2">Order</div>
          </div>
        </div>

        {/* Scrollable Body */}
        <div
          className="overflow-auto bg-white"
          style={{ maxHeight, minHeight: "120px" }}
        >
          {headings.length === 0 ? (
            <div className="text-center text-gray-400 py-8">
              No headings found
            </div>
          ) : (
            <div className="divide-y divide-gray-100">
              {headings.map((heading) => (
                <div
                  key={heading.id}
                  className="grid grid-cols-12 gap-0 hover:bg-gray-50 transition-colors text-xs"
                >
                  <div className="col-span-2 px-3 py-2 flex items-center">
                    <Badge
                      variant="info"
                      className="bg-purple-100 text-purple-800"
                    >
                      {heading.level}
                    </Badge>
                  </div>
                  <div
                    className="col-span-8 px-3 py-2 flex items-center truncate text-gray-900"
                    title={heading.text}
                  >
                    {heading.text}
                  </div>
                  <div className="col-span-2 px-3 py-2 flex items-center text-gray-900">
                    {heading.order}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
