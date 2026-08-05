import { useCallback, useEffect, useRef, useState } from "react"
import {
	type CursorReq,
	getAllMetadata,
} from "@/app/features/files/handlers"
import type { Metadata } from "@/app/features/files/types"

export interface UseFilesOptions {
	limit?: number
}

export interface UseFilesReturn {
	files: Metadata[]
	isLoading: boolean
	isLoadingMore: boolean
	hasMore: boolean
	loadMore: () => Promise<void>
	refresh: () => Promise<void>
}

export function useFiles({ limit = 20 }: UseFilesOptions = {}): UseFilesReturn {
	const [files, setFiles] = useState<Metadata[]>([])
	const [cursor, setCursor] = useState<CursorReq | null>(null)
	const [isLoading, setIsLoading] = useState(true)
	const [isLoadingMore, setIsLoadingMore] = useState(false)
	const [hasMore, setHasMore] = useState(true)

	const loadingRef = useRef(false)

	const reset = useCallback(async () => {
		setIsLoading(true)
		setFiles([])
		setCursor(null)
		setHasMore(true)

		const resp = await getAllMetadata({ cursor: { modified_at: null, id: null }, limit })
		const data = resp.data

		setFiles(data)
		if (data.length < limit) {
			setHasMore(false)
		}
		if (data.length > 0) {
			const last = data[data.length - 1]
			setCursor({ modified_at: last.modified_at, id: last.id })
		}
		setIsLoading(false)
	}, [limit])

	useEffect(() => { reset() }, [reset])

	const loadMore = useCallback(async () => {
		if (!hasMore || loadingRef.current || !cursor) return

		loadingRef.current = true
		setIsLoadingMore(true)

		try {
			const resp = await getAllMetadata({ cursor, limit })
			const data = resp.data

			if (data.length === 0) {
				setHasMore(false)
			} else {
				setFiles(prev => [...prev, ...data])
				if (data.length < limit) {
					setHasMore(false)
				}
				const last = data[data.length - 1]
				setCursor({ modified_at: last.modified_at, id: last.id })
			}
		} finally {
			setIsLoadingMore(false)
			loadingRef.current = false
		}
	}, [cursor, hasMore, limit])

	return { files, isLoading, isLoadingMore, hasMore, loadMore, refresh: reset }
}
