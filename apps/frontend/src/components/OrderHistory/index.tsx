import {
  useEffect,
  useState,
  useCallback,
  useMemo,
  useRef,
  useLayoutEffect,
} from "react";
import {
  OrderbookRecord,
  OrderbookResponse,
  TradeRecord,
  TransformedOrder,
} from "@/types/market";
import { useStore } from "@/state/store";
import { FormatNumber } from "../FormatNumber";
import {
  cancelOrder,
  getOrderbook,
  getTrades,
  submitOrder,
} from "@/services/api";
import { mergeUniqueTrades, resolveCoreumExplorer } from "@/utils";
import "./order-history.scss";
import { DEX } from "coreum-js";
import { TxRaw } from "coreum-js/dist/main/cosmos";
import { Side } from "coreum-js/dist/main/coreum/dex/v1/order";
import { fromByteArray } from "base64-js";
import { UpdateStrategy, wsManager, NetworkToEnum } from "@/services/websocket";
import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import debounce from "lodash/debounce";
import { FixedSizeList as List } from "react-window";
import { OrderStatus } from "coredex-api-types/order";
import { Method } from "coredex-api-types/update";
import AutoSizer from "react-virtualized-auto-sizer";
dayjs.extend(duration);

const TABS = {
  OPEN_ORDERS: "OPEN_ORDERS",
  ORDER_HISTORY: "ORDER_HISTORY",
};

const ROW_HEIGHT = 26;

