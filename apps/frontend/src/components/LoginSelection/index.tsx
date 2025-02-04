import { useState } from "react";
import { useStore } from "@/state/store";
import { ExtensionWallets } from "coreum-js-nightly/dist/main/types";
import { quantum } from "ldrs";
import { Client } from "coreum-js-nightly";
quantum.register();
import "./login-selection.scss";

const LoginSelection = ({
  closeModal,
}: {
  closeModal?: (arg: boolean) => void;
}) => {
  const { setCoreum, network, setWallet, pushNotification } = useStore();
  const [isLoading, setIsLoading] = useState(false);

  const resolveOption = (option: ExtensionWallets) => {
    switch (option) {
      case "leap":
        return "Leap";
      case "keplr":
        return "Keplr";
      case "cosmostation":
        return "Cosmostation";
      default:
        return option;
    }
  };

  const connectWallet = async (option: ExtensionWallets) => {
    try {
      setIsLoading(true);
      const client: Client = new Client({ network: network });
      const connectOptions = {
        withWS: false,
      };
      console.log("Attempting to connect with:", option);
      await client.connectWithExtension(option, connectOptions);
      console.log("Connected with:", option);
      setWallet({ address: client!.address });
      setCoreum(client);
      closeModal && closeModal(false);
      pushNotification({
        message: "Connected",
        type: "success",
      });
    } catch (e: any) {
      console.error("Connection failed:", e);
      setIsLoading(false);
    }
  };

  return isLoading ? (
    <div className="loader">
      <p className="connecting">Connecting</p>
      <l-quantum size="40" speed="6" color="white"></l-quantum>
    </div>
  ) : (
    <div className="wallet-options">
      {
        // SUPPORTED_WALLETS
        // Keplr is the only supported wallet for devnet
        Object.values(ExtensionWallets).map((option, idx) => {
          return (
            <div
              className="wallet-option"
              style={{
                display: "flex",
                gap: "10px",
                alignItems: "center",
                cursor: "pointer",
              }}
              onClick={() => {
                if (
                  option === ExtensionWallets.KEPLR ||
                  option === ExtensionWallets.LEAP ||
                  option === ExtensionWallets.COSMOSTATION
                ) {
                  connectWallet(option);
                }
              }}
              key={idx}
            >
              <img
                src={`/trade/images/${option}.svg`}
                alt={option}
                width={24}
                height={24}
              />
              <p>{resolveOption(option)}</p>
            </div>
          );
        })
      }
    </div>
  );
};

export default LoginSelection;
