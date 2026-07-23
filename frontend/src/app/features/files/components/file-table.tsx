import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { useEffect, useMemo, useState } from "react"
import { Checkbox } from "@/components/ui/checkbox"
import { UploadDialog } from "@/app/features/shared/components/dialog/upload-dialog"
import { Button } from "@/components/ui/button"
import { Clock, FolderPlus, Star } from "lucide-react"
import {
    type CursorReq,
    getAllMetadata,
} from "@/app/features/files/handlers.ts"
import type { Metadata } from "@/app/features/files/types.ts"
import { useAuthStore } from "@/security/auth/authstore/auth-store"
import { getIconForFile } from "@react-symbols/icons/utils"
import { FileDialog } from "@/app/features/files/components/file-dialog.tsx";
import FileDropdown from "@/app/features/files/components/file-dropdown.tsx";

/**
 * Main file table component in files page.
 * */
export function FileTable() {
    const [dialogOpen, setDialogOpen] = useState(false)
    const [activeFile, setActiveFile] = useState<Metadata | null>(null)
    const userId = useAuthStore((s) => s.userId)
    const [files, setFiles] = useState<Metadata[]>([])
    const [selected, setSelected] = useState<string[]>([])
    const { request } = usePagination(userId, 20);

    /**
     * Main refresh file function.
     * */
    async function refreshFiles() {
        if (!userId || !request.user_id) return;
        const resp = await getAllMetadata(request);
        setFiles(resp.data);
    }

    useEffect(() => {
        refreshFiles();
    }, [request, userId]);

    /**
     * Handles opening FileDialog component on file.
     * */
    function openDialog(file: Metadata) {
        setActiveFile(file)
        setDialogOpen(true)
    }

    /**
     * Handles checkbox toggle in file table.
     * */
    const toggleSelect = (id: string) => {
        setSelected((prev) =>
            prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]
        )
    }

    /**
     * Handles the selection of all file checkboxes within the file table.
     * */
    const selectAll = (checked: boolean) => {
        setSelected(checked ? files.map((f) => f.id) : [])
    }

    /**
     * Upon file deletion, refreshes files from backend.
     * This allows for instant UI update post event.
     * */
    async function handleFileDeleted() {
        await refreshFiles();
    }

    return (
        <div>
            <h1 className="text-2xl font-semibold pb-4 pt-4 m-1">All files</h1>

            <nav className="w-full flex gap-3">
                <Button variant="outline"><Clock /> Recents</Button>

                <UploadDialog
                    onUploaded={async () => {
                        await refreshFiles();
                    }}
                />

                <Button variant="outline"><FolderPlus /> New Folder</Button>
                <Button variant="outline"><Star /> Favorites</Button>
            </nav>

            <Table className="mt-2 w-full table-fixed">
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-[30px]">
                            <Checkbox
                                checked={files.length > 0 && selected.length === files.length}
                                onCheckedChange={(v) => selectAll(v === true)}
                            />
                        </TableHead>
                        <TableHead className="w-[30px]" />
                        <TableHead>Name</TableHead>
                        <TableHead className="w-[150px]">Last Modified</TableHead>
                        <TableHead className="w-[50px]" />
                    </TableRow>
                </TableHeader>

                <TableBody>
                    {files.map((file) => (
                        <TableRow key={file.id} className="cursor-pointer">
                            <TableCell>
                                <Checkbox
                                    checked={selected.includes(file.id)}
                                    onCheckedChange={() => toggleSelect(file.id)}
                                />
                            </TableCell>

                            <TableCell onClick={() => openDialog(file)}>
                                <div className="w-4">
                                    {getIconForFile({ fileName: file.file_name })}
                                </div>
                            </TableCell>

                            <TableCell onClick={() => openDialog(file)}>
                                <p className="truncate">{file.file_name}</p>
                            </TableCell>

                            <TableCell onClick={() => openDialog(file)}>
                                <p className="truncate">{formatDate(file.modified_at)}</p>
                            </TableCell>

                            <TableCell>
                                <FileDropdown
                                    fileId={file.id}
                                    onDeleted={handleFileDeleted}
                                />
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>

            {activeFile && (
                <FileDialog
                    open={dialogOpen}
                    onOpenChange={setDialogOpen}
                    metadata={activeFile}
                    ipfsLink=""
                    spaceName=""
                    spaceDid=""
                />
            )}
        </div>
    )
}

/**
 * Hook for file pagination on scroll.
 * */
function usePagination(userId: string | null, limit: number = 20) {
    const [cursor] = useState<CursorReq | null>(null);
    const request = useMemo(() => ({ user_id: userId, cursor, limit }), [userId, cursor, limit]);
    return { request };
}

/**
 * Helper function to format date to human-readable format, backend
 * returns a timestamp.
 * @param date - date as a string.
 * @return string - date as string.
 * */
function formatDate(date: string): string {
    return new Date(date).toLocaleString()
}
