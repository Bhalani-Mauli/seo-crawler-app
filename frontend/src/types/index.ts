export interface CrawlData {
  html_version: string | null;
  title: string | null;
  headings: Record<string, number>;
  internal_links: number;
  external_links: number;
  inaccessible_links: number;
  has_login_form: boolean;
  status: CrawlStatus;
  created_at: string | null;
  updated_at: string | null;
}

export interface URLData {
  id: number;
  url: string;
  crawl_data: CrawlData;
}

export interface PaginationResponse {
  data: URLData[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface LinkData {
  id: number;
  url: string;
  text: string | null;
  type: "internal" | "external";
  status_code: number | null;
  is_accessible: boolean;
}

export interface HeadingData {
  id: number;
  level: string;
  text: string;
  order: number;
}

export interface LinksResponse {
  links: LinkData[];
}

export interface HeadingsResponse {
  headings: HeadingData[];
}

export type CrawlStatus = "pending" | "running" | "done" | "error" | "stopped";

export interface FetchResultsParams {
  page?: number;
  pageSize?: number;
  search?: string;
  status?: string;
  sortBy?: string;
  sortOrder?: "asc" | "desc";
}

export interface ApiResponse<T = any> {
  message?: string;
  data?: T;
}
