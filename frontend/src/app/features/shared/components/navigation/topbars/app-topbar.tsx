import * as React from "react"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
    DropdownMenu,
    DropdownMenuTrigger,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu"
import { ModeToggle } from "@/components/mode-toggle"
import { SettingsToggle } from "../../buttons/settings-toggle"

interface AppTopbarProps extends React.HTMLAttributes<HTMLDivElement> {
    children?: React.ReactNode;

}

export function AppTopbar({ children, className, ...props }: AppTopbarProps) {
    return (
        <header
            className={`flex w-full items-center px-4 py-2 ${className}`}
            {...props}
        >
            <div className="flex items-center gap-2">
                {children}
            </div>

            <div className="ml-auto flex items-center gap-2">
                <ModeToggle />
                {/*<SettingsToggle />*/}
            </div>
        </header>
    )
}

