import { useState, useEffect, useCallback, useRef } from "react";
import type { PaginationResponse, FetchResultsParams } from "../types";
import { apiService } from "../services/api";
import { POLLING_INTERVAL, MAX_RESULTS_PER_PAGE } from "../constants";

export const useCrawlData = () => {
  const [crawlData, setCrawlData] = useState<PaginationResponse>({
    data: [],
    total: 0,
    page: 1,
    page_size: 10,
    total_pages: 0,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isPolling, setIsPolling] = useState(false);
  const pollingIntervalRef = useRef<number | null>(null);

  // Check if there are any pending or running crawls
  const hasActiveCrawls = useCallback((data: PaginationResponse) => {
    return data.data.some(
      (item) =>
        item.crawl_data.status === "pending" ||
        item.crawl_data.status === "running"
    );
  }, []);

  // Start polling
  const startPolling = useCallback(() => {
    if (pollingIntervalRef.current) {
      clearInterval(pollingIntervalRef.current);
    }
    setIsPolling(true);
    pollingIntervalRef.current = setInterval(() => {
      fetchData();
    }, POLLING_INTERVAL);
  }, []);

  // Stop polling
  const stopPolling = useCallback(() => {
    if (pollingIntervalRef.current) {
      clearInterval(pollingIntervalRef.current);
      pollingIntervalRef.current = null;
    }
    setIsPolling(false);
  }, []);

  const fetchData = useCallback(
    async (params: FetchResultsParams = {}) => {
      try {
        setLoading(true);
        setError(null);
        const response = await apiService.fetchResults({
          page: 1,
          pageSize: MAX_RESULTS_PER_PAGE,
          sortBy: "created_at",
          sortOrder: "desc",
          ...params,
        });
        setCrawlData(response);

        if (hasActiveCrawls(response)) {
          if (!isPolling) {
            startPolling();
          }
        } else {
          stopPolling();
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to fetch data");
        console.error("Error fetching results:", err);
      } finally {
        setLoading(false);
      }
    },
    [hasActiveCrawls, isPolling, startPolling, stopPolling]
  );

  const addUrl = useCallback(
    async (url: string) => {
      try {
        setLoading(true);
        setError(null);
        const response = await apiService.addUrl(url);

        await fetchData();

        return response;
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Failed to add URL";
        setError(errorMessage);
        throw new Error(errorMessage);
      } finally {
        setLoading(false);
      }
    },
    [fetchData]
  );

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  useEffect(() => {
    return () => {
      if (pollingIntervalRef.current) {
        clearInterval(pollingIntervalRef.current);
      }
    };
  }, []);

  return {
    crawlData,
    loading,
    error,
    fetchData,
    addUrl,
    isPolling,
  };
};
