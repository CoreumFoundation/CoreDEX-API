import { WalletType } from "graz";

export interface ToasterProps {
  type: "error" | "success" | "warning" | "info";
  message: string;
}

export const SUPPORTED_WALLETS = [
  WalletType.LEAP,
  WalletType.COSMOSTATION,
  WalletType.METAMASK_SNAP_LEAP,
  WalletType.KEPLR,
];
