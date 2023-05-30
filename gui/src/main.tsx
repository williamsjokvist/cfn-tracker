import React from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from "react-router-dom";
import { CFNMachineContext } from './state-machine/machine';
import { router } from "./router";

import './styles/globals.sass'
import './i18n';

const app = (
  <React.StrictMode>
    <CFNMachineContext.Provider>
      <RouterProvider router={router} />
    </CFNMachineContext.Provider>
  </React.StrictMode>
)

const root = document.getElementById("root")!

createRoot(root).render(app);