import { useCallback, useEffect, useMemo, useState } from "react";
import { FormatNumber } from "../FormatNumber";
import { useStore } from "@/state/store";
import { getTickers } from "@/services/api";
import { Method, NetworkToEnum, WebSocketMessage } from "@/services/websocket";
import { useWebSocket } from "@/hooks/websocket";
import "./tickers.scss";

const Tickers = () => {
  const { setTickers, tickers, market, network } = useStore();
  const [change, setChange] = useState<number>(0);

  // initial tickers
  useEffect(() => {
    const fetchTickers = async () => {
      try {
        const base = market.base.Denom.Denom;
        const counter = market.counter.Denom.Denom;
        const symbols = btoa(JSON.stringify([`${base}_${counter}`]));
        const response = await getTickers(symbols);

        if (response.status === 200 && response.data.Tickers) {
          const data = response.data;
          const ticker = data.Tickers[market.pair_symbol];
          setTickers(ticker);
        }
      } catch (e) {
        console.log("ERROR GETTING TICKERS DATA >>", e);
        setTickers(null);
      }
    };

    fetchTickers();
  }, [market.pair_symbol]);

  const handleTickerUpdate = useCallback(
    (message: WebSocketMessage) => {
      const tickerContent =
        message.Subscription?.Content.Tickers[market.pair_symbol];
      if (!tickerContent) return;
      setTickers(tickerContent);
    },
    [setTickers]
  );

  const subscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.TICKER,
      ID: market.pair_symbol,
    }),
    [market.pair_symbol, network]
  );

  useWebSocket(subscription, handleTickerUpdate);

  useEffect(() => {
    if (tickers) {
      const { OpenPrice, LastPrice } = tickers;
      const difference = Number(LastPrice) - Number(OpenPrice);
      const change = 100 * (difference / Number(OpenPrice));
      setChange(change);
    }
  }, [tickers, market]);

  return (
    <div className="tickers-container">
      <div className="price-container">
        <div className="price">
          <FormatNumber number={tickers ? tickers.LastPrice : 0} />
        </div>
        <div
          className={`change ${Number(change) > 0 ? "positive" : "negative"}`}
        >
          <span>{Number(change) >= 0 ? "+" : ""}</span>
          <FormatNumber number={change} precision={2} />
        </div>
      </div>

      <div className="volume-base">
        <div className="label">{`24h Volume (${market.base.Denom.Name})`}</div>
        <div className="volume">
          <FormatNumber number={tickers ? tickers.Volume : 0} />
        </div>
      </div>

      <div className="volume-counter">
        <div className="label">{`24h Volume (${market.counter.Denom.Name})`}</div>
        <div className="volume">
          <FormatNumber number={tickers ? tickers.Invertedvolume : 0} />
        </div>
      </div>

      <div className="high">
        <div className="label">{`24h High`}</div>
        <div className="volume">
          <FormatNumber number={tickers ? tickers.HighPrice : 0} />
        </div>
      </div>

      <div className="low">
        <div className="label">{`24h Low`}</div>
        <div className="volume">
          <FormatNumber number={tickers ? tickers.LowPrice : 0} />
        </div>
      </div>
    </div>
  );
};

export default Tickers;
