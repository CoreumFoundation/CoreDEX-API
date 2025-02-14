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
import duration from "dayjs/plugin/duration";
dayjs.extend(duration);

const ONE_MINUTE = dayjs.duration(1, "minutes").asSeconds();

const ExchangeHistory = () => {
  const { market, network, exchangeHistory, setExchangeHistory } = useStore();
  const [timeRange, setTimeRange] = useState({
    from: dayjs().subtract(3, "day").subtract(5, "minutes").unix(),
    to: dayjs().subtract(3, "day").unix(),
  });
  const [isFetchingMore, setIsFetchingMore] = useState(false);
  const historyRef = useRef<HTMLDivElement>(null);
  const sentinelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchExchangeHistory = async () => {
      const { from, to } = timeRange;

      try {
        const response = await getTrades(market.pair_symbol, from, to);
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

  useEffect(() => {
    if (!historyRef.current) return;
    const container = historyRef.current;

    const handleScroll = async () => {
      const threshold = 50;
      if (
        container.scrollTop + container.clientHeight >=
        container.scrollHeight - threshold
      ) {
        if (exchangeHistory && exchangeHistory.length > 0 && !isFetchingMore) {
          setIsFetchingMore(true);
          const distanceFromBottom =
            container.scrollHeight -
            container.scrollTop -
            container.clientHeight;

          await loadOlderHistory();

          requestAnimationFrame(() => {
            const newScrollHeight = container.scrollHeight;
            container.scrollTop =
              newScrollHeight - container.clientHeight - distanceFromBottom;
            setIsFetchingMore(false);
          });
        }
      }
    };

    container.addEventListener("scroll", handleScroll);
    return () => {
      container.removeEventListener("scroll", handleScroll);
    };
  }, [exchangeHistory, timeRange, market.pair_symbol, isFetchingMore]);

  const mergeUniqueTrades = (
    prevHistory: TradeRecord[],
    newTrades: TradeRecord[]
  ): TradeRecord[] => {
    const merged = [...prevHistory, ...newTrades];
    // Filter duplicates by TXID
    const unique = merged.filter(
      (trade, index, self) =>
        index === self.findIndex((t) => t.TXID === trade.TXID)
    );
    return unique.sort((a, b) => b.BlockTime.seconds - a.BlockTime.seconds);
  };

  const loadOlderHistory = async () => {
    try {
      const currentOldest = timeRange.from;
      const newFrom = currentOldest - ONE_MINUTE;
      const response = await getTrades(
        market.pair_symbol,
        newFrom,
        currentOldest
      );
      if (response.status === 200) {
        const olderData = response.data;
        const prevHistory = exchangeHistory || [];
        const mergedHistory = mergeUniqueTrades(prevHistory, olderData);
        wsManager.setInitialState(subscription, mergedHistory);
        setExchangeHistory(mergedHistory);
        setTimeRange((prev) => ({ ...prev, from: newFrom }));
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
