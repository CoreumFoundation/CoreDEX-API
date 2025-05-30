interface CoreumEnv {
  VITE_ENV_BASE_API: string;
  VITE_ENV_WS: string;
  VITE_ENV_DEFAULT_MARKET_CONFIGS: string;
}

const env: CoreumEnv =
  import.meta.env.VITE_ENV_MODE === "development"
    ? {
        VITE_ENV_BASE_API: import.meta.env.VITE_ENV_BASE_API,
        VITE_ENV_WS: import.meta.env.VITE_ENV_WS,
        VITE_ENV_DEFAULT_MARKET_CONFIGS: import.meta.env
          .VITE_ENV_DEFAULT_MARKET_CONFIGS,
      }
    : (window as any).COREUM?.env;

export const BASE_API_URL = env.VITE_ENV_BASE_API;
export const WS_URL = env.VITE_ENV_WS;
export const DEFAULT_MARKET_CONFIGS = JSON.parse(
  env.VITE_ENV_DEFAULT_MARKET_CONFIGS
) as Record<string, any>;
