import type { ColumnDef } from "@tanstack/react-table";
import type { URLData } from "../../types";
import { Badge } from "../ui/Badge";
import { getStatusColor } from "../../utils/status";

export const tableColumns: ColumnDef<URLData>[] = [
  {
    accessorKey: "url",
    header: "URL",
    cell: ({ row }) => (
      <a
        href={row.original.url}
        target="_blank"
        rel="noopener noreferrer"
        className="text-blue-600 hover:text-blue-800 underline transition-colors"
      >
        {row.original.url}
      </a>
    ),
  },
  {
    accessorKey: "crawl_data.title",
    header: "Title",
    cell: ({ row }) => (
      <span
        className="max-w-xs truncate block text-gray-900"
        title={row.original.crawl_data.title || undefined}
      >
        {row.original.crawl_data.title || "N/A"}
      </span>
    ),
  },
  {
    accessorKey: "crawl_data.html_version",
    header: "HTML Version",
    cell: ({ row }) => (
      <span className="text-gray-900">
        {row.original.crawl_data.html_version || "N/A"}
      </span>
    ),
  },
  {
    accessorKey: "crawl_data.internal_links",
    header: "Internal Links",
    cell: ({ row }) => (
      <span className="text-gray-900">
        {row.original.crawl_data.internal_links}
      </span>
    ),
  },
  {
    accessorKey: "crawl_data.external_links",
    header: "External Links",
    cell: ({ row }) => (
      <span className="text-gray-900">
        {row.original.crawl_data.external_links}
      </span>
    ),
  },
  {
    accessorKey: "crawl_data.inaccessible_links",
    header: "Broken Links",
    cell: ({ row }) => (
      <span className="text-gray-900">
        {row.original.crawl_data.inaccessible_links}
      </span>
    ),
  },
  {
    accessorKey: "crawl_data.has_login_form",
    header: "Login Form",
    cell: ({ row }) => (
      <span
        className={
          row.original.crawl_data.has_login_form
            ? "text-green-600"
            : "text-gray-500"
        }
      >
        {row.original.crawl_data.has_login_form ? "Yes" : "No"}
      </span>
    ),
  },
  {
    accessorKey: "crawl_data.status",
    header: "Status",
    cell: ({ row }) => {
      const status = row.original.crawl_data.status;
      return (
        <Badge
          variant={
            status === "done"
              ? "success"
              : status === "error"
              ? "error"
              : "warning"
          }
          className={getStatusColor(status)}
        >
          {status}
        </Badge>
      );
    },
  },
];
