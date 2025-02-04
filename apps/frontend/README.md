# Coreum Dex UI

This is the front-end UI for the Coreum Dex demo.

Note: this webapp is not fully responsive and is developed for a desktop screen of 1440px - 1024px. Some screen sizes may not look the best as of now.

## To Run

- `cd` into `/apps/frontend/` and run `./bin/start.sh`
- You might need to run `chmod +x ./bin/start.sh` first

## Websocket

This app uses primarily React and Zustand on the front-end. But for your application you may choose a different framework. Within `/src/services/websocket.ts` is a TypeScript utility file that manages the websocket connection and should be reusable in other frameworks.

In this demo app, we use the websocket manager inside of a React Hook (see `hooks/websocket.tsx) to make things easier.

An alternative would be to use the websocket manager in a state manager of your choice and handle the connections within your components (eg.`state/websocket.ts`)

Below is an example of how you might use the hook to connect, subscribe and listen to messages.

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

## Other

You may want to refer to `/utils/index.ts` for some helper functions. `getAvgPriceFromOBbyVolume` in particular contains logic that calculates a volume weighted price based on the current orderbook.
