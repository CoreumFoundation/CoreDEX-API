# Coreum Dex UI

This is the front-end UI for the Coreum Dex demo.

Note: this webapp is not fully responsive and is developed for a desktop screen of 1440px - 1024px. Some screen sizes may not look the best as of now.

There are two utility files in this project that are intended for you to use for your own Dex: `services/websocket.ts` and `services/api.ts`. See Websocket and API sections below.

## To Run

- `cd` into `/apps/frontend/` and run `./bin/start.sh`
- You might need to run `chmod +x ./bin/start.sh` first

### Start Parameters

```
VITE_ENV_BASE_API: API server URL used in services/api.ts
VITE_ENV_WS: Websocket URL used in services/websocket.ts
VITE_ENV_MODE: Development flag for how the app imports envs when in development vs when in production
VITE_ENV_DEFAULT_MARKET_CONFIGS: object that stores the pair id for market depending on selected network {"devnet":"123..", "testnet": "abc..."}
```

In this project, envs are imported in `/config/envs.ts`.

If you are using a different framework or build tool, you will need to adjust how envs are named and imported in these files:

```
// for vite
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
```

## Websocket

This app uses primarily React and Zustand on the front-end. But for your application you may choose a different framework. Within `services/websocket.ts` is a TypeScript utility file that manages the websocket connection and is reusable in other frameworks/projects.

The websocket manager takes care of data handling and manages its own state.

You may choose to create a wrapper like a React hook for the websocket manager in your project. But below is an example of how you would connect and subscribe using only the websocket manager:

```
import { useEffect, useMemo } from 'react';
import {
  Method,
  NetworkToEnum,
  UpdateStrategy,
  wsManager,
} from "@/services/websocket";

const Test: React.FC = () => {
  const walletSubscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.WALLET,
      ID: `${wallet ? wallet.address : ""}`,
    }),
    [market.pair_symbol, wallet]
  );

  const handleWalletUpdate = (message: WalletBalances) => {
    if (message.length > 0) {
      setWalletBalances(message);
    }
  };

  useEffect(() => {
    wsManager.connected().then(() => {
      wsManager.subscribe(
        walletSubscription,
        handleWalletUpdate,
        UpdateStrategy.REPLACE
      );
    });
    return () => {
      wsManager.unsubscribe(walletSubscription, setWalletBalances);
    };
  }, [walletSubscription]);
  return (
    <div>
      ...
    </div>
  );
};

export default Test;
```

- you can pass a custom handler function or just a setter. The above example uses `handleWalletUpdate`
- there are a few basic update strategies like replace, merge, etc that describe how data should be added to the ws state (append, replace, etc)

## API

This project also handles all backend API calls in `services/api.ts` and can be ported to your project as well.

## Other

You may want to refer to `/utils/index.ts` for some helper functions. `getAvgPriceFromOBbyVolume` in particular contains logic that calculates a volume weighted price based on the current orderbook.
