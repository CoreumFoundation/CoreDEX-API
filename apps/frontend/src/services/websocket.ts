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

export const NetworkToEnum = (network: CoreumNetwork) => {
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

export interface WebSocketMessage {
  Action: Action;
  Subscription?: Subscription;
}

type MessageHandler = (message: WebSocketMessage) => void;

class WebSocketManager {
  private static instance: WebSocketManager;
  private ws: WebSocket | null = null;
  private handlers: Map<string, MessageHandler[]> = new Map();
  private pendingSubscriptions: Subscription[] = [];
  private pendingUnsubscriptions: Subscription[] = [];
  private isConnected = false;

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

    this.pendingUnsubscriptions.forEach((sub) => this.sendUnsubscription(sub));
    this.pendingUnsubscriptions = [];

    this.pendingSubscriptions.forEach((sub) => this.sendSubscription(sub));
    this.pendingSubscriptions = [];
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

  public subscribe(subscription: Subscription, handler: MessageHandler) {
    const key = this.getSubscriptionKey(subscription);

    if (!this.handlers.has(key)) {
      this.handlers.set(key, []);
    }
    this.handlers.get(key)?.push(handler);

    if (this.isConnected) {
      this.sendSubscription(subscription);
    } else {
      this.pendingSubscriptions.push(subscription);
    }
  }

  public unsubscribe(subscription: Subscription, handler: MessageHandler) {
    const key = this.getSubscriptionKey(subscription);
    const handlers = this.handlers.get(key) || [];
    this.handlers.set(
      key,
      handlers.filter((h) => h !== handler)
    );

    if (this.isConnected) {
      this.sendUnsubscription(subscription);
    } else {
      this.pendingUnsubscriptions.push(subscription);
    }
  }

  private handleMessage(data: string) {
    try {
      if (data === "Connected") return;

      const message: WebSocketMessage = JSON.parse(data);
      this.notifyHandlers(message);
    } catch (error) {
      console.error("Error handling message:", error, "Raw data:", data);
    }
  }

  private notifyHandlers(message: WebSocketMessage) {
    if (!message.Subscription) return;

    const key = this.getSubscriptionKey(message.Subscription);
    this.handlers.get(key)?.forEach((handler) => handler(message));
  }

  private getSubscriptionKey(sub: Subscription): string {
    return `${sub.Network}-${sub.Method}-${sub.ID}`;
  }

  public close() {
    this.ws?.close();
    this.handlers.clear();
    this.pendingSubscriptions = [];
    this.pendingUnsubscriptions = [];
  }
}

export const wsManager = WebSocketManager.getInstance();
