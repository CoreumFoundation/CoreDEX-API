import { useEffect, useMemo, useRef, useState } from "react";
import dayjs from "dayjs";
import { FormatNumber } from "../FormatNumber";
import { useStore } from "@/state/store";
import { getTrades } from "@/services/api";
import { SideBuy, TradeRecord } from "@/types/market";
import "./exchange-history.scss";
import {
  UpdateStrategy,
  wsManager,
  Method,
  NetworkToEnum,
} from "@/services/websocket";

const FIVE_MINUTES = 1 * 60;

const ExchangeHistory = () => {
  const { market, network, exchangeHistory, setExchangeHistory } = useStore();
  const [timeRange, setTimeRange] = useState({
    from: Math.floor(Date.now() / 1000),
    to: Math.floor(Date.now() / 1000) - FIVE_MINUTES,
  });
  const [isFetchingMore, setIsFetchingMore] = useState(false);
  const historyRef = useRef<HTMLDivElement>(null);
  const sentinelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchExchangeHistory = async () => {
      const { from, to } = timeRange;

      try {
        const response = await getTrades(market.pair_symbol, to, from);
        if (response.status === 200) {
          const data = response.data;
          setExchangeHistory(data);

          wsManager.setInitialState(subscription, data);
        }
      } catch (e) {
        console.log("ERROR GETTING EXCHANGE HISTORY DATA >>", e);
        setExchangeHistory([]);
      }
    };
    fetchExchangeHistory();
  }, [market.pair_symbol, timeRange]);

  const subscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.TRADES_FOR_SYMBOL,
      ID: market.pair_symbol,
    }),
    [market.pair_symbol, network]
  );

  const handler = (data: any) => {
    setExchangeHistory(data);
  };

  useEffect(() => {
    wsManager.connected().then(() => {
      wsManager.subscribe(subscription, handler, UpdateStrategy.MERGE);
    });
    return () => {
      wsManager.unsubscribe(subscription, handler);
    };
  }, [subscription]);

  // useEffect(() => {
  //   if (!historyRef.current) return;

  //   const handleScroll = async () => {
  //     const container = historyRef.current;
  //     if (!container) return;

  //     // Check if user has scrolled to the bottom (within a small threshold)
  //     const threshold = 50; // pixels from the bottom
  //     if (
  //       container.scrollTop + container.clientHeight >=
  //       container.scrollHeight - threshold
  //     ) {
  //       if (exchangeHistory && exchangeHistory.length > 0) {
  //         // Capture the current scroll height before loading more
  //         const oldScrollHeight = container.scrollHeight;

  //         // Calculate new boundaries based on timeRange
  //         const newFrom = timeRange.to; // current oldest timestamp
  //         const newTo = newFrom - FIVE_MINUTES; // 5 minutes older

  //         // Load older history and wait for it to finish
  //         await loadOlderHistory(newTo, newFrom);

  //         // After the new data is rendered, compute the change in scroll height.
  //         const newScrollHeight = container.scrollHeight;
  //         const deltaHeight = newScrollHeight - oldScrollHeight;

  //         // Adjust scrollTop by adding the delta so that the user's view stays at the same item.
  //         container.scrollTop = container.scrollTop + deltaHeight;
  //       }
  //     }
  //   };

  //   const container = historyRef.current;
  //   container.addEventListener("scroll", handleScroll);

  //   return () => {
  //     container.removeEventListener("scroll", handleScroll);
  //   };
  // }, [exchangeHistory, timeRange, market.pair_symbol]);

  const mergeUniqueTrades = (
    prevHistory: TradeRecord[],
    newTrades: TradeRecord[]
  ): TradeRecord[] => {
    const merged = [...prevHistory, ...newTrades];

    const unique = merged.filter(
      (trade, index, self) =>
        index === self.findIndex((t) => t.TXID === trade.TXID)
    );
    return unique;
  };

  const loadOlderHistory = async (newTo: number, newFrom: number) => {
    try {
      const response = await getTrades(market.pair_symbol, newTo, newFrom);
      if (response.status === 200) {
        const olderData = response.data;

        // Get the previous state from the store (using useStore.getState())
        const prevHistory = exchangeHistory || [];
        const mergedHistory = mergeUniqueTrades(prevHistory, olderData);

        wsManager.setInitialState(subscription, mergedHistory);
        setExchangeHistory(mergedHistory);
        setTimeRange({ ...timeRange, to: newTo });
      }
    } catch (e) {
      console.log("ERROR FETCHING OLDER HISTORY >>", e);
    }
  };

  return (
    <div className="exchange-history-container">
      <div className="exchange-history-title">Exchange History</div>

      <div className="header">
        <div className="exchange-history-body-row label">Price</div>
        <div className="exchange-history-body-row label">Volume</div>
        <div className="exchange-history-body-row label time">Time</div>
      </div>

      {exchangeHistory && exchangeHistory.length > 0 ? (
        <div ref={historyRef} className="exchange-history-body-rows">
          {exchangeHistory.map((trade, index: number) => {
            return (
              <div className="exchange-history-body-row" key={index}>
                <div
                  className={`exchange-history-body-value  ${
                    trade.Side === SideBuy.BUY ? "positive" : "negative"
                  }`}
                >
                  <FormatNumber number={trade.HumanReadablePrice} />
                </div>
                <div className="exchange-history-body-value volume">
                  <FormatNumber number={trade.SymbolAmount} />
                </div>
                <div className="exchange-history-body-value time">
                  {dayjs.unix(trade.BlockTime.seconds).format("h:mm A")}
                </div>
              </div>
            );
          })}
          <div ref={sentinelRef} style={{ height: "1px" }}></div>
        </div>
      ) : (
        <div className="no-data-container">
          <img src="/trade/images/warning.png" alt="warning" />
          <p className="no-data">No Data Found</p>
        </div>
      )}
    </div>
  );
};

export default ExchangeHistory;
