import React from "react";
import type { LinkData } from "../../types";
import { Badge } from "../ui/Badge";
import { LINK_TYPE_COLORS } from "../../constants";
import {
  getStatusCodeColor,
  getAccessibilityColor,
  getAccessibilityText,
} from "../../utils/status";

interface LinksTableProps {
  links: LinkData[];
  maxHeight?: string; // Optional prop to control max height
}

export const LinksTable: React.FC<LinksTableProps> = ({
  links,
  maxHeight = "300px", // Reduced default height
}) => {
  return (
    <div className="flex flex-col">
      <h4 className="font-semibold mb-3 text-lg flex-shrink-0">
        Links ({links.length})
      </h4>
      <div className="flex flex-col border border-gray-200 rounded-lg overflow-hidden">
        {/* Fixed Header */}
        <div className="flex-shrink-0 bg-gray-50 border-b border-gray-200">
          <div className="grid grid-cols-12 gap-0 text-xs font-medium text-gray-700">
            <div className="col-span-2 px-3 py-2">Type</div>
            <div className="col-span-4 px-3 py-2">URL</div>
            <div className="col-span-3 px-3 py-2">Text</div>
            <div className="col-span-2 px-3 py-2">Status</div>
            <div className="col-span-1 px-3 py-2">Accessible</div>
          </div>
        </div>

        {/* Scrollable Body */}
        <div
          className="overflow-auto bg-white"
          style={{ maxHeight, minHeight: "120px" }}
        >
          {links.length === 0 ? (
            <div className="text-center text-gray-400 py-8">No links found</div>
          ) : (
            <div className="divide-y divide-gray-100">
              {links.map((link) => (
                <div
                  key={link.id}
                  className="grid grid-cols-12 gap-0 hover:bg-gray-50 transition-colors text-xs"
                >
                  <div className="col-span-2 px-3 py-2 flex items-center">
                    <Badge
                      variant={link.type === "internal" ? "info" : "success"}
                      className={LINK_TYPE_COLORS[link.type]}
                    >
                      {link.type}
                    </Badge>
                  </div>
                  <div className="col-span-4 px-3 py-2 flex items-center">
                    <a
                      href={link.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-blue-600 hover:text-blue-800 underline truncate block transition-colors"
                      title={link.url}
                    >
                      {link.url}
                    </a>
                  </div>
                  <div
                    className="col-span-3 px-3 py-2 flex items-center truncate text-gray-900"
                    title={link.text || undefined}
                  >
                    {link.text || "-"}
                  </div>
                  <div className="col-span-2 px-3 py-2 flex items-center">
                    <Badge
                      variant={
                        link.status_code && link.status_code >= 400
                          ? "error"
                          : "success"
                      }
                      className={getStatusCodeColor(link.status_code)}
                    >
                      {link.status_code || "N/A"}
                    </Badge>
                  </div>
                  <div className="col-span-1 px-3 py-2 flex items-center">
                    <Badge
                      variant={link.is_accessible ? "success" : "error"}
                      className={getAccessibilityColor(link.is_accessible)}
                    >
                      {getAccessibilityText(link.is_accessible)}
                    </Badge>
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
