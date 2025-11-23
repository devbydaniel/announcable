import path from "path";
import { defineConfig } from "vite";
import cssInjectedByJsPlugin from "vite-plugin-css-injected-by-js";

export default defineConfig({
  // PostCSS config is loaded from postcss.config.js
  plugins: [cssInjectedByJsPlugin()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  build: {
    lib: {
      entry: "src/main.ts",
      formats: ["umd"],
      name: "ReleaseBeaconWidget",
      fileName: () => `widget.js`,
    },
    rollupOptions: {
      external: [],
    },
    outDir: "dist",
    emptyOutDir: true,
  },
  define: {
    "process.env": {},
  },
});
