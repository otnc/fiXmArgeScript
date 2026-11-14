import cloudflare from "@astrojs/cloudflare";
import starlight from "@astrojs/starlight";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "astro/config";

import { fixmargeGrammars } from "./src/grammar.mjs";

export default defineConfig({
  site: "https://fixmargescript.oto.im",
  // 紹介とドキュメントは作っておき、/pkg/* だけを Worker で開かれたときに作る
  // 画像はビルドのときに整え、セッションは使わない (どちらも既定では Cloudflare の Images と KV を足してしまう)。
  // 作っておくページは Node で作る。Markdown を読む satteri が workerd 向けには wasm 版を探して失敗するため
  adapter: cloudflare({ imageService: "compile", prerenderEnvironment: "node" }),
  session: false,
  integrations: [
    starlight({
      title: "fiXmArge Script",
      description: "「魔法少女ノ魔女裁判」のフィクスマージ語をキーワードにした、非公式のファンメイドのプログラミング言語です。",
      logo: { dark: "../assets/logo-dark.svg", light: "../assets/logo-light.svg", replacesTitle: true },
      favicon: "/favicon.svg",
      defaultLocale: "root",
      locales: { root: { label: "日本語", lang: "ja" } },
      social: [{ icon: "github", label: "GitHub", href: "https://github.com/otnc/fiXmArgeScript" }],
      editLink: { baseUrl: "https://github.com/otnc/fiXmArgeScript/edit/main/site/" },
      customCss: ["./src/styles/global.css"],
      components: { Footer: "./src/components/Footer.astro" },
      // ドキュメントのコードの色分けに、エディタと同じ TextMate 文法を使う
      expressiveCode: {
        shiki: { langs: fixmargeGrammars },
        themes: ["github-dark", "github-light"],
        // 背景と枠はサイトの色に合わせる (テーマの青みがかった灰色は黒地から浮くため)
        styleOverrides: {
          codeBackground: "var(--sl-color-gray-6)",
          borderColor: "var(--sl-color-hairline)",
          borderRadius: "0",
          frames: {
            editorTabBarBackground: "var(--sl-color-gray-6)",
            editorActiveTabBackground: "var(--sl-color-gray-5)",
            editorActiveTabIndicatorTopColor: "var(--sl-color-accent)",
            terminalTitlebarBackground: "var(--sl-color-gray-6)",
            terminalBackground: "var(--sl-color-gray-6)",
          },
        },
      },
      sidebar: [{ label: "はじめに", items: [{ autogenerate: { directory: "docs/start" } }] }],
    }),
  ],
  vite: {
    plugins: [tailwindcss()],
    // fixmarge は new URL("./fixmarge.wasm", import.meta.url) で wasm を探すので、まとめ直さずにそのまま使う
    optimizeDeps: { exclude: ["fixmarge"] },
  },
});
