import { useEffect, useState, useCallback, useMemo } from "react";
import {
  OrderbookRecord,
  OrderbookResponse,
  SIDE_BUY,
  TradeHistoryResponse,
  TransformedOrder,
} from "@/types/market";
import { useStore } from "@/state/store";
import { FormatNumber } from "../FormatNumber";
import { cancelOrder, getOrderbook, getTrades } from "@/services/general";
import { Method, NetworkToEnum, WebSocketMessage } from "@/services/websocket";
import { useWebSocket } from "@/hooks/websocket";
import { resolveCoreumExplorer } from "@/utils";
import "./order-history.scss";
import Modal from "../Modal";
import Button, { ButtonVariant } from "../Button";

const TABS = {
  OPEN_ORDERS: "OPEN_ORDERS",
  ORDER_HISTORY: "ORDER_HISTORY",
};

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
  } = useStore();

  const [activeTab, setActiveTab] = useState(TABS.OPEN_ORDERS);
  const [cancelOrderModal, setCancelOrderModal] = useState(false);
  const [cancelOrderId, setCancelOrderId] = useState("");

  // fetch order history
  useEffect(() => {
    const fetchExchangeHistory = async () => {
      if (!wallet?.address) return;
      const from = new Date().getTime();
      const to = from - 2592000000; // 30 days ago
      try {
        const response = await getTrades(
          market.pair_symbol,
          to,
          from,
          wallet?.address
        );
        if (response.status === 200) {
          const data = response.data;
          setOrderHistory(data);
        }
      } catch (e) {
        console.log("ERROR GETTING ORDER HISTORY DATA >>", e);
        setOrderHistory(null);
      }
    };
    fetchExchangeHistory();
  }, [market.pair_symbol, wallet]);

  // fetch open orders filtered from orderbook. transform to formatted data
  useEffect(() => {
    const fetchOpenOrders = async () => {
      try {
        const response = await getOrderbook(
          market.pair_symbol,
          wallet?.address
        );
        const data = response.data;
        if (data) {
          const openOrders = transformOrderbook(data);
          setOpenOrders(openOrders);
        }
      } catch (e) {
        console.log("ERROR GETTING OPEN ORDERS DATA >>", e);
        setOpenOrders(null);
      }
    };

    fetchOpenOrders();
  }, [market.pair_symbol, wallet]);

  const transformOrderbook = (
    orderbook: OrderbookResponse
  ): TransformedOrder[] => {
    const transformSide = (
      orders: OrderbookRecord[],
      side: SIDE_BUY.BUY | SIDE_BUY.SELL
    ) =>
      orders.map(
        (order) =>
          ({
            Side: side,
            Price: order.HumanReadablePrice,
            Volume: order.SymbolAmount,
            Total:
              Number(order.HumanReadablePrice) * Number(order.SymbolAmount),
            Account: order.Account,
            Sequence: order.Sequence,
            OrderID: order.OrderID,
          } as TransformedOrder)
      );

    return [
      ...transformSide(orderbook.Buy, SIDE_BUY.BUY),
      ...transformSide(orderbook.Sell, SIDE_BUY.SELL),
    ].sort((a, b) => a.Sequence - b.Sequence);
  };

  // TODO: move to ws service
  const handleOrderHistory = useCallback(
    (message: WebSocketMessage) => {
      const data = message.Subscription?.Content;
      if (data.length > 0) {
        if (orderHistory) {
          const updatedHistory: TradeHistoryResponse =
            orderHistory.concat(data);
          setOrderHistory(updatedHistory);
        } else {
          setOrderHistory(data);
        }
      }
    },
    [setOpenOrders]
  );

  const orderHistorySubscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.TRADES_FOR_ACCOUNT_AND_SYMBOL,
      ID: `${wallet ? wallet.address : ""}_${market.pair_symbol}`,
    }),
    [market.pair_symbol, wallet]
  );

  const handleOpenOrders = useCallback(
    (message: WebSocketMessage) => {
      const data = message.Subscription?.Content;
      if (data) {
        const updatedHistory = transformOrderbook(data);
        setOpenOrders(updatedHistory);
      }
    },
    [setOpenOrders, transformOrderbook]
  );

  const openOrderSubscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT,
      ID: `${wallet ? wallet.address : ""}_${market.pair_symbol}`,
    }),
    [market.pair_symbol, wallet]
  );

  useWebSocket(orderHistorySubscription, handleOrderHistory);
  useWebSocket(openOrderSubscription, handleOpenOrders);

  const handleCancelOrder = async (id: string) => {
    if (!wallet?.address) return;
    try {
      const response = await cancelOrder(wallet.address, id);
      if (response.status === 200 && response.data) {
        try {
          const tx = response.data.TXBytes;
          await navigator.clipboard.writeText(tx);

          pushNotification({
            type: "success",
            message: `Order Cancelled! TXHash copied to clipboard: ${tx.slice(
              0,
              6
            )}...${tx.slice(-4)}`,
          });
        } catch (copyError) {
          console.error("Copy failed:", copyError);
        }
      }
    } catch (e) {
      console.log("ERROR CANCELLING ORDER >>", e);
      pushNotification({
        type: "error",
        message: "Error cancelling order",
      });
    }
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
              <>
                {openOrders && openOrders.length > 0 ? (
                  openOrders.map((order: TransformedOrder, index) => {
                    return (
                      <div key={index} className="open-row">
                        <div
                          className={
                            order.Side === SIDE_BUY.BUY ? `buy` : "sell"
                          }
                        >
                          {order.Side === SIDE_BUY.BUY
                            ? "Buy"
                            : order.Side === SIDE_BUY.SELL
                            ? "Sell"
                            : "Unspecified"}
                        </div>
                        <div className="order-id"> {order.Sequence}</div>
                        <FormatNumber
                          number={order.Price}
                          precision={5}
                          className="price"
                        />
                        <FormatNumber
                          number={order.Volume}
                          precision={4}
                          className="volume"
                        />
                        <FormatNumber
                          number={order.Total}
                          precision={4}
                          className="total"
                        />
                        <div
                          className="cancel-order-container"
                          onClick={() => {
                            setCancelOrderModal(true);
                            setCancelOrderId(order.OrderID);
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
              </>
            ) : (
              <>
                {orderHistory && orderHistory.length > 0 ? (
                  orderHistory.map((order, index) => {
                    return (
                      <a
                        key={index}
                        className="history-row"
                        href={`${resolveCoreumExplorer(network)}/transactions/${
                          order.TXID
                        }`}
                      >
                        <div
                          className={
                            order.Side === SIDE_BUY.BUY ? `buy` : "sell"
                          }
                        >
                          {order.Side === SIDE_BUY.BUY ? "Buy" : "Sell"}
                        </div>
                        <div className="order-id"> {order.Sequence}</div>
                        <div className="status">TODO</div>
                        <FormatNumber
                          number={order.HumanReadablePrice}
                          precision={5}
                          className="price"
                        />
                        <FormatNumber
                          number={order.SymbolAmount}
                          precision={4}
                          className="volume"
                        />
                        <FormatNumber
                          number={
                            Number(order.HumanReadablePrice) *
                            Number(order.SymbolAmount)
                          }
                          precision={4}
                          className="total"
                        />
                        <p className="date">
                          {new Date(
                            order.BlockTime.seconds * 1000
                          ).toLocaleString()}
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
                  })
                ) : (
                  <div className="no-orders">
                    <img src="/trade/images/planet-graphic.svg" alt="" />
                    You have no orders!
                  </div>
                )}
              </>
            )}
          </div>
        </>
      )}
      <Modal
        isOpen={cancelOrderModal}
        onClose={() => {
          setCancelOrderModal(false);
          setCancelOrderId("");
        }}
        title="Cancel Open Order"
        children={
          <div className="cancel-order">
            <p className="cancel-order-description">
              Do you want to cancel this open order?
            </p>
            <div className="cancel-order-btns">
              <Button
                variant={ButtonVariant.PRIMARY}
                onClick={() => {
                  handleCancelOrder(cancelOrderId);
                  setCancelOrderModal(false);
                  setCancelOrderId("");
                }}
                width={"100%"}
                height={37}
                label="Confirm"
              />
              <Button
                variant={ButtonVariant.DANGER}
                onClick={() => {
                  setCancelOrderModal(false);
                  setCancelOrderId("");
                }}
                width={"100%"}
                height={37}
                label="Cancel"
              />
            </div>
          </div>
        }
      />
    </div>
  );
};

export default OrderHistory;
