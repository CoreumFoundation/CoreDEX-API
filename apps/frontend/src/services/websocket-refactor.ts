// websocket.ts
import { CoreumNetwork } from "coreum-js-nightly";

export enum Action {
  SUBSCRIBE = 0,
  UNSUBSCRIBE = 1,
  CLOSE = 2,
  RESPONSE = 3,
}

export enum Method {
  METHOD_DO_NOT_USE = 0,
  TRADES_FOR_SYMBOL = 1,
  TRADES_FOR_ACCOUNT = 2,
  TRADES_FOR_ACCOUNT_AND_SYMBOL = 3,
  OHLC = 4,
  TICKER = 5,
  ORDERBOOK = 6,
  ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT = 7,
}

export enum Network {
  NETWORK_DO_NOT_USE = 0,
  MAINNET = 1,
  TESTNET = 2,
  DEVNET = 3,
}

export const NetworkToEnum = (network: CoreumNetwork): Network => {
  switch (network) {
    case CoreumNetwork.MAINNET:
      return Network.MAINNET;
    case CoreumNetwork.TESTNET:
      return Network.TESTNET;
    case CoreumNetwork.DEVNET:
      return Network.DEVNET;
  }
};

export interface Subscription {
  Network: Network;
  Method: Method;
  ID: string;
  Content?: any;
}

/**
 * The message that comes over the WebSocket.
 */
export interface WebSocketMessage {
  Action: Action;
  Subscription?: Subscription;
}

/**
 * A handler that is called when a message for a subscription is received.
 * The `data` parameter is the updated state (after applying the update logic).
 */
export type MessageHandler = (data: any) => void;

/**
 * A function that defines how to update state.
 * It receives the current state (which can be null on first call)
 * and the new message content, then returns the new state.
 */
export type UpdateFunction = (prevState: any, newContent: any) => any;

interface SubscriptionConfig {
  subscription: Subscription;
  handlers: MessageHandler[];
  updateFn?: UpdateFunction; // If provided, use this to merge new messages into the state.
}

class WebSocketManager {
  private static instance: WebSocketManager;
  private ws: WebSocket | null = null;
  private isConnected = false;
  private pendingSubscriptions: SubscriptionConfig[] = [];
  private pendingUnsubscriptions: Subscription[] = [];

  // A mapping from subscription key to a subscription config (handlers and update function)
  private subscriptions: Map<string, SubscriptionConfig> = new Map();
  // Internal state per subscription key.
  private stateStore: Map<string, any> = new Map();

  private constructor() {}

  public static getInstance(): WebSocketManager {
    if (!WebSocketManager.instance) {
      WebSocketManager.instance = new WebSocketManager();
    }
    return WebSocketManager.instance;
  }

  public connect(url: string) {
    if (this.ws) return;

    this.ws = new WebSocket(url);
    this.ws.onopen = () => this.handleConnectionOpen();
    this.ws.onmessage = (event) => this.handleMessage(event.data);
    this.ws.onclose = () => {
      this.isConnected = false;
      this.ws = null;
      console.log("WebSocket disconnected");
    };
  }

  private handleConnectionOpen() {
    this.isConnected = true;
    console.log("WebSocket connected");

    // Send all pending subscriptions.
    this.pendingSubscriptions.forEach((config) => {
      this.sendSubscription(config.subscription);
      // Also store the subscription in our active registry.
      const key = this.getSubscriptionKey(config.subscription);
      this.subscriptions.set(key, config);
    });
    this.pendingSubscriptions = [];

    // Process pending unsubscriptions.
    this.pendingUnsubscriptions.forEach((sub) => this.sendUnsubscription(sub));
    this.pendingUnsubscriptions = [];
  }

  private sendSubscription(subscription: Subscription) {
    this.ws?.send(
      JSON.stringify({
        Action: Action.SUBSCRIBE,
        Subscription: subscription,
      })
    );
  }

  private sendUnsubscription(subscription: Subscription) {
    this.ws?.send(
      JSON.stringify({
        Action: Action.UNSUBSCRIBE,
        Subscription: subscription,
      })
    );
  }

  /**
   * Subscribe with a subscription object, a message handler, and optionally an update function.
   *
   * The update function determines how new messages modify the current state.
   * For example, an orderbook update may simply replace the state,
   * while an OHLC subscription may append new data to the existing state.
   */
  public subscribe(
    subscription: Subscription,
    handler: MessageHandler,
    updateFn?: UpdateFunction
  ) {
    const key = this.getSubscriptionKey(subscription);
    let config = this.subscriptions.get(key);
    if (config) {
      // If a config already exists, simply add the handler.
      config.handlers.push(handler);
    } else {
      // Otherwise, create a new config.
      config = { subscription, handlers: [handler], updateFn };
      this.subscriptions.set(key, config);
    }

    if (this.isConnected) {
      this.sendSubscription(subscription);
    } else {
      this.pendingSubscriptions.push(config);
    }
  }

  /**
   * Unsubscribe a particular handler for the given subscription.
   */
  public unsubscribe(subscription: Subscription, handler: MessageHandler) {
    const key = this.getSubscriptionKey(subscription);
    const config = this.subscriptions.get(key);
    if (config) {
      config.handlers = config.handlers.filter((h) => h !== handler);
      // If no more handlers are registered for this subscription, remove it completely.
      if (config.handlers.length === 0) {
        this.subscriptions.delete(key);
        if (this.isConnected) {
          this.sendUnsubscription(subscription);
        } else {
          this.pendingUnsubscriptions.push(subscription);
        }
        // Optionally remove the state.
        this.stateStore.delete(key);
      }
    }
  }

  /**
   * Handle an incoming message. If it includes a subscription,
   * update the internal state (using the update function if provided) and notify the handlers.
   */
  private handleMessage(data: string) {
    try {
      if (data === "Connected") return;
      const message: WebSocketMessage = JSON.parse(data);
      if (!message.Subscription) return;

      const key = this.getSubscriptionKey(message.Subscription);
      const config = this.subscriptions.get(key);
      if (!config) return;

      // Get the new content from the incoming message.
      const newContent = message.Subscription.Content;

      // Update state if an update function is provided, otherwise simply replace.
      const prevState = this.stateStore.get(key) || null;
      const newState = config.updateFn
        ? config.updateFn(prevState, newContent)
        : newContent;
      this.stateStore.set(key, newState);

      // Notify all handlers with the updated state.
      config.handlers.forEach((handler) => handler(newState));
    } catch (error) {
      console.error("Error handling message:", error, "Raw data:", data);
    }
  }

  private getSubscriptionKey(sub: Subscription): string {
    return `${sub.Network}-${sub.Method}-${sub.ID}`;
  }

  public close() {
    this.ws?.close();
    this.subscriptions.clear();
    this.pendingSubscriptions = [];
    this.pendingUnsubscriptions = [];
    this.stateStore.clear();
  }
}

export const wsManager = WebSocketManager.getInstance();
