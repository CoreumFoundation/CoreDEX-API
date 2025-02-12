import { useEffect } from "react";
import {
  wsManager,
  Subscription,
  WebSocketMessage,
} from "../services/websocket";
import { WS_URL } from "@/config/envs";

export const useWebSocket = (
  subscription: Subscription,
  callback: (message: WebSocketMessage) => void
) => {
  useEffect(() => {
    wsManager.connect(WS_URL);

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
