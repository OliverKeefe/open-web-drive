export type Commit = {
    id: string;
    author: `${string} <${string}>`;;
    date: Date
    message: string;
}