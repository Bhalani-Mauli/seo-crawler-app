import React from "react";
import { APP_CONFIG } from "../../constants/index";

export const Header: React.FC = () => {
  return (
    <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
      <h1 className="text-3xl font-bold text-gray-900">{APP_CONFIG.NAME}</h1>
      <p className="text-gray-600 mt-2">Manage and monitor your SEO crawls</p>
    </div>
  );
};
