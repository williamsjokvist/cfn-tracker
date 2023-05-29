import { createHashRouter } from "react-router-dom";
import { RootPage } from "./routes/root";
import { HistoryPage } from "./routes/history";
import { PageWrapper } from '@/components/Wrapper';

export const router = createHashRouter([
  {
    path: "/",
    element: <PageWrapper><RootPage /></PageWrapper>,
  },
  {
    path: "/history",
    element: <PageWrapper><HistoryPage /></PageWrapper>,
  },
]);