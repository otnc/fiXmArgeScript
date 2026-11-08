import { defineConfig } from "vite";

import { cloudflare } from "@cloudflare/vite-plugin";

export default defineConfig({
  // install.sh はインストーラと一緒に installer/ で管理し、同じサイトの /install.sh で配る
  publicDir: "../installer/public",
  // fxma は new URL("./fxma.wasm", import.meta.url) で wasm を探すので、まとめ直さずにそのまま使う
  optimizeDeps: { exclude: ["fxma"] },
  plugins: [cloudflare()],
});
