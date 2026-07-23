import {RestHandler} from "@/app/features/shared/api/rest/rest-handler.ts";
import type { Metadata } from "@/app/features/files/types.ts";
import {UploadForm} from "@/app/features/shared/files/upload.ts";

export interface CursorReq {
    modified_at: string | null;
    id: string | null;
}

export interface GetAllMetadataReq {
    // TODO: remove user_id from request body — backend should derive
    // this from the JWT to prevent IDOR. Currently the backend trusts
    // the client-sent value, which is a security risk.
    user_id: string | null;
    cursor: CursorReq;
    limit: number;
}

export interface GetAllMetadataResp {
    status: string;
    data: Metadata[];
}

const api = new RestHandler(`http://127.0.0.1:8081`);

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
