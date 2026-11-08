-- 言語サーバー fxma trench の設定 (Neovim 0.11 から)。vim.lsp.enable('fxma') でつなぐ。
---@type vim.lsp.Config
return {
  cmd = { 'fxma', 'trench' },
  filetypes = { 'fxma', 'fxmac' },
  root_markers = { 'fxma.mod', '.git' },
}
