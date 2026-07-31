import { Button } from "@/components/ui/button"
import {Flame, FolderClosed, GitBranch, LogOut, Settings2, Snowflake, SquareTerminal, Users} from "lucide-react"
import React from "react"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {Link} from "@tanstack/react-router";
import {useAuthStore} from "@/security/auth/authstore/auth-store.ts";
import {logout} from "@/security/auth/keycloak/keycloak.ts";

interface SecondarySidebarProps {
    children?: React.ReactNode
}

export function SecondarySidebar({ children }: SecondarySidebarProps) {
    const handleLogout = async() => {
        useAuthStore.getState().setToken(null);
        useAuthStore.getState().setUserId(null);
        await logout();
    }

    return (
        <aside className="fixed left-0 top-0 z-30 h-screen w-16 border-r bg-sidebar flex flex-col">
            <div className="h-12" />

            <div className="flex-1 grid grid-rows-12 place-items-center">
                <Link to={"/files"}>
                    <Button variant="ghost" className="cursor-pointer">
                        <FolderClosed />
                    </Button>
                </Link>
                <Button variant="ghost" className="cursor-pointer">
                    <Link to={"/archive"}>
                        <Users className={"h-8 w-8"}/>
                    </Link>
                </Button>
                <Button variant="ghost" className="cursor-pointer">
                    <Link to={"/archive"}>
                        <Snowflake className={"h-8 w-8"}/>
                    </Link>
                </Button>
                <Button variant="ghost" className="cursor-pointer" onClick={handleLogout}>
                        <LogOut className={"h-8 w-8"}/>
                </Button>
                {children}
            </div>
            <div className="h-12 flex items-center justify-center border-t">
                <Link to={"/profile"}>
                    <Avatar>
                        <AvatarImage />
                        <AvatarFallback>OK</AvatarFallback>
                    </Avatar>
                </Link>
            </div>
        </aside>
    )
}

