const test = require("node:test");
const assert = require("node:assert/strict");
const { findCommand } = require("../find-command");

const only = (...files) => (file) => files.includes(file);

test("PATH にあればそれを使う", () => {
  const found = findCommand("fixmarge", {
    platform: "linux",
    env: { PATH: "/usr/bin:/opt/fx/bin" },
    home: "/home/a",
    isFile: only("/opt/fx/bin/fixmarge", "/home/a/.fixmarge/bin/fixmarge"),
  });
  assert.equal(found, "/opt/fx/bin/fixmarge");
});

test("PATH になければ install.sh の場所を探す", () => {
  const found = findCommand("fixmarge", {
    platform: "darwin",
    env: { PATH: "/usr/bin" },
    home: "/Users/a",
    isFile: only("/Users/a/.fixmarge/bin/fixmarge"),
  });
  assert.equal(found, "/Users/a/.fixmarge/bin/fixmarge");
});

test("FIXMARGE_INSTALL があればその下の bin を探す", () => {
  const found = findCommand("fixmarge", {
    platform: "linux",
    env: { PATH: "", FIXMARGE_INSTALL: "/srv/fx" },
    home: "/home/a",
    isFile: only("/srv/fx/bin/fixmarge"),
  });
  assert.equal(found, "/srv/fx/bin/fixmarge");
});

test("Windows ではセットアップの場所の fixmarge.exe を探す", () => {
  const found = findCommand("fixmarge", {
    platform: "win32",
    env: { PATH: "C:\\Windows", LOCALAPPDATA: "C:\\Users\\a\\AppData\\Local" },
    home: "C:\\Users\\a",
    isFile: only("C:\\Users\\a\\AppData\\Local\\Programs\\fixmarge\\fixmarge.exe"),
  });
  assert.equal(found, "C:\\Users\\a\\AppData\\Local\\Programs\\fixmarge\\fixmarge.exe");
});

test("go install の場所も探す", () => {
  const found = findCommand("fixmarge", {
    platform: "linux",
    env: { PATH: "" },
    home: "/home/a",
    isFile: only("/home/a/go/bin/fixmarge"),
  });
  assert.equal(found, "/home/a/go/bin/fixmarge");
});

test("パスを書いたときはそのファイルだけを見る", () => {
  const options = { platform: "linux", env: { PATH: "/usr/bin" }, home: "/home/a", isFile: only("/usr/bin/fixmarge") };
  assert.equal(findCommand("/opt/fixmarge", options), undefined);
  assert.equal(findCommand("/usr/bin/fixmarge", options), "/usr/bin/fixmarge");
});

test("Windows ではパスに .exe を省いてもよい", () => {
  const found = findCommand("D:\\tools\\fixmarge", {
    platform: "win32",
    env: { PATH: "" },
    home: "C:\\Users\\a",
    isFile: only("D:\\tools\\fixmarge.exe"),
  });
  assert.equal(found, "D:\\tools\\fixmarge.exe");
});

test("どこにもなければ undefined", () => {
  const found = findCommand("fixmarge", { platform: "linux", env: { PATH: "/usr/bin" }, home: "/home/a", isFile: () => false });
  assert.equal(found, undefined);
});
