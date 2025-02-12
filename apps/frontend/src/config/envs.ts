interface CoreumEnv {
  VITE_ENV_BASE_API: string;
  VITE_ENV_WS: string;
}

const env: CoreumEnv =
  import.meta.env.VITE_ENV_MODE === "development"
    ? {
        VITE_ENV_BASE_API: import.meta.env.VITE_ENV_BASE_API,
        VITE_ENV_WS: import.meta.env.VITE_ENV_WS,
      }
    : (window as any).COREUM?.env;

export const BASE_API_URL = env.VITE_ENV_BASE_API;
export const WS_URL = env.VITE_ENV_WS;
