" 言語サーバー fixmarge trench に自動でつなぐ。止めるときは、このプラグインが読み込まれる前に g:fixmarge_lsp を 0 にする。
scriptencoding utf-8
if exists('g:loaded_fixmarge')
  finish
endif
let g:loaded_fixmarge = 1

if has('nvim-0.11')
  " 設定は lsp/fixmarge.lua。fixmarge が見つからなければ、そこの root_dir で起動を止めて知らせる
  if get(g:, 'fixmarge_lsp', 1)
    lua vim.lsp.enable('fixmarge')
  endif
  finish
endif

function! s:Start() abort
  if !get(g:, 'fixmarge_lsp', 1)
    return
  endif
  if has('nvim')
    let command = fixmarge#lsp#Command()
    if command ==# ''
      call fixmarge#lsp#NotFound()
      return
    endif
    call luaeval('vim.lsp.start(_A)', {'name': 'fixmarge', 'cmd': [command, 'trench'], 'root_dir': fixmarge#lsp#Root()})
    return
  endif
  " 遅れて読み込まれたとき (ファイルタイプで読み込むプラグインマネージャーなど) は、ここで登録する
  call fixmarge#lsp#Register()
  if fixmarge#lsp#Client() ==# 'none'
    call fixmarge#lsp#NotFound()
  endif
endfunction

augroup fixmarge_lsp
  autocmd!
  autocmd FileType fixmarge,fixmargec call s:Start()
  if !has('nvim')
    " vim-lsp は有効になるときに lsp_setup を出す。ほかのクライアントは、全部のプラグインが読み込まれてから探す
    autocmd User lsp_setup call fixmarge#lsp#Register('vimlsp')
    autocmd VimEnter * call fixmarge#lsp#Register()
  endif
augroup END
