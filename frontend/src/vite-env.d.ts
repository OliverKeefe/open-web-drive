/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_KEYCLOAK_URL: string;
    readonly VITE_KEYCLOAK_REALM: string;
    readonly VITE_KEYCLOAK_CLIENT_ID: string;
    readonly VITE_BACKEND_BASE_URL: string;
    readonly VITE_IPFS_GATEWAY_URL: string;
    readonly VITE_MAX_UPLOAD_SIZE_BYTES: string;
    readonly VITE_STORACHA_DID: string;
    readonly VITE_STORACHA_ACCOUNT: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
