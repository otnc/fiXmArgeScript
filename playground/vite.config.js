import { defineConfig } from "vite";

import { cloudflare } from "@cloudflare/vite-plugin";

export default defineConfig({
  // fxma は new URL("./fxma.wasm", import.meta.url) で wasm を探すので、まとめ直さずにそのまま使う
  optimizeDeps: { exclude: ["fxma"] },
  plugins: [cloudflare()],
});
