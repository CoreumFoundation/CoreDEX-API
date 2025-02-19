import { useEffect, useMemo, useState } from "react";
import { useStore } from "@/state/store";
import {
  OrderType,
  TradeType,
  OrderbookAction,
  WalletAsset,
  TimeInForceString,
  TimeSelection,
  TimeInForceStringToEnum,
  WalletBalances,
} from "@/types/market";
import { getAvgPriceFromOBbyVolume, multiply, noExponents } from "@/utils";
import { FormatNumber } from "../FormatNumber";
import { Input, InputType } from "../Input";
import Button, { ButtonVariant } from "../Button";
import BigNumber from "bignumber.js";
import { submitOrder, getWalletAssets, createOrder } from "@/services/api";
import { DEX } from "coreum-js-nightly";
import { TxRaw } from "coreum-js-nightly/dist/main/cosmos";
import "./order-actions.scss";
import {
  Side,
  OrderType as OT,
  TimeInForce,
} from "coreum-js-nightly/dist/main/coreum/dex/v1/order";
import { MsgPlaceOrder } from "coreum-js-nightly/dist/main/coreum/dex/v1/tx";
import { fromByteArray } from "base64-js";
import Dropdown, { DropdownVariant } from "../Dropdown";
import { DatetimePicker } from "../DatetimePicker";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import advancedFormat from "dayjs/plugin/advancedFormat";
import {
  Method,
  NetworkToEnum,
  UpdateStrategy,
  wsManager,
} from "@/services/websocket";

dayjs.extend(utc);
dayjs.extend(advancedFormat);
BigNumber.config({ DECIMAL_PLACES: 30, EXPONENTIAL_AT: 0 });

