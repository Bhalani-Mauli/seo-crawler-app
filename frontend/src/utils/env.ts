export const validateEnv = () => {
  const requiredVars = ["VITE_API_BASE_URL", "VITE_API_KEY"];

  const missingVars = requiredVars.filter(
    (varName) => !import.meta.env[varName]
  );

  if (missingVars.length > 0) {
    console.warn(
      `Missing environment variables: ${missingVars.join(
        ", "
      )}. Using fallback values.`
    );
  }

  return {
    isValid: missingVars.length === 0,
    missingVars,
  };
};

export const getEnvVar = (key: string, fallback?: string): string => {
  return import.meta.env[key] || fallback || "";
};

export const getEnvVarAsNumber = (key: string, fallback: number): number => {
  const value = import.meta.env[key];
  if (!value) return fallback;

  const num = Number(value);
  return isNaN(num) ? fallback : num;
};

export const getEnvVarAsBoolean = (key: string, fallback: boolean): boolean => {
  const value = import.meta.env[key];
  if (!value) return fallback;

  return value.toLowerCase() === "true";
};
