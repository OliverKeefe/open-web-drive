/**
* Handles client-side REST requests and responses (POST, GET, PUT, DELETE).
* */
import { useAuthStore } from "@/security/auth/authstore/auth-store.ts";

export class RestHandler {
    private readonly baseURL: string;

    /**
     * @constructor
     * @param baseURL The base URL for REST api calls to the Gestalto Control Plane Backend.
     * */
    constructor(baseURL: string) {
        this.baseURL = baseURL;
    }

    /**
     * Gets the user UUID from keycloak kc.token.sub in AuthStore.
     * @return string | null of userId (uuid as string or null val).
     * */
    private get userId(): string | null {
        const { userId } = useAuthStore.getState();
        return userId;
    }

    /**
     * Gets the JWT Token from keycloak kc.token in AuthStore.
     *
     * @return string | null of kc.token (as string or null val).
     * */
    private get token(): string | null {
        const { token } = useAuthStore.getState();
        return token;
    }

    /**
     * Handles logic for failed HTTP request.
     * @return Promise<void>
     * @throws Error containing HTTP status and information backend allows frontned to see.
     * */
    private async handleFailedRequest(response: Response): Promise<void> {
        if (!response.ok) {
            const errorText = await response.text();
            console.error(`Request failed: ${response.status} ${response.statusText}`, errorText);
            throw new Error(`Request failed with status ${response.status}`);
        }
    }

    /**
     * Handles HTTP GET request.
     * @param endpoint string of the api endpoint uri (e.g. `/files`).
     * @return Promise<R> response - R being a generic representing the response type.
     * */
    public async handleGet<R = unknown>(endpoint: string): Promise<R> {
        const userId = this.userId;
        const url = `${this.baseURL}/${endpoint}${userId}`;
        const options: RequestInit = { method: "GET" };
        const response = await fetch(url, options);
        await this.handleFailedRequest(response);
        return await response.json();
    }

    public async handlePost<T, R = unknown>(endpoint: string, payload: T): Promise<R> {
        const token = useAuthStore.getState().token;
        const url = `${this.baseURL}/${endpoint}`;
        const options: RequestInit = {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(payload)
        };

        const response = await fetch(url, options);
        await this.handleFailedRequest(response);
        return await response.json();
    }

    public async handleDelete<T, R = unknown>(endpoint: string, payload: T): Promise<R> {
        const token = useAuthStore.getState().token;
        const url = `${this.baseURL}/${endpoint}`;
        const options: RequestInit = {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(payload)
        };

        const response = await fetch(url, options);
        await this.handleFailedRequest(response);
        return await response.json();
    }

    public async handleDownload<T>(endpoint: string, payload: T): Promise<void> {
        const token = useAuthStore.getState().token;
        const url = `${this.baseURL}/${endpoint}`;
        const options: RequestInit = {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(payload)
        };

        const response = await fetch(url, options);
        await this.handleFailedRequest(response);

        const blob = await response.blob();
        const objectUrl = URL.createObjectURL(blob);

        const a = document.createElement("a");
        a.href = objectUrl;
        a.download = this.getFileNameFromResponse(response) || "download";
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(objectUrl);
    }

    private getFileNameFromResponse(response: Response): string | null {
        const disposition = response.headers.get("Content-Disposition");
        if (!disposition) return null;
        const match = disposition.match(/filename="(.+)"/);
        return match ? match[1] : null;
    }
}