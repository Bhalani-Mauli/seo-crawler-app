import type {
  PaginationResponse,
  FetchResultsParams,
  LinksResponse,
  HeadingsResponse,
  URLData,
  ApiResponse,
} from "../types";
import { API_CONFIG } from "../constants";
import {
  createApiHeaders,
  handleApiResponse,
  buildQueryString,
} from "../utils/api";

class ApiService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = API_CONFIG.BASE_URL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const currentHeaders = createApiHeaders();
    const response = await fetch(url, {
      headers: { ...currentHeaders, ...options.headers },
      ...options,
    });

    return handleApiResponse<T>(response);
  }

  async fetchResults(
    params: FetchResultsParams = {}
  ): Promise<PaginationResponse> {
    const queryString = buildQueryString({
      page: params.page || 1,
      page_size: params.pageSize || 10,
      search: params.search || "",
      status: params.status || "",
      sort_by: params.sortBy || "created_at",
      sort_order: params.sortOrder || "desc",
    });

    return this.request<PaginationResponse>(`/results?${queryString}`);
  }

  async addUrl(url: string): Promise<ApiResponse> {
    return this.request<ApiResponse>("/crawl", {
      method: "POST",
      headers: createApiHeaders({ "Content-Type": "application/json" }),
      body: JSON.stringify({ url }),
    });
  }

  async rerunUrls(urls: string[]): Promise<ApiResponse> {
    return this.request<ApiResponse>("/bulk/rerun", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ urls }),
    });
  }

  async deleteUrls(urls: string[]): Promise<ApiResponse> {
    return this.request<ApiResponse>("/bulk/delete", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ urls }),
    });
  }

  async stopCrawl(id: number): Promise<ApiResponse> {
    return this.request<ApiResponse>(`/stop/${id}`, {
      method: "POST",
    });
  }

  async fetchResultById(id: number): Promise<URLData> {
    return this.request<URLData>(`/results/${id}`);
  }

  async fetchLinksByResultId(id: number): Promise<LinksResponse> {
    return this.request<LinksResponse>(`/results/${id}/links`);
  }

  async fetchHeadingsByResultId(id: number): Promise<HeadingsResponse> {
    return this.request<HeadingsResponse>(`/results/${id}/headings`);
  }
}

export const apiService = new ApiService();
