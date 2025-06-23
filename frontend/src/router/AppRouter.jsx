import AuthorizationPage from "@/pages/AuthorizationPage";
import MainPage from "@/pages/MainPage";
import { Route, Routes } from "react-router-dom";
import {routes} from "@/router/"
import HomePage from "@/pages/HomePage";

function AppRouter() {
  return (
    <Routes>
      <Route path={routes.user.home} element={<HomePage />} />
      <Route path={routes.public.main} element={<MainPage />} />
      <Route path={routes.public.auth.route} element={<AuthorizationPage />} />
    </Routes>
  );
}

export default AppRouter;