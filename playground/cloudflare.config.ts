import { defineConfig } from "cf/config";

// プレイグラウンドを Cloudflare Workers の静的アセットとして https://fixmargescript.oto.im/ に置く。
// 後で install.sh も同じ場所に置くので、Worker の名前はプレイグラウンドに限らない名前にする
export default defineConfig({
  worker: {
    name: "fixmargescript",
    compatibilityDate: "2026-11-04",
    domains: ["fixmargescript.oto.im"],
    observability: { enabled: true },
    assets: { notFoundHandling: "single-page-application" },
  },
});
