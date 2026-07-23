import './App.css'
import { ThemeProvider } from "@/components/theme-provider"
import { RouterProvider } from "@tanstack/react-router";
import { router } from "@/routes/app-routes";

interface AppProps {
    isAuthenticated: boolean;
}

function App({ isAuthenticated }: AppProps) {
    return (
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
            <RouterProvider router={router} />
        </ThemeProvider>
    );
}

export default App;
