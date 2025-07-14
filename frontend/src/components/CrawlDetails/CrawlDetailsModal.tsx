import React from "react";
import { Modal } from "../ui/Modal";
import { CrawlDetailsHeader } from "./CrawlDetailsHeader";
import { CrawlDetailsContent } from "./CrawlDetailsContent";
import type { URLData, LinkData, HeadingData } from "../../types";

interface CrawlDetailsModalProps {
  isOpen: boolean;
  onClose: () => void;
  detailMeta: URLData | null;
  detailLinks: LinkData[];
  detailHeadings: HeadingData[];
  detailLoading: boolean;
}

export const CrawlDetailsModal: React.FC<CrawlDetailsModalProps> = ({
  isOpen,
  onClose,
  detailMeta,
  detailLinks,
  detailHeadings,
  detailLoading,
}) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose} size="xl">
      <CrawlDetailsHeader detailMeta={detailMeta} onClose={onClose} />
      <CrawlDetailsContent
        detailLinks={detailLinks}
        detailHeadings={detailHeadings}
        detailLoading={detailLoading}
      />
    </Modal>
  );
};
