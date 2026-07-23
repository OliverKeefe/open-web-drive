import {
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router";
import Login from "@/app/features/auth/pages/login-page";
import { Files } from "@/app/features/files/pages/files-page";
import Documents from "@/app/features/documents/pages/documents-page.tsx";
import Photos from "@/app/features/photos/pages/photos-page";
import Settings from "@/app/features/settings/pages/settings-page";
import GitHome from "@/app/features/git/pages/git-home-page";
import HotStoragePage from "@/app/features/storage/pages/hot-storage-page";
import ArchivePage from "@/app/features/archive/pages/archive-page";
import ProfilePage from "@/app/features/profile/pages/profile-page";
import Layout from "@/app/features/shared/components/layout/layout";

const rootRoute = createRootRoute({
  component: Layout,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: Files,
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/login",
  component: Login,
});

const filesRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/files",
  component: Files,
});

const documentsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/documents",
  component: Documents,
});

const photosRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/photos",
  component: Photos,
});

const settingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/settings",
  component: Settings,
});

const gitRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/git",
  component: GitHome,
});

const hotStorageRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/hot-storage-settings",
  component: HotStoragePage,
});

const archiveRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/archive",
  component: ArchivePage,
});

const profileRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/profile",
  component: ProfilePage,
});

const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  filesRoute,
  documentsRoute,
  photosRoute,
  settingsRoute,
  gitRoute,
  hotStorageRoute,
  archiveRoute,
  profileRoute,
]);

export const router = createRouter({ routeTree });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
