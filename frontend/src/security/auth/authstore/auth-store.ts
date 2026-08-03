import { create } from 'zustand'

interface AuthState {
    userId: string | null;
    token: string | null;
    username: string | null;
    setUserId: (id: string | null) => void;
    setToken: (token: string | null) => void;
    setUsername: (username: string | null) => void;
}

export const useAuthStore = create<AuthState>((set) => ({
    userId: null,
    token: null,
    username: null,
    setUserId: (id) => set({ userId: id }),
    setToken: (token) => set({ token: token }),
    setUsername: (username) => set({ username })
}));