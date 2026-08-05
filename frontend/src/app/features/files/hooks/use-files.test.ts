import { renderHook, act, waitFor } from "@testing-library/react"
import { describe, it, expect, vi, beforeEach } from "vitest"
import { useFiles } from "./use-files"
import type { GetAllMetadataResp } from "@/app/features/files/handlers"
import type { Metadata } from "@/app/features/files/types"

vi.mock("@/app/features/files/handlers", () => ({ getAllMetadata: vi.fn() }))

import { getAllMetadata } from "@/app/features/files/handlers"
const mockGetAll = vi.mocked(getAllMetadata)

function makeMeta(ts: string, id: string, name?: string): Metadata {
	return {
		id,
		file_id: `file-${id}`,
		version: 1,
		file_name: name ?? `${id}.txt`,
		path: `/${id}.txt`,
		relative_path: `${id}.txt`,
		size: 100,
		file_type: "text/plain",
		hash: `hash-${id}`,
		created_at: ts,
		modified_at: ts,
		uploaded_at: ts,
	}
}

function page(n: number, offset = 0): Metadata[] {
	return Array.from({ length: n }, (_, i) => {
		const idx = offset + i
		return makeMeta(
			new Date(2024, 0, 30 - idx).toISOString(),
			`uuid-${idx}`,
			`file-${idx}.txt`,
		)
	})
}

describe("useFiles", () => {
	beforeEach(() => {
		vi.clearAllMocks()
	})

	it("fetches the first page on mount", async () => {
		mockGetAll.mockResolvedValue({ status: "ok", data: page(5) } as GetAllMetadataResp)

		const { result } = renderHook(() => useFiles({ limit: 20 }))

		expect(result.current.isLoading).toBe(true)

		await waitFor(() => expect(result.current.isLoading).toBe(false))

		expect(mockGetAll).toHaveBeenCalledTimes(1)
		expect(mockGetAll).toHaveBeenCalledWith({
			cursor: { modified_at: null, id: null },
			limit: 20,
		})
		expect(result.current.files).toHaveLength(5)
		expect(result.current.hasMore).toBe(false)
	})

	it("sets hasMore when response equals the page limit", async () => {
		mockGetAll.mockResolvedValue({ status: "ok", data: page(20) } as GetAllMetadataResp)

		const { result } = renderHook(() => useFiles({ limit: 20 }))

		await waitFor(() => expect(result.current.isLoading).toBe(false))

		expect(result.current.hasMore).toBe(true)
	})

	it("appends data on loadMore and derives cursor from the last item", async () => {
		mockGetAll
			.mockResolvedValueOnce({ status: "ok", data: page(20) } as GetAllMetadataResp)
			.mockResolvedValueOnce({ status: "ok", data: page(5, 20) } as GetAllMetadataResp)

		const { result } = renderHook(() => useFiles({ limit: 20 }))

		await waitFor(() => expect(result.current.isLoading).toBe(false))

		await act(async () => {
			await result.current.loadMore()
		})

		expect(mockGetAll).toHaveBeenCalledTimes(2)
		expect(result.current.files).toHaveLength(25)

		const afterFirstPage = result.current.files[19]
		expect(mockGetAll).toHaveBeenNthCalledWith(2, {
			cursor: { modified_at: afterFirstPage.modified_at, id: afterFirstPage.id },
			limit: 20,
		})
	})

	it("loadMore is no-op when hasMore is false", async () => {
		mockGetAll.mockResolvedValue({ status: "ok", data: page(3) } as GetAllMetadataResp)

		const { result } = renderHook(() => useFiles({ limit: 20 }))

		await waitFor(() => expect(result.current.isLoading).toBe(false))
		expect(result.current.hasMore).toBe(false)

		await act(async () => {
			await result.current.loadMore()
		})

		expect(mockGetAll).toHaveBeenCalledTimes(1)
	})

	it("loadMore is no-op when cursor is null (empty initial response)", async () => {
		mockGetAll.mockResolvedValue({ status: "ok", data: [] } as GetAllMetadataResp)

		const { result } = renderHook(() => useFiles({ limit: 20 }))

		await waitFor(() => expect(result.current.isLoading).toBe(false))

		await act(async () => {
			await result.current.loadMore()
		})

		expect(mockGetAll).toHaveBeenCalledTimes(1)
	})

	it("prevents concurrent loadMore calls", async () => {
		let resolveDeferred!: (v: GetAllMetadataResp) => void
		const deferred = new Promise<GetAllMetadataResp>((r) => { resolveDeferred = r })

		mockGetAll
			.mockResolvedValueOnce({ status: "ok", data: page(20) } as GetAllMetadataResp)
			.mockReturnValueOnce(deferred)

		const { result } = renderHook(() => useFiles({ limit: 20 }))

		await waitFor(() => expect(result.current.isLoading).toBe(false))

		act(() => { void result.current.loadMore() })
		await waitFor(() => expect(result.current.isLoadingMore).toBe(true))

		await act(async () => {
			await result.current.loadMore()
		})

		await act(async () => {
			resolveDeferred({ status: "ok", data: page(5, 20) } as GetAllMetadataResp)
		})
		await waitFor(() => expect(result.current.isLoadingMore).toBe(false))

		expect(mockGetAll).toHaveBeenCalledTimes(2)
	})

	it("refresh resets state and re-fetches the first page", async () => {
		mockGetAll
			.mockResolvedValueOnce({ status: "ok", data: page(5) } as GetAllMetadataResp)
			.mockResolvedValueOnce({ status: "ok", data: page(3) } as GetAllMetadataResp)

		const { result } = renderHook(() => useFiles({ limit: 20 }))

		await waitFor(() => expect(result.current.isLoading).toBe(false))
		expect(result.current.files).toHaveLength(5)

		await act(async () => {
			await result.current.refresh()
		})

		expect(mockGetAll).toHaveBeenCalledTimes(2)
		expect(result.current.files).toHaveLength(3)
		expect(mockGetAll).toHaveBeenNthCalledWith(2, {
			cursor: { modified_at: null, id: null },
			limit: 20,
		})
	})

	it("accumulates three pages correctly", async () => {
		mockGetAll
			.mockResolvedValueOnce({ status: "ok", data: page(20) } as GetAllMetadataResp)
			.mockResolvedValueOnce({ status: "ok", data: page(20, 20) } as GetAllMetadataResp)
			.mockResolvedValueOnce({ status: "ok", data: page(5, 40) } as GetAllMetadataResp)

		const { result } = renderHook(() => useFiles({ limit: 20 }))

		await waitFor(() => expect(result.current.isLoading).toBe(false))
		expect(result.current.files).toHaveLength(20)

		await act(async () => { await result.current.loadMore() })
		expect(result.current.files).toHaveLength(40)
		expect(result.current.hasMore).toBe(true)

		await act(async () => { await result.current.loadMore() })
		expect(result.current.files).toHaveLength(45)
		expect(result.current.hasMore).toBe(false)
	})
})
