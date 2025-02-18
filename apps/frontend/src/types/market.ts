// TODO - remove enums and types that can be replaced from backend protos
import { TimeInForce } from "coreum-js-nightly/dist/main/coreum/dex/v1/order";

export type Market = {
  base: Token;
  counter: Token;
  pair_symbol: string;
  reversed_pair_symbol: string;
};
export type Token = {
  Denom: {
    Currency: string;
    Issuer?: string;
    Precision: number;
    Denom: string;
    Name: string;
    Description: string;
  };
  SendCommission?: any;
  BurnRate?: any;
  InitialAmount?: any;
  Description: string;
  MetaData?: any;
};
export type CurrencyResponse = {
  Currencies: Token[];
};

// order actions
export enum OrderType {
  LIMIT = 1,
  MARKET = 2,
}

export enum TimeInForceString {
  goodTilCancel = "Good till Cancel",
  goodTilTime = "Good till Time",
  immediateOrCancel = "Immediate or Cancel",
  fillOrKill = "Fill or Kill",
}

export enum TimeInForceStringToEnum {
  "Good till Cancel" = TimeInForce.TIME_IN_FORCE_GTC,
  "Good till Time" = TimeInForce.TIME_IN_FORCE_GTC,
  "Immediate or Cancel" = TimeInForce.TIME_IN_FORCE_IOC,
  "Fill or Kill" = TimeInForce.TIME_IN_FORCE_FOK,
}

export enum TimeSelection {
  "5M" = "5 Mins",
  "15M" = "15 Mins",
  "30M" = "30 Mins",
  "1H" = "1 Hour",
  "6H" = "6 Hours",
  "12H" = "12 Hours",
  "1D" = "1 Day",
  CUSTOM = "Custom",
}

export enum SideBuy {
  BUY = 1,
  SELL = 2,
}

export type CreateOrderObject = {
  sender: string;
  type: OrderType;
  id?: number;
  base_denom: string;
  quote_denom: string;
  price: {
    exp: number;
    num: number;
  };
  quantity: number;
  side: SideBuy;
  good_til: {
    good_til_block_height: number;
    good_til_block_time: string;
  };
  timeInForce: TimeInForce;
};
export enum OrderType {
  BUY = "buy",
  SELL = "sell",
}
export type Order = {
  direction: OrderType;
  quantity: {
    currency: string;
    value: string;
    issuer?: string;
  };
  totalPrice: {
    currency: string;
    value: string;
    issuer?: string;
  };
  fee?: string;
  passive?: boolean;
  fillOrKill?: boolean;
  immediateOrCancel?: boolean;
  expirationTime?: string;
};
export enum TradeType {
  MARKET = "market",
  LIMIT = "limit",
}
export type OrderbookAction = {
  type: OrderType;
  price: number;
  volume: number;
};

// tradingview
export type BarPeriodParams = {
  from: number;
  to: number;
  countBack?: number;
  firstDataRequest?: boolean;
};
export type DataFeedAsset = {
  id: string;
  name: string;
};
export type ChartSubscription = {
  key: string;
  symbolInfo: BarSymbolInfo;
  resolution: string;
  onRealtimeCallback?: any;
};
export type BarSymbolInfo = {
  id: string;
  name: string;
  exchange?: string;
  session: string;
  timezone: string;
  has_intraday: boolean;
  has_weekly_and_monthly: boolean;
  supported_resolutions: string[];
  pricescale?: number;
  minmov: number;
  base_name?: string[];
  legs?: string[];
  full_name?: string;
  pro_name?: string;
  data_status?: string;
  ticker?: string;
};
export type OhlcRecord = [number, string, string, string, string, string];
export type OhlcResponse = OhlcRecord[];

export type ChartFeedBarsParams = {
  symbol: string;
  from: string;
  to: string;
  period: string;
};

// order history
export type FormattedOpenOrder = {
  side: string;
  price: string;
  volume: string;
  total: string;
  sequence: number;
};
export type Exchange = {
  id: string;
  txid: string;
  symbol: string;
  buyer: string;
  seller: string;
  is_seller_taker: boolean;
  amount: string;
  price: string;
  quote_amount: string;
  executed_at: string;
  time: string;
};

// ticker
export interface Ticker {
  OpenTime: number;
  CloseTime: number;
  OpenPrice: number;
  HighPrice: number;
  LowPrice: number;
  LastPrice: number;
  FirstPrice: number;
  Volume: number;
  Invertedvolume: number;
  Inverted: boolean;
}
export type TickerResponse = {
  Tickers: Record<string, Ticker>;
  USDTickers: Record<string, Ticker>;
};

// orderbook
export interface OrderbookRecord {
  Price: string;
  HumanReadablePrice: string;
  Amount: string;
  SymbolAmount: string;
  Sequence: number;
  Account?: string;
  OrderID: string;
}
export type OrderbookResponse = {
  Buy: OrderbookRecord[];
  Sell: OrderbookRecord[];
};

// exchange/trade history
export type TradeHistoryResponse = TradeRecord[];

export type TradeRecord = {
  Account: string;
  OrderID: string;
  Sequence: number;
  Amount: number;
  Price: number;
  HumanReadablePrice: string;
  SymbolAmount: string;
  Denom1: {
    Currency: string;
    Issuer: string;
    Precision: number;
    IsIBC: boolean;
    Denom: string;
    Name: string;
    Description: string;
    Icon: string;
  };
  Denom2: {
    Currency: string;
    Issuer: string;
    Precision: number;
    IsIBC: boolean;
    Denom: string;
    Name: string;
    Description: string;
    Icon: string;
  };
  Side: SideBuy;
  BlockTime: {
    seconds: number;
    nanos: number;
  };
  TradingFee: number;
  MetaData: {
    Network: string;
    CreatedAt: {
      seconds: number;
      nanos: number;
    };
    UpdatedAt: {
      seconds: number;
      nanos: number;
    };
  };
  TXID: string;
  BlockHeight: number;
  USD: number;
  FeeUSD: number;
  Status: OrderHistoryStatus;
};

export enum OrderHistoryStatus {
  // order_status_unspecified reserves the default value, to protect against unexpected settings.
  OrderStatus_ORDER_STATUS_UNSPECIFIED = 0,
  // order_status_open means that the order is open with any remaining quantity
  OrderStatus_ORDER_STATUS_OPEN = 1,
  // order_status_cancelled means the user has canceled the order.
  OrderStatus_ORDER_STATUS_CANCELED = 2,
  // order_status_filled means that the order is filled (quantity remaining is 0)
  OrderStatus_ORDER_STATUS_FILLED = 3,
  // order_status_expired means that the order is expired (e.g. a block event has passed the good til block height/time).
  OrderStatus_ORDER_STATUS_EXPIRED = 4,
}

export interface TransformedOrder {
  Side: number;
  Price: string;
  HumanReadablePrice: string;
  SymbolAmount: string;
  Amount: string;
  Total: number;
  Sequence: number;
  Account: string;
  OrderID: string;
}

export type CreateOrder = {
  Sender: string;
  Type: OrderType;
  OrderID: number;
  BaseDenom: string;
  QuoteDenom: string;
  Price: string;
  Quantity: string;
  Side: SideBuy;
  GoodTil: {
    GoodTilBlockHeight: number;
    GoodTilBlockTime: string;
  };
  TimeInForce: TimeInForce;
};

export type WalletAsset = {
  Amount: string;
  Denom: string;
  SymbolAmount: string;
};

export type WalletBalances = WalletAsset[];

export type CancelOrderResponse = {
  TXBytes: string;
};
