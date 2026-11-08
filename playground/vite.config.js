import { defineConfig } from "vite";

export default defineConfig({
  // GitHub Pages の https://otnc.github.io/fiXmArgeScript/ のように、どのパスに置いても動くようにする
  base: "./",
  // fxma は new URL("./fxma.wasm", import.meta.url) で wasm を探すので、まとめ直さずにそのまま使う
  optimizeDeps: { exclude: ["fxma"] },
});
