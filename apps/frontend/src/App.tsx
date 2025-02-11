import { RouterProvider } from "react-router-dom";
import router from "./router";

import { Toaster } from "./components/Toaster";
import { useEffect } from "react";
import { wsManager } from "./services/websocket-refactor";

const WEBSOCKET_URL =
  import.meta.env.VITE_ENV_MODE === "development"
    ? import.meta.env.VITE_ENV_WS
    : (window as any).COREUM.env.VITE_ENV_WS;

function App() {
  useEffect(() => {
    wsManager.connect(WEBSOCKET_URL);
  }, []);

  return (
    <>
      <Toaster />
      <RouterProvider router={router} />
    </>
  );
}

export default App;
