import React from 'react'
import {createRoot} from 'react-dom/client'
import {createBrowserRouter, RouterProvider,} from "react-router-dom";
import {ConfigProvider, theme} from "antd";
import 'antd/dist/reset.css';
import ErrorPage from "./error-page";
import Compress from "./routes/Compress";
import App from "./App";
import About from "./routes/About";
import Setting from "./routes/Setting";
import {QueryClient, QueryClientProvider} from "react-query";

const router = createBrowserRouter([
    {
        path: "/",
        element: <App/>,
        errorElement: <ErrorPage/>,
        children: [
            {
                index: true,
                element: <Compress/>,
            },
            {
                path: "/setting",
                element: <Setting/>,
            },
            {
                path: "/about",
                element: <About/>,
            },
        ],
    },
]);

const container = document.getElementById('root')

const root = createRoot(container!);
let queryClient = new QueryClient();

root.render(
    <React.StrictMode>
        <ConfigProvider
            theme={{
                algorithm: theme.defaultAlgorithm,
            }}
        >
            <QueryClientProvider client={queryClient}>
                <RouterProvider router={router}/>
            </QueryClientProvider>
        </ConfigProvider>
    </React.StrictMode>
)
