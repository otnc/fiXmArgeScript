// vsce package の前に、リポジトリ直下の2つのライセンスを LICENSE にまとめる。vsce は LICENSE という名前のファイルを探すため。
// できたものはコミットしない (.gitignore)。
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = join(import.meta.dirname, "..", "..");
const texts = ["LICENSE-MIT", "LICENSE-APACHE"].map((name) => readFileSync(join(root, name), "utf8"));
writeFileSync(
  join(import.meta.dirname, "LICENSE"),
  "MIT と Apache-2.0 のどちらかを選べる二重ライセンスです。\n\n" + texts.join("\n\n"),
);
