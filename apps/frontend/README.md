# Coreum Dex UI

This is the front-end UI for the Coreum Dex demo.

Note: this webapp is not fully responsive and is developed for a desktop screen of 1440px - 1024px. Some screen sizes may not look the best as of now.

There are two utility files in this project that are intended for you to use for your own Dex: `services/websocket.ts` and `services/api.ts`. See Websocket and API sections below.

## To Run

- `cd` into `/apps/frontend/` and run `./bin/start.sh`
- You might need to run `chmod +x ./bin/start.sh` first

### Start Parameters

//TODO better description

```
VITE_ENV_BASE_API: API server URL used in services/api.ts
VITE_ENV_WS: Websocket URL used in services/websocket.ts
VITE_ENV_MODE: Development flag for how the app imports envs when in development vs when in production
```

In this project, envs are imported in `services/general.ts` and `services/websocket.ts`.

If you are using a different framework or build tool, you will need to adjust how envs are named and imported in these files:

```
// for vite
const API_URL =
  import.meta.env.VITE_ENV_MODE === "development"
    ? import.meta.env.VITE_ENV_BASE_API
    : (window as any).COREUM.env.VITE_ENV_BASE_API;
```

## Websocket

This app uses primarily React and Zustand on the front-end. But for your application you may choose a different framework. Within `services/websocket.ts` is a TypeScript utility file that manages the websocket connection and is reusable in other frameworks/projects.

The websocket manager takes care of data handling and manages its own state.

You may choose to create a wrapper like a React hook for the websocket manager in your project. But below is an example of how you would connect and subscribe using only the websocket manager:

// TODO update examples after refactoring websocket

```
import { useEffect } from 'react';
import useWebSocketStore from './websocketStore';

const TestOrderBook: React.FC = () => {
  const { isConnected, messages, connect, subscribe } = useWebSocketStore();

  // memoize callback and subscription to prevent unnecessary subscriptions/unsubscriptions
  const handleOrderbookUpdate = useCallback(
    (message: WebSocketMessage) => {
      setOrderbook(message.Subscription?.Content);
    },
    [setOrderbook]
  );

  const subscription = useMemo(
    () => ({
      Network: Network.DEVNET,
      Method: Method.ORDERBOOK,
      ID: market.pair_symbol,
    }),
    [market.pair_symbol]
  );

  useWebSocket(subscription, handleOrderbookUpdate);

  return (
    <div>
      ...
    </div>
  );
};

export default TestOrderBook;
```

## API

This project also handles all backend API calls in `services/api.ts` and can be ported to your project as well.

## Other

You may want to refer to `/utils/index.ts` for some helper functions. `getAvgPriceFromOBbyVolume` in particular contains logic that calculates a volume weighted price based on the current orderbook.
