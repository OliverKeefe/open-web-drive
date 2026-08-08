import { EllipsisVertical, Info, Settings, Trash2, CloudDownload } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuthStore } from "@/security/auth/authstore/auth-store";
import { RestHandler } from "@/app/features/shared/api/rest/rest-handler";
import { ConfirmAlertDialog } from "@/app/features/shared/components/alerts/confirm-alert";
import {enqueueSnackbar} from "notistack";

interface FileDropdownProps {
    fileId: string;
    onDeleted: () => void;
    onInfo: () => void;
}

const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;
if (!backendBaseUrl) throw new Error('VITE_BACKEND_BASE_URL is not set. Set it in your .env file or pass it at build time.');

function FileDropdown({ fileId, onDeleted, onInfo }: FileDropdownProps) {
    const userId = useAuthStore((s) => s.userId);

    async function handleDelete(e: React.MouseEvent) {
        e.stopPropagation();

        if (!userId) {
            console.error("No user ID found");
            return;
        }

        try {
            const api = new RestHandler(backendBaseUrl);

            await api.handleDelete<
                { id: string; },
                void
            >("api/files/delete", {
                id: fileId,
            });

            onDeleted();
            const message = "File: " + fileId + " deleted successfully.";
            enqueueSnackbar(message, { autoHideDuration: 1000 })


        } catch (error) {
            console.error("Failed to delete file on server:", error);
            const message = "File: " + fileId + " could not be deleted.";
            enqueueSnackbar(message, { autoHideDuration: 5000 })
        }
    }

    async function handleDownload(e: React.MouseEvent) {
        e.stopPropagation();

        try {
            const api = new RestHandler(backendBaseUrl);
            await api.handleDownload<{ id: string }>("api/files/download", {
                id: fileId,
            });
            enqueueSnackbar("Download started", { autoHideDuration: 1000 });
        } catch (error) {
            console.error("Failed to download file:", error);
            enqueueSnackbar("Download failed", { autoHideDuration: 5000 });
        }
    }

    function handleInfo(e: React.MouseEvent) {
        e.stopPropagation();
        onInfo();
    }

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button
                    variant="ghost"
                    size="icon"
                    className="ms-1 cursor-pointer"
                >
                    <EllipsisVertical className="h-4 w-4" />
                </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent align="end">
                <ConfirmAlertDialog
                    onConfirm={handleDelete}
                    message={"This file will be removed, permanently."}
                    icon={Trash2}
                >
                    <DropdownMenuItem
                        className="cursor-pointer text-destructive focus:text-destructive"
                        onSelect={(e) => e.preventDefault()}
                    >
                        <Trash2 className="mr-2 h-4 w-4" />
                        <span>Delete</span>
                    </DropdownMenuItem>
                </ConfirmAlertDialog>
                <DropdownMenuSeparator />

                <DropdownMenuItem className="cursor-pointer" onClick={handleInfo}>
                    <Info className="mr-2 h-4 w-4" />
                    <span>File Info</span>
                </DropdownMenuItem>

                <DropdownMenuItem className="cursor-pointer">
                    <Settings className="mr-2 h-4 w-4" />
                    <span>File Settings</span>
                </DropdownMenuItem>

                <DropdownMenuItem className="cursor-pointer" onClick={handleDownload}>
                    <CloudDownload className="mr-2 h-4 w-4"/>
                    <span>Download</span>
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    );
}

export default FileDropdown
