import {
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router";
import Login from "@/app/features/auth/pages/login-page";
import { Files } from "@/app/features/files/pages/files-page";
import Settings from "@/app/features/settings/pages/settings-page";
import HotStoragePage from "@/app/features/storage/pages/hot-storage-page";
import ArchivePage from "@/app/features/archive/pages/archive-page";
import ProfilePage from "@/app/features/profile/pages/profile-page";
import Layout from "@/app/features/shared/components/layout/layout";
import VersionControlPage from "@/app/features/versions/pages/version-control-page";

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

const settingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/settings",
  component: Settings,
});

const versionControlRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/version-control",
  component: VersionControlPage,
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
  settingsRoute,
  versionControlRoute,
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
