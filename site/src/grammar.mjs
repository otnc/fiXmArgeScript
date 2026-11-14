// ドキュメントのコードの色分け (Shiki) に、VS Code 拡張と同じ TextMate 文法を渡す。
// 文法は Go の語の表から internal/grammar が作っているので、ここには語を書かない
import { readFileSync } from "node:fs";

function load(name, aliases) {
  const url = new URL(`../../packaging/vscode/syntaxes/${name}.tmLanguage.json`, import.meta.url);
  return { ...JSON.parse(readFileSync(url, "utf8")), name, aliases };
}

export const fixmargeGrammars = [load("fixmarge", ["fxma"]), load("fixmargec", ["fxmac"]), load("fixmarge-mod", [])];
