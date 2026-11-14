" fiXmArge Script のプラグインの入口。ファイルタイプを付け、バッファを設定し、言語サーバー fixmarge trench に自動でつなぐ。
" ftdetect・ftplugin・indent・lsp のフォルダに分けないのは、リポジトリの直下のフォルダを減らすため。
" 言語サーバーへの自動の接続を止めるときは、このプラグインが読み込まれる前に g:fixmarge_lsp を 0 にする。
scriptencoding utf-8
if exists('g:loaded_fixmarge')
  finish
endif
let g:loaded_fixmarge = 1

" 下書きと清書は語の書き方が違うので、別のファイルタイプにする (VS Code 拡張と同じ)。
" fixmarge.mod と fixmarge.sum は、Vim が *.mod を modsim3、*.sum を hercules とするので、決まった名前で上書きする。
if has('nvim')
  lua vim.filetype.add({
        \ extension = { fxma = 'fixmarge', fxmac = 'fixmargec' },
        \ filename = { ['fixmarge.mod'] = 'fixmargemod', ['fixmarge.sum'] = 'fixmargesum' },
        \ })
else
  augroup fixmarge_filetype
    autocmd!
    " 中身が # で始まるファイルは Vim が conf と見なすので、setfiletype ではなく上書きする
    autocmd BufRead,BufNewFile *.fxma set filetype=fixmarge
    autocmd BufRead,BufNewFile *.fxmac set filetype=fixmargec
    autocmd BufRead,BufNewFile fixmarge.mod set filetype=fixmargemod
    autocmd BufRead,BufNewFile fixmarge.sum set filetype=fixmargesum
  augroup END
endif

augroup fixmarge
  autocmd!
  autocmd FileType fixmarge,fixmargec call fixmarge#Setup()
augroup END

if !get(g:, 'fixmarge_lsp', 1)
  finish
endif

if has('nvim-0.11')
  " vim.lsp.config で先に書かれた設定 (cmd を変えるなど) を優先し、書かれていない所だけをここで埋める。
  " root_dir で on_dir を呼ばなければ起動しないので、fixmarge が見つからないときは Neovim の英語のエラーを出さず、こちらの知らせだけにする
  lua vim.lsp.config('fixmarge', vim.tbl_extend('keep', vim.lsp.config.fixmarge or {}, {
        \   cmd = function(dispatchers)
        \     return vim.lsp.rpc.start({ vim.fn['fixmarge#lsp#Command'](), 'trench' }, dispatchers)
        \   end,
        \   filetypes = { 'fixmarge', 'fixmargec' },
        \   root_dir = function(bufnr, on_dir)
        \     if vim.fn['fixmarge#lsp#Command']() == '' then
        \       vim.fn['fixmarge#lsp#NotFound']()
        \       return
        \     end
        \     on_dir(vim.fn['fixmarge#lsp#Root'](vim.api.nvim_buf_get_name(bufnr)))
        \   end,
        \ }))
  lua vim.lsp.enable('fixmarge')
  finish
endif

augroup fixmarge_lsp
  autocmd!
  autocmd FileType fixmarge,fixmargec call fixmarge#lsp#Start()
  if !has('nvim')
    " vim-lsp は有効になるときに lsp_setup を出す。ほかのクライアントは、全部のプラグインが読み込まれてから探す
    autocmd User lsp_setup call fixmarge#lsp#Register('vimlsp')
    autocmd VimEnter * call fixmarge#lsp#Register()
  endif
augroup END
