import MediaInfo from "mediainfo.js";


/**
 * Metadata interface maps extracted metadata to object.
 * Fields `path`, `relativePath`, `lastModified`, `lastModifiedDate`
 * `size` and `fileType` are set using MediaInfo.js metadata
 * extraction. The `id` field is a client-generated UUID used as the
 * multipart form part key to associate metadata with its file data.
 *
 * Property names match the Go multipartMetadata json tags exactly
 * so JSON.stringify() produces the correct wire format.
 */
export type Metadata = ExtractedMetadata & {
    id: string
    file_name: string
    created_at: number
    uploadedAt: number
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

    const result = await mediaInfo.analyzeData(
        () => file.size,
        async (chunkSize, offset) => {
            const buffer = await file.slice(offset, offset + chunkSize).arrayBuffer()
            return new Uint8Array(buffer)
        }
    )

    mediaInfo.close()

    if (!result) {
        throw new Error("MediaInfo returned void")
    }

    //const media = JSON.parse(result);

    console.log('Metadata:', result);
    return {
        path: file.name,
        relativePath: file.webkitRelativePath || file.name,
        lastModified: file.lastModified,
        lastModifiedDate: new Date(file.lastModified).toISOString(),
        size: file.size,
        fileType: file.type,

        //media,

        id: crypto.randomUUID(),
        file_name: file.name,
        created_at: file.lastModified,
        uploadedAt: Date.now(),
    }
}

// TODO: Compute this in backend instead, bad security and also >100mb browser dies.
async function getCheckSum(file: File): Promise<string> {
    const arrayBuffer = await file.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);

    return Array.from(new Uint8Array(hashBuffer))
        .map(b => b.toString(16).padStart(2, '0'))
        .join('');
}