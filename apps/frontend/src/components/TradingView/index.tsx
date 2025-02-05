import { useCallback, useEffect, useMemo, useState } from "react";
import themes from "./tools/theme";
import { widget as Widget } from "../../vendor/tradingview/charting_library";
import { CoreumDataFeed } from "./tools/api";
import { DEFAULT_CONFIGS, getOverrides } from "./tools/config";
import { useSaveAndClear, useMountChart, useChartTheme } from "@/hooks";
import { useStore } from "@/state/store";
import "./tradingview.scss";
import {
  Action,
  Method,
  NetworkToEnum,
  WebSocketMessage,
} from "@/services/websocket";
import { useWebSocket } from "@/hooks/websocket";
import { resolveResolution } from "./tools/utils";
import { OhlcRecord } from "@/types/market";

declare global {
  interface Window {
    tvWidget: any;
  }
}

const resolutions: { [key: string]: string } = {
  "1": "1",
  "3": "3",
  "5": "5",
  "15": "15",
  "30": "30",
  "1h": "1h",
  "3h": "3h",
  "6h": "6h",
  "12h": "12h",
  "1D": "1D",
  "3D": "3D",
  "1W": "1W",
  "60": "1h",
  "180": "3h",
  "360": "6h",
  "720": "12h",
};

const TradingView = ({ height }: { height: number | string }) => {
  const { market, chartPeriod, setChartPeriod, network } = useStore();
  const [resolution, setResolution] = useState<string>(chartPeriod);
  const [dataFeed, setDataFeed] = useState<CoreumDataFeed | null>(null);
  const [lastUpdate, setLastUpdate] = useState<any>(null);

  const ohlcSubscription = useMemo(() => {
    const base = market.base.Denom;
    const counter = market.counter.Denom;

    return {
      Network: NetworkToEnum(network),
      Method: Method.OHLC,
      ID: `${base.Denom}_${counter.Denom}_${resolveResolution(chartPeriod)}`,
    };
  }, [market, chartPeriod]);

  const handleDataFeedUpdate = useCallback(
    (message: WebSocketMessage) => {
      if (message.Action === Action.RESPONSE && message.Subscription?.Content) {
        setLastUpdate(message.Subscription.Content);
      }
    },
    [setLastUpdate]
  );

  useWebSocket(ohlcSubscription, handleDataFeedUpdate);

  useEffect(() => {
    mountChart();

    return () => {
      if (window.tvWidget) {
        window.tvWidget.remove();
        window.tvWidget = null;
      }
      if (dataFeed) {
        dataFeed.subscriptions = [];
      }
    };
  }, [market.pair_symbol]);

  // data feed updates
  useEffect(() => {
    if (!lastUpdate || !dataFeed) return;
    const bars = lastUpdate.map((el: OhlcRecord) => ({
      time: Number(el[0] * 1000),
      open: Number(el[1]),
      high: Number(el[2]),
      low: Number(el[3]),
      close: Number(el[4]),
      volume: Number(el[5]),
    }));

    dataFeed.subscriptions.forEach((sub) => {
      if (bars.length > 0) {
        sub.onRealtimeCallback(bars[bars.length - 1]);
      }
    });
  }, [lastUpdate, dataFeed]);

  const coreumConstructorFeed = () => {
    const { base: baseDenom, counter: counterDenom } = market;

    const getDenomString = (denom: typeof baseDenom) =>
      `${denom.Denom.Currency}${
        denom.Denom.Issuer ? `-${denom.Denom.Issuer}` : ""
      }`;

    const base = getDenomString(baseDenom);
    const counter = getDenomString(counterDenom);

    return {
      id: `${base}_${counter}`,
      name: `${baseDenom.Denom.Name?.toUpperCase()} / ${counterDenom.Denom.Name?.toUpperCase()}`,
    };
  };

  const mountChart = () => {
    const symbol = coreumConstructorFeed();

    if (window.tvWidget) {
      window.tvWidget.remove();
    }
    const dataFeedInstance = new CoreumDataFeed(symbol);
    setDataFeed(dataFeedInstance);

    const widgetOptions = {
      // debug: true, // TV logs
      symbol: symbol.name,
      datafeed: dataFeedInstance,
      height: height,
      interval: resolution,
      theme: "Dark",
      loading_screen: {
        backgroundColor: themes["dark"].colors.chart.background,
        foregroundColor: "#D81D3C",
      },
      locale: "en",
      ...DEFAULT_CONFIGS,
    };
    const widget = (window.tvWidget = new Widget(widgetOptions as any));

    widget.onChartReady(() => {
      setTimeout(() => {
        widget
          .activeChart()
          .onIntervalChanged()
          .subscribe(null, (interval) => {
            updateResolution(interval);
          });
        widget.applyOverrides(getOverrides("dark"));
        setReady(true);
      });
    });
  };

  const updateResolution = (res: string) => {
    const widget = window.tvWidget;
    widget.chart().setResolution(res);

    // resolve the resolution or default to '1W' if invalid
    const validResolutions = Object.keys(resolutions);
    const resolvedRes = validResolutions.includes(res)
      ? resolutions[res]
      : "1W";
    localStorage.chart_resolution = resolvedRes;
    setResolution(resolvedRes);

    // resubscribe to bars with the new resolution
    if (dataFeed) {
      dataFeed.subscriptions.forEach((sub) => {
        dataFeed.unsubscribeBars(sub.key);
        dataFeed.subscribeBars(
          sub.symbolInfo,
          res,
          sub.onRealtimeCallback,
          sub.key
        );
      });
    }
  };

  const openIndicators = () => {
    const widget = window.tvWidget;
    widget?.chart().executeActionById("insertIndicator");
  };

  const { chartReady, setReady } = useMountChart(mountChart);
  useChartTheme(chartReady);
  const { clearable, saveChart, clearChart } = useSaveAndClear(
    mountChart,
    setReady
  );

  return (
    <div className="chart-wrapper">
      <div className="top-toolbar">
        <div className="intervals">
          <div
            className={`interval ${resolution === "1D" ? "active" : ""}`}
            onClick={() => {
              updateResolution("1D");
              setChartPeriod("1D");
            }}
          >
            1D
          </div>

          <div
            onClick={() => {
              updateResolution("3D");
              setChartPeriod("3D");
            }}
            className={`interval ${resolution === "3D" ? "active" : ""}`}
          >
            3D
          </div>

          <div
            onClick={() => {
              updateResolution("1W");
              setChartPeriod("1W");
            }}
            className={`interval ${resolution === "1W" ? "active" : ""}`}
          >
            1W
          </div>

          <div onClick={() => openIndicators()} className="interval indicators">
            <p className="interval">Technical Ind</p>
          </div>
        </div>
      </div>
      {!chartReady && (
        <div
          style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            height: "100%",
            width: "100%",
          }}
        >
          loading...
        </div>
      )}
      <div
        id="chartContainer"
        className="chart-container"
        style={{ display: chartReady ? "initial" : "none" }}
      />
    </div>
  );
};

export default TradingView;
