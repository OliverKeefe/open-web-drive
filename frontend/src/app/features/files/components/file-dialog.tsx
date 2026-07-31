import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog.tsx";
import type {Metadata} from "@/app/features/files/types.ts";
import {useState} from "react";
import {Button} from "@/components/ui/button.tsx";
import {EllipsisVertical} from "lucide-react";
import {DialogDescription} from "@radix-ui/react-dialog";

interface FileDialogProps{
    open: boolean,
    onOpenChange: (open: boolean) => void,
    metadata: Metadata,
    ipfsLink: string,
    spaceName: string,
    spaceDid: string
}

export function FileDialog({
                               open,
                               onOpenChange,
                               metadata,
                               ipfsLink,
                               spaceName,
                               spaceDid,
                           }: FileDialogProps) {
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>{metadata.file_name}</DialogTitle>
                    <DialogDescription>
                        File info and backup uris
                    </DialogDescription>
                </DialogHeader>

                <h2>File size: {metadata.size / 1024} KB</h2>
                <h2>Last modified at: {formatDate(metadata.modified_at)}</h2>
                <h2>
                    Uploaded at:{" "}
                    {metadata.uploaded_at
                        ? formatDate(metadata.uploaded_at)
                        : "11/02/2026 11:38:42"}
                </h2>
                <h2>{metadata.id}</h2>
                <h2>{metadata.owner_id}</h2>

                <h2>Backup</h2>
                <h3>Space: test-space</h3>
                <p>Visibility: PUBLIC</p>
                <p className={"truncate"}>DID: bafybeia7wkemsgryogneimjafwwkb33ifwh2oo3djba3lqfeg3lkrqn464</p>
                <p className={"truncate"}>Shards: bagbaierahrldusuunn3mt3xcgfue3aav6zcynk7ynpwzhgyi4l6muyp4hjhq</p>
                <h1>Recovery </h1>
                    <Button variant={"outline"}>
                        <a href={import.meta.env.VITE_IPFS_GATEWAY_URL}> Fetch from IPFS</a>
                    </Button>
            </DialogContent>
        </Dialog>
    );
}

//TODO: remove this duplcate function at some point
function formatDate(date: string): string {
    const d = new Date(date);

    if (isNaN(d.getTime())) {
        return "11/02/2026 11:38AM";
    }

    return d.toLocaleString();
}
