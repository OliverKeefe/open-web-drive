import path from "path"
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import {defineConfig, type Plugin} from "vite"
import fs from "fs"

function copyWasmModules(): Plugin {
    const filePaths = [
        'node_modules/mediainfo.js/dist/MediaInfoModule.wasm',
    ]

    const copy = () => {
        const publicDir = path.resolve(__dirname, 'public')

        for (const file of filePaths) {
            const src = path.resolve(__dirname, file)
            const dest = path.join(publicDir, path.basename(file))

            fs.mkdirSync(publicDir, { recursive: true })
            fs.copyFileSync(src, dest)
        }
    }

    return {
        name: 'copy-wasm',
        apply: 'serve',
        configureServer() {
            copy()
        },
        buildStart() {
            copy()
        }
    }
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [
      react(),
      tailwindcss(),
      copyWasmModules()
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
})

