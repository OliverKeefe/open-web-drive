import {
    Tooltip,
    TooltipTrigger,
    TooltipContent
} from "@/components/ui/tooltip.tsx";
import { Button } from "@/components/ui/button";
import type {ReactNode} from "react";

interface IconButtonProps {
    label: string
    icon: ReactNode
    onClick?: () => void
    className?: string
}

export function IconButton({ label, icon, onClick, className }: IconButtonProps) {
    return (
        <Tooltip>
            <TooltipTrigger asChild>
                <Button variant="outline" size="icon" className={className} onClick={onClick}>
                    {icon}
                    <span className="sr-only">{label}</span>
                </Button>
            </TooltipTrigger>
            <TooltipContent>
                <p>{label}</p>
            </TooltipContent>
        </Tooltip>
    )
}