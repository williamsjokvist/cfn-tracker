import React from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from "react-router-dom";

import { CFNMachineContext } from './machine';
import { LanguageProvider } from './context/LanguageProvider';
import { router } from "./router";

import './styles/globals.sass'
import './styles/sidebar.sass'

import './i18n';

const app = (
  <React.StrictMode>
    <CFNMachineContext.Provider>
      <LanguageProvider>
        <RouterProvider router={router} />
      </LanguageProvider>
    </CFNMachineContext.Provider>
  </React.StrictMode>
)

const root = document.getElementById("root")!

createRoot(root).render(app);