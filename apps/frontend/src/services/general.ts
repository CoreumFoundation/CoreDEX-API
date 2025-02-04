import {
  CurrencyResponse,
  OrderbookResponse,
  TickerResponse,
  OhlcResponse,
  TradeHistoryResponse,
  CreateOrder,
  WalletAsset,
} from "@/types/market";
import { APIMethod, request } from "@/utils/api";
import { AxiosResponse } from "axios";

const API_URL =
  import.meta.env.VITE_MODE === "development"
    ? import.meta.env.VITE_ENV_BASE_API
    : (window as any).SOLOGENIC.env.VITE_ENV_BASE_API;

export const getOHLC = async (
  symbol: string,
  period: string,
  from: string,
  to: string
): Promise<AxiosResponse<OhlcResponse>> => {
  try {
    const data = request(
      {},
      `${API_URL}/ohlc?symbol=${symbol}&period=${period}&from=${from}&to=${to}`,
      APIMethod.GET
    );
    if (!data) {
      throw new Error("No data received from API");
    }
    return data;
  } catch (e) {
    console.log("ERROR GETTING OHLC DATA >>", e);
    throw e;
  }
};

export const getTrades = async (
  symbol: string,
  from: number,
  to: number,
  account?: string
): Promise<AxiosResponse<TradeHistoryResponse>> => {
  try {
    const data = await request(
      {},
      `${API_URL}/trades?symbol=${symbol}&from=${from}&to=${to}${
        account ? `&account=${account}` : ""
      }`,
      APIMethod.GET
    );
    if (!data) {
      throw new Error("No data received from API");
    }
    return data;
  } catch (e) {
    console.log("ERROR GETTING TRADES DATA >>", e);
    throw e;
  }
};

export const getTickers = async (
  symbols: string
): Promise<AxiosResponse<TickerResponse>> => {
  try {
    const data = await request(
      {},
      `${API_URL}/tickers?symbols=${symbols}`,
      APIMethod.GET
    );

    if (!data) {
      throw new Error("No data received from API");
    }

    return data;
  } catch (e) {
    console.log("ERROR GETTING TICKERS DATA >>", e);
    throw e;
  }
};

export const getCurrencies = async (): Promise<
  AxiosResponse<CurrencyResponse>
> => {
  try {
    const data = await request({}, `${API_URL}/currencies`, APIMethod.GET);
    if (!data) {
      throw new Error("No data received from API");
    }
    return data;
  } catch (e) {
    console.log("ERROR GETTING CURRENCIES DATA >>", e);
    throw e;
  }
};

export const createOrder = async (order: CreateOrder) => {
  try {
    const data = request(order, `${API_URL}/order/create`, APIMethod.POST);
    if (!data) {
      throw new Error("No data received from API");
    }
    return data;
  } catch (e) {
    console.log("ERROR CREATING ORDER >>", e);
    throw e;
  }
};

export const submitOrder = async (order: { TX: string }) => {
  try {
    const data = request(order, `${API_URL}/order/submit`, APIMethod.POST);
    if (!data) {
      throw new Error("No data received from API");
    }
    return data;
  } catch (e) {
    console.log("ERROR SUBMITTING ORDER >>", e);
    throw e;
  }
};

export const getOrderbook = async (
  symbol: string,
  account?: string
): Promise<AxiosResponse<OrderbookResponse>> => {
  try {
    const data = request(
      {},
      `${API_URL}/order/orderbook?symbol=${symbol}${
        account ? `&account=${account}` : ""
      }`,
      APIMethod.GET
    );
    if (!data) {
      throw new Error("No data received from API");
    }
    return data;
  } catch (e) {
    console.log("ERROR GETTING ORDERBOOK DATA >>", e);
    throw e;
  }
};

export const getWalletAssets = async (
  address: string
): Promise<AxiosResponse<WalletAsset[]>> => {
  try {
    const data = await request(
      {},
      `${API_URL}/wallet/assets?address=${address}`,
      APIMethod.GET
    );
    if (!data) {
      throw new Error("No data received from API");
    }
    return data;
  } catch (e) {
    console.log("ERROR GETTING WALLET ASSETS DATA >>", e);
    throw e;
  }
};

export const cancelOrder = async (
  address: string,
  id: string
): Promise<AxiosResponse<any>> => {
  try {
    const data = await request(
      {
        Sender: address,
        OrderID: id,
      },
      `${API_URL}/order/cancel`,
      APIMethod.POST
    );
    if (!data) {
      throw new Error("No data received from API");
    }
    return data;
  } catch (e) {
    console.log("ERROR CANCELING ORDER >>", e);
    throw e;
  }
};
