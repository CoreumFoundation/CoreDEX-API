import { CoreumNetwork } from "coreum-js";
import { Action, Subscription as Sub } from "coredex-api-types/update";
import { Network } from "coredex-api-types/metadata";

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

export type Subscription = Omit<Sub, "Content"> & {
  Content?: any;
};

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
  public stateStore: Map<string, any> = new Map();

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
      if (!("Content" in message.Subscription)) return;

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

  public getSubscriptionKey(sub: Subscription): string {
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

  public clearState() {
    this.stateStore.clear();
    this.subscriptions.clear();
    this.pendingSubscriptions = [];
    this.pendingUnsubscriptions = [];
  }

  public getState(): Map<string, any> {
    return this.stateStore;
  }
}

export const wsManager = WebSocketManager.getInstance();
