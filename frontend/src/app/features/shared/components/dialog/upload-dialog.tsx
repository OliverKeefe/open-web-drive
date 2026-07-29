import React, { useState } from "react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import {
    Dropzone,
    DropzoneContent,
    DropzoneEmptyState,
} from "@/components/ui/shadcn-io/dropzone";
import { Upload } from "lucide-react";
import { UploadForm } from "@/app/features/shared/files/upload";

type UploadDialogProps = {
    onUploaded?: () => void;
};

export function UploadDialog({ onUploaded }: UploadDialogProps) {
    const [open, setDialogOpen] = useState(false);
    const [files, setFiles] = useState<File[] | null>(null);
    const [isUploading, setIsUploading] = useState(false);

    const handleDrop = (newFiles: File[]) => {
        setFiles((prev) => {
            const existing = prev ?? [];
            const merged = [...existing, ...newFiles];

            return Array.from(
                new Map(
                    merged.map((f) => [
                        `${f.name}-${f.size}-${f.lastModified}`,
                        f,
                    ])
                ).values()
            );
        });
    };

    async function handleDialogUpload(): Promise<void> {
        if (!files || files.length === 0) return;

        try {
            setIsUploading(true);

            const upload = new UploadForm(files);
            await upload.prepare();
            await upload.send();

            onUploaded?.();

            setDialogOpen(false);
            setFiles(null);
        } catch (err) {
            console.error("Upload failed:", err);
        } finally {
            setIsUploading(false);
        }
    }

    return (
        <Dialog
            open={open}
            onOpenChange={(isOpen) => {
                setDialogOpen(isOpen);
                if (!isOpen) setFiles(null);
            }}
        >
            <DialogTrigger asChild>
                <Button variant="default">
                    <Upload className="mr-2 h-4 w-4" />
                    Upload
                </Button>
            </DialogTrigger>

            <DialogContent className="max-w-3xl">
                <DialogHeader>
                    <DialogTitle>Upload File</DialogTitle>
                    <DialogDescription>
                        Upload a file or folder.
                    </DialogDescription>
                </DialogHeader>

                <Dropzone
                    maxFiles={10}
                    maxSize={Number(import.meta.env.VITE_MAX_UPLOAD_SIZE_BYTES)}
                    minSize={1}
                    onDrop={handleDrop}
                    onError={console.error}
                    src={files}
                >
                    <DropzoneEmptyState />
                    <DropzoneContent />
                </Dropzone>

                <div className="flex flex-col space-y-2 mt-4">
                    <Button
                        onClick={handleDialogUpload}
                        disabled={!files || isUploading}
                    >
                        {isUploading ? "Uploading..." : "Upload"}
                    </Button>
                </div>
            </DialogContent>
        </Dialog>
    );
}
