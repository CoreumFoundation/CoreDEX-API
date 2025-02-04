import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import * as path from "path";
import nodePolyfills from "vite-plugin-node-stdlib-browser";
// import mkcert from "vite-plugin-mkcert";
import { fileURLToPath } from "url";

// Get the equivalent of __dirname
const __dirname = path.dirname(fileURLToPath(import.meta.url));

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), nodePolyfills()],
  build: {
    assetsDir: "trade",
  },
  resolve: {
    alias: {
      react: path.resolve("./node_modules/react"),
      "react-dom": path.resolve("./node_modules/react-dom"),
      "@": path.resolve(__dirname, "src"),
    },
  },
  server: {
    host: "localhost",
    port: 3000,
  },
  base: "/",
});
