import { create } from "zustand";
import { Client, CoreumNetwork } from "coreum-js-nightly";
import {
  ICoreumWallet,
  Market,
  OrderbookResponse,
  TickerResponse,
  TradeHistoryResponse,
  TransformedOrder,
} from "@/types/market";
import { ToasterProps } from "@/types";
import { toast } from "react-toastify";

const defaultMarket: Market = {
  base: {
    Denom: {
      Currency: "nor",
      Issuer: "devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
      Precision: 6,
      Denom: "nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
      Name: "NOR",
      Description: "NOR",
      IsIBC: false,
    },
    Description: "NOR",
  },
  counter: {
    Denom: {
      Currency: "alb",
      Issuer: "devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
      Precision: 6,
      Denom: "alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
      Name: "ALB",
      Description: "ALB",
      IsIBC: false,
    },
    Description: "ALB",
  },
  pair_symbol:
    "nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43_alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
  reversed_pair_symbol:
    "alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43_nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
};

const storedMarketString = localStorage.getItem("market");
const storedMarket: Market = storedMarketString
  ? JSON.parse(storedMarketString)
  : defaultMarket;

export type State = {
  fetching: boolean;
  wallet: any;
  setWallet: (wallet: any) => Promise<void>;
  network: CoreumNetwork;
  setNetwork: (network: CoreumNetwork) => void;
  coreum: Client | null;
  setCoreum: (client: Client | null) => Promise<void>;
  pushNotification: (object: ToasterProps) => void;
  isLoading: boolean;
  setIsLoading: (isLoading: boolean) => void;
  market: Market;
  setMarket: (market: Market) => void;
  currencies: any;
  setCurrencies: (currencies: any) => void;
  orderbook: OrderbookResponse | null;
  setOrderbook: (orderbook: OrderbookResponse | null) => void;
  openOrders: TransformedOrder[] | null;
  setOpenOrders: (openOrders: TransformedOrder[] | null) => void;
  loginModal: boolean;
  setLoginModal: (loginModal: boolean) => void;
  tickers: TickerResponse | null;
  setTickers: (tickers: TickerResponse | null) => void;
  chartPeriod: string;
  setChartPeriod: (period: string) => void;
  exchangeHistory: TradeHistoryResponse | [];
  setExchangeHistory: (
    exchangeHistory:
      | TradeHistoryResponse
      | []
      | ((prev: TradeHistoryResponse | []) => TradeHistoryResponse | [])
  ) => void;
  orderHistory: TradeHistoryResponse | [];
  setOrderHistory: (
    orderHistory:
      | TradeHistoryResponse
      | []
      | ((prev: TradeHistoryResponse | []) => TradeHistoryResponse | [])
  ) => void;
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
  setWallet: async (wallet: ICoreumWallet) => {
    localStorage.wallet = JSON.stringify({
      address: wallet.address,
      method: wallet.method,
    });
    set(() => ({
      wallet,
    }));
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
  isLoading: false,
  setIsLoading: (isLoading: boolean) => {
    set({ isLoading });
  },

  market: storedMarket,
  setMarket: (market: Market) => {
    localStorage.setItem("market", JSON.stringify(market));
    set({ market });
  },
  currencies: null,
  setCurrencies: (currencies: any) => {
    set({ currencies });
  },
  tickers: null,
  setTickers: (tickers: TickerResponse | null) => {
    set({ tickers });
  },
  orderbook: null,
  setOrderbook: (orderbook: OrderbookResponse | null) => {
    set({ orderbook });
  },
  openOrders: null,
  setOpenOrders: (openOrders: TransformedOrder[] | null) => {
    set({ openOrders });
  },
  orderHistory: [],
  setOrderHistory: (orderHistory) => {
    set((state) => ({
      orderHistory:
        typeof orderHistory === "function"
          ? orderHistory(state.orderHistory)
          : orderHistory,
    }));
  },
  loginModal: false,
  setLoginModal: (loginModal: boolean) => {
    set({ loginModal });
  },
  chartPeriod: "1",
  setChartPeriod: (period: string) => {
    set({ chartPeriod: period });
  },
  exchangeHistory: [],
  setExchangeHistory: (exchangeHistory) =>
    set((state) => ({
      exchangeHistory:
        typeof exchangeHistory === "function"
          ? exchangeHistory(state.exchangeHistory)
          : exchangeHistory,
    })),
}));
