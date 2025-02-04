import { useEffect } from "react";
import {
  wsManager,
  Subscription,
  WebSocketMessage,
} from "../services/websocket";

const WEBSOCKET_URL =
  import.meta.env.VITE_MODE === "development"
    ? import.meta.env.VITE_ENV_WS
    : (window as any).SOLOGENIC.env.VITE_ENV_WS;

export const useWebSocket = (
  subscription: Subscription,
  callback: (message: WebSocketMessage) => void
) => {
  useEffect(() => {
    wsManager.connect(WEBSOCKET_URL);

    const handler = (message: WebSocketMessage) => {
      if (!message.Subscription) return;

      const isMatch =
        message.Subscription.Network === subscription.Network &&
        message.Subscription.Method === subscription.Method &&
        message.Subscription.ID === subscription.ID;

      if (isMatch) {
        try {
          const content = message.Subscription.Content
            ? JSON.parse(message.Subscription.Content)
            : null;

          callback({
            ...message,
            Subscription: { ...message.Subscription, Content: content },
          });
        } catch (error) {
          console.error(
            `Failed to process ${subscription.Method} message:`,
            error
          );
        }
      }
    };

    wsManager.subscribe(subscription, handler);

    return () => {
      wsManager.unsubscribe(subscription, handler);
    };
  }, [subscription, callback]);
};
