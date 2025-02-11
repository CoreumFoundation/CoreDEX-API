import { useCallback, useEffect, useMemo, useRef } from "react";
import dayjs from "dayjs";
import { FormatNumber } from "../FormatNumber";
import { useStore } from "@/state/store";
import { getTrades } from "@/services/api";
import { Method, NetworkToEnum, WebSocketMessage } from "@/services/websocket";
import { useWebSocket } from "@/hooks/websocket";
import { SideBuy, TradeHistoryResponse } from "@/types/market";
import "./exchange-history.scss";

const ExchangeHistory = () => {
  const { market, setExchangeHistory, exchangeHistory, network } = useStore();
  const historyRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchExchangeHistory = async () => {
      const from = new Date().getTime();
      const to = from - 86400000; //  1 day ago

      try {
        const response = await getTrades(market.pair_symbol, to, from);
        if (response.status === 200) {
          const data = response.data;
          console.log(data);
          setExchangeHistory(data);
        }
      } catch (e) {
        console.log("ERROR GETTING EXCHANGE HISTORY DATA >>", e);
        setExchangeHistory(null);
      }
    };
    fetchExchangeHistory();
  }, [market.pair_symbol]);

  const handleExchangeHistoryUpdate = useCallback(
    (message: WebSocketMessage) => {
      const data = message.Subscription?.Content;
      // TODO move this to ws service

      if (data.length > 0) {
        if (exchangeHistory) {
          const updatedHistory: TradeHistoryResponse =
            exchangeHistory.concat(data);
          setExchangeHistory(updatedHistory);
        } else {
          setExchangeHistory(data);
        }
      }
    },
    [setExchangeHistory]
  );

  const subscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.TRADES_FOR_SYMBOL,
      ID: market.pair_symbol,
    }),
    [market]
  );

  useWebSocket(subscription, handleExchangeHistoryUpdate);

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
                  {dayjs(trade.BlockTime.seconds).format("HH:mm:ss")}
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
