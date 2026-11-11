// 下書きの色分け。語の一覧は fixmarge の words() (Go の語の表) から受け取り、ここには書かない。
import { StreamLanguage } from "@codemirror/language";

// 語の使い道 → CodeMirror の字句の種類。先にあるものを優先する
const styles = [
  ["keyword", "keyword"],
  ["error", "atom"],
  ["package", "builtin"],
  ["function", "builtin"],
  ["builtin", "builtin"],
];

export function fxmaLanguage(words) {
  const styleOf = new Map();
  for (const w of words) {
    const style = styles.find(([use]) => w.uses.includes(use));
    if (style) styleOf.set(w.name, style[1]);
  }
  return StreamLanguage.define({
    name: "fxma",
    startState: () => ({ tripleString: false }),
    token(stream, state) {
      if (state.tripleString) return readTripleString(stream, state);
      if (stream.eatSpace()) return null;
      if (stream.match('"""')) {
        state.tripleString = true;
        return readTripleString(stream, state);
      }
      if (stream.peek() === '"') return readString(stream);
      if (stream.match("#")) {
        stream.skipToEnd();
        return "comment";
      }
      if (stream.match(/^\d[\d_]*(\.\d+)?([eE][+-]?\d+)?/)) return "number";
      const word = stream.match(/^[A-Za-z_][A-Za-z0-9_]*/);
      if (word) return styleOf.get(word[0]) ?? "variableName";
      stream.next();
      return "operator";
    },
    languageData: { commentTokens: { line: "#" } },
  });
}

function readString(stream) {
  stream.next();
  while (!stream.eol()) {
    const c = stream.next();
    if (c === "\\") stream.next();
    else if (c === '"') break;
  }
  return "string";
}

function readTripleString(stream, state) {
  while (!stream.eol()) {
    if (stream.match('"""')) {
      state.tripleString = false;
      break;
    }
    stream.next();
  }
  return "string";
}
