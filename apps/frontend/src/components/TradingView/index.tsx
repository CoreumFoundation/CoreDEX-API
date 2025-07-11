import { useEffect, useMemo, useRef, useState } from "react";
import themes from "./tools/theme";
import { widget as Widget } from "../../vendor/tradingview/charting_library";
import { CoreumDataFeed } from "./tools/api";
import { DEFAULT_CONFIGS, getOverrides } from "./tools/config";
import { useSaveAndClear, useMountChart } from "@/hooks";
import { useStore } from "@/state/store";
import "./tradingview.scss";
import { resolveResolution } from "./tools/utils";
import { OhlcRecord } from "@/types/market";
import { NetworkToEnum, UpdateStrategy, wsManager } from "@/services/websocket";
import { Method } from "coredex-api-types/update";
import dayjs from "dayjs";
import timezone from "dayjs/plugin/timezone";
dayjs.extend(timezone);
import { mirage } from "ldrs";
mirage.register();

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

// three minute threshold
const REMOUNT_THRESHOLD = 3 * 60 * 1000;

const TradingView = ({ height }: { height: number | string }) => {
  const { market, chartPeriod, setChartPeriod, network } = useStore();
  const [resolution, setResolution] = useState<string>(chartPeriod);
  const [dataFeed, setDataFeed] = useState<CoreumDataFeed | null>(null);
  const [lastUpdate, setLastUpdate] = useState<any>(null);
  const lastRemountTime = useRef(0);

  const ohlcSubscription = useMemo(() => {
    const base = market.base.Denom;
    const counter = market.counter.Denom;

    return {
      Network: NetworkToEnum(network),
      Method: Method.OHLC,
      ID: `${base.Denom}_${counter.Denom}_${resolveResolution(chartPeriod)}`,
    };
  }, [market, chartPeriod]);

  useEffect(() => {
    wsManager.connected().then(() => {
      wsManager.subscribe(
        ohlcSubscription,
        setLastUpdate,
        UpdateStrategy.REPLACE
      );
    });
    return () => {
      wsManager.unsubscribe(ohlcSubscription, setLastUpdate);
    };
  }, [ohlcSubscription]);

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
  }, [market.pair_symbol, network]);

  // remount chart if away from tab for a while
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (document.visibilityState === "visible") {
        const now = Date.now();
        if (now - lastRemountTime.current > REMOUNT_THRESHOLD) {
          mountChart();
          lastRemountTime.current = now;
        }
      }
    };

    document.addEventListener("visibilitychange", handleVisibilityChange);

    return () => {
      document.removeEventListener("visibilitychange", handleVisibilityChange);
    };
  }, []);

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
        bars.forEach((bar: any) => {
          handleWebsocketTick(sub, bar);
        });
      }
    });
  }, [lastUpdate, dataFeed]);

  const handleWebsocketTick = (sub: any, newTick: any) => {
    const lastBar = sub.lastBar;

    if (!lastBar) {
      sub.lastBar = newTick;
      sub.onRealtimeCallback(newTick);
      return;
    }

    if (newTick.time >= lastBar.time) {
      sub.lastBar = newTick;
      sub.onRealtimeCallback(newTick);
    } else {
      return;
    }
  };

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
      timezone: dayjs.tz.guess(),
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

    const validResolutions = Object.keys(resolutions);
    const resolvedRes = validResolutions.includes(res)
      ? resolutions[res]
      : "1W";
    localStorage.chart_resolution = resolvedRes;
    setResolution(resolvedRes);

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
  useSaveAndClear(mountChart, setReady);

  return (
    <div className="chart-wrapper">
      <div className="top-toolbar">
        {dataFeed?.errorMessage && (
          <p className="tradingview-error">
            Error: Something went wrong with chart data
          </p>
        )}
        <div className="intervals">
          <div
            className={`interval ${resolution === "1" ? "active" : ""}`}
            onClick={() => {
              updateResolution("1");
              setChartPeriod("1");
            }}
          >
            1m
          </div>
          <div
            className={`interval ${resolution === "5" ? "active" : ""}`}
            onClick={() => {
              updateResolution("5");
              setChartPeriod("5");
            }}
          >
            5m
          </div>
          <div
            className={`interval ${resolution === "15" ? "active" : ""}`}
            onClick={() => {
              updateResolution("15");
              setChartPeriod("15");
            }}
          >
            15m
          </div>
          <div
            className={`interval ${resolution === "30" ? "active" : ""}`}
            onClick={() => {
              updateResolution("30");
              setChartPeriod("30");
            }}
          >
            30m
          </div>
          <div
            className={`interval ${resolution === "1h" ? "active" : ""}`}
            onClick={() => {
              updateResolution("60");
              setChartPeriod("60");
            }}
          >
            1h
          </div>
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
          <l-mirage size="40" speed="6" color="#25d695"></l-mirage>
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
