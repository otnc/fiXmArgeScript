// 言語サーバーとして動かす fixmarge の実行ファイルを探す。
const fs = require("node:fs");
const path = require("node:path");
const os = require("node:os");

/**
 * installDirs は、セットアップ・install.sh・go install が fixmarge を置く場所。
 * どれも PATH に足すが、入れる前から動いている VS Code は古い PATH のままなので、ここも探す。
 * @param {NodeJS.Platform} platform
 * @param {NodeJS.ProcessEnv} env
 * @param {string} home
 */
function installDirs(platform, env, home) {
  const p = platform === "win32" ? path.win32 : path.posix;
  const goBin = env.GOBIN || p.join(env.GOPATH || p.join(home, "go"), "bin");
  if (platform === "win32") {
    return [
      env.LOCALAPPDATA && p.join(env.LOCALAPPDATA, "Programs", "fixmarge"),
      env.ProgramFiles && p.join(env.ProgramFiles, "fixmarge"),
      goBin,
    ].filter(Boolean);
  }
  return [p.join(env.FIXMARGE_INSTALL || p.join(home, ".fixmarge"), "bin"), goBin];
}

/**
 * findCommand は command (設定の fixmarge.path) を実行ファイルのパスにする。見つからなければ undefined。
 * パスを含む書き方ならそのファイルだけを見て、名前だけなら PATH とインストール先を探す。
 * @param {string} command
 * @param {{ platform?: NodeJS.Platform, env?: NodeJS.ProcessEnv, home?: string, isFile?: (file: string) => boolean }} [options]
 * @returns {string | undefined}
 */
function findCommand(command, options = {}) {
  const { platform = process.platform, env = process.env, home = os.homedir(), isFile = isRegularFile } = options;
  const p = platform === "win32" ? path.win32 : path.posix;
  // Windows では拡張子を省いた名前でも .exe を探す (spawn は shell を通さないので .cmd などは動かせない)
  const names = platform === "win32" && !p.extname(command) ? [command + ".exe", command] : [command];
  if (command.includes("/") || (platform === "win32" && command.includes("\\"))) {
    return names.find(isFile);
  }
  const dirs = [...(env.PATH ?? env.Path ?? "").split(p.delimiter).filter(Boolean), ...installDirs(platform, env, home)];
  for (const dir of dirs) {
    for (const name of names) {
      const file = p.join(dir, name);
      if (isFile(file)) return file;
    }
  }
  return undefined;
}

/** @param {string} file */
function isRegularFile(file) {
  try {
    return fs.statSync(file).isFile();
  } catch {
    return false;
  }
}

module.exports = { findCommand };
