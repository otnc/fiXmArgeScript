// vsce package の前に、リポジトリ直下の2つのライセンスをこのフォルダに写す。どちらもコミットしない (.gitignore)。
// npm や Releases の配布物と同じく2つのファイルのまま入れる。LICENSE という1つのファイルはないので、vsce は --skip-license で動かす。
import { copyFileSync } from "node:fs";
import { join } from "node:path";

const root = join(import.meta.dirname, "..", "..");
for (const name of ["LICENSE-MIT", "LICENSE-APACHE"]) {
  copyFileSync(join(root, name), join(import.meta.dirname, name));
}
