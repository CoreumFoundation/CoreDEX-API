import BigNumber from "bignumber.js";
import classNames from "classnames";
import { CSSProperties, ReactNode } from "react";
import "./format-number.scss";

export interface FormatNumberProps {
  /** Sets the number to be formatted. */
  number: number | string;
  className?: string;
  fontSize?: number | string;
  insideElement?: boolean;
  smallDecimals?: boolean;
  customStyle?: CSSProperties;
  precision?: number;
  prefix?: string | ReactNode;
  suffix?: string | ReactNode;
}

export function FormatNumber(props: FormatNumberProps) {
  const {
    number,
    className = "",
    fontSize,
    insideElement = true,
    smallDecimals = true,
    customStyle,
    precision,
    prefix,
    suffix,
  } = props;

  const formatDecimals = (decimals: string, prec: number) => {
    if (/^0+$/.test(decimals)) {
      return "0".repeat(prec);
    }

    const match = decimals.match(/^0{4,}/);

    if (match) {
      const [zeros] = match;
      const zeroCount = zeros.length;
      const remainingDigits = decimals.slice(zeroCount);

      const zeroStr = zeroCount.toString().padStart(2, "0");
      const significantSpace = prec - 3;
      const significant = remainingDigits.slice(0, significantSpace);
      const trailingZeros = "0".repeat(
        Math.max(0, significantSpace - significant.length)
      );

      return (
        <>
          {zeroStr && <span className="decimal">0</span>}
          <span className="subscript">{zeroStr}</span>
          {significant}
          {trailingZeros}
        </>
      );
    }

    return decimals.padEnd(prec, "0").slice(0, prec);
  };

  const num = typeof number === "string" ? number : new BigNumber(number);
  const [ints, decimals] = (
    typeof num === "string" ? num : num.toFormat()
  ).split(".");

  let bigInts: string | BigNumber = ints;
  if (!ints.includes(",")) {
    bigInts = ints === "-0" ? ints : new BigNumber(ints).toFormat();
  }

  const finalDecimals = decimals || "";
  let formattedDecimals = null;

  if (precision) {
    const processed = formatDecimals(finalDecimals, precision);
    formattedDecimals = (
      <span className={`${smallDecimals ? "decimal" : ""}`}>.{processed}</span>
    );
  } else if (finalDecimals) {
    formattedDecimals = (
      <span className={`${smallDecimals ? "decimal" : ""}`}>
        .{finalDecimals}
      </span>
    );
  }

  const content = (
    <>
      {prefix && prefix}
      {bigInts === "NaN" ? "--" : bigInts}
      {formattedDecimals}
      {suffix && suffix}
    </>
  );

  return insideElement ? (
    <p
      className={classNames("format__number", className)}
      style={{ fontSize, ...customStyle }}
    >
      {content}
    </p>
  ) : (
    content
  );
}
