// fiXmArge Script の WebAssembly 版。Node.js でもブラウザでも動く。
// fxma.wasm は最初に呼ばれたときに1回だけ読み込み、後の呼び出しで使い回す。
import "./wasm_exec.js";

let loading;

function load() {
  loading ??= (async () => {
    // Go の main が run などを用意し終えたら、この関数を呼ぶ
    const ready = new Promise((resolve) => {
      globalThis.__fxmaReady = resolve;
    });
    const go = new globalThis.Go();
    const bytes = await readWasm(new URL("./fxma.wasm", import.meta.url));
    const { instance } = await WebAssembly.instantiate(bytes, go.importObject);
    go.run(instance);
    await ready;
    return globalThis.__fxma;
  })();
  return loading;
}

async function readWasm(url) {
  if (url.protocol === "file:") {
    // import 文で書くと、ブラウザ向けにまとめるツール (Vite など) が node:fs を探しに行くので、実行中に取り出す
    return process.getBuiltinModule("node:fs/promises").readFile(url);
  }
  const response = await fetch(url);
  return response.arrayBuffer();
}

export async function run(source, options = {}) {
  return (await load()).run(source, options);
}

export async function build(source) {
  return (await load()).build(source);
}

export async function words() {
  return (await load()).words();
}