const OrderActions = ({
  orderbookAction,
}: {
  orderbookAction?: OrderbookAction;
}) => {
  const {
    orderbook,
    wallet,
    setLoginModal,
    pushNotification,
    market,
    coreum,
    network,
  } = useStore();

  const [orderType, setOrderType] = useState<OrderType>(OrderType.BUY);
  const [totalPrice, setTotalPrice] = useState(0);
  const [limitPrice, setLimitPrice] = useState("");
  const [volume, setVolume] = useState<string>("");
  const [tradeType, setTradeType] = useState(TradeType.MARKET);
  const [walletBalances, setWalletBalances] = useState<WalletBalances | null>(
    null
  );
  const [marketBalances, setMarketBalances] = useState({
    base: "0",
    counter: "0",
  });
  const [advSettingsOpen, setAdvSetting] = useState<boolean>(false);
  const [timeInForce, setTimeInForce] = useState<TimeInForceString>(
    TimeInForceString.goodTilCancel
  );
  const [goodTilValue, setGoodTilValue] = useState<number>(1);
  const [goodTilUnit, setGoodTilUnit] = useState<string>("Minutes");
  const [expirationTime, setExpirationTime] = useState<Date>();
  const [customTime, setCustomTime] = useState<string>("");
  const [blockHeight, setBlockHeight] = useState<number>(0);

  useEffect(() => {
    fetchWalletAssets();
  }, [wallet, market]);

  const walletSubscription = useMemo(
    () => ({
      Network: NetworkToEnum(network),
      Method: Method.WALLET,
      ID: `${wallet ? wallet.address : ""}`,
    }),
    [market.pair_symbol, wallet]
  );

  const handleWalletUpdate = (message: WalletBalances) => {
    if (message.length > 0) {
      setWalletBalances(message);
    }
  };

  useEffect(() => {
    wsManager.connected().then(() => {
      wsManager.subscribe(
        walletSubscription,
        handleWalletUpdate,
        UpdateStrategy.REPLACE
      );
    });
    return () => {
      wsManager.unsubscribe(walletSubscription, setWalletBalances);
    };
  }, [walletSubscription]);

  useEffect(() => {
    if (!walletBalances) return;

    const baseBalanceObject = walletBalances.find(
      (asset: WalletAsset) => asset.Denom === market.base.Denom.Denom
    );
    const counterBalanceObject = walletBalances.find(
      (asset: WalletAsset) => asset.Denom === market.counter.Denom.Denom
    );

    setMarketBalances({
      base: baseBalanceObject?.SymbolAmount || "0",
      counter: counterBalanceObject?.SymbolAmount || "0",
    });
  }, [market, walletBalances]);

  const fetchWalletAssets = async () => {
    if (!wallet?.address) return;
    try {
      const response = await getWalletAssets(wallet?.address);
      if (response.status === 200 && response.data.length > 0) {
        const data = response.data;
        setWalletBalances(data);
        wsManager.setInitialState(walletSubscription, data);
      }
    } catch (e) {
      console.log("ERROR GETTING WALLET ASSETS DATA >>", e);
    }
  };

  // trigger when click on orderbook
  useEffect(() => {
    if (orderbookAction?.price) {
      setTradeType(TradeType.LIMIT);
      setOrderType(orderbookAction.type);

      const volumeBN = new BigNumber(orderbookAction.volume);
      const priceBN = new BigNumber(orderbookAction.price);

      setVolume(volumeBN.toString());
      setLimitPrice(priceBN.toString());
      setTotalPrice(priceBN.times(volumeBN).toNumber());
    }
  }, [orderbookAction]);

  useEffect(() => {
    if (tradeType === TradeType.LIMIT) {
      const vol = volume ? Number(volume) : 0;

      const total = multiply(Number(limitPrice), vol);
      BigNumber(volume)
        .multipliedBy(
          new BigNumber(10).exponentiatedBy(market.base.Denom.Precision)
        )
        .toFixed();
      setTotalPrice(
        !total.isNaN()
          ? Number(noExponents(Number(total)).replaceAll(",", ""))
          : 0
      );
    }

    if (orderbook) {
      if (tradeType === TradeType.MARKET) {
        const avgPrice = Number(
          getAvgPriceFromOBbyVolume(
            orderType === OrderType.BUY ? orderbook.Buy : orderbook.Sell,
            volume
          )
        );

        const vol = volume ? Number(volume) : 0;

        const expRegex = /[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)/g;
        const total = multiply(avgPrice, vol).toNumber();
        setTotalPrice(
          !expRegex.test(total.toString()) && !isNaN(total) ? total : 0
        );
      }
    }
  }, [volume, limitPrice, orderbook, tradeType, orderType]);

  useEffect(() => {
    const tomorrow = dayjs.utc().add(1, "days");

    tomorrow.second(0);
    tomorrow.minute(0);
    tomorrow.hour(12);

    setCustomTime(new Date(tomorrow.format()).toString());
  }, []);

  useEffect(() => {
    if (timeInForce === TimeInForceString.goodTilTime) {
      let now = dayjs.utc();

      switch (goodTilUnit) {
        case "Minutes":
          now = dayjs.utc().add(goodTilValue, "minutes");
          break;
        case "Hours":
          now = dayjs.utc().add(goodTilValue, "hours");
          break;
        case "Days":
          now = dayjs.utc().add(goodTilValue, "days");
          break;
        case "Custom":
          now = dayjs.utc(customTime);
          break;
      }
      setExpirationTime(now.toDate());
    }
  }, [timeInForce, goodTilUnit, customTime, goodTilValue]);

  // format price for regex according to coreum backend
  // 1.5 -> 15e-1 or 1e+1 -> 10
  const formatPriceForRegex = (value: BigNumber): string => {
    let [mantissa, exponent = ""] = value.toExponential().split("e");
    exponent = exponent.replace(/^\+/, "");

    if (mantissa.includes(".")) {
      const decimalIndex = mantissa.indexOf(".");
      const decimalPlaces = mantissa.length - decimalIndex - 1;
      mantissa = mantissa.replace(".", "");
      const adjustedExponent = (parseInt(exponent, 10) || 0) - decimalPlaces;
      exponent = adjustedExponent.toString();
    }

    let processedExponent = "";
    if (exponent) {
      const exponentMatch = exponent.match(/^(-?)(\d+)$/);
      if (!exponentMatch) {
        throw new Error(`Invalid exponent: ${exponent}`);
      }

      const [_, sign, digits] = exponentMatch;
      const trimmedDigits = digits.replace(/^0+/, "") || "0";
      const isZeroExponent = trimmedDigits === "0";

      if (!isZeroExponent) {
        processedExponent = `${sign}${trimmedDigits}`;
      }
    }

    let result = mantissa;
    if (processedExponent) {
      result += `e${processedExponent}`;
    }

    return result;
  };

  const handleSubmit = async () => {
    try {
      const orderCreate: MsgPlaceOrder = {
        sender: wallet.address,
        type:
          tradeType === TradeType.LIMIT
            ? OT.ORDER_TYPE_LIMIT
            : OT.ORDER_TYPE_MARKET,
        id: crypto.randomUUID(),
        baseDenom: market.base.Denom.Denom,
        quoteDenom: market.counter.Denom.Denom,
        price:
          tradeType === TradeType.LIMIT
            ? formatPriceForRegex(BigNumber(limitPrice))
            : "",
        quantity: BigNumber(volume)
          .multipliedBy(
            new BigNumber(10).exponentiatedBy(market.base.Denom.Precision)
          )
          .toFixed(),
        side: orderType === OrderType.BUY ? Side.SIDE_BUY : Side.SIDE_SELL,
        goodTil:
          tradeType === TradeType.LIMIT &&
          timeInForce === TimeInForceString.goodTilTime
            ? {
                goodTilBlockTime: expirationTime,
                goodTilBlockHeight: blockHeight,
              }
            : undefined,
        timeInForce:
          tradeType === TradeType.LIMIT
            ? (TimeInForceStringToEnum[timeInForce] as any)
            : TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
      };

      const test = createOrder(orderCreate);
      console.log(test);

      // const orderMessage = DEX.PlaceOrder(orderCreate);
      // const signedTx = await coreum?.signTx([orderMessage]);
      // const encodedTx = TxRaw.encode(signedTx!).finish();
      // const base64Tx = fromByteArray(encodedTx);
      // const submitResponse = await submitOrder({ TX: base64Tx });

      // if (submitResponse.status !== 200) {
      //   pushNotification({
      //     type: "error",
      //     message: "There was an issue submitting your order",
      //   });
      //   throw new Error("Error submitting order");
      // }
      // const txHash = submitResponse.data.TXHash;
      // pushNotification({
      //   type: "success",
      //   message: `Order Placed! TXHash: ${txHash.slice(0, 6)}...${txHash.slice(
      //     -4
      //   )}`,
      // });

      // setTimeInForce(TimeInForceString.goodTilCancel);
      // setGoodTilValue(1);
      // setVolume("");
      // setLimitPrice("");
    } catch (e: any) {
      console.log("ERROR HANDLING SUBMIT ORDER >>", e);
      pushNotification({
        type: "error",
        message: e.error.message || e.message,
      });
      throw e;
    }
  };

  return (
    <div className="order-actions-container">
      <div className="order-actions-content" style={{ padding: "16px" }}>
        <div className="order-top">
          <div className="order-switch">
            <div
              className={`switch switch-buy ${
                orderType === OrderType.BUY ? "active" : ""
              }`}
              onClick={() => setOrderType(OrderType.BUY)}
            >
              <p>Buy</p>
            </div>

            <div
              className={`switch switch-sell ${
                orderType === OrderType.SELL ? "active" : ""
              }`}
              onClick={() => setOrderType(OrderType.SELL)}
            >
              <p>Sell</p>
            </div>
          </div>
          <div className="order-trade">
            <div className="order-trade-types">
              <div
                className={`type-item ${
                  tradeType === TradeType.MARKET ? "active" : ""
                }`}
                onClick={() => {
                  setTradeType(TradeType.MARKET);
                }}
              >
                Market
              </div>
              <div
                className={`type-item ${
                  tradeType === TradeType.LIMIT ? "active" : ""
                }`}
                onClick={() => {
                  setTradeType(TradeType.LIMIT);
                }}
              >
                Limit
              </div>
            </div>
          </div>

          <div className="order-trade">
            {tradeType === TradeType.LIMIT ? (
              <div className="limit-type-wrapper">
                <Input
                  maxLength={16}
                  placeholder="Enter Amount"
                  type={InputType.NUMBER}
                  onValueChange={(val: string) => {
                    setVolume(val);
                  }}
                  value={volume}
                  inputName="volume"
                  label="Amount"
                  customCss={{
                    fontSize: 14,
                  }}
                  inputWrapperClassname="order-input"
                  decimals={market.base.Denom.Precision}
                  adornmentRight={market.base.Denom.Currency}
                />
                <Input
                  maxLength={16}
                  placeholder="Enter Limit Price"
                  type={InputType.NUMBER}
                  onValueChange={(val: string) => {
                    setLimitPrice(val);
                  }}
                  value={limitPrice}
                  inputName="limit-price"
                  label="Price"
                  customCss={{
                    fontSize: 14,
                  }}
                  inputWrapperClassname="order-input"
                  decimals={market.counter.Denom.Precision}
                  adornmentRight={market.counter.Denom.Currency}
                />

                <div className="advanced-settings-header">
                  <div
                    className="advanced-accordion"
                    onClick={() => {
                      setAdvSetting(!advSettingsOpen);
                    }}
                  >
                    <p
                      className={`advanced-label ${
                        advSettingsOpen ? "active" : ""
                      }`}
                    >
                      Advanced Settings
                    </p>
                    <img
                      className={`advanced-arrow ${
                        advSettingsOpen ? "active" : ""
                      }`}
                      src="/trade/images/arrow.svg"
                      alt=""
                    />
                  </div>

                  <div
                    className={`advanced-settings-content ${
                      advSettingsOpen ? "open" : ""
                    }`}
                  >
                    <div className="time-in-force">
                      <p className="time-in-force-label">Time in Force</p>
                      <Dropdown
                        variant={DropdownVariant.OUTLINED}
                        items={(
                          Object.keys(TimeInForceString) as Array<
                            keyof typeof TimeInForceString
                          >
                        ).map((key) => [TimeInForceString[key]])}
                        value={timeInForce}
                        onClick={(item) => {
                          setTimeInForce(item[0] as TimeInForceString);
                        }}
                        renderItem={(item) => <div>{item}</div>}
                      />

                      {timeInForce === TimeInForceString.goodTilTime && (
                        <>
                          {
                            <div className="good-til-time">
                              <div className="time-selector">
                                <img
                                  src="/trade/images/arrow.svg"
                                  alt=""
                                  style={{
                                    transform: "rotate(90deg)",
                                  }}
                                  onClick={() =>
                                    goodTilUnit !== TimeSelection.CUSTOM &&
                                    setGoodTilValue((prev) =>
                                      Math.max(prev - 1, 1)
                                    )
                                  }
                                />

                                <div className="good-til-values">
                                  {goodTilUnit !== TimeSelection.CUSTOM && (
                                    <span className="time-value">
                                      {goodTilValue}
                                    </span>
                                  )}

                                  <span className="time-unit">
                                    {goodTilUnit}
                                  </span>
                                </div>

                                <img
                                  src="/trade/images/arrow.svg"
                                  alt=""
                                  style={{
                                    transform: "rotate(-90deg)",
                                  }}
                                  onClick={() =>
                                    goodTilUnit !== TimeSelection.CUSTOM &&
                                    setGoodTilValue((prev) => prev + 1)
                                  }
                                />
                              </div>

                              <div className="unit-selector">
                                <div
                                  className="unit"
                                  onClick={() => setGoodTilUnit("Minutes")}
                                >
                                  mins
                                </div>
                                <div
                                  className="unit"
                                  onClick={() => setGoodTilUnit("Hours")}
                                >
                                  hrs
                                </div>
                                <div
                                  className="unit"
                                  onClick={() => setGoodTilUnit("Days")}
                                >
                                  day
                                </div>
                                <div
                                  className="unit"
                                  onClick={() => setGoodTilUnit("Custom")}
                                >
                                  custom
                                </div>
                              </div>
                            </div>
                          }

                          {goodTilUnit === TimeSelection.CUSTOM && (
                            <div className="custom-time">
                              <DatetimePicker
                                selectedDate={customTime}
                                onChange={(val: any) => setCustomTime(val)}
                                width="100%"
                                minDate={
                                  new Date(dayjs.utc().add(1, "day").format())
                                }
                              />
                              <Input
                                maxLength={16}
                                placeholder="Block Height"
                                type={InputType.NUMBER}
                                onValueChange={(val: any) => {
                                  setBlockHeight(val);
                                }}
                                value={blockHeight}
                                inputName="limit-price"
                                label="Block Height"
                                customCss={{
                                  fontSize: 14,
                                }}
                                inputWrapperClassname="order-input"
                                decimals={0}
                              />
                            </div>
                          )}
                        </>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ) : (
              <div className="market-type-wrapper">
                <Input
                  maxLength={16}
                  placeholder="Enter Amount"
                  label="Amount"
                  type={InputType.NUMBER}
                  onValueChange={(val: string) => {
                    setVolume(val);
                  }}
                  value={volume}
                  inputName="volume"
                  customCss={{
                    fontSize: 16,
                  }}
                  decimals={13}
                  adornmentRight={market.base.Denom.Currency}
                />
              </div>
            )}
          </div>
        </div>

        <div className="order-bottom">
          <div className="order-total">
            <p className="order-total-label">Total:</p>
            <div className="right">
              <FormatNumber
                number={totalPrice || 0}
                className="order-total-number"
              />
              <p className="order-total-currency">
                {market.counter.Denom.Currency}
              </p>
            </div>
          </div>

          {!wallet ? (
            <div className="connect-wallet">
              <Button
                variant={ButtonVariant.PRIMARY}
                onClick={() => {
                  setLoginModal(true);
                }}
                image="/trade/images/wallet.svg"
                width={"100%"}
                height={37}
                label="Connect Wallet"
              />
            </div>
          ) : (
            <>
              <Button
                variant={ButtonVariant.PRIMARY}
                onClick={() => {
                  handleSubmit();
                }}
                image="/trade/images/arrow-right.svg"
                width={"100%"}
                height={37}
                label="Confirm Order"
                disabled={
                  !volume ||
                  volume === "0" ||
                  (tradeType === TradeType.LIMIT && !limitPrice) ||
                  (orderType === OrderType.BUY &&
                    totalPrice > Number(marketBalances.counter)) ||
                  (orderType === OrderType.SELL &&
                    totalPrice > Number(marketBalances.base))
                }
              />
            </>
          )}
        </div>
      </div>

      <div className="available-balances">
        <p className="title">Assets</p>
        <div className="balance-row">
          <p className="balance-label">{market.base.Denom.Currency} Balance</p>

          <FormatNumber number={marketBalances.base} />
        </div>

        <div className="balance-row">
          <p className="balance-label">
            {market.counter.Denom.Currency} Balance
          </p>

          <FormatNumber number={marketBalances.counter} />
        </div>
      </div>
    </div>
  );
};

export default OrderActions;
