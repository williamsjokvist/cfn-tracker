import React from 'react'
import ReactDOM from 'react-dom/client'
import {
  createHashRouter,
  RouterProvider,
} from "react-router-dom";
import Root from "./routes/root";
import History from "./routes/history";
import './i18n';
import './index.css'
import Wrapper from './components/Wrapper';

const router = createHashRouter([
  {
    path: "/",
    element: <Root />,
  },
  {
    path: "/tracking",
    element: <Root />,
  },
  {
    path: "/history",
    element: <History />,
  },
]);



ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Wrapper>
      <RouterProvider router={router} />
    </Wrapper>
  </React.StrictMode>
);