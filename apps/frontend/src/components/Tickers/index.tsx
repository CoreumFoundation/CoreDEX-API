import { useEffect, useMemo, useState } from "react";
import { FormatNumber } from "../FormatNumber";
import { useStore } from "@/state/store";
import { getTickers } from "@/services/api";

import {
  wsManager,
  UpdateStrategy,
  NetworkToEnum,
  Method,
} from "@/services/websocket-refactor";
import "./tickers.scss";
import { TickerResponse } from "@/types/market";

const Tickers = () => {
  const { market, network } = useStore();
  const [change, setChange] = useState<number>(0);
  const [tickers, setTickers] = useState<TickerResponse | null>(null);

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
          setTickers(data);
        }
      } catch (e) {
        console.log("ERROR GETTING TICKERS DATA >>", e);
        setTickers(null);
      }
    };

    fetchTickers();
  }, [market.pair_symbol]);

  const subscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.TICKER,
      ID: market.pair_symbol,
    }),
    [market.pair_symbol, network]
  );

  useEffect(() => {
    wsManager.connected().then(() => {
      wsManager.subscribe(subscription, setTickers, UpdateStrategy.REPLACE);
    });
    return () => {
      wsManager.unsubscribe(subscription, setTickers);
    };
  }, [subscription]);

  useEffect(() => {
    if (tickers) {
      const { OpenPrice, LastPrice } = tickers.Tickers[market.pair_symbol];
      const difference = Number(LastPrice) - Number(OpenPrice);
      const change = 100 * (difference / Number(OpenPrice));
      setChange(change);
    }
  }, [tickers, market]);

  return (
    <div className="tickers-container">
      <div className="price-container">
        <div className="price">
          <FormatNumber
            number={tickers ? tickers.Tickers[market.pair_symbol].LastPrice : 0}
          />
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
          <FormatNumber
            number={tickers ? tickers.Tickers[market.pair_symbol].Volume : 0}
          />
        </div>
      </div>

      <div className="volume-counter">
        <div className="label">{`24h Volume (${market.counter.Denom.Name})`}</div>
        <div className="volume">
          <FormatNumber
            number={
              tickers ? tickers.Tickers[market.pair_symbol].Invertedvolume : 0
            }
          />
        </div>
      </div>

      <div className="high">
        <div className="label">{`24h High`}</div>
        <div className="volume">
          <FormatNumber
            number={tickers ? tickers.Tickers[market.pair_symbol].HighPrice : 0}
          />
        </div>
      </div>

      <div className="low">
        <div className="label">{`24h Low`}</div>
        <div className="volume">
          <FormatNumber
            number={tickers ? tickers.Tickers[market.pair_symbol].LowPrice : 0}
          />
        </div>
      </div>
    </div>
  );
};

export default Tickers;
