/**
 * Helper function to format date to human-readable format, backend
 * returns a timestamp.
 * @param date - date as a string.
 * @return string - date as string.
 * */
export function formatDate(date: string): string {
    return new Date(date).toLocaleString()
}