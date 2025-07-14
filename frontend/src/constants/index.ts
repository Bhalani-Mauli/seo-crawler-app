import { getEnvVar, getEnvVarAsNumber, getEnvVarAsBoolean } from "../utils/env";

export const API_CONFIG = {
  BASE_URL: getEnvVar("VITE_API_BASE_URL", "http://localhost:8080/api"),
  API_KEY: getEnvVar("VITE_API_KEY", "seo-crawler-api-key-2025"),
} as const;

export const APP_CONFIG = {
  NAME: getEnvVar("VITE_APP_NAME", "SEO Crawler"),
  VERSION: getEnvVar("VITE_APP_VERSION", "1.0.0"),
} as const;

export const POLLING_INTERVAL = getEnvVarAsNumber(
  "VITE_POLLING_INTERVAL",
  10000
);

export const MAX_RESULTS_PER_PAGE = getEnvVarAsNumber(
  "VITE_MAX_RESULTS_PER_PAGE",
  50
);

export const FEATURE_FLAGS = {
  DEBUG_MODE: getEnvVarAsBoolean("VITE_ENABLE_DEBUG_MODE", false),
  ANALYTICS: getEnvVarAsBoolean("VITE_ENABLE_ANALYTICS", false),
} as const;

export const STATUS_COLORS = {
  pending: "text-yellow-600 bg-yellow-100",
  running: "text-blue-600 bg-blue-100",
  done: "text-green-600 bg-green-100",
  error: "text-red-600 bg-red-100",
  stopped: "text-gray-600 bg-gray-100",
} as const;

export const LINK_TYPE_COLORS = {
  internal: "bg-blue-100 text-blue-800",
  external: "bg-green-100 text-green-800",
} as const;

export const STATUS_CODE_COLORS = {
  success: "bg-green-100 text-green-800",
  error: "bg-red-100 text-red-800",
  unknown: "bg-gray-100 text-gray-800",
} as const;

export const ACCESSIBILITY_COLORS = {
  accessible: "bg-green-100 text-green-800",
  inaccessible: "bg-red-100 text-red-800",
} as const;
