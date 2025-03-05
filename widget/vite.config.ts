import path from "path";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
// @ts-ignore
import namespace from "postcss-plugin-namespace";
import autoprefixer from "autoprefixer";
import tailwindcss from "tailwindcss";

export default defineConfig({
  css: {
    postcss: {
      plugins: [tailwindcss(), namespace(".announcable-widget"), autoprefixer],
    },
  },
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  build: {
    lib: {
      entry: "src/main.tsx",
      formats: ["umd"],
      name: "ReleaseBeaconWidget",
      fileName: () => `widget.js`,
    },
    rollupOptions: {
      external: [], // This ensures React is included in the bundle
    },
    outDir: "dist",
    emptyOutDir: true,
  },
  define: {
    "process.env": {},
  },
});
