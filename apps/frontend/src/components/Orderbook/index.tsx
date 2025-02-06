import { useState, useEffect, useRef, useCallback, useMemo } from "react";
import BigNumber from "bignumber.js";
import { useStore } from "@/state/store";
import { TooltipPosition, useTooltip } from "@/hooks";
import { toFixedDown } from "@/utils";
import { FormatNumber } from "../FormatNumber";
import { OrderType, OrderbookAction, OrderbookRecord } from "@/types/market";
import { getOrderbook } from "@/services/general";
import { Method, NetworkToEnum, WebSocketMessage } from "@/services/websocket";
import { useWebSocket } from "@/hooks/websocket";
import "./orderbook.scss";

enum ORDERBOOK_TYPE {
  BUY = "buy",
  SELL = "sell",
  BOTH = "both",
}

export default function Orderbook({
  setOrderbookAction,
}: {
  setOrderbookAction: (action: OrderbookAction) => void;
}) {
  const { setOrderbook, orderbook, market, network } = useStore();

  const [spread, setSpread] = useState<BigNumber>(new BigNumber(0));
  const [topBuyVolume, setTopBuyVolume] = useState<number>(0);
  const [topSellVolume, setTopSellVolume] = useState<number>(0);

  // tooltip state
  const { showTooltip, hideTooltip } = useTooltip();
  const [leftPos, _] = useState<number>(0);
  const [orderbookType, __] = useState<ORDERBOOK_TYPE>(ORDERBOOK_TYPE.BOTH);
  const componentRef = useRef<HTMLDivElement>(null);

  // note: because we want Sells to display as descending, we reverse the order
  // we have to do some index manipulation to maintain proper hover and tooltip
  // in calcaulteGroupData and renderOrderRow
  // the alternative is to reverse the Sell array before setting in state but that
  // would require a lot of array manipulation on every message

  useEffect(() => {
    const fetchOrderbook = async () => {
      try {
        const response = await getOrderbook(market.pair_symbol);
        if (response.status === 200 && response.data) {
          const data = response.data;
          setOrderbook(data);
        }
      } catch (e) {
        console.log("ERROR GETTING ORDERBOOK DATA >>", e);
        setOrderbook(null);
      }
    };
    fetchOrderbook();
  }, [market.base, market.counter]);

  const handleOrderbookUpdate = useCallback(
    (message: WebSocketMessage) => {
      setOrderbook(message.Subscription?.Content);
    },
    [setOrderbook]
  );

  const subscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.ORDERBOOK,
      ID: market.pair_symbol,
    }),
    [market.pair_symbol]
  );

  useWebSocket(subscription, handleOrderbookUpdate);

  // calculate spread
  useEffect(() => {
    if (!orderbook) return;

    const calculateSpread = () => {
      const bestBid = orderbook.Buy[0]?.HumanReadablePrice;
      const bestAsk = orderbook.Sell[0]?.HumanReadablePrice;

      if (!bestBid && !bestAsk) return new BigNumber(0);

      if (bestBid && bestAsk) {
        return new BigNumber(bestAsk).minus(bestBid);
      }

      return new BigNumber(bestAsk || bestBid || 0);
    };

    setSpread(calculateSpread());
  }, [orderbook, market.pair_symbol]);

  // scroll buys to bottom
  useEffect(() => {
    const buysOb = document.getElementById("buys_ob");
    if (buysOb === null) {
      const timer = setInterval(function () {
        const buysObook = document.getElementById("buys_ob");

        if (buysObook) {
          clearInterval(timer);
        }
      }, 200);
    } else {
      buysOb.scrollTop = buysOb.scrollHeight;
    }
  }, []);

  // find the highest volume in the orderbook
  useEffect(() => {
    if (orderbook?.Buy?.length) {
      const highestBuyVolume = orderbook.Buy.reduce((max, buy) => {
        const current = new BigNumber(buy.SymbolAmount);
        return current.isGreaterThan(max) ? current : max;
      }, new BigNumber(0)).toNumber();
      setTopBuyVolume(highestBuyVolume);
    }

    if (orderbook?.Sell?.length) {
      const highestSellVolume = orderbook.Sell.reduce((max, sell) => {
        const current = new BigNumber(sell.SymbolAmount);
        return current.isGreaterThan(max) ? current : max;
      }, new BigNumber(0)).toNumber();
      setTopSellVolume(highestSellVolume);
    }
  }, [orderbook]);

  const calculateGroupData = useCallback(
    (lines: OrderbookRecord[], index: number, orderType: ORDERBOOK_TYPE) => {
      let avgPriceSum = 0;
      let totalVolume = 0;
      let sum = 0;
      const lineGroup = [];

      // buys are in descending order, sells are in ascending order
      const increment = orderType === ORDERBOOK_TYPE.BUY ? -1 : 1;
      const endCondition = orderType === ORDERBOOK_TYPE.BUY ? 0 : lines.length;

      for (
        let i = index;
        orderType === ORDERBOOK_TYPE.BUY ? i >= endCondition : i < endCondition;
        i += increment
      ) {
        const line = lines[i];
        if (!line) break;

        lineGroup.push(line);
        avgPriceSum += Number(line.HumanReadablePrice);
        totalVolume += Number(line.SymbolAmount);
        sum += Number(line.HumanReadablePrice) * Number(line.SymbolAmount);
      }

      return {
        avgPrice: lineGroup.length ? avgPriceSum / lineGroup.length : 0,
        sum,
        lineGroup,
        totalVolume,
      };
    },
    []
  );

  const handleTooltip = useCallback(
    (
      e: React.MouseEvent<HTMLDivElement>,
      index: number,
      orderType: ORDERBOOK_TYPE,
      isHovering: boolean
    ) => {
      if (!orderbook || !isHovering) {
        removeHoverClasses();
        hideTooltip();
        return;
      }

      const lines =
        orderType === ORDERBOOK_TYPE.BUY ? orderbook.Buy : orderbook.Sell;

      const { avgPrice, sum, lineGroup, totalVolume } = calculateGroupData(
        lines,
        index,
        orderType
      );

      lineGroup.forEach((g) => {
        const element = document.querySelector(
          `.orderbook-row[data-value="${g.HumanReadablePrice}_${orderType}_${g.Sequence}"]`
        );
        element?.classList.add(
          `${orderType === ORDERBOOK_TYPE.BUY ? "hovered-buy" : "hovered-sell"}`
        );
      });

      const tooltipContent = `
      <div class="orderbook-tooltip">
        <div class="inline-item">
          <p class="inline-item-label">Avg. Price:</p> ~ ${toFixedDown(
            avgPrice,
            12
          )}
        </div>
        
        <div class="inline-item">
          <p class="inline-label">Sum (${
            market.counter.Denom.Name
          }):</p> ${toFixedDown(totalVolume, 12)}
        </div>
        
        <div class="inline-item">
          <p class="inline-label">Total Volume:</p> ${toFixedDown(sum, 12)}
        </div>
      </div>
    `;

      showTooltip(
        e.currentTarget,
        tooltipContent,
        leftPos >= 245 ? TooltipPosition.LEFT : TooltipPosition.RIGHT,
        "204px"
      );
    },
    [
      orderbook,
      calculateGroupData,
      market.base,
      market.counter,
      showTooltip,
      hideTooltip,
      leftPos,
    ]
  );

  const removeHoverClasses = useCallback(() => {
    const elements = document.getElementsByClassName("orderbook-row");
    Array.from(elements).forEach((el) => {
      el.classList.remove("hovered-buy", "hovered-sell");
    });
  }, []);
  const renderOrderRow = (
    order: OrderbookRecord,
    index: number,
    type: ORDERBOOK_TYPE
  ) => {
    const isBuy = type === ORDERBOOK_TYPE.BUY;
    const volBar = Math.max(
      2,
      (Number(order.SymbolAmount) * 100) /
        (isBuy ? topBuyVolume : topSellVolume)
    );

    return (
      <div
        key={index}
        className="orderbook-row"
        data-value={`${order.HumanReadablePrice}_${type}_${order.Sequence}`}
        onMouseEnter={(e) => handleTooltip(e, index, type, true)}
        onMouseLeave={() => {
          removeHoverClasses();
          hideTooltip();
        }}
        onClick={() => {
          setOrderbookAction({
            type: isBuy ? OrderType.BUY : OrderType.SELL,
            price: Number(order.HumanReadablePrice),
            volume: Number(order.SymbolAmount),
          });
        }}
      >
        <div
          style={{ width: `${volBar}%` }}
          className={`volume-bar ${isBuy ? "buys" : "sells"}`}
        />
        <div className="orderbook-numbers-wrapper">
          <FormatNumber
            number={Number(order.HumanReadablePrice)}
            className={`orderbook-number price-${isBuy ? "buys" : "sells"}`}
          />
          <FormatNumber
            number={Number(order.SymbolAmount)}
            className="orderbook-number"
          />
          <FormatNumber
            number={
              Number(order.HumanReadablePrice) * Number(order.SymbolAmount)
            }
            className="orderbook-number"
          />
        </div>
      </div>
    );
  };

  return (
    <div className="orderbook-container" ref={componentRef}>
      <div className="orderbook-body">
        <div className="orderbook-header">
          <div className="title">Orderbook</div>
        </div>

        {orderbook ? (
          <>
            <div className="orderbook-header-wrapper">
              <div className="orderbook-header-cell">Price</div>
              <div className="orderbook-header-cell">Volume</div>
              <div className="orderbook-header-cell">
                Total Amount
                <span className="tooltip-total">Price * Volume</span>
              </div>
            </div>

            <div className="orderbook-sections">
              {(orderbookType === ORDERBOOK_TYPE.BUY ||
                orderbookType === ORDERBOOK_TYPE.BOTH) && (
                <div
                  className="orderbook-wrapper"
                  style={{ flexDirection: "column-reverse" }}
                  id="buys_ob"
                >
                  {orderbook.Buy.slice(0, 50).map((buy, i) =>
                    renderOrderRow(buy, i, ORDERBOOK_TYPE.BUY)
                  )}
                </div>
              )}

              {spread && (
                <div className="orderbook-spread">
                  <p className="spread-label">{`${market.counter.Denom.Name} Spread:`}</p>
                  <FormatNumber number={spread?.valueOf()} />
                </div>
              )}

              {(orderbookType === ORDERBOOK_TYPE.SELL ||
                orderbookType === ORDERBOOK_TYPE.BOTH) && (
                <div className="orderbook-wrapper" id="sells_ob">
                  {/* reverse, create array of original indices, take top 5, 
                  then map back to data to maintain proper hover and tooltip */}
                  {orderbook.Sell.map((_, idx) => idx)
                    .slice(0, 50)
                    .reverse()
                    .map((originalIndex) => {
                      const sell = orderbook.Sell[originalIndex];
                      return renderOrderRow(
                        sell,
                        originalIndex,
                        ORDERBOOK_TYPE.SELL
                      );
                    })}
                </div>
              )}
            </div>
          </>
        ) : (
          <div className="no-data-container">
            <img src="/trade/images/warning.png" alt="warning" />
            <p className="no-data">No Data Found</p>
          </div>
        )}
      </div>
    </div>
  );
}
