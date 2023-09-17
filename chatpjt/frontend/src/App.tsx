import { createBrowserRouter, createRoutesFromElements, Route, RouterProvider } from "react-router-dom"
import Home from "./pages/Home"
import React from "react";
import { ThemeProvider } from "@mui/material";
import createMuiTheme from "./theme/theme";

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route>
      <Route path="/" element={<Home />} />
    </Route>
  )
);


const App: React.FC = () => {
  const theme = createMuiTheme();

  return (
  <ThemeProvider theme={theme}>
    <RouterProvider router={router} />
  </ThemeProvider>
  );
};

export default App;
