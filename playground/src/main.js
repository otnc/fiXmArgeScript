import { basicSetup, EditorView } from "codemirror";
import { HighlightStyle, indentUnit, syntaxHighlighting } from "@codemirror/language";
import { EditorState } from "@codemirror/state";
import { tags } from "@lezer/highlight";
import { build, run, words } from "fixmarge";

import sample from "../../examples/fizzbuzz.fxma?raw";
import { fxmaLanguage } from "./language.js";

const mainFile = "main.fxma";

// 止まらないプログラムを書くとタブごと固まるので、文と繰り返しの数に上限を付ける。
// 空の繰り返しなら1秒足らず、関数を呼び続けるものでも10秒ほどで届く数
const maxSteps = 5_000_000;

// 開いているファイル。main.fxma は消せない。ほかのファイルは werer でも :memoid でも読める
const files = new Map([[mainFile, sample]]);
let current = mainFile;

const tabs = document.getElementById("tabs");
const output = document.getElementById("output");
const stdin = document.getElementById("stdin");
const runButton = document.getElementById("run");
const buildButton = document.getElementById("build");

const language = fxmaLanguage(await words());

// 色はページの CSS 変数から取り、明るい画面でも暗い画面でもページと合わせる
const theme = EditorView.theme({
  "&": { backgroundColor: "var(--panel)", color: "var(--text)", height: "100%" },
  ".cm-gutters": { backgroundColor: "var(--panel)", color: "var(--muted)", borderRight: "1px solid var(--line)" },
  ".cm-activeLine, .cm-activeLineGutter": { backgroundColor: "var(--active-line)" },
  ".cm-cursor": { borderLeftColor: "var(--text)" },
  "&.cm-focused .cm-selectionBackground, .cm-selectionBackground": { backgroundColor: "var(--selection)" },
});
const highlight = HighlightStyle.define([
  { tag: tags.keyword, color: "var(--syntax-keyword)", fontWeight: "bold" },
  { tag: tags.atom, color: "var(--syntax-error)" },
  { tag: tags.standard(tags.variableName), color: "var(--syntax-builtin)" },
  { tag: tags.string, color: "var(--syntax-string)" },
  { tag: tags.number, color: "var(--syntax-number)" },
  { tag: tags.comment, color: "var(--muted)", fontStyle: "italic" },
]);
const view = new EditorView({
  parent: document.getElementById("editor"),
  state: stateFor(sample),
});

function stateFor(doc) {
  return EditorState.create({
    doc,
    extensions: [basicSetup, language, theme, syntaxHighlighting(highlight), indentUnit.of("  "), EditorState.tabSize.of(2)],
  });
}

function open(name) {
  files.set(current, view.state.doc.toString());
  current = name;
  view.setState(stateFor(files.get(name)));
  renderTabs();
}

function close(name) {
  files.delete(name);
  if (current === name) {
    current = mainFile;
    view.setState(stateFor(files.get(mainFile)));
  }
  renderTabs();
}

function renderTabs() {
  tabs.replaceChildren(
    ...[...files.keys()].map((name) => {
      const tab = document.createElement("span");
      tab.className = "tab";
      tab.classList.toggle("active", name === current);
      const label = document.createElement("button");
      label.textContent = name;
      label.addEventListener("click", () => open(name));
      tab.append(label);
      if (name !== mainFile) {
        const x = document.createElement("button");
        x.className = "close";
        x.textContent = "×";
        x.setAttribute("aria-label", `${name} を閉じる`);
        x.addEventListener("click", () => close(name));
        tab.append(x);
      }
      return tab;
    }),
  );
}

document.getElementById("new-file").addEventListener("submit", (event) => {
  event.preventDefault();
  const input = document.getElementById("new-file-name");
  const name = input.value.trim();
  if (name === "" || files.has(name)) return;
  files.set(name, "");
  input.value = "";
  open(name);
  view.focus();
});

function show(stdout, stderr, note) {
  const parts = [stdout];
  if (stderr) parts.push(span("stderr", stderr));
  if (note) parts.push(span("note", note));
  output.replaceChildren(...parts);
}

function span(className, text) {
  const s = document.createElement("span");
  s.className = className;
  s.textContent = text;
  return s;
}

async function busy(work) {
  runButton.disabled = buildButton.disabled = true;
  try {
    await work();
  } catch (error) {
    // fixmarge.wasm を読み込めなかったときなど。何も出さないと、ボタンを押しても動かないように見える
    show("", `SaiVI ${error instanceof Error ? error.message : error}\n`);
  } finally {
    runButton.disabled = buildButton.disabled = false;
  }
}

runButton.addEventListener("click", () =>
  busy(async () => {
    files.set(current, view.state.doc.toString());
    const others = Object.fromEntries([...files].filter(([name]) => name !== mainFile));
    const r = await run(files.get(mainFile), { stdin: stdin.value, files: others, maxSteps });
    show(r.stdout, r.stderr, `shid ${r.exitCode}`);
    // :memoid で書いたファイルもタブに出す
    for (const [name, content] of Object.entries(r.files)) files.set(name, content);
    // 開いているファイルが書き換わったときだけ作り直す (取り消しの履歴を消さないため)
    if (files.get(current) !== view.state.doc.toString()) view.setState(stateFor(files.get(current)));
    renderTabs();
  }),
);

buildButton.addEventListener("click", () =>
  busy(async () => {
    files.set(current, view.state.doc.toString());
    const b = await build(files.get(current));
    show(b.output ?? "", b.stderr);
  }),
);

renderTabs();
runButton.disabled = buildButton.disabled = false;
