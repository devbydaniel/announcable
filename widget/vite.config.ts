import path from "path";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import cssInjectedByJsPlugin from "vite-plugin-css-injected-by-js";
// @ts-ignore
import namespace from "postcss-plugin-namespace";
import autoprefixer from "autoprefixer";
import tailwindcss from "tailwindcss";

export default defineConfig({
  css: {
    postcss: {
      plugins: [
        tailwindcss(),
        namespace(".release-beacon-widget"),
        autoprefixer,
      ],
    },
  },
  plugins: [react(), cssInjectedByJsPlugin()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  build: {
    lib: {
      entry: "src/main.tsx",
      name: "ReleaseBeaconWidget",
      fileName: (format) => `release-beacon-widget.${format}.js`,
    },
    rollupOptions: {
      external: [], // This ensures React is included in the bundle
    },
  },
  define: {
    "process.env": {},
  },
});
