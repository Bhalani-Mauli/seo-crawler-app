export { apiService } from "./services/api";

import { apiService } from "./services/api";

export const {
  fetchResults,
  addUrl,
  rerunUrls,
  deleteUrls,
  stopCrawl,
  fetchResultById,
  fetchLinksByResultId,
  fetchHeadingsByResultId,
} = apiService;
