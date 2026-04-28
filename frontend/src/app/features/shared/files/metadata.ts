import MediaInfo from "mediainfo.js";
import {uploadObject} from "@/app/features/shared/ipfs/upload.ts";


/**
 * Metadata interface maps extracted metadata to object.
 * Fields `path`, `relativePath`, `lastModified`, `lastModifiedDate`
 * `size` and `fileType` are set using MediaInfo.js metadata
 * extraction. The rest, `id`, `ownerId` and `checkSum` are
 * injected once extraction is complete. `id` serves as a key to
 * associate each metadata object with its respective ByteData.
 * */
export type Metadata = ExtractedMetadata& {
    id: string
    ownerId: string
    checkSum: string
    cids: string[]
}

type ExtractedMetadata = {
    path: string
    relativePath: string
    lastModified: number
    lastModifiedDate: string
    size: number
    fileType: string
}

/**
 * <p>
 * extractMetadata uses MediaInfo to extract metadata from a file uploaded
 * via an UploadDialog's Dropzone component.
 * </p>
 * <p>
 * Metadata is extracted client-side and sent to backend via REST api due to
 * metadata loss that occurs with multipart form uploads.
 * </p>
 *
 * @param file - a File, uploaded via UploadDialog's Dropzone component.
 * @returns result - extracted metadata parsed to JSON.
 * */
export async function extractMetadata(file: File): Promise<Metadata> {
    const mediaInfo = await MediaInfo({
        format: "JSON",
        //locateFile: (path) => `/mediainfo/dist/${path}`,
        locateFile: () => `/MediaInfoModule.wasm`,
    })

    const [result, checksum] = await Promise.all([
        mediaInfo.analyzeData(
            () => file.size,
            async (chunkSize, offset) => {
                const buffer = await file.slice(offset, offset + chunkSize).arrayBuffer()
                return new Uint8Array(buffer)
        }
        ),
        getCheckSum(file)
    ]);

    mediaInfo.close()

    if (!result) {
        throw new Error("MediaInfo returned void")
    }

    //const media = JSON.parse(result);

    console.log('Metadata:', result);
    console.log('SHA-256:', checksum);
    return {
        path: file.name,
        relativePath: file.webkitRelativePath || file.name,
        lastModified: file.lastModified,
        lastModifiedDate: new Date(file.lastModified).toISOString(),
        size: file.size,
        fileType: file.type,

        //media,

        id: crypto.randomUUID(),
        ownerId: "29b6f168-03f6-4801-81d5-603b52f2c932",
        checkSum: checksum
    }
}

function generateFileId(): string {
    //TODO: Check UUID exist in db.
    return crypto.randomUUID();
}

function getOwnerId(): string {
    //TODO: Check UUID exist in db.
    return crypto.randomUUID();
}

// TODO: Compute this in backend instead, bad security and also >100mb browser dies.
async function getCheckSum(file: File): Promise<string> {
    const arrayBuffer = await file.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);

    return Array.from(new Uint8Array(hashBuffer))
        .map(b => b.toString(16).padStart(2, '0'))
        .join('');
}