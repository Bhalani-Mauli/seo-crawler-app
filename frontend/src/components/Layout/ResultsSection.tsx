import React from "react";
import { CrawlTable } from "../CrawlTable/CrawlTable";
import type { URLData } from "../../types";

interface ResultsSectionProps {
  data: URLData[];
  total: number;
  onDetails: (id: number) => void;
}

export const ResultsSection: React.FC<ResultsSectionProps> = ({
  data,
  total,
  onDetails,
}) => {
  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold text-gray-900">Crawl Results</h2>
        <span className="text-sm text-gray-500">Total Results: {total}</span>
      </div>
      <CrawlTable data={data} onDetails={onDetails} />
    </div>
  );
};
