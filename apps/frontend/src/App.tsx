import { RouterProvider } from "react-router-dom";
import router from "./router";
import { Toaster } from "./components/Toaster";
import { useEffect } from "react";
import { wsManager } from "./services/websocket";
import { WS_URL } from "./config/envs";
import { useStore } from "./state/store";
import { Client } from "coreum-js-nightly";
import { ICoreumWallet } from "./types/market";

function App() {
  const { wallet, network, setWallet, setCoreum, pushNotification } =
    useStore();

  useEffect(() => {
    fetch("/build.version")
      .then((response) => {
        if (!response.ok) {
          throw new Error(
            `Failed to load build.version: ${response.statusText}`
          );
        }
        return response.text();
      })
      .then((versionText) => {
        console.log("Build version:", versionText);
      })
      .catch((error) => {
        console.error("Error fetching build version:", error);
      });
  }, []);

  useEffect(() => {
    wsManager.connect(WS_URL);
  }, []);

  useEffect(() => {
    const localWallet = getLocalWallet();
    if (localWallet && !wallet) {
      connectWalletOnLoad(localWallet as ICoreumWallet);
    }
  }, []);

  const getLocalWallet = () => {
    if (localStorage.wallet) {
      return JSON.parse(localStorage.wallet);
    }

    return null;
  };

  const connectWalletOnLoad = async (wallet: ICoreumWallet) => {
    try {
      const client = new Client({ network: network });
      await client.connectWithExtension(wallet.method, { withWS: false });
      setWallet({ address: wallet.address, method: wallet.method });
      setCoreum(client);
      pushNotification({
        message: "Connected",
        type: "success",
      });
    } catch (e: any) {
      console.error("Connection failed:", e);
    }
  };

  return (
    <>
      <Toaster />
      <RouterProvider router={router} />
    </>
  );
}

export default App;
