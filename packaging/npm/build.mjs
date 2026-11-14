// fixmarge.wasm を作り、Go 付属の wasm_exec.js とライセンスをこのフォルダに写す。どちらもコミットしない (.gitignore)。
import { execFileSync } from "node:child_process";
import { copyFileSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const here = import.meta.dirname;
const root = join(here, "..", "..");

execFileSync("go", ["build", "-trimpath", "-ldflags=-s -w", "-o", join(here, "fixmarge.wasm"), "./cmd/fixmarge-wasm"], {
  cwd: root,
  env: { ...process.env, GOOS: "js", GOARCH: "wasm" },
  stdio: "inherit",
});
const goroot = execFileSync("go", ["env", "GOROOT"], { encoding: "utf8" }).trim();
// mise などで入れた Go のファイルは読み取り専用で、copyFileSync はその属性まで写すので、次に作り直すときに上書きできなくなる。
// 中身だけを書き、前に写した読み取り専用のファイルは先に消す
rmSync(join(here, "wasm_exec.js"), { force: true });
writeFileSync(join(here, "wasm_exec.js"), readFileSync(join(goroot, "lib", "wasm", "wasm_exec.js")));
for (const name of ["LICENSE-MIT", "LICENSE-APACHE"]) {
  copyFileSync(join(root, name), join(here, name));
}
