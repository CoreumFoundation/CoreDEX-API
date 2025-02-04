import { CoreumNetwork } from "coreum-js-nightly";

export interface Features {
  Flag: string;
  Status: boolean;
  Description?: string;
}

interface Environment {
  KycAPI: string;
  AccountAPI: string;
}
export interface NetworkEnvs {
  Auth: string;
  KYC: string;
}
interface Networks {
  Network: CoreumNetwork;
  Envs: NetworkEnvs;
}
interface Envs extends Environment {
  Networks: Networks[];
}

const ENV_NAMES: Array<string> = [""];

export const getAppEnv = () => {
  const env: { [key: string]: string } = {};
  const prefix = "VITE_";
  const config =
    (import.meta as any).env.VITE_MODE === "development"
      ? (import.meta as any).env
      : (window as any).ATG.env;

  ENV_NAMES.forEach((e) => {
    env[e] = config[`${prefix}${e}`];
    if (!env[e])
      throw new Error(`FATAL: Env ${e} is not set. Please add env ${e}`);
  });

  return env;
};

export const getNetworkEnvs = (
  env: Envs,
  network: CoreumNetwork
): NetworkEnvs => {
  const networkEnv = env.Networks.find((netEnv) => netEnv.Network === network);
  return networkEnv!.Envs;
};

export const custom_nodes = {
  mainnet: "https://full-node-californium.mainnet-1.coreum.dev:26657",
};

const getConfig = () => {
  try {
    const app_envs = getAppEnv();
    console.log("app envs", { app_envs });

    return app_envs && app_envs;
  } catch (error) {
    console.log("PARSING ENVS ERROR", error);
  }
};

export const config = getConfig();
