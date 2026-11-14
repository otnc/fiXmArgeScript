// spec/ の「言語」の章をドキュメント (/docs/spec/ と /docs/std/) に写す。仕様の元は spec/ だけにし、サイト用に書き直さない。
// 載せる章とその並びは spec/README.md の「言語」の一覧、標準パッケージは spec/std/README.md の一覧で決める。
// 写すときに、最初の見出しを Starlight の題にし、章どうしのリンクをサイトの URL に直す。
import { mkdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { dirname, join, posix } from "node:path";
import { fileURLToPath } from "node:url";

const root = fileURLToPath(new URL("../..", import.meta.url));
const out = fileURLToPath(new URL("../src/content/docs/docs", import.meta.url));
const repo = "https://github.com/otnc/fiXmArgeScript/blob/main";

// links は Markdown の一覧の行 ("- [題](path)" や表の "[`:meme`](meme.md)") から、リンク先を順に取り出す
function links(markdown) {
  return [...markdown.matchAll(/\]\(([^)#]+\.md)\)/g)].map((m) => m[1]);
}

// 「## 言語」の節だけを読む (「## 開発」の章はサイトに出さない)
function languageSection(markdown) {
  const start = markdown.indexOf("\n## 言語");
  const end = markdown.indexOf("\n## ", start + 1);
  return markdown.slice(start, end < 0 ? undefined : end);
}

const chapters = links(languageSection(readFileSync(join(root, "spec/README.md"), "utf8")));
// 一覧の表のリンクは同じフォルダのファイルだけ (頭の ../README.md は言語の目次へのリンク)
const packages = links(readFileSync(join(root, "spec/std/README.md"), "utf8"))
  .filter((p) => !p.includes("/"))
  .map((p) => posix.join("std", p));

// url は spec/ からのパス (std/meme.md) をサイトの URL にする。サイトに出さないものは GitHub のファイルにする
const published = new Map();
for (const p of chapters) published.set(p, p === "std/README.md" ? "/docs/std/" : `/docs/spec/${p.replace(/\.md$/, "")}/`);
for (const p of packages) published.set(p, `/docs/${p.replace(/\.md$/, "")}/`);

function url(specPath, hash) {
  const page = published.get(specPath);
  return page ? page + hash : `${repo}/${posix.normalize(posix.join("spec", specPath))}${hash}`;
}

function convert(specPath, order) {
  const source = readFileSync(join(root, "spec", specPath), "utf8");
  const [, title, body] = /^# (.+)\n([\s\S]*)$/.exec(source) ?? [];
  if (!title) throw new Error(`spec/${specPath} の1行目に「# 題」がありません`);
  // コード (``` の塊と `...`) の中の ](...) はリンクではないので、それ以外の所だけ直す
  const fixed = body
    .split(/(```[\s\S]*?```|`[^`\n]*`)/)
    .map((part, i) =>
      i % 2 === 1
        ? part
        : part.replace(/\]\((?!https?:|#|\/)([^)#]+)(#[^)]*)?\)/g, (_, target, hash = "") => {
            return `](${url(posix.join(posix.dirname(specPath), target), hash)})`;
          }),
    )
    .join("");
  const front = [
    "---",
    `title: ${JSON.stringify(title)}`,
    `editUrl: ${JSON.stringify(`https://github.com/otnc/fiXmArgeScript/edit/main/spec/${specPath}`)}`,
    `sidebar: { order: ${order} }`,
    "---",
    "",
  ].join("\n");
  const page = published.get(specPath).replace(/^\/docs\//, "").replace(/\/$/, "");
  const file = join(out, specPath === "std/README.md" ? "std/index.md" : `${page}.md`);
  mkdirSync(dirname(file), { recursive: true });
  writeFileSync(file, front + fixed);
}

rmSync(join(out, "spec"), { recursive: true, force: true });
rmSync(join(out, "std"), { recursive: true, force: true });
chapters.forEach((p, i) => convert(p, i));
packages.forEach((p, i) => convert(p, i + 1));
