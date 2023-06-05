import { Route, createHashRouter, createRoutesFromElements } from "react-router-dom";
import { RootPage } from "./routes/root";
import { HistoryPage } from "./routes/history";
import { PageWrapper } from '@/components/Wrapper';

export const router = createHashRouter(
  createRoutesFromElements(
    <Route element={<PageWrapper/>}>
      <Route path='/' element={<RootPage />}/>
      <Route path='/history' element={<HistoryPage />}/>
    </Route>
  )
);