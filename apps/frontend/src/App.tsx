import { RouterProvider } from "react-router-dom";
import router from "./router";

import { Toaster } from "./components/Toaster";
import { useEffect } from "react";

function App() {
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

  return (
    <>
      <Toaster />
      <RouterProvider router={router} />
    </>
  );
}

export default App;
