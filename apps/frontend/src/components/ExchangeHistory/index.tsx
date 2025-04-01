import { useEffect, useLayoutEffect, useMemo, useRef, useState } from "react";
import dayjs from "dayjs";
import { FormatNumber } from "../FormatNumber";
import { useStore } from "@/state/store";
import { getTrades } from "@/services/api";
import { TradeRecord } from "@/types/market";
import { Side } from "coredex-api-types/order-properties";

import "./exchange-history.scss";
import {
  UpdateStrategy,
  wsManager,
  NetworkToEnum,
  Subscription,
} from "@/services/websocket";
import { Method } from "coredex-api-types/update";
import duration from "dayjs/plugin/duration";
import debounce from "lodash/debounce";
import { FixedSizeList as List } from "react-window";
import { mirage } from "ldrs";
import { mergeUniqueTrades } from "@/utils";
mirage.register();

dayjs.extend(duration);

const ROW_HEIGHT = 26;
const containerHeight = 212;

const ExchangeHistory = () => {
  const { market, network, exchangeHistory, setExchangeHistory } = useStore();

  // initial window
  const [timeRange, setTimeRange] = useState({
    from: dayjs().subtract(1, "hour").unix(),
    to: dayjs().unix(),
  });
  const [isFetchingMore, setIsFetchingMore] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [isLoading, setIsLoading] = useState(false);

  const historyRef = useRef<HTMLDivElement>(null);

  const subscription: Subscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.TRADES_FOR_SYMBOL,
      ID: market.pair_symbol,
    }),
    [market.pair_symbol, network]
  );

  const handler = (newTrades: TradeRecord[]) => {
    const merged = mergeUniqueTrades(exchangeHistory, newTrades);
    setExchangeHistory(merged);
  };

  useEffect(() => {
    const initFetch = async () => {
      setIsLoading(true);
      await fetchTrades();
      setIsLoading(false);
    };
    initFetch();
  }, []);

  useEffect(() => {
    if (exchangeHistory && exchangeHistory.length > 0) {
      wsManager.connected().then(() => {
        wsManager.subscribe(subscription, handler, UpdateStrategy.MERGE);
      });
      return () => {
        wsManager.unsubscribe(subscription, handler);
      };
    }
  }, [subscription, exchangeHistory]);

  const fetchTrades = async (): Promise<boolean> => {
    try {
      const response = await getTrades({
        symbol: market.pair_symbol,
        side: Side.SIDE_BUY,
      });
      if (
        response.status === 200 &&
        response.data &&
        response.data.length > 0
      ) {
        setExchangeHistory(response.data);
        wsManager.setInitialState(subscription, response.data);
        return true;
      }
      return false;
    } catch (e) {
      console.log("ERROR GETTING ORDER HISTORY DATA >>", e);
      return false;
    }
  };

  const loadOlderHistory = async (): Promise<number> => {
    try {
      const currentWindow = timeRange.to - timeRange.from;
      const newTo = timeRange.from;
      const newFrom = newTo - currentWindow;
      const response = await getTrades({
        symbol: market.pair_symbol,
        from: newFrom,
        to: newTo,
      });
      if (response.status === 200) {
        const olderData = response.data;
        if (!olderData || olderData.length === 0) {
          setHasMore(false);
          return 0;
        }
        const prevHistory = exchangeHistory || [];
        const mergedHistory = mergeUniqueTrades(prevHistory, olderData);
        wsManager.setInitialState(subscription, mergedHistory);
        setExchangeHistory(mergedHistory);
        setTimeRange({ from: newFrom, to: newTo });
        return olderData.length;
      }
    } catch (e) {
      console.log("ERROR FETCHING OLDER HISTORY >>", e);
    }
    return 0;
  };

  useLayoutEffect(() => {
    if (!historyRef.current) return;
    const container = historyRef.current;

    const handleScroll = async () => {
      const threshold = 50;
      if (
        container.scrollTop + container.clientHeight >=
        container.scrollHeight - threshold
      ) {
        if (
          exchangeHistory &&
          exchangeHistory.length > 0 &&
          !isFetchingMore &&
          hasMore
        ) {
          setIsFetchingMore(true);
          const anchorEl = container.querySelector(
            ".exchange-history-body-rows"
          );
          const anchorRect = anchorEl ? anchorEl.getBoundingClientRect() : null;
          const previousScrollTop = container.scrollTop;

          await loadOlderHistory();

          requestAnimationFrame(() => {
            if (anchorEl && anchorRect) {
              const newAnchorRect = anchorEl.getBoundingClientRect();
              const delta = newAnchorRect.top - anchorRect.top;
              container.scrollTop = previousScrollTop + delta;
            } else {
              container.scrollTop = previousScrollTop;
            }
            setIsFetchingMore(false);
          });
        }
      }
    };

    const debouncedHandleScroll = debounce(handleScroll, 300);
    container.addEventListener("scroll", debouncedHandleScroll);
    return () => {
      container.removeEventListener("scroll", debouncedHandleScroll);
      debouncedHandleScroll.cancel();
    };
  }, [exchangeHistory, timeRange, market.pair_symbol, isFetchingMore, hasMore]);

  const Row = ({
    index,
    style,
  }: {
    index: number;
    style: React.CSSProperties;
  }) => {
    const trade = exchangeHistory[index];
    return (
      <div style={style} className="exchange-history-body-row">
        <div className={`exchange-history-body-value`}>
          <FormatNumber number={trade?.HumanReadablePrice} />
        </div>
        <div className="exchange-history-body-value volume">
          <FormatNumber number={trade?.SymbolAmount} />
        </div>
        <div className="exchange-history-body-value time">
          {dayjs.unix(trade?.BlockTime.seconds).format("MM/DD HH:mm")}
        </div>
      </div>
    );
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
        <div className="exchange-history-body-rows">
          <List
            height={containerHeight}
            itemCount={exchangeHistory.length}
            itemSize={ROW_HEIGHT}
            width={"100%"}
            outerRef={historyRef}
          >
            {Row}
          </List>
          <div style={{ height: "1px" }}></div>
        </div>
      ) : (
        <div className="no-data-container">
          {isLoading ? (
            <>
              <l-mirage size="40" speed="6" color="#25d695"></l-mirage>
            </>
          ) : (
            <>
              <img src="/trade/images/warning.png" alt="warning" />
              <p className="no-data">No Data Found</p>
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default ExchangeHistory;
