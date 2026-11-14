import assert from "node:assert/strict";
import { test } from "node:test";

import { build, run, words } from "../index.js";

test("run", async () => {
  const r = await run('name izu zanfession\n"weRer, " + name redogh\n', { stdin: "fxma\n" });
  assert.deepEqual([r.stdout, r.stderr, r.exitCode], ["weRer, fxma\n", "", 0]);
});

test("run reports errors", async () => {
  const r = await run("cnt redogh\n");
  assert.equal(r.stderr, "Nii eXi' main.fxma:1:1 cnt\n");
  assert.equal(r.exitCode, 1);
});

test("files are shared by werer and memoid", async () => {
  const src = '"./util" no double werer\n":memoid" aste memoid werer\nmemoid no surc("out.txt", "#{double(21)}")\n';
  const r = await run(src, { files: { "util.fxma": "oqun weho double(n)\n  n * 2 itswa\nfinend\n" } });
  assert.equal(r.stderr, "");
  assert.equal(r.files["out.txt"], "42");
});

test("exit code and args", async () => {
  const src = '":ieviljio" aste p werer\ndv izu p no diviljio()\ndv no "conwa" redogh\ndv no "wid" redogh\ndv no "shid"(3)\n';
  const r = await run(src, { args: ["a"] });
  assert.deepEqual([r.stdout, r.exitCode], ['["a"]\nnode\n', 3]);
});

test("maxSteps stops an endless loop", async () => {
  const r = await run("tadact azagara\nfinend\n", { maxSteps: 1000 });
  assert.deepEqual([r.stderr, r.exitCode], ["Bavellabion main.fxma:1:1 LAikIll\n", 1]);
});

test("build", async () => {
  assert.deepEqual(await build("x izu tadact\nx redogh\n"), { output: "x Iz'u TadACt\nx ReDogh'\n", stderr: "" });
  const bad = await build('x izu "abc\n');
  assert.equal(bad.output, null);
});

test("words", async () => {
  const izu = (await words()).find((w) => w.name === "izu");
  assert.deepEqual(izu, { name: "izu", spelling: "Iz'u", uses: ["keyword"] });
});
