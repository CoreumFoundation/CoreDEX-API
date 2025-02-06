import { useEffect, useState } from "react";
import { useStore } from "@/state/store";
import {
  OrderType,
  TradeType,
  OrderbookAction,
  WalletAsset,
} from "@/types/market";
import { getAvgPriceFromOBbyVolume, multiply, noExponents } from "@/utils";
import { FormatNumber } from "../FormatNumber";
import { Input, InputType } from "../Input";
import Button, { ButtonVariant } from "../Button";
import BigNumber from "bignumber.js";
import { submitOrder, getWalletAssets } from "@/services/general";
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
  const [baseBalance, setBaseBalance] = useState<number | string>(0);
  const [counterBalance, setCounterBalance] = useState<number | string>(0);

  useEffect(() => {
    const fetchWalletAssets = async () => {
      if (!wallet?.address) return;
      try {
        const response = await getWalletAssets(wallet?.address);
        if (response.status === 200 && response.data.length > 0) {
          const data = response.data;
          console.log("ASSET BALANCES", data);
          setBalances(data);
        }
      } catch (e) {
        console.log("ERROR GETTING WALLET ASSETS DATA >>", e);
      }
    };
    fetchWalletAssets();
  }, [wallet?.address]);

  useEffect(() => {
    if (balances && balances.length > 0) {
      const baseBalance: WalletAsset = balances.find(
        (asset: WalletAsset) => asset.Denom === market.base.Denom.Denom
      );
      const counterBalance: WalletAsset = balances.find(
        (asset: WalletAsset) => asset.Denom === market.counter.Denom.Denom
      );
      if (baseBalance) {
        setBaseBalance(baseBalance.SymbolAmount);
      }
      if (counterBalance) {
        setCounterBalance(counterBalance.SymbolAmount);
      }
    }
  }, [market.pair_symbol, balances]);

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

  // format price for regex according to coreum backend
  // 1.5 -> 15e-1 or 1e+1 -> 10
  function formatPriceForRegex(value: BigNumber): string {
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
  }

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

      try {
        const txHash = submitResponse.data.TXHash;
        await navigator.clipboard.writeText(txHash);

        pushNotification({
          type: "success",
          message: `Order Placed! TXHash copied to clipboard: ${txHash.slice(
            0,
            6
          )}...${txHash.slice(-4)}`,
        });
      } catch (copyError) {
        console.error("Copy failed:", copyError);
      }
    } catch (e: any) {
      console.log("ERROR HANDLING SUBMIT ORDER >>", e.error.message);
      throw e;
    }
  };

  return (
    <div className="order-actions-container">
      <div className="order-actions-content" style={{ padding: "16px" }}>
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
              />
            </div>
          )}
        </div>

        <div className="order-bottom">
          <div className="order-total">
            <p className="order-total-label">Total:</p>
            <div className="right">
              <FormatNumber
                number={totalPrice || 0}
                className="order-total-number"
                precision={7}
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
                  (orderType === OrderType.SELL && totalPrice > Number(baseBalance))
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
          <p className="balance-value">{baseBalance}</p>
        </div>

        <div className="balance-row">
          <p className="balance-label">
            {market.counter.Denom.Currency} Balance
          </p>
          <p className="balance-value">{counterBalance}</p>
        </div>
      </div>
    </div>
  );
};

export default OrderActions;
