import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import dayjs from "dayjs";
import { FormatNumber } from "../FormatNumber";
import { useStore } from "@/state/store";
import { getTrades } from "@/services/api";
import { Method, NetworkToEnum, WebSocketMessage } from "@/services/websocket";
import { useWebSocket } from "@/hooks/websocket";
import { SideBuy, TradeHistoryResponse } from "@/types/market";
import "./exchange-history.scss";
import { UpdateStrategy, wsManager } from "@/services/websocket-refactor";

const ExchangeHistory = () => {
  const { market, network, exchangeHistory, setExchangeHistory } = useStore();
  const historyRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchExchangeHistory = async () => {
      const from = new Date().getTime();
      const to = from - 86400000; //  1 day ago

      try {
        const response = await getTrades(market.pair_symbol, to, from);
        if (response.status === 200) {
          const data = response.data;
          setExchangeHistory(data);

          wsManager.setInitialState(subscription, data);
        }
      } catch (e) {
        console.log("ERROR GETTING EXCHANGE HISTORY DATA >>", e);
        setExchangeHistory(null);
      }
    };
    fetchExchangeHistory();
  }, [market.pair_symbol]);

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
      wsManager.unsubscribe(subscription, setExchangeHistory);
    };
  }, [subscription]);

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
