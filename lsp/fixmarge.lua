-- 言語サーバー fixmarge trench の設定 (Neovim 0.11 から)。vim.lsp.enable('fixmarge') でつなぐ。
---@type vim.lsp.Config
return {
  cmd = { 'fixmarge', 'trench' },
  filetypes = { 'fixmarge', 'fixmargec' },
  root_markers = { 'fixmarge.mod', '.git' },
}
