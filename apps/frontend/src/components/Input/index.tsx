import React, { CSSProperties, ReactNode } from "react";
import classNames from "classnames";
import { NumericFormat, NumericFormatProps } from "react-number-format";
import "./input.scss";

export enum InputType {
  TEXT = "text",
  NUMBER = "number",
  TEXTAREA = "textarea",
}

export enum ErrorInputType {
  BORDER = "error-border",
  UNDERLINE = "error-underline",
}

export enum InputSize {
  SMALL = "small",
  MEDIUM = "medium",
  LARGE = "large",
}

export interface InputAdornmentClass {
  left?: string;
  right?: string;
}

interface InputProps {
  inputName: string;
  onValueChange: (value: string) => void;
  value: string | number;
  placeholder: string;
  type?: InputType;
  label?: string;
  sublabel?: string;
  error?: ReactNode | string;
  isError?: boolean;
  customCss?: CSSProperties;
  adornmentClassname?: string | InputAdornmentClass;
  inputClassname?: string;
  inputWrapperClassname?: string;
  wrapperClassname?: string;
  autocomplete?: "on" | "off";
  disabled?: boolean;
  decimals?: number;
  errorType?: ErrorInputType;
  height?: string | number;
  width?: string | number;
  size?: InputSize;
  adornmentLeft?: ReactNode;
  adornmentRight?: ReactNode;
  labelTooltip?: ReactNode | string;
  tooltipWidth?: string;
  step?: number;
  enforceStep?: boolean;
  onBlur?: (e: any) => void;
}
type Input =
  | (InputProps & React.InputHTMLAttributes<HTMLInputElement>)
  | (InputProps & React.TextareaHTMLAttributes<HTMLTextAreaElement>)
  | (InputProps & NumericFormatProps);

export const Input = ({
  inputName,
  value,
  onValueChange,
  placeholder,
  type,
  label,
  sublabel,
  error,
  isError,
  customCss,
  adornmentClassname,
  inputClassname,
  wrapperClassname,
  inputWrapperClassname,
  autocomplete = "off",
  disabled = false,
  decimals = 6,
  width = "100%",
  height = "fit-content",
  size = InputSize.MEDIUM,
  errorType = ErrorInputType.BORDER,
  adornmentLeft,
  adornmentRight,
  labelTooltip,
  tooltipWidth = "200px",
  ...rest
}: Input) => {
  ``;
  const inputClasses = classNames(
    "input",
    {
      "input-error": isError && errorType === ErrorInputType.UNDERLINE,
    },
    inputClassname
  );

  const inputElement =
    type === InputType.TEXTAREA ? (
      <textarea
        autoComplete={autocomplete}
        className={inputClasses}
        onChange={(e) => onValueChange(e.target.value)}
        style={customCss}
        name={inputName}
        value={value}
        placeholder={placeholder}
        disabled={disabled}
        {...(rest as React.TextareaHTMLAttributes<HTMLTextAreaElement>)}
      />
    ) : type === InputType.TEXT ? (
      <input
        autoComplete={autocomplete}
        className={inputClasses}
        onChange={(e) => onValueChange(e.target.value)}
        style={customCss}
        type={type}
        name={inputName}
        value={value}
        placeholder={placeholder}
        disabled={disabled}
        {...(rest as React.InputHTMLAttributes<HTMLInputElement>)}
      />
    ) : (
      <NumericFormat
        autoComplete={autocomplete}
        className={inputClasses}
        style={customCss}
        allowLeadingZeros
        thousandSeparator
        allowNegative={false}
        decimalScale={decimals}
        name={inputName}
        placeholder={placeholder}
        value={value}
        disabled={disabled}
        onValueChange={({ value }) => {
          if (value === "") return onValueChange("0");
          const sanitizedValue = value.replace(/^0+(?=\d)/, "");
          onValueChange(sanitizedValue || "0");
        }}
        {...(rest as NumericFormatProps)}
      />
    );

  return (
    <div
      className={classNames("input__wrapper", wrapperClassname, {
        "input__wrapper--error": isError,
      })}
    >
      {label && (
        <div className="input__label">
          <div className="input__label__label">
            <p>{label}</p>
          </div>
          {sublabel && <p className="input__label__sublabel">{sublabel}</p>}
        </div>
      )}
      <div
        className={classNames(
          "input__container",
          `input__container--${size}`,
          inputWrapperClassname,
          {
            "input__container--disabled": disabled,
            "input__container--error":
              isError && errorType === ErrorInputType.BORDER,
          }
        )}
        style={{ width, height }}
      >
        {adornmentLeft && (
          <div
            className={classNames(
              "input__adornment input__adornment--left",
              {
                [`${adornmentClassname}`]:
                  typeof adornmentClassname === "string",
              },
              (adornmentClassname as InputAdornmentClass)?.left
            )}
          >
            {adornmentLeft}
          </div>
        )}
        {inputElement}

        {adornmentRight && (
          <div
            className={classNames(
              "input__adornment input__adornment--right",
              {
                [`${adornmentClassname}`]:
                  typeof adornmentClassname === "string",
              },
              (adornmentClassname as InputAdornmentClass)?.right
            )}
          >
            {adornmentRight}
          </div>
        )}
      </div>
      {isError && error ? (
        <>
          {typeof error === "string" ? (
            <div className="input__error">
              <span className="input__error__txt">{error}</span>
            </div>
          ) : (
            error
          )}
        </>
      ) : (
        <></>
      )}
    </div>
  );
};
