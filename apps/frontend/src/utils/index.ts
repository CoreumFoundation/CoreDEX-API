import { OrderbookRecord } from "@/types/market";
import BigNumber from "bignumber.js";
import { CoreumNetwork } from "coreum-js-nightly";

export const toFixedDown = (
  float: number | BigNumber,
  decimals: number
): string => {
  let zeros: string[] = [];

  for (let i = 0; i < decimals; i++) {
    zeros.push("0");
  }

  let zerosString = zeros.join("");
  let factor = new BigNumber(`1${zerosString}`);

  let newAmount: BigNumber;
  let numeralFloat = BigNumber.isBigNumber(float)
    ? float
    : new BigNumber(float);
  if (factor.gt(0)) {
    newAmount = numeralFloat.decimalPlaces(decimals, BigNumber.ROUND_DOWN);
  } else {
    newAmount = new BigNumber(0);
  }

  return newAmount.toFormat();
};

export const multiply = (a: number | string, b: number | string): BigNumber => {
  let x = new BigNumber(a);
  let y = new BigNumber(b);
  return x.multipliedBy(y);
};

export const divide = (a: number | string, b: number | string): BigNumber => {
  let x = new BigNumber(a);
  let y = new BigNumber(b);
  return x.dividedBy(y);
};

export const noExponents = (number: number) => {
  let expRegex = /[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)/g;
  if (!expRegex.test(String(number))) return resolveAndFixPrecision(number);

  var data = String(number).split(/[eE]/);

  var z = "",
    sign = number < 0 ? "-" : "",
    str = data[0].replace(".", ""),
    mag = Number(data[1]) + 1;

  if (mag < 0) {
    z = sign + "0.";
    while (mag++) z += "0";
    return z + str.replace(/^\-/, "");
  }
  mag -= str.length;

  while (mag--) z += "0";

  return str + z;
};

export const resolveAndFixPrecision = (num: string | number): string => {
  let precision = 2;
  let amount = typeof num === "number" ? num : Number(num);
  if (amount > 10000) {
    precision = 2;
  } else if (amount > 100) {
    precision = 4;
  } else if (amount >= 1) {
    precision = 6;
  } else if (amount < 1) {
    precision = 8;
  } else if (amount < 0.00000001) {
    precision = 12;
  }
  const fixed = toFixedDown(amount, precision);
  if (fixed === "NaN" || Number.isNaN(fixed) || fixed === "0") {
    return num.toString();
  }
  return fixed;
};

// calculates the volume weighted average price from the orderbook
export const getAvgPriceFromOBbyVolume = (
  ob: OrderbookRecord[],
  targetVolume: string
): number => {
  if (!Array.isArray(ob) || ob.length === 0) return 0;

  const parsedTargetVolume = Math.max(
    Number.parseFloat(targetVolume) || 0,
    Number.EPSILON
  );

  let remainingVolume = parsedTargetVolume;
  let totalVolumeUsed = 0;
  let totalPriceVolumeProduct = 0;

  if (parsedTargetVolume <= Number.EPSILON) return 0;

  for (const order of ob) {
    if (remainingVolume <= 0) break;

    const orderPrice = Number.parseFloat(order.HumanReadablePrice) || 0;
    const orderVolume = Number.parseFloat(order.SymbolAmount) || 0;

    if (orderPrice <= 0 || orderVolume <= 0) continue;

    const fillAmount = Math.min(remainingVolume, orderVolume);

    totalVolumeUsed += fillAmount;
    totalPriceVolumeProduct += orderPrice * fillAmount;
    remainingVolume -= fillAmount;
  }

  return totalVolumeUsed > 0
    ? Number((totalPriceVolumeProduct / totalVolumeUsed).toPrecision(8))
    : 0;
};

export const resolveCoreumExplorer = (network: CoreumNetwork) => {
  switch (network) {
    case CoreumNetwork.TESTNET:
      return "https://explorer.testnet-1.coreum.dev/coreum";
    case CoreumNetwork.DEVNET:
      return "https://explorer.devnet-1.coreum.dev/coreum";
    default:
      return "https://explorer.coreum.com/coreum";
  }
};
