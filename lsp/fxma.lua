-- 言語サーバー fixmarge trench の設定 (Neovim 0.11 から)。vim.lsp.enable('fxma') でつなぐ。
---@type vim.lsp.Config
return {
  cmd = { 'fixmarge', 'trench' },
  filetypes = { 'fxma', 'fxmac' },
  root_markers = { 'fixmarge.mod', '.git' },
}
