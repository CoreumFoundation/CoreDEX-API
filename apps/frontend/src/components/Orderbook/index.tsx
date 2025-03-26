import { useState, useEffect, useRef, useCallback, useMemo } from "react";
import BigNumber from "bignumber.js";
import { useStore } from "@/state/store";
import { useTooltip, TooltipPosition } from "@/hooks";
import { toFixedDown } from "@/utils";
import { FormatNumber } from "../FormatNumber";
import { OrderbookAction, OrderbookRecord } from "@/types/market";
import { getOrderbook } from "@/services/api";
import { wsManager, UpdateStrategy, NetworkToEnum } from "@/services/websocket";
import "./orderbook.scss";
import { Side } from "coreum-js-nightly/dist/main/coreum/dex/v1/order";
import { Method } from "coredex-api-types/update";
import { mirage } from "ldrs";
mirage.register();

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
  const { market, network, orderbook, setOrderbook } = useStore();
  const { showTooltip, hideTooltip } = useTooltip();
  const [isLoading, setIsLoading] = useState(false);
  const [spread, setSpread] = useState<BigNumber>(new BigNumber(0));
  const [topBuyVolume, setTopBuyVolume] = useState<number>(0);
  const [topSellVolume, setTopSellVolume] = useState<number>(0);
  const [totalVolume, setTotalVolume] = useState<number>(0);
  const componentRef = useRef<HTMLDivElement>(null);
  const sellsObRef = useRef<HTMLDivElement>(null);

  const subscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.ORDERBOOK,
      ID: market.pair_symbol,
    }),
    [market.pair_symbol, network]
  );

  useEffect(() => {
    wsManager.connected().then(() => {
      wsManager.subscribe(subscription, setOrderbook, UpdateStrategy.REPLACE);
    });
    return () => {
      wsManager.unsubscribe(subscription, setOrderbook);
    };
  }, [subscription]);

  useEffect(() => {
    const fetchOrderbook = async () => {
      try {
        setIsLoading(true);
        const response = await getOrderbook(market.pair_symbol);
        if (response.status === 200 && response.data) {
          const data = response.data;
          setOrderbook(data);
        }
        setIsLoading(false);
      } catch (e) {
        console.log("ERROR GETTING ORDERBOOK DATA >>", e);
        setOrderbook(null);
        setIsLoading(false);
      }
    };
    fetchOrderbook();
  }, [market.base, market.counter]);

  // calculate spread
  useEffect(() => {
    if (!orderbook || !orderbook.Buy || !orderbook.Sell) return;

    const calculateSpread = () => {
      const highestBuy = orderbook.Buy[0]?.HumanReadablePrice;
      const lowestSell =
        orderbook.Sell[orderbook.Sell.length - 1]?.HumanReadablePrice;

      if (!highestBuy && !lowestSell) return new BigNumber(0);

      if (highestBuy && lowestSell) {
        return new BigNumber(lowestSell).minus(highestBuy);
      }

      return new BigNumber(lowestSell || highestBuy || 0);
    };

    setSpread(calculateSpread());
  }, [orderbook, market.pair_symbol]);

  // scroll sells to bottom
  useEffect(() => {
    if (sellsObRef.current && orderbook?.Sell) {
      setTimeout(() => {
        sellsObRef.current!.scrollTop = sellsObRef.current!.scrollHeight;
      }, 50);
    }
  }, [orderbook?.Sell, market]);

  // find the highest volume in the orderbook
  useEffect(() => {
    if (orderbook?.Buy?.length) {
      const highestBuyVolume = orderbook.Buy.reduce((max, buy) => {
        const current = new BigNumber(buy.RemainingSymbolAmount);
        return current.isGreaterThan(max) ? current : max;
      }, new BigNumber(0)).toNumber();
      setTopBuyVolume(highestBuyVolume);
    }

    if (orderbook?.Sell?.length) {
      const highestSellVolume = orderbook.Sell.reduce((max, sell) => {
        const current = new BigNumber(sell.RemainingSymbolAmount);
        return current.isGreaterThan(max) ? current : max;
      }, new BigNumber(0)).toNumber();
      setTopSellVolume(highestSellVolume);
    }
  }, [orderbook]);

  const calculateGroupData = useCallback(
    (lines: OrderbookRecord[], index: number, orderType: ORDERBOOK_TYPE) => {
      let avgPriceSum = 0;
      let totalVolumeCalc = 0;
      let sum = 0;
      const lineGroup = [];

      if (orderType === ORDERBOOK_TYPE.BUY) {
        for (let i = index; i >= 0; i--) {
          const line = lines[i];
          if (!line) break;
          lineGroup.push(line);
          avgPriceSum += Number(line.HumanReadablePrice);
          totalVolumeCalc += Number(line.RemainingSymbolAmount);
          sum +=
            Number(line.HumanReadablePrice) *
            Number(line.RemainingSymbolAmount);
        }
      } else {
        for (let i = index; i < lines.length; i++) {
          const line = lines[i];
          if (!line) break;
          lineGroup.push(line);
          avgPriceSum += Number(line.HumanReadablePrice);
          totalVolumeCalc += Number(line.RemainingSymbolAmount);
          sum +=
            Number(line.HumanReadablePrice) *
            Number(line.RemainingSymbolAmount);
        }
      }

      setTotalVolume(totalVolumeCalc);

      return {
        avgPrice: lineGroup.length ? avgPriceSum / lineGroup.length : 0,
        sum,
        lineGroup,
        totalVolume: totalVolumeCalc,
      };
    },
    [setTotalVolume]
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
            <p class="inline-item-label">Sum (${
              market.counter.Denom.Name
            }):</p> ${toFixedDown(totalVolume, 12)}
          </div>
          <div class="inline-item">
            <p class="inline-item-label">Total Amount:</p> ${toFixedDown(
              sum,
              12
            )}
          </div>
        </div>
      `;

      showTooltip(
        e.currentTarget,
        tooltipContent,
        TooltipPosition.RIGHT,
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
      (Number(order.RemainingSymbolAmount) * 100) /
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
            type: isBuy ? Side.SIDE_SELL : Side.SIDE_BUY,
            price: Number(order.HumanReadablePrice),
            volume: totalVolume,
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
            number={Number(order.RemainingSymbolAmount)}
            className="orderbook-number"
          />
          <FormatNumber
            number={
              Number(order.HumanReadablePrice) *
              Number(order.RemainingSymbolAmount)
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
              <div className="orderbook-wrapper" id="sells_ob" ref={sellsObRef}>
                {orderbook.Sell &&
                  orderbook.Sell.slice(0, 50).map((order, index) =>
                    renderOrderRow(order, index, ORDERBOOK_TYPE.SELL)
                  )}
              </div>

              {spread && (
                <div className="orderbook-spread">
                  <p className="spread-label">{`${market.counter.Denom.Name} Spread:`}</p>
                  <FormatNumber number={spread?.toNumber()} />
                </div>
              )}

              <div className="orderbook-wrapper" id="buys_ob">
                {orderbook.Buy &&
                  orderbook.Buy.slice(0, 50).map((order, index) =>
                    renderOrderRow(order, index, ORDERBOOK_TYPE.BUY)
                  )}
              </div>
            </div>
          </>
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
    </div>
  );
}
