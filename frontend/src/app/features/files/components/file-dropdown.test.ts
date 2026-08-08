import { createElement } from "react"
import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"
import FileDropdown from "./file-dropdown"

/**
 * Renders the FileDropdown with spied callbacks.
 * @return an object with the onInfo spy for assertions.
 * */
function renderDropdown() {
    const onInfo = vi.fn()
    render(
        createElement(FileDropdown, {
            fileId: "file-1",
            onDeleted: vi.fn(),
            onInfo,
        })
    )
    return { onInfo }
}

describe("FileDropdown", () => {
    /**
     * Verifies selecting File Info in the dropdown calls onInfo.
     * */
    it("calls onInfo when File Info is selected", async () => {
        const user = userEvent.setup()
        const { onInfo } = renderDropdown()

        await user.click(screen.getByRole("button"))
        await user.click(await screen.findByText("File Info"))

        expect(onInfo).toHaveBeenCalledTimes(1)
    })
})
