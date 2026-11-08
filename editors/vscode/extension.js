// fiXmArge Script の VS Code 拡張。色分けの土台は TextMate 文法 (syntaxes/、Go の internal/grammar から生成) で、
// 診断・補完・ホバー・定義へ移動・整形・色分けの上書きは言語サーバー (fxma trench) に任せる。
const vscode = require("vscode");
const { LanguageClient } = require("vscode-languageclient/node");

/** @type {LanguageClient | undefined} */
let client;

async function start() {
  const command = vscode.workspace.getConfiguration("fxma").get("path") || "fxma";
  client = new LanguageClient(
    "fxma",
    "fiXmArge Script",
    { command, args: ["trench"] },
    { documentSelector: [{ language: "fxma" }, { language: "fxmac" }] },
  );
  try {
    await client.start();
  } catch {
    // 表示はコンソールと同じくフィクスマージ語だけにする。意味は README に書いてある
    vscode.window.showErrorMessage(`Nii eXi' ${command}`);
    client = undefined;
  }
}

async function stop() {
  await client?.stop();
  client = undefined;
}

/** @param {vscode.ExtensionContext} context */
async function activate(context) {
  context.subscriptions.push(
    vscode.workspace.onDidChangeConfiguration(async (e) => {
      if (e.affectsConfiguration("fxma.path")) {
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
