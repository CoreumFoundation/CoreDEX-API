import { ExtensionWallets } from "coreum-js";
import { Side, TimeInForce } from "coreum-js/dist/main/coreum/dex/v1/order";
import { Trade } from "coredex-api-types/trade";
import { OrderStatus } from "coredex-api-types/order";
import { Denom } from "coredex-api-types/denom";

export type ICoreumWallet = {
  address: string;
  method: ExtensionWallets;
};

export type Market = {
  base: Token;
  counter: Token;
  pair_symbol: string;
  reversed_pair_symbol: string;
};

export type Token = {
  Denom: Denom;
  SendCommission?: any;
  BurnRate?: any;
  InitialAmount?: any;
  Description: string;
  MetaData?: any;
};

export type CurrencyResponse = {
  Currencies: Token[];
};

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

export type OrderbookAction = {
  type: Side;
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
  InvertedVolume: number;
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
  RemainingAmount: string;
  RemainingSymbolAmount: string;
}
export type OrderbookResponse = {
  Buy: OrderbookRecord[];
  Sell: OrderbookRecord[];
};

// exchange/trade history
export type TradeHistoryResponse = TradeRecord[];
export type TradeRecord = Trade & {
  HumanReadablePrice: string;
  SymbolAmount: string;
  Status: OrderStatus;
  BlockTime: {
    seconds: number;
    nanos: number;
  };
};

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
  RemainingAmount: string;
  RemainingSymbolAmount: string;
}

export type WalletAsset = {
  Amount: string;
  Denom: string;
  SymbolAmount: string;
};

export type WalletBalances = WalletAsset[];

export type CancelOrderResponse = {
  TXBytes: string;
};

export type MarketData = {
  Denom1: {
    Currency: string;
    Issuer: string;
    Precision: number;
    Denom: string;
  };
  Denom2: {
    Currency: string;
    Issuer: string;
    Precision: number;
    Denom: string;
  };
  MetaData: {
    Network: number;
    UpdatedAt: {
      seconds: number;
      nanos: number;
    };
    CreatedAt: {
      seconds: number;
      nanos: number;
    };
  };
  PriceTick: {
    Value: number;
    Exp: number;
  };
  QuantityStep: number;
};
