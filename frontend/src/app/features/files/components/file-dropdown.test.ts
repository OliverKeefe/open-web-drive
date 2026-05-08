import { render, screen, fireEvent } from "@testing-library/react";
import FileDropdown from "@/app/features/files/components/file-dropdown"; // Adjust path as needed
import { useAuthStore } from "@/security/auth/authstore/auth-store";
import {jest} from "globals";

// 1. Mock the Auth Store
jest.mock("@/security/auth/authstore/auth-store", () => ({
    useAuthStore: jest.fn(),
}));

// 2. Mock notistack (since it requires a Provider otherwise)
jest.mock("notistack", () => ({
    enqueueSnackbar: jest.fn(),
}));

describe("FileDropdown", () => {
    const mockProps = {
        fileId: "123",
        fileName: "test-file.txt",
        onDeleted: jest.fn(),
    };

    beforeEach(() => {
        // Mock userId return value
        (useAuthStore as unknown as jest.Mock).mockReturnValue("user-456");
    });

    it("should display the menu items after clicking the trigger", () => {
        render(<FileDropdown {...mockProps} />);

        // Find and click the trigger button
        const trigger = screen.getByRole("button");
        fireEvent.click(trigger);

        // Verify items appear
        expect(screen.getByText("Delete")).toBeInTheDocument();
        expect(screen.getByText("Download")).toBeInTheDocument();
        expect(screen.getByText("File Info")).toBeInTheDocument();
        expect(screen.getByText("File Settings")).toBeInTheDocument();
    });
});
