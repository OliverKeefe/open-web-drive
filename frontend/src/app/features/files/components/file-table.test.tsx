import { render, screen, within } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, it, expect, vi, beforeEach } from "vitest"
import { FileTable } from "./file-table"
import { TooltipProvider } from "@/components/ui/tooltip"
import type { Metadata } from "@/app/features/files/types"
import { getAllMetadata } from "@/app/features/files/handlers"


vi.mock("@/app/features/files/handlers", () => ({
    getAllMetadata: vi.fn(),
}))

vi.mock("@/app/features/files/components/file-dropdown", () => ({
    default: ({ onInfo }: { onInfo: () => void }) => (
        <button onClick={onInfo}>file info</button>
    ),
}))

vi.mock("@/app/features/shared/components/dialog/upload-dialog", () => ({
    UploadDialog: () => <div data-testid="upload-dialog" />
}))

const mockGetAll = vi.mocked(getAllMetadata)

function makeMetadata(id: string, file_name: string, size: number): Metadata {
    return {
        id,
        file_id: `file-${id}`,
        version: 1,
        file_name,
        path: `/${file_name}`,
        relative_path: file_name,
        size,
        file_type: "application/octet-stream",
        hash: `hash-${id}`,
        created_at: "2026-01-01T10:00:00.000Z",
        modified_at: "2026-02-02T10:00:00.000Z",
        uploaded_at: "2026-03-03T10:00:00.000Z",
    }
}

async function renderWithFiles(files: Metadata[]) {
    mockGetAll.mockResolvedValue({ status: "ok", data: files })
    render(
        <TooltipProvider>
            <FileTable />
        </TooltipProvider>
    )
    await screen.findByText(files[0].file_name)
}

describe("FileTable + FileDialog", () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it("renders no dialog before File Info is requested", async() => {
        await renderWithFiles([makeMetadata("m1", "report.pdf", 2048)])

        expect(screen.queryByRole("dialog")).not.toBeInTheDocument()
    })

    it("opens the dialog when File Info is clicked in the dropdown", async () => {
        await renderWithFiles([makeMetadata("m1", "report.pdf", 2048)])
        const user = userEvent.setup()

        await user.click(screen.getByRole("button", { name: "file info" }))

        expect(await screen.findByRole("dialog")).toBeInTheDocument()
        expect(screen.getByRole("heading", { name: "report.pdf" })).toBeInTheDocument()
        expect(screen.getByText("File size: 2 KB")).toBeInTheDocument()
    })

    it("does not open the dialog when a file name is clicked", async () => {
        await renderWithFiles([makeMetadata("m1", "report.pdf", 2048)])
        const user = userEvent.setup()

        await user.click(screen.getByText("report.pdf"))

        expect(screen.queryByRole("dialog")).not.toBeInTheDocument()
    })

    it("closes the dialog when the close button is clicked", async () => {
        await renderWithFiles([makeMetadata("m1", "report.pdf", 2048)])
        const user = userEvent.setup()

        await user.click(screen.getByRole("button", { name: "file info" }))
        await user.click(await screen.findByRole("button", { name: "Close" }))

        expect(screen.queryByRole("dialog")).not.toBeInTheDocument()
    })

    describe("selection", () => {
        it("toggles selection when the file row is clicked", async () => {
            await renderWithFiles([makeMetadata("m1", "report.pdf", 2048)])
            const user = userEvent.setup()
            const row = screen.getByRole("row", { name: /report\.pdf/ })
            const checkbox = within(row).getByRole("checkbox")

            await user.click(screen.getByText("report.pdf"))
            expect(checkbox).toBeChecked()

            await user.click(screen.getByText("report.pdf"))
            expect(checkbox).not.toBeChecked()
        })

        it("does not double-toggle when the checkbox itself is clicked", async () => {
            await renderWithFiles([makeMetadata("m1", "report.pdf", 2048)])
            const user = userEvent.setup()
            const row = screen.getByRole("row", { name: /report\.pdf/ })

            await user.click(within(row).getByRole("checkbox"))

            expect(within(row).getByRole("checkbox")).toBeChecked()
        })

        it("does not toggle selection when the dropdown is interacted with", async () => {
            await renderWithFiles([makeMetadata("m1", "report.pdf", 2048)])
            const user = userEvent.setup()
            const row = screen.getByRole("row", { name: /report\.pdf/ })
            const checkbox = within(row).getByRole("checkbox")

            await user.click(screen.getByRole("button", { name: "file info" }))

            expect(checkbox).not.toBeChecked()
        })
    })
})