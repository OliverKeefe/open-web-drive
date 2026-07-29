import {extractMetadata, type Metadata} from "@/app/features/shared/files/metadata";
import {useAuthStore} from "@/security/auth/authstore/auth-store.ts";
import { uploadObject } from "@/app/features/shared/ipfs/upload";
import { type ObjectType } from "@/app/features/shared/ipfs/types";

type PayloadItem = {
    id: string
    metadata: Metadata
    file: File
}

// type Payload = Record<string, PayloadItem>

export class UploadForm {
    private readonly files: File[];
    private payload: PayloadItem[];
    private readonly formData: FormData;

    constructor(files: File[]) {
        this.files = files;
        this.formData = new FormData();
        this.payload = [];
    }

    public async prepare() {
        for (let i = 0; i < this.files.length; i++) {
            const metadata = await extractMetadata(this.files[i]);
            const file = this.files[i];
            const id = metadata.id;
            this.payload.push({ id, metadata, file })
        }

    }

    private buildFormData() {
        Object.values(this.payload).forEach(({ metadata, file }) => {
            this.formData.append(
                `metadata-${metadata.id}`,
                JSON.stringify(metadata)
            )

            this.formData.append(
                `filedata-${metadata.id}`,
                file,
                file.name
            )
        })
    }

    public async send(): Promise<any> {
        this.buildFormData();

        const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;
        if (!backendBaseUrl) throw new Error('VITE_BACKEND_BASE_URL is not set');
        const url = `${backendBaseUrl}/api/files/upload`;
        const token = useAuthStore.getState().token;

        const response = await fetch(url, {
            method: "POST",
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: this.formData,
        });

        const contentType = response.headers.get("content-type");

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || `Upload failed (${response.status})`);
        }

        if (contentType?.includes("application/json")) {
            return await response.json();
        }

        return null;
    }

    private async handleFailedUpload(response: Response): Promise<void> {
        if (!response.ok) {
            const errorText = await response.text();
            console.error(`Upload failed: ${response.status} ${response.statusText}`, errorText);
            throw new Error(`Upload failed with status ${response.status}`);
        }
    }

    private async ipfsUpload(): Promise<Response> {
        for (const { id, metadata, file } of this.payload) {
            try {
                if (metadata.size > 4 * 1024 * 1024 * 1024) {
                    await uploadObject(id, file, metadata, ObjectType.CAR_FILE_SHARDS);
                }

                else if (metadata.fileType === "dir") {
                    await uploadObject(id, file, metadata, ObjectType.DIRECTORY);
                }

                else {
                    await uploadObject(id, file, metadata, ObjectType.FILE);
                }
            } catch(err) {
                console.log(`unable to upload file: ${metadata.id} to IPFS`, err);
            }
        }
    }
}

