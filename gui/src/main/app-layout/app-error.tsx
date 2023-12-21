import React from "react";
import { useRouteError } from "react-router-dom";

import { AppTitleBar } from "./app-titlebar";

export const AppErrorBoundary = () => {
  const error = useRouteError();
  const [errMsg, setErrMsg] = React.useState("")
  React.useEffect(() => {
    console.error(error);
    if (error instanceof Error) {
      setErrMsg(error.message)
    }
  }, [])
  return (
    <>
      <AppTitleBar />
      <main className="grid justify-center items-center h-screen w-full text-white">
        <div className="bg-black bg-opacity-25 px-6 py-3 rounded-md">
          <h1 className="text-2xl text-center font-bold">An application error occured</h1>
          <p className="text-xl">{errMsg}</p>
        </div>
      </main>
    </>
  )
}