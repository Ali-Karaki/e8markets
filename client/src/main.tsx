import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider } from "@tanstack/react-router";
import ReactDOM from "react-dom/client";

import { ApiError } from "./lib/api";
import { getRouter } from "./router";

const ONE_DAY = 1 * 24 * 60 * 60 * 1000;

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: (failureCount, error) => {
        if (error instanceof ApiError && error.status >= 400) return false;
        return failureCount < 1;
      },
      refetchOnWindowFocus: false,
      staleTime: 7 * ONE_DAY,
      gcTime: 7 * ONE_DAY,
    },
  },
});
const router = getRouter(queryClient);

const rootElement = document.getElementById("app");

if (rootElement && !rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  );
}
