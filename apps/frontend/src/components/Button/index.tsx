import React from "react";
import "./button.scss";
import { useStore } from "@/state/store";
import { mirage } from "ldrs";
mirage.register();

interface ButtonProps {
  label: string;
  image?: string;
  onClick?: () => void;
  className?: string;
  variant?: ButtonVariantType;
  height?: string | number;
  width?: string | number;
  disabled?: boolean;
}

export const ButtonVariant = {
  PRIMARY: "primary",
  SECONDARY: "secondary",
  TERTIARY: "tertiary",
  DANGER: "danger",
};

export type ButtonVariantType =
  (typeof ButtonVariant)[keyof typeof ButtonVariant];

const Button: React.FC<ButtonProps> = ({
  label,
  image,
  onClick,
  className,
  variant = ButtonVariant.PRIMARY,
  height,
  width,
  disabled,
}) => {
  const { isLoading } = useStore();

  return (
    <button
      className={`button button-${variant}  ${className || ""}`}
      onClick={onClick}
      style={{ height, width }}
      disabled={disabled}
    >
      {isLoading ? (
        <l-mirage size="40" speed="6" color="#25d695"></l-mirage>
      ) : (
        <>
          {image && <img src={image} alt={label} className="button-image" />}
          <span className="button-label">{label}</span>
        </>
      )}
    </button>
  );
};

export default Button;
