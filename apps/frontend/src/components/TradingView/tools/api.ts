import { getBars } from "./utils";
import {
  BarPeriodParams,
  BarSymbolInfo,
  ChartSubscription,
  DataFeedAsset,
} from "@/types/market";
import { SUPPORTED_RESOLUTIONS } from "./utils";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";
import { useStore } from "@/state/store";

dayjs.extend(utc);
dayjs.extend(timezone);
declare global {
  interface Window {
    chart_int_feed: any;
  }
}

export class CoreumDataFeed {
  subscriptions: ChartSubscription[];
  asset: DataFeedAsset;
  searchSymbols: any;
  errorMessage: boolean | null;
  constructor(asset: DataFeedAsset) {
    this.subscriptions = [];
    this.asset = asset;
    this.searchSymbols = null;
    this.errorMessage = null;
    return this;
  }

  onReady(cb: any) {
    clearInterval(window.chart_int_feed);

    setTimeout(() => {
      cb({ SUPPORTED_RESOLUTIONS });
    }, 0);
  }

  resolveSymbol(
    symbolName: string,
    onSymbolResolvedCallback: any,
    onResolveErrorCallback: any
  ) {
    setTimeout(() => {
      const tickers = useStore.getState().tickers;
      const market = useStore.getState().market;
      const ticker = tickers && tickers.Tickers[market.pair_symbol];

      if (ticker && ticker.LastPrice) {
        const zeros = () => {
          let string = "";
          const decimalLength = String(ticker.LastPrice).includes(".")
            ? String(ticker.LastPrice).split(".")[1].length
            : 10;

          const len = decimalLength < 15 ? decimalLength : 15;

          for (let i = 0; i < len; i++) string += "0";

          return string;
        };

        const symbol_stub: BarSymbolInfo = {
          ...this.asset,
          // exchange: "COREUMDEX",
          session: "24x7",
          timezone: dayjs.tz.guess(),
          has_intraday: true,
          has_weekly_and_monthly: true,
          supported_resolutions: SUPPORTED_RESOLUTIONS,
          pricescale:
            tickers && Number(ticker.LastPrice) < 0.000001
              ? Number(`1${zeros()}`)
              : 1000000,
          minmov: 1,
        };

        setTimeout(() => {
          onSymbolResolvedCallback(symbol_stub);
        }, 0);
      } else {
        this.errorMessage = true;
        onResolveErrorCallback(`Could not resolve symbol ${symbolName}`);
      }
    }, 1500);
  }

  getBars(
    symbolInfo: BarSymbolInfo,
    resolution: string,
    periodParams: BarPeriodParams,
    onHistoryCallback: any,
    onErrorCallback: any
  ) {
    if (periodParams.firstDataRequest) periodParams.to = Date.now() / 1000;
    getBars(symbolInfo.id, resolution, periodParams.from, periodParams.to)
      .then((bars) => {
        const sortedBars = bars
          .map((el) => ({
            time: el.time,
            close: el.close,
            open: el.open,
            high: el.high,
            low: el.low,
            volume: el.volume,
          }))
          .sort((a, b) => a.time - b.time);

        const uniqueBars = [];
        for (let bar of sortedBars) {
          if (
            uniqueBars.length === 0 ||
            uniqueBars[uniqueBars.length - 1].time !== bar.time
          ) {
            uniqueBars.push(bar);
          }
        }

        onHistoryCallback(uniqueBars, { noData: uniqueBars.length === 0 });
      })
      .catch((err) => {
        onErrorCallback(err);
        console.log(err);
      });
  }

  subscribeBars(
    symbolInfo: BarSymbolInfo,
    resolution: string,
    onRealtimeCallback: any,
    key: string
  ) {
    this.unsubscribeBars(key);
    this.subscriptions.push({
      key: `${key}`,
      symbolInfo,
      resolution,
      onRealtimeCallback,
    });
  }

  unsubscribeBars(key: string) {
    this.subscriptions = this.subscriptions.filter((s) => s.key !== key);
  }

  reset() {
    this.subscriptions = [];
  }
  
  getErrorMessage() {
    return this.errorMessage;
  }
}
