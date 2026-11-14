-- Neovim のプラグイン (リポジトリの直下) を確かめる。リポジトリの直下で動かす:
--   FIXMARGE=dist/vimtest/fixmarge nvim --headless -u NONE -l scripts/check/nvim.lua
-- 本物の fixmarge trench につなぎ、診断が届くことまで見る。通れば終了コード 0、落ちれば 1。
local root = vim.fn.getcwd()
local errors = {}
local function check(ok, message)
  if not ok then
    table.insert(errors, message)
  end
end

vim.opt.runtimepath:prepend(root)
vim.g.fixmarge_path = os.getenv('FIXMARGE')
vim.cmd('filetype plugin indent on')
vim.cmd('syntax on')
vim.cmd('runtime plugin/fixmarge.vim')

-- 構文の誤りがあると、言語サーバーが AyAmis の診断を出す
local file = vim.fn.tempname() .. '.fxma'
vim.fn.writefile({ 'x izu (1' }, file)
vim.cmd.edit(file)
check(vim.bo.filetype == 'fixmarge', 'filetype: ' .. vim.bo.filetype)
check(vim.bo.indentexpr == 'fixmarge#indent#Expr()', 'indentexpr: ' .. vim.bo.indentexpr)

local attached = vim.wait(10000, function()
  return #vim.lsp.get_clients({ name = 'fixmarge', bufnr = 0 }) > 0
end)
check(attached, 'fixmarge trench did not attach')
local diagnosed = vim.wait(10000, function()
  return #vim.diagnostic.get(0) > 0
end)
check(diagnosed, 'no diagnostics')
if diagnosed then
  local message = vim.diagnostic.get(0)[1].message
  check(message:find('AyAmis', 1, true) ~= nil, 'diagnostic: ' .. message)
end

-- 見つからないときは起動せず、知らせだけを出す
vim.g.fixmarge_path = 'fixmarge-not-installed'
vim.cmd('silent! messages clear')
local other = vim.fn.tempname() .. '.fxma'
vim.fn.writefile({ '1 redogh' }, other)
vim.cmd.edit(other)
vim.wait(500)
check(#vim.lsp.get_clients({ name = 'fixmarge', bufnr = 0 }) == 0, 'started without fixmarge')
check(vim.fn.execute('messages'):find("Nii eXi' fixmarge-not-installed", 1, true) ~= nil, 'not-found message')

if #errors == 0 then
  io.stdout:write('PASS nvim\n')
  os.exit(0)
end
io.stdout:write(table.concat(errors, '\n') .. '\n')
os.exit(1)
