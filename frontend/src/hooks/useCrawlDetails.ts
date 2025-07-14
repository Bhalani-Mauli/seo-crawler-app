import { useState, useCallback } from "react";
import type { LinkData, HeadingData, URLData } from "../types";
import { apiService } from "../services/api";

export const useCrawlDetails = () => {
  const [detailId, setDetailId] = useState<number | null>(null);
  const [detailLinks, setDetailLinks] = useState<LinkData[]>([]);
  const [detailHeadings, setDetailHeadings] = useState<HeadingData[]>([]);
  const [detailMeta, setDetailMeta] = useState<URLData | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);

  const openDetails = useCallback(async (id: number) => {
    setDetailId(id);
    setDetailLoading(true);

    try {
      const [linksRes, headingsRes, metaRes] = await Promise.all([
        apiService.fetchLinksByResultId(id),
        apiService.fetchHeadingsByResultId(id),
        apiService.fetchResultById(id),
      ]);

      setDetailLinks(linksRes.links || []);
      setDetailHeadings(headingsRes.headings || []);
      setDetailMeta(metaRes || null);
    } catch (error) {
      console.error("Error fetching details:", error);
      setDetailLinks([]);
      setDetailHeadings([]);
      setDetailMeta(null);
    } finally {
      setDetailLoading(false);
    }
  }, []);

  const closeDetails = useCallback(() => {
    setDetailId(null);
    setDetailLinks([]);
    setDetailHeadings([]);
    setDetailMeta(null);
  }, []);

  return {
    detailId,
    detailLinks,
    detailHeadings,
    detailMeta,
    detailLoading,
    openDetails,
    closeDetails,
  };
};
