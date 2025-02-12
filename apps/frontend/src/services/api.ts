import {
  CurrencyResponse,
  OrderbookResponse,
  TickerResponse,
  OhlcResponse,
  TradeHistoryResponse,
  CreateOrder,
  WalletAsset,
  CancelOrderResponse,
} from "@/types/market";
import { APIMethod, request } from "@/utils/api";
import { AxiosResponse } from "axios";
import { BASE_API_URL } from "@/config/envs";

export const getOHLC = async (
  symbol: string,
  period: string,
  from: string,
  to: string
): Promise<AxiosResponse<OhlcResponse>> => {
  const params = new URLSearchParams({
    symbol,
    period,
    from,
    to,
  });

  const response = await request(
    {},
    `${BASE_API_URL}/ohlc?${params}`,
    APIMethod.GET
  );

  if (!response.data) {
    throw new Error("No data received from OHLC API");
  }

  return response;
};

export const getTrades = async (
  symbol: string,
  from: number,
  to: number,
  account?: string
): Promise<AxiosResponse<TradeHistoryResponse>> => {
  const params = new URLSearchParams({
    symbol,
    from: from.toString(),
    to: to.toString(),
    ...(account && { account }),
  });

  const response = await request(
    {},
    `${BASE_API_URL}/trades?${params}`,
    APIMethod.GET
  );

  if (!response.data) {
    throw new Error("No data received from Trades API");
  }

  return response;
};

export const getTickers = async (
  symbols: string
): Promise<AxiosResponse<TickerResponse>> => {
  const params = new URLSearchParams({
    symbols,
  });

  const response = await request(
    {},
    `${BASE_API_URL}/tickers?${params}`,
    APIMethod.GET
  );

  if (!response.data) {
    throw new Error("No data received from Tickers API");
  }

  return response;
};

export const getCurrencies = async (): Promise<
  AxiosResponse<CurrencyResponse>
> => {
  const response = await request(
    {},
    `${BASE_API_URL}/currencies`,
    APIMethod.GET
  );

  if (!response.data) {
    throw new Error("No data received from Currencies API");
  }

  return response;
};

export const createOrder = async (order: CreateOrder) => {
  const response = await request(
    order,
    `${BASE_API_URL}/order/create`,
    APIMethod.POST
  );

  if (!response.data) {
    throw new Error("No data received from CreateOrder API");
  }

  return response;
};

export const submitOrder = async (order: { TX: string }) => {
  const response = await request(
    order,
    `${BASE_API_URL}/order/submit`,
    APIMethod.POST
  );

  if (!response.data) {
    throw new Error("No data received from SubmitOrder API");
  }

  return response;
};

export const getOrderbook = async (
  symbol: string,
  account?: string
): Promise<AxiosResponse<OrderbookResponse>> => {
  const params = new URLSearchParams({
    symbol,
    ...(account && { account }),
  });

  const response = await request(
    {},
    `${BASE_API_URL}/order/orderbook?${params}`,
    APIMethod.GET
  );

  if (!response.data) {
    throw new Error("No data received from Orderbook API");
  }

  return response;
};

export const getWalletAssets = async (
  address: string
): Promise<AxiosResponse<WalletAsset[]>> => {
  const params = new URLSearchParams({
    address,
  });

  const response = await request(
    {},
    `${BASE_API_URL}/wallet/assets?${params}`,
    APIMethod.GET
  );

  if (!response.data) {
    throw new Error("No data received from Wallet Assets API");
  }

  return response;
};

export const cancelOrder = async (
  address: string,
  id: string
): Promise<CancelOrderResponse> => {
  const response = await request(
    {
      Sender: address,
      OrderID: id,
    },
    `${BASE_API_URL}/order/cancel`,
    APIMethod.POST
  );

  if (!response.data) {
    throw new Error("No data received from CancelOrder API");
  }

  return response.data;
};
