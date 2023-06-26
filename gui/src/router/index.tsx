import { Route, createHashRouter, createRoutesFromElements } from "react-router-dom";
import { PageWrapper } from '@/components/Wrapper';
import { RootPage } from "./routes/root";
import { OutputPage } from "./routes/output";
import { HistoryPage } from "./routes/history";

export const router = createHashRouter(
  createRoutesFromElements(
    <Route element={<PageWrapper/>}>
      <Route path='/' element={<RootPage />}/>
      <Route path='/history' element={<HistoryPage />}/>
      <Route path='/output' element={<OutputPage />}/>
    </Route>
  )
);