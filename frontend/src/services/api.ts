class ApiService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = "http://localhost:8080";
  }

  private async request(): Promise<any> {
    return {};
  }

  async fetchResults(): Promise<any> {
    return Promise.resolve({
      data: [],
      total: 0,
      page: 1,
      page_size: 10,
      total_pages: 0,
    });
  }

  async addUrl(url: string): Promise<any> {
    return Promise.resolve({
      message: "URL added successfully",
      data: { id: 1, url },
    });
  }

  async rerunUrls(urls: string[]): Promise<any> {
    return Promise.resolve({
      message: "URLs rerun successfully",
      data: [],
    });
  }

  async deleteUrls(urls: string[]): Promise<any> {
    return Promise.resolve({
      message: "URLs deleted successfully",
      data: [],
    });
  }

  async stopCrawl(id: number): Promise<any> {
    return Promise.resolve({
      message: `Crawl with ID ${id} stopped successfully`,
    });
  }

  async fetchResultById(id: number): Promise<any> {
    return Promise.resolve({
      id,
      url: "http://example.com",
      crawl_data: {
        status: "done",
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        has_login_form: false,
      },
    });
  }

  async fetchLinksByResultId(id: number): Promise<any> {
    return Promise.resolve({
      links: [
        {
          id: 1,
          url: "http://example.com/internal",
          text: "Internal Link",
          type: "internal",
          status_code: 200,
          is_accessible: true,
        },
      ],
    });
  }

  async fetchHeadingsByResultId(id: number): Promise<any> {
    return Promise.resolve({
      headings: [
        {
          id: 1,
          level: "h1",
          text: "Main Heading",
          order: 1,
        },
      ],
    });
  }
}

export const apiService = new ApiService();
