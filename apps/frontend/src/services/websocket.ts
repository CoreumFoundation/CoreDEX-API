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

export interface WebSocketMessage {
  Action: Action;
  Subscription?: Subscription;
}

export type MessageHandler = (data: any) => void;

export type UpdateFunction = (prevState: any, newContent: any) => any;

interface SubscriptionConfig {
  subscription: Subscription;
  handlers: MessageHandler[];
  updateFn?: UpdateFunction;
}

export enum UpdateStrategy {
  REPLACE = "REPLACE", // replaces state
  MERGE = "MERGE", // adds new data to beginning of state
  APPEND = "APPEND", // adds new data to end of state
}

const updateFunctions: Record<UpdateStrategy, UpdateFunction> = {
  [UpdateStrategy.REPLACE]: (_prev: any, newContent: any) => newContent,
  [UpdateStrategy.MERGE]: (prev: any, newContent: any) => {
    const prevArr = Array.isArray(prev) ? prev : [];
    if (Array.isArray(newContent) && newContent.length === 0) {
      return prevArr;
    }
    if (Array.isArray(newContent)) {
      return [...newContent, ...prevArr];
    }
    return [newContent, ...prevArr];
  },
  [UpdateStrategy.APPEND]: (prev: any, newContent: any) => {
    const prevArr = Array.isArray(prev) ? prev : [];
    if (Array.isArray(newContent) && newContent.length === 0) {
      return prevArr;
    }
    if (Array.isArray(newContent)) {
      return [...prevArr, ...newContent];
    }
    return [...prevArr, newContent];
  },
};

class WebSocketManager {
  private static instance: WebSocketManager;
  private ws: WebSocket | null = null;
  private isConnected = false;
  private connectedPromise: Promise<void>;
  private connectedResolver!: () => void;
  private pendingSubscriptions: SubscriptionConfig[] = [];
  private pendingUnsubscriptions: Subscription[] = [];
  private subscriptions: Map<string, SubscriptionConfig> = new Map();
  private stateStore: Map<string, any> = new Map();

  private constructor() {
    this.connectedPromise = new Promise((resolve) => {
      this.connectedResolver = resolve;
    });
  }

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

    if (this.connectedResolver) {
      this.connectedResolver();
    }

    this.pendingSubscriptions.forEach((config) => {
      console.log(
        "Sending pending subscription for key:",
        this.getSubscriptionKey(config.subscription)
      );
      this.sendSubscription(config.subscription);
      this.subscriptions.set(
        this.getSubscriptionKey(config.subscription),
        config
      );
    });
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

  /**
   * Subscribe with a subscription object, a message handler, and optionally an update strategy.
   *
   * The update strategy determines how new messages modify the current state.
   * For example, an orderbook update may simply replace the state,
   * while an OHLC subscription may append new data to the existing state.
   */
  public subscribe(
    subscription: Subscription,
    handler: MessageHandler,
    strategy: UpdateStrategy
  ) {
    const key = this.getSubscriptionKey(subscription);
    let config = this.subscriptions.get(key);
    const updateFn = updateFunctions[strategy];

    if (config) {
      config.handlers.push(handler);
    } else {
      config = { subscription, handlers: [handler], updateFn };
      this.subscriptions.set(key, config);
    }

    if (this.isConnected) {
      this.sendSubscription(subscription);
    } else {
      this.pendingSubscriptions.push(config);
    }
  }

  public unsubscribe(subscription: Subscription, handler: MessageHandler) {
    const key = this.getSubscriptionKey(subscription);
    const config = this.subscriptions.get(key);
    if (config) {
      config.handlers = config.handlers.filter((h) => h !== handler);
      if (config.handlers.length === 0) {
        this.subscriptions.delete(key);
        if (this.isConnected) {
          this.sendUnsubscription(subscription);
        } else {
          this.pendingUnsubscriptions.push(subscription);
        }
        this.stateStore.delete(key);
      }
    }
  }

  private handleMessage(data: string) {
    try {
      if (data === "Connected") return;
      const message: WebSocketMessage = JSON.parse(data);
      if (!message.Subscription) return;

      const key = this.getSubscriptionKey(message.Subscription);

      const config = this.subscriptions.get(key);
      if (!config) return;
      const newContent = JSON.parse(message.Subscription.Content);
      const prevState = this.stateStore.get(key) || [];
      const newState = config.updateFn
        ? config.updateFn(prevState, newContent)
        : newContent;
      this.stateStore.set(key, newState);

      config.handlers.forEach((handler) => handler(newState));
    } catch (error) {
      console.error("Error handling message:", error, "Raw data:", data);
    }
  }

  private getSubscriptionKey(sub: Subscription): string {
    return `${sub.Network}-${sub.Method}-${sub.ID}`;
  }

  public setInitialState(subscription: Subscription, initialData: any) {
    const key = this.getSubscriptionKey(subscription);
    this.stateStore.set(key, initialData);
  }

  public close() {
    this.ws?.close();
    this.subscriptions.clear();
    this.pendingSubscriptions = [];
    this.pendingUnsubscriptions = [];
    this.stateStore.clear();
  }

  public connected(): Promise<void> {
    return this.connectedPromise;
  }
}

export const wsManager = WebSocketManager.getInstance();
