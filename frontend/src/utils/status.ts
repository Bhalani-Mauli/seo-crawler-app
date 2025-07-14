import type { CrawlStatus } from "../types";
import {
  STATUS_COLORS,
  STATUS_CODE_COLORS,
  ACCESSIBILITY_COLORS,
} from "../constants";

export const getStatusColor = (status: CrawlStatus): string => {
  return STATUS_COLORS[status] || STATUS_COLORS.stopped;
};

export const getStatusCodeColor = (statusCode: number | null): string => {
  if (!statusCode) return STATUS_CODE_COLORS.unknown;

  if (statusCode >= 200 && statusCode < 400) {
    return STATUS_CODE_COLORS.success;
  }

  return STATUS_CODE_COLORS.error;
};

export const getAccessibilityColor = (isAccessible: boolean): string => {
  return isAccessible
    ? ACCESSIBILITY_COLORS.accessible
    : ACCESSIBILITY_COLORS.inaccessible;
};

export const getAccessibilityText = (isAccessible: boolean): string => {
  return isAccessible ? "Yes" : "No";
};