const OrderHistory = () => {
  const {
    setOpenOrders,
    openOrders,
    market,
    wallet,
    network,
    orderHistory,
    setOrderHistory,
    pushNotification,
    coreum,
  } = useStore();

  const [activeTab, setActiveTab] = useState(TABS.OPEN_ORDERS);
  const [timeRange, setTimeRange] = useState({
    from: dayjs().subtract(1, "hour").unix(),
    to: dayjs().unix(),
  });
  const [isFetchingMore, setIsFetchingMore] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const listOuterRef = useRef<HTMLDivElement>(null);
  const currentMarketRef = useRef(market.pair_symbol);

  const resolveOrderStatus = (status: OrderStatus) => {
    switch (status) {
      case OrderStatus.ORDER_STATUS_OPEN:
        return "Open";
      case OrderStatus.ORDER_STATUS_EXPIRED:
        return "Expired";
      case OrderStatus.ORDER_STATUS_CANCELED:
        return "Cancelled";
      case OrderStatus.ORDER_STATUS_FILLED:
        return "Filled";
      default:
        return "Unspecified";
    }
  };

  // reset state on market change and init fetch
  useEffect(() => {
    wsManager.clearState();

    currentMarketRef.current = market.pair_symbol;
    setOrderHistory([]);
    setOpenOrders(null);
    setHasMore(true);
    setTimeRange({
      from: dayjs().subtract(1, "hour").unix(),
      to: dayjs().unix(),
    });

    if (wallet?.address) {
      const initFetch = async () => {
        await fetchOrderHistory();
      };
      initFetch();
    }
  }, [market.pair_symbol, wallet]);

  const fetchOrderHistory = async (): Promise<boolean> => {
    try {
      const response = await getTrades({
        symbol: market.pair_symbol,
        account: wallet?.address,
      });

      if (response.status === 200) {
        const data = response.data || [];
        setOrderHistory(data);
        wsManager.setInitialState(orderHistorySubscription, data);
        return data.length > 0;
      }

      setOrderHistory([]);
      wsManager.setInitialState(orderHistorySubscription, []);
      return false;
    } catch (e) {
      console.log("ERROR GETTING ORDER HISTORY DATA >>", e);
      setOrderHistory([]);
      wsManager.setInitialState(orderHistorySubscription, []);
      return false;
    }
  };

  const transformOrderbook = useCallback(
    (orderbook: OrderbookResponse): TransformedOrder[] => {
      const transformSide = (orders: OrderbookRecord[], side: Side) =>
        orders.map((order) => {
          return {
            Side: side,
            HumanReadablePrice: order.HumanReadablePrice,
            Price: order.Price,
            Amount: order.Amount,
            SymbolAmount: order.SymbolAmount,
            Total:
              Number(order.HumanReadablePrice) * Number(order.SymbolAmount),
            Account: order.Account,
            Sequence: order.Sequence,
            OrderID: order.OrderID,
            RemainingAmount: order.RemainingAmount,
            RemainingSymbolAmount: order.RemainingSymbolAmount,
          } as TransformedOrder;
        });

      return [
        ...transformSide(orderbook.Buy, Side.SIDE_BUY),
        ...transformSide(orderbook.Sell, Side.SIDE_SELL),
      ].sort((a, b) => Number(b.Sequence) - Number(a.Sequence));
    },
    []
  );

  // fetch open orders filtered from orderbook. transform to formatted data
  useEffect(() => {
    const fetchOpenOrders = async () => {
      try {
        if (!wallet?.address) return;
        const response = await getOrderbook(
          market.pair_symbol,
          wallet?.address
        );
        const data = response.data;
        if (data) {
          const openOrders = transformOrderbook(data);
          setOpenOrders(openOrders);
          wsManager.setInitialState(openOrderSubscription, openOrders);
        }
      } catch (e) {
        console.log("ERROR GETTING OPEN ORDERS DATA >>", e);
        setOpenOrders(null);
      }
    };

    fetchOpenOrders();
  }, [market.pair_symbol, wallet]);

  const openOrderSubscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT,
      ID: `${wallet ? wallet.address : ""}_${market.pair_symbol}`,
    }),
    [market.pair_symbol, wallet, network]
  );

  const orderHistorySubscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.TRADES_FOR_ACCOUNT_AND_SYMBOL,
      ID: `${wallet ? wallet.address : ""}_${market.pair_symbol}`,
    }),
    [market.pair_symbol, wallet, network]
  );

  const handleOpenOrders = useCallback(
    (message: OrderbookResponse) => {
      if (currentMarketRef.current === market.pair_symbol) {
        const updatedHistory = transformOrderbook(message);
        setOpenOrders(updatedHistory);
      }
    },
    [setOpenOrders, transformOrderbook, market.pair_symbol]
  );

  useEffect(() => {
    if (!wallet) return;
    wsManager.connected().then(() => {
      wsManager.subscribe(
        openOrderSubscription,
        handleOpenOrders,
        UpdateStrategy.REPLACE
      );
    });
    return () => {
      wsManager.unsubscribe(openOrderSubscription, handleOpenOrders);
    };
  }, [openOrderSubscription, wallet, handleOpenOrders, market.pair_symbol]);

  const orderHistoryHandler = useCallback(
    (newTrades: TradeRecord[]) => {
      if (currentMarketRef.current === market.pair_symbol) {
        const validTrades = newTrades.filter((trade) => {
          return (
            !trade ||
            `${trade.Denom1?.Denom}_${trade.Denom2?.Denom}` ===
              market.pair_symbol
          );
        });

        if (validTrades.length > 0) {
          setOrderHistory((prev) => mergeUniqueTrades(prev, validTrades));
        }
      }
    },
    [market.pair_symbol, setOrderHistory]
  );

  useEffect(() => {
    if (!wallet) return;
    wsManager.connected().then(() => {
      wsManager.subscribe(
        orderHistorySubscription,
        orderHistoryHandler,
        UpdateStrategy.MERGE
      );
    });
    return () => {
      wsManager.unsubscribe(orderHistorySubscription, orderHistoryHandler);
    };
  }, [
    orderHistorySubscription,
    wallet,
    orderHistoryHandler,
    market.pair_symbol,
  ]);

  // uses a 1hour window to not overload api with possibly huge requests
  // but this means we will not get new data when large gaps in between trades (eg. multiple days, weeks, etc)
  // should update api to use offset and limit if we want to load more data without setting a retry limit
  const loadOlderHistory = useCallback(async (): Promise<number> => {
    try {
      const currentWindow = timeRange.to - timeRange.from;
      const oldestTrade = orderHistory[orderHistory.length - 1];
      const oldestTimestamp = oldestTrade?.BlockTime?.seconds ?? timeRange.from;
      const newTo = oldestTimestamp - 1;
      const newFrom = newTo - currentWindow;
      const response = await getTrades({
        symbol: market.pair_symbol,
        from: newFrom,
        to: newTo,
        account: wallet?.address,
      });
      if (response.status === 200) {
        const olderData = response.data;
        if (!olderData || olderData.length === 0) {
          setHasMore(false);
          return 0;
        }
        const prevHistory = orderHistory || [];
        const mergedHistory = mergeUniqueTrades(prevHistory, olderData);
        wsManager.setInitialState(orderHistorySubscription, mergedHistory);
        setOrderHistory(mergedHistory);
        setTimeRange({ from: newFrom, to: newTo });
        return olderData.length;
      }
    } catch (e) {
      console.log("ERROR FETCHING OLDER HISTORY >>", e);
    }
    return 0;
  }, [
    timeRange,
    orderHistory,
    market.pair_symbol,
    wallet?.address,
    setHasMore,
    setOrderHistory,
    setTimeRange,
    orderHistorySubscription,
  ]);

  useLayoutEffect(() => {
    if (activeTab !== TABS.ORDER_HISTORY) return;
    if (!listOuterRef.current) return;
    const container = listOuterRef.current;
    const handleScroll = async () => {
      const threshold = 50;

      if (
        container.scrollTop + container.clientHeight >=
        container.scrollHeight - threshold
      ) {
        if (
          orderHistory &&
          orderHistory.length > 0 &&
          !isFetchingMore &&
          hasMore
        ) {
          setIsFetchingMore(true);
          const anchorEl = container.querySelector(".history-row");
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
    requestAnimationFrame(() => {
      container.addEventListener("scroll", debouncedHandleScroll);
    });
    return () => {
      container.removeEventListener("scroll", debouncedHandleScroll);
      debouncedHandleScroll.cancel();
    };
  }, [
    orderHistory,
    timeRange,
    market.pair_symbol,
    isFetchingMore,
    hasMore,
    activeTab,
    loadOlderHistory,
  ]);

  const handleCancelOrder = async (id: string) => {
    if (!wallet?.address) return;
    try {
      const orderCancel = {
        Sender: wallet.address,
        OrderID: id,
      };
      const orderCancelResponse = await cancelOrder(orderCancel);
      const cancelMessage = DEX.CancelOrder(
        orderCancelResponse.data.OrderCancel
      );
      const signedTx = await coreum?.signTx(
        [cancelMessage],
        undefined,
        orderCancelResponse.data.Sequence
      );
      const encodedTx = TxRaw.encode(signedTx!).finish();
      const base64Tx = fromByteArray(encodedTx);
      const submitResponse = await submitOrder({ TX: base64Tx });

      if (submitResponse.status !== 200) {
        pushNotification({
          type: "error",
          message: "There was an issue cancelling your order",
        });
        throw new Error("Error submitting order");
      }

      const txHash = submitResponse.data.TXHash;
      pushNotification({
        type: "success",
        message: `Order Cancelled! TXHash: ${txHash.slice(
          0,
          6
        )}...${txHash.slice(-4)}`,
      });
    } catch (e: any) {
      console.log("ERROR CANCELLING ORDER >>", e);
      pushNotification({
        type: "error",
        message: "Error cancelling order",
      });
    }
  };

  const Row = ({
    index,
    style,
  }: {
    index: number;
    style: React.CSSProperties;
  }) => {
    if (!orderHistory) return null;
    const order = orderHistory[index];
    return (
      <a
        style={style}
        className="history-row"
        href={`${resolveCoreumExplorer(network)}/transactions/${order.TXID}`}
        target="_blank"
      >
        <div className={order.Side === Side.SIDE_BUY ? "buy" : "sell"}>
          {order.Side === Side.SIDE_BUY ? "Buy" : "Sell"}
        </div>
        <div className="order-id">{order.Sequence}</div>
        <div className="status">{resolveOrderStatus(order.Status)}</div>
        <FormatNumber number={order.HumanReadablePrice} className="price" />
        <FormatNumber number={order.SymbolAmount} className="volume" />
        <FormatNumber
          number={Number(order.HumanReadablePrice) * Number(order.SymbolAmount)}
          className="total"
        />
        <p className="date">
          {dayjs.unix(order.BlockTime?.seconds ?? 0).format("MM/DD/YY HH:mm")}
        </p>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="21"
          height="21"
          viewBox="0 0 21 21"
          fill="none"
          className="external-link"
        >
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M2 5.941C2.00498 5.83929 2.03774 5.74086 2.09471 5.65645C2.15166 5.57191 2.23083 5.50453 2.32345 5.46193L10.2992 2.04202C10.43 1.98599 10.578 1.98599 10.7087 2.04202L18.6849 5.46208C18.7777 5.50308 18.8566 5.56994 18.9126 5.65471C18.9685 5.73949 18.9988 5.83885 19 5.94037L19 15.0602C18.9999 15.162 18.9698 15.2617 18.9135 15.3465C18.8572 15.4314 18.7771 15.4981 18.6834 15.5378L10.703 18.958C10.5728 19.014 10.4253 19.014 10.2952 18.958L2.31494 15.5378C2.22147 15.4978 2.14189 15.4312 2.08584 15.3463C2.02981 15.2614 2 15.162 2 15.0602V5.941ZM3.03967 14.7178L9.9805 17.6924V9.70242L3.03967 6.72761V14.7178ZM3.83948 5.93982L10.5002 8.79507L17.161 5.93982L10.5002 3.0851L3.83948 5.93982ZM11.0199 17.6924L17.9608 14.7178V6.72761L11.0199 9.70229V17.6924Z"
            fill="#5E6773"
          />
        </svg>
      </a>
    );
  };

  return (
    <div className="order-history-container">
      <div className="order-history-tabs">
        <div className="options">
          <div
            className={activeTab === TABS.OPEN_ORDERS ? "tab active" : "tab"}
            onClick={() => setActiveTab(TABS.OPEN_ORDERS)}
          >
            Open Orders
          </div>
          <div
            className={activeTab === TABS.ORDER_HISTORY ? "tab active" : "tab"}
            onClick={() => setActiveTab(TABS.ORDER_HISTORY)}
          >
            Order History
          </div>
        </div>
      </div>

      {!wallet?.address ? (
        <div className="no-orders">
          <img src="/trade/images/planet-graphic.svg" alt="" />
          Sign in with your wallet to view your orders.
        </div>
      ) : (
        <>
          <div
            className={
              activeTab === TABS.OPEN_ORDERS
                ? `open-orders-labels`
                : `order-history-labels`
            }
          >
            <div className="order-label">Side</div>
            <div className="order-label">OrderId</div>
            {activeTab === TABS.ORDER_HISTORY && (
              <div className="order-label">Status</div>
            )}
            <div className="order-label">Price</div>
            <div className="order-label">Volume</div>
            <div className="order-label">Total</div>
            {activeTab === TABS.OPEN_ORDERS && <div></div>}
            {activeTab === TABS.ORDER_HISTORY ? (
              <div className="order-label date">
                {activeTab === TABS.ORDER_HISTORY ? "Date" : "Time"}
              </div>
            ) : (
              <div></div>
            )}
          </div>

          <div className="order-history-body">
            {activeTab === TABS.OPEN_ORDERS ? (
              <div className="open-orders">
                {openOrders && openOrders.length > 0 ? (
                  openOrders.map((order: TransformedOrder, index) => {
                    return (
                      <div key={index} className="open-row">
                        <div
                          className={
                            order.Side === Side.SIDE_BUY ? `buy` : "sell"
                          }
                        >
                          {order.Side === Side.SIDE_BUY
                            ? "Buy"
                            : order.Side === Side.SIDE_SELL
                            ? "Sell"
                            : "Unspecified"}
                        </div>
                        <div className="order-id"> {order.Sequence}</div>
                        <FormatNumber
                          number={order.HumanReadablePrice}
                          className="price"
                        />
                        <FormatNumber
                          number={order.RemainingSymbolAmount}
                          className="volume"
                        />
                        <FormatNumber number={order.Total} className="total" />
                        <div
                          className="cancel-order-container"
                          onClick={() => {
                            handleCancelOrder(order.OrderID);
                          }}
                        >
                          <svg
                            className="cancel-order"
                            xmlns="http://www.w3.org/2000/svg"
                            width="21"
                            height="21"
                            viewBox="0 0 21 21"
                            fill="none"
                          >
                            <path
                              fillRule="evenodd"
                              clipRule="evenodd"
                              d="M4.61205 3.27658C4.24327 2.90781 3.64536 2.90781 3.27658 3.27658C2.90781 3.64536 2.90781 4.24327 3.27658 4.61205L9.16453 10.5L3.27658 16.3879C2.90781 16.7567 2.90781 17.3546 3.27658 17.7234C3.64536 18.0922 4.24327 18.0922 4.61205 17.7234L10.5 11.8355L16.3877 17.7232C16.7565 18.092 17.3544 18.092 17.7232 17.7232C18.092 17.3544 18.092 16.7565 17.7232 16.3877L11.8355 10.5L17.7232 4.61228C18.092 4.2435 18.092 3.64559 17.7232 3.27681C17.3544 2.90803 16.7565 2.90803 16.3877 3.27681L10.5 9.16453L4.61205 3.27658Z"
                              fill="#5E6773"
                            />
                          </svg>
                        </div>
                      </div>
                    );
                  })
                ) : (
                  <div className="no-orders">
                    <img src="/trade/images/planet-graphic.svg" alt="" />
                    You have no orders!
                  </div>
                )}
              </div>
            ) : (
              <div className="order-history">
                {orderHistory && orderHistory.length > 0 ? (
                  <AutoSizer>
                    {({ height, width }) => (
                      <List
                        height={height}
                        itemCount={orderHistory.length}
                        itemSize={ROW_HEIGHT}
                        width={width}
                        outerRef={listOuterRef}
                      >
                        {Row}
                      </List>
                    )}
                  </AutoSizer>
                ) : (
                  <div className="no-orders">
                    <img src="/trade/images/planet-graphic.svg" alt="" />
                    You have no orders!
                  </div>
                )}
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
};

export default OrderHistory;
