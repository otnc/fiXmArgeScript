-- 言語サーバー fixmarge trench の設定 (Neovim 0.11 から)。plugin/fixmarge.vim が vim.lsp.enable('fixmarge') でつなぐ。
-- fixmarge の場所は g:fixmarge_path と PATH、インストール先から探す (autoload/fixmarge/lsp.vim)。
---@type vim.lsp.Config
return {
  cmd = function(dispatchers)
    return vim.lsp.rpc.start({ vim.fn['fixmarge#lsp#Command'](), 'trench' }, dispatchers)
  end,
  filetypes = { 'fixmarge', 'fixmargec' },
  -- on_dir を呼ばなければ起動しない。見つからないときに Neovim の英語のエラーを出さず、こちらの知らせだけにする
  root_dir = function(bufnr, on_dir)
    if vim.fn['fixmarge#lsp#Command']() == '' then
      vim.fn['fixmarge#lsp#NotFound']()
      return
    end
    on_dir(vim.fn['fixmarge#lsp#Root'](vim.api.nvim_buf_get_name(bufnr)))
  end,
}
