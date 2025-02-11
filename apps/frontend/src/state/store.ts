import { create } from "zustand";
import { Client, CoreumNetwork } from "coreum-js-nightly";
import {
  Market,
  OrderbookResponse,
  Ticker,
  Token,
  TradeHistoryResponse,
  TransformedOrder,
} from "@/types/market";
import { ToasterProps } from "@/types";
import { toast } from "react-toastify";

export type State = {
  fetching: boolean;
  wallet: any;
  setWallet: (wallet: any) => Promise<void>;
  account: string;
  setAccount: (accountID: string) => void;
  network: CoreumNetwork;
  setNetwork: (network: CoreumNetwork) => void;
  coreum: Client | null;
  setCoreum: (client: Client | null) => Promise<void>;
  pushNotification: (object: ToasterProps) => void;

  // market
  market: Market;
  setMarket: (market: Market) => void;
  currencies: any;
  setCurrencies: (currencies: Token[] | null) => void;
  orderbook: OrderbookResponse | null;
  setOrderbook: (orderbook: OrderbookResponse | null) => void;
  openOrders: TransformedOrder[] | null;
  setOpenOrders: (openOrders: TransformedOrder[] | null) => void;
  loginModal: boolean;
  setLoginModal: (loginModal: boolean) => void;
  tickers: Ticker | null;
  setTickers: (tickers: Ticker | null) => void;
  chartPeriod: string;
  setChartPeriod: (period: string) => void;
  exchangeHistory: TradeHistoryResponse | null;
  setExchangeHistory: (exchangeHistory: TradeHistoryResponse | null) => void;
  orderHistory: TradeHistoryResponse | null;
  setOrderHistory: (orderHistory: TradeHistoryResponse | null) => void;
};

export const useStore = create<State>((set) => ({
  fetching: false,
  network: sessionStorage.network || CoreumNetwork.DEVNET,
  setNetwork: (network: CoreumNetwork) => {
    if (
      sessionStorage.network &&
      sessionStorage.network !== network &&
      localStorage.token
    ) {
      localStorage.removeItem("token");
    }
    sessionStorage.network = network;
    set(() => ({ network }));
  },
  wallet: null,
  setWallet: async (wallet: any) => {
    set(() => ({
      wallet,
    }));
  },
  account: "",
  setAccount: (accountID: string) => {
    set({ account: accountID });
  },
  coreum: null,
  setCoreum: async (client: Client | null) => {
    set({ coreum: client });
  },

  pushNotification: ({ type, message }) => {
    if (type === "success") {
      toast.success(message);
    } else if (type === "error") {
      toast.error(message);
    } else {
      toast.warning(message);
    }
  },

  market: {
    base: {
      Denom: {
        Currency: "dextestdenom0",
        Issuer: "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
        Precision: 6,
        Denom: "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
        Name: "DexTestDenom0",
        Description: "Dex Test Denom",
      },
      Description: "Dex Test Denom",
    },
    counter: {
      Denom: {
        Currency: "dextestdenom1",
        Issuer: "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
        Precision: 6,
        Denom: "dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
        Name: "DexTestDenom1",
        Description: "Dex Test Denom",
      },

      Description: "Dex Test Denom",
    },
    pair_symbol:
      "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
    reversed_pair_symbol:
      "dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygsdextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
  },
  setMarket: (market: Market) => {
    set({ market: market });
  },
  currencies: null,
  setCurrencies: (currencies: Token[] | null) => {
    set({ currencies: currencies });
  },
  tickers: null,
  setTickers: (tickers: Ticker | null) => {
    set({ tickers: tickers });
  },
  orderbook: null,
  setOrderbook: (orderbook: OrderbookResponse | null) => {
    set({ orderbook: orderbook });
  },
  openOrders: null,
  setOpenOrders: (openOrders: TransformedOrder[] | null) => {
    set({ openOrders: openOrders });
  },
  orderHistory: null,
  setOrderHistory: (orderHistory: TradeHistoryResponse | null) => {
    set({ orderHistory: orderHistory });
  },
  loginModal: false,
  setLoginModal: (loginModal: boolean) => {
    set({ loginModal: loginModal });
  },
  chartPeriod: "1W",
  setChartPeriod: (period: string) => {
    set({ chartPeriod: period });
  },
  exchangeHistory: null,
  setExchangeHistory: (exchangeHistory: TradeHistoryResponse | null) => {
    set({ exchangeHistory: exchangeHistory });
  },
}));
