import { Market } from "@/types/market";
import { CoreumNetwork } from "coreum-js-nightly";

// default markets
export const DEVNET_DEFAULT: Market = {
  base: {
    Denom: {
      Currency: "nor",
      Issuer: "devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
      Precision: 6,
      Denom: "nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
      Name: "NOR",
      Description: "NOR",
      IsIBC: false,
    },
    Description: "NOR",
  },
  counter: {
    Denom: {
      Currency: "alb",
      Issuer: "devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
      Precision: 6,
      Denom: "alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
      Name: "ALB",
      Description: "ALB",
      IsIBC: false,
    },
    Description: "ALB",
  },
  pair_symbol:
    "nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43_alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
  reversed_pair_symbol:
    "alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43_nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
};

export const TESTNET_DEFAULT: Market = {
  base: {
    Denom: {
      Currency: "nor",
      Issuer: "testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57",
      Precision: 6,
      Denom: "nor-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57",
      Name: "NOR",
      Description: "NOR",
      IsIBC: false,
    },
    Description: "NOR",
  },
  counter: {
    Denom: {
      Currency: "alb",
      Issuer: "testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57",
      Precision: 6,
      Denom: "alb-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57",
      Name: "ALB",
      Description: "ALB",
      IsIBC: false,
    },
    Description: "ALB",
  },
  pair_symbol:
    "nor-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57_alb-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57",
  reversed_pair_symbol:
    "alb-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57_nor-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57",
};

export const PREDEFINED_MARKETS: Record<string, Market> = {
  [DEVNET_DEFAULT.pair_symbol]: DEVNET_DEFAULT,
  [TESTNET_DEFAULT.pair_symbol]: TESTNET_DEFAULT,
};

export const GLOBAL_FALLBACK_MARKET = DEVNET_DEFAULT;

// function to get the default market based on the current network
// reads from the environment variable VITE_DEFAULT_MARKET_CONFIGS
// fallsback to devnet market if not found
export const getDefaultMarket = (currentNetwork: CoreumNetwork): Market => {
  const defaultConfigString = import.meta.env.VITE_DEFAULT_MARKET_CONFIGS as
    | string
    | undefined;

  if (!defaultConfigString) {
    console.warn(
      "VITE_DEFAULT_MARKET_CONFIGS is not set. Using global fallback market."
    );
    return GLOBAL_FALLBACK_MARKET;
  }

  try {
    const networkToMarketSymbolMap = JSON.parse(defaultConfigString) as Record<
      string,
      string
    >;
    const marketSymbol = networkToMarketSymbolMap[currentNetwork];

    if (!marketSymbol) {
      console.warn(
        `No default market symbol found for network "${currentNetwork}" in VITE_DEFAULT_MARKET_CONFIGS. Using global fallback.`
      );
      return GLOBAL_FALLBACK_MARKET;
    }

    const marketConfig = PREDEFINED_MARKETS[marketSymbol];
    if (!marketConfig) {
      console.warn(
        `Market configuration for symbol "${marketSymbol}" (for network "${currentNetwork}") not found in PREDEFINED_MARKETS. Using global fallback.`
      );
      return GLOBAL_FALLBACK_MARKET;
    }

    return marketConfig;
  } catch (error) {
    console.error(
      "Error parsing VITE_DEFAULT_MARKET_CONFIGS or looking up market. Using global fallback.",
      error
    );
    return GLOBAL_FALLBACK_MARKET;
  }
};
