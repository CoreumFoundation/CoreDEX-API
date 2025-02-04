import { Suspense } from "react";
import { createBrowserRouter } from "react-router-dom";
import Home from "@/containers";

const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <Suspense>
        <Home />
      </Suspense>
    ),
  },
]);

export default router;
