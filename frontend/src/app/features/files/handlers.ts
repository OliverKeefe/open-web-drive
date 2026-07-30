import {RestHandler} from "@/app/features/shared/api/rest/rest-handler.ts";
import type { Metadata } from "@/app/features/files/types.ts";
import {UploadForm} from "@/app/features/shared/files/upload.ts";

export interface CursorReq {
    modified_at: string | null;
    id: string | null;
}

export interface GetAllMetadataReq {
    cursor: CursorReq;
    limit: number;
}

export interface GetAllMetadataResp {
    status: string;
    data: Metadata[];
}

const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;
if (!backendBaseUrl) throw new Error('VITE_BACKEND_BASE_URL is not set. Set it in your .env file or pass it at build time.');
const api = new RestHandler(backendBaseUrl);

export async function getAllMetadata(request: GetAllMetadataReq): Promise<GetAllMetadataResp> {
        const resp: any = await api.handlePost<GetAllMetadataReq, any>(`api/files/get-all`, request);
        console.log("RAW RESPONSE:", resp);
        return {
            status: resp.status || "fetched",
            data: Array.isArray(resp.data) ? resp.data : [],
        };
}

export async function uploadFiles(files: File[]) {
    const form = new UploadForm(files);
    await form.prepare();
    return await form.send();
}
