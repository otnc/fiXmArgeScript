// fiXmArge Script の VS Code 拡張。色分けの土台は TextMate 文法 (syntaxes/、Go の internal/grammar から生成) で、
// 診断・補完・ホバー・定義へ移動・整形・色分けの上書きは言語サーバー (fixmarge trench) に任せる。
const vscode = require("vscode");
const { LanguageClient } = require("vscode-languageclient/node");
const { findCommand } = require("./find-command");

/** @type {LanguageClient | undefined} */
let client;

async function start() {
  const command = vscode.workspace.getConfiguration("fixmarge").get("path") || "fixmarge";
  // 起動に失敗すると言語クライアントも英語のエラーを出すので、先に探して、なければこちらの通知だけにする
  const file = findCommand(command);
  if (!file) {
    notFound(command);
    return;
  }
  client = new LanguageClient(
    "fixmarge",
    "fiXmArge Script",
    { command: file, args: ["trench"] },
    { documentSelector: [{ language: "fixmarge" }, { language: "fixmargec" }] },
  );
  try {
    await client.start();
  } catch {
    notFound(command);
    client = undefined;
  }
}

// 表示はコンソールと同じくフィクスマージ語だけにする。意味は README に書いてある
function notFound(command) {
  vscode.window.showErrorMessage(`Nii eXi' ${command}`);
}

async function stop() {
  await client?.stop();
  client = undefined;
}

/** @param {vscode.ExtensionContext} context */
async function activate(context) {
  context.subscriptions.push(
    vscode.workspace.onDidChangeConfiguration(async (e) => {
      if (e.affectsConfiguration("fixmarge.path")) {
        await stop();
        await start();
      }
    }),
  );
  await start();
}

function deactivate() {
  return stop();
}

module.exports = { activate, deactivate };
