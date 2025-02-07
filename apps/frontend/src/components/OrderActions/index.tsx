import { useEffect, useState } from "react";
import { useStore } from "@/state/store";
import {
  OrderType,
  TradeType,
  OrderbookAction,
  WalletAsset,
  TIME_IN_FORCE_STRING,
  TIME_SELECTION,
} from "@/types/market";
import { getAvgPriceFromOBbyVolume, multiply, noExponents } from "@/utils";
import { FormatNumber } from "../FormatNumber";
import { Input, InputType } from "../Input";
import Button, { ButtonVariant } from "../Button";
import BigNumber from "bignumber.js";
import { submitOrder, getWalletAssets } from "@/services/api";
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

BigNumber.config({ DECIMAL_PLACES: 30, EXPONENTIAL_AT: 0 });

const OrderActions = ({
  orderbookAction,
}: {
  orderbookAction?: OrderbookAction;
}) => {
  const { orderbook, wallet, setLoginModal, pushNotification, market, coreum } =
    useStore();

  const [orderType, setOrderType] = useState<OrderType>(OrderType.BUY);
  const [totalPrice, setTotalPrice] = useState(0);
  const [limitPrice, setLimitPrice] = useState("");
  const [volume, setVolume] = useState<string>("");
  const [tradeType, setTradeType] = useState(TradeType.MARKET);
  const [balances, setBalances] = useState<any>(null);
  const [baseBalance, setBaseBalance] = useState<string>("0");
  const [counterBalance, setCounterBalance] = useState<string>("0");
  const [advSettingsOpen, setAdvSetting] = useState<boolean>(false);
  const [timeInForce, setTimeInForce] = useState<TIME_IN_FORCE_STRING>(
    TIME_IN_FORCE_STRING.goodTilCancel
  );
  const [timeToCancel, setTimeToCancel] = useState<TIME_SELECTION>(
    TIME_SELECTION["5M"]
  );

  useEffect(() => {
    fetchWalletAssets();
  }, [wallet, market]);

  useEffect(() => {
    if (!balances) return;

    const baseBalanceObject = balances.find(
      (asset: WalletAsset) => asset.Denom === market.base.Denom.Denom
    );
    const counterBalanceObject = balances.find(
      (asset: WalletAsset) => asset.Denom === market.counter.Denom.Denom
    );

    setBaseBalance(baseBalanceObject ? baseBalanceObject.SymbolAmount : "0");
    setCounterBalance(
      counterBalanceObject ? counterBalanceObject.SymbolAmount : "0"
    );
  }, [market, balances]);

  // trigger when click on orderbook
  useEffect(() => {
    if (orderbookAction?.price) {
      setTradeType(TradeType.LIMIT);
      setOrderType(orderbookAction.type);

      const volumeBN = new BigNumber(orderbookAction.volume);
      const priceBN = new BigNumber(orderbookAction.price);

      setVolume(volumeBN.toFixed(18));
      setLimitPrice(priceBN.toFixed(18));
      setTotalPrice(priceBN.times(volumeBN).toNumber());
    }
  }, [orderbookAction]);

  useEffect(() => {
    if (tradeType === TradeType.LIMIT) {
      const vol = volume ? Number(volume) : 0;

      const total = multiply(Number(limitPrice), vol);
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

  const fetchWalletAssets = async () => {
    if (!wallet?.address) return;
    try {
      const response = await getWalletAssets(wallet?.address);
      if (response.status === 200 && response.data.length > 0) {
        const data = response.data;
        setBalances(data);
      }
    } catch (e) {
      console.log("ERROR GETTING WALLET ASSETS DATA >>", e);
    }
  };

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
        ...(tradeType === TradeType.LIMIT
          ? {
              price: formatPriceForRegex(
                BigNumber(limitPrice).dividedBy(
                  BigNumber(10).pow(
                    market.counter.Denom.Denom === "udevcore"
                      ? 6
                      : market.counter.Denom.Precision
                  )
                )
              ),
            }
          : { price: "" }),
        quantity: BigNumber(volume)
          .multipliedBy(
            BigNumber(10).pow(
              market.counter.Denom.Denom === "udevcore"
                ? 6
                : market.counter.Denom.Precision
            )
          )
          .toFixed(0),
        side: orderType === OrderType.BUY ? Side.SIDE_BUY : Side.SIDE_SELL,
        goodTil: undefined,
        timeInForce:
          tradeType === TradeType.LIMIT
            ? TimeInForce.TIME_IN_FORCE_GTC
            : TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
      };

      const orderMessage = DEX.PlaceOrder(orderCreate);
      const signedTx = await coreum?.signTx([orderMessage]);
      const encodedTx = TxRaw.encode(signedTx!).finish();
      const base64Tx = fromByteArray(encodedTx);
      const submitResponse = await submitOrder({ TX: base64Tx });

      if (submitResponse.status !== 200) {
        pushNotification({
          type: "error",
          message: "There was an issue submitting your order",
        });
        throw new Error("Error submitting order");
      }
      const txHash = submitResponse.data.TXHash;
      pushNotification({
        type: "success",
        message: `Order Placed! TXHash: ${txHash.slice(0, 6)}...${txHash.slice(
          -4
        )}`,
      });
    } catch (e: any) {
      console.log("ERROR HANDLING SUBMIT ORDER >>", e.error.message);
      pushNotification({
        type: "error",
        message: e.error.message,
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
                  decimals={13}
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
                  decimals={13}
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
                          Object.keys(TIME_IN_FORCE_STRING) as Array<
                            keyof typeof TIME_IN_FORCE_STRING
                          >
                        ).map((key) => [TIME_IN_FORCE_STRING[key]])}
                        value={timeInForce}
                        renderItem={(item) => (
                          <div
                            className="network-item"
                            onClick={() => {
                              setTimeInForce(item[0] as TIME_IN_FORCE_STRING);
                            }}
                          >
                            {item}
                          </div>
                        )}
                      />

                      {timeInForce === TIME_IN_FORCE_STRING.goodTilTime && (
                        <Dropdown
                          variant={DropdownVariant.OUTLINED}
                          items={(
                            Object.keys(TIME_SELECTION) as Array<
                              keyof typeof TIME_SELECTION
                            >
                          ).map((key) => [TIME_SELECTION[key]])}
                          value={timeToCancel}
                          renderItem={(item) => (
                            <div
                              className="network-item"
                              onClick={() => {
                                setTimeToCancel(item[0] as TIME_SELECTION);
                              }}
                            >
                              {item}
                            </div>
                          )}
                        />
                      )}

                      {timeInForce === TIME_IN_FORCE_STRING.goodTilTime && timeToCancel === TIME_SELECTION.CUSTOM (
                         <DatetimePicker
                         selectedDate={customTime}
                         onChange={(val: any) => {
                           setCustomTime(val?.toString() || "");
                         }}
                         width={"100%"}
                         minDate={new Date(dayjs.utc().add(1, "day").format())}
                       />
                     
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
            <p className="order-total-label">Total Cost:</p>
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
                    totalPrice > Number(counterBalance)) ||
                  (orderType === OrderType.SELL &&
                    totalPrice > Number(baseBalance))
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
          <p className="balance-value">
            {Number(baseBalance).toLocaleString()}
          </p>
        </div>

        <div className="balance-row">
          <p className="balance-label">
            {market.counter.Denom.Currency} Balance
          </p>
          <p className="balance-value">
            {Number(counterBalance).toLocaleString()}
          </p>
        </div>
      </div>
    </div>
  );
};

export default OrderActions;
