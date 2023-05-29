import React from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from "react-router-dom";
import { router } from "./router";

import './globals.sass'
import './i18n';

const app = (
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
)

const root = document.getElementById("root")!

createRoot(root).render(app);