import { useState, useEffect } from "react";
import { Navbar } from "../components/Layout/Navbar";
import { Header } from "../components/Layout/Header";
import { UrlForm } from "../components/forms/UrlForm";
import { MessageDisplay } from "../components/Layout/MessageDisplay";
import { ResultsSection } from "../components/Layout/ResultsSection";
import { CrawlDetailsModal } from "../components/CrawlDetails/CrawlDetailsModal";
import { useCrawlData } from "../hooks/useCrawlData";
import { useCrawlDetails } from "../hooks/useCrawlDetails";
import { validateEnv } from "../utils/env";

export const DashboardPage: React.FC = () => {
  const [message, setMessage] = useState<string | null>(null);
  const { crawlData, loading, addUrl, isPolling } = useCrawlData();
  const {
    detailId,
    detailLinks,
    detailHeadings,
    detailMeta,
    detailLoading,
    openDetails,
    closeDetails,
  } = useCrawlDetails();

  useEffect(() => {
    const envValidation = validateEnv();
    if (!envValidation.isValid) {
      console.warn("Environment validation failed:", envValidation.missingVars);
    }
  }, []);

  const handleUrlSubmit = async (url: string) => {
    try {
      const response = await addUrl(url);
      setMessage(response.message || "URL added successfully");
    } catch (error) {
      setMessage("Error submitting URL");
    }
  };

  const handleDetails = (id: number) => {
    openDetails(id);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <div className="py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <Header />
          <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
            <UrlForm onSubmit={handleUrlSubmit} loading={loading} />
            <MessageDisplay message={message} />
            {isPolling && (
              <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
                <div className="flex items-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600 mr-2"></div>
                  <span className="text-sm text-blue-700">
                    Monitoring crawl progress...
                  </span>
                </div>
              </div>
            )}
          </div>
          <ResultsSection
            data={crawlData.data}
            total={crawlData.total}
            onDetails={handleDetails}
          />
        </div>
      </div>
      <CrawlDetailsModal
        isOpen={detailId !== null}
        onClose={closeDetails}
        detailMeta={detailMeta}
        detailLinks={detailLinks}
        detailHeadings={detailHeadings}
        detailLoading={detailLoading}
      />
    </div>
  );
};

export default DashboardPage;
