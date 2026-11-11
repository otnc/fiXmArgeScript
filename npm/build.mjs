// fixmarge.wasm を作り、Go 付属の wasm_exec.js とライセンスをこのフォルダに写す。どちらもコミットしない (.gitignore)。
import { execFileSync } from "node:child_process";
import { copyFileSync } from "node:fs";
import { join } from "node:path";

const here = import.meta.dirname;
const root = join(here, "..");

execFileSync("go", ["build", "-trimpath", "-ldflags=-s -w", "-o", join(here, "fixmarge.wasm"), "./cmd/fixmarge-wasm"], {
  cwd: root,
  env: { ...process.env, GOOS: "js", GOARCH: "wasm" },
  stdio: "inherit",
});
const goroot = execFileSync("go", ["env", "GOROOT"], { encoding: "utf8" }).trim();
copyFileSync(join(goroot, "lib", "wasm", "wasm_exec.js"), join(here, "wasm_exec.js"));
for (const name of ["LICENSE-MIT", "LICENSE-APACHE"]) {
  copyFileSync(join(root, name), join(here, name));
}
