import { RouterProvider } from "react-router-dom";
import router from "./router";

import { Toaster } from "./components/Toaster";

function App() {
  return (
    <>
      <Toaster />
      <RouterProvider router={router} />
    </>
  );
}

export default App;
