" 言語サーバー fixmarge trench の起動に使う値と、Vim の LSP クライアントへの登録。
" Neovim は lsp/fixmarge.lua と plugin/fixmarge.vim でつなぐ (ここの Command と Root も使う)。
scriptencoding utf-8

" fixmarge#lsp#Command は fixmarge の実行ファイルの場所。見つからなければ ''。
" 名前だけなら PATH のほかに、セットアップ・install.sh・go install の入れる場所も探す (VS Code 拡張と同じ)。
" それらは PATH に足されるが、入れる前から動いている Vim は古い PATH のままなので。
function! fixmarge#lsp#Command() abort
  let command = get(g:, 'fixmarge_path', 'fixmarge')
  if command =~# '[/\\]'
    " ~ や $HOME を書けるように展開する
    let command = expand(command)
    return executable(command) ? exepath(command) : ''
  endif
  let found = exepath(command)
  if found !=# ''
    return found
  endif
  let name = has('win32') && fnamemodify(command, ':e') ==# '' ? command . '.exe' : command
  for dir in s:InstallDirs()
    let file = dir . '/' . name
    if executable(file)
      return file
    endif
  endfor
  return ''
endfunction

function! s:InstallDirs() abort
  let home = expand('~')
  let gopath = empty($GOPATH) ? home . '/go' : split($GOPATH, has('win32') ? ';' : ':')[0]
  let gobin = empty($GOBIN) ? gopath . '/bin' : $GOBIN
  if has('win32')
    return filter([
          \ empty($LOCALAPPDATA) ? '' : $LOCALAPPDATA . '/Programs/fixmarge',
          \ empty($ProgramFiles) ? '' : $ProgramFiles . '/fixmarge',
          \ gobin], 'v:val !=# ""')
  endif
  return [(empty($FIXMARGE_INSTALL) ? home . '/.fixmarge' : $FIXMARGE_INSTALL) . '/bin', gobin]
endfunction

" fixmarge#lsp#Root はパッケージの根 (fixmarge.mod のあるフォルダ)。なければ Git の根、それもなければファイルのフォルダ。
function! fixmarge#lsp#Root(...) abort
  let dir = fnamemodify(a:0 ? a:1 : expand('%:p'), ':p:h')
  let mod = findfile('fixmarge.mod', fnameescape(dir) . ';')
  if mod !=# ''
    return fnamemodify(mod, ':p:h')
  endif
  " .git はふつうフォルダだが、git worktree ではファイルになる
  let git = finddir('.git', fnameescape(dir) . ';')
  if git !=# ''
    return fnamemodify(git, ':p:h:h')
  endif
  let git = findfile('.git', fnameescape(dir) . ';')
  return git !=# '' ? fnamemodify(git, ':p:h') : dir
endfunction

" fixmarge#lsp#NotFound は fixmarge が見つからないことを1度だけ知らせる。表示はコンソールと同じくフィクスマージ語だけ。
function! fixmarge#lsp#NotFound() abort
  if get(s:, 'warned', 0)
    return
  endif
  let s:warned = 1
  echohl WarningMsg
  echomsg "Nii eXi' " . get(g:, 'fixmarge_path', 'fixmarge')
  echohl None
endfunction

let s:filetypes = ['fixmarge', 'fixmargec']

" fixmarge#lsp#Client は登録した Vim のクライアントの名前。まだなら ''。
function! fixmarge#lsp#Client() abort
  return get(s:, 'client', '')
endfunction

" fixmarge#lsp#Register は、入っているクライアントのうち最初に見つかった1つに登録する。
" 2つに登録すると、同じ診断が二重に出るため。client を渡すと、そのクライアントだけを試す (vim-lsp の lsp_setup から呼ぶ)。
function! fixmarge#lsp#Register(...) abort
  if get(s:, 'client', '') !=# '' || !get(g:, 'fixmarge_lsp', 1)
    return
  endif
  let clients = a:0 ? [a:1] : ['yegappan', 'vimlsp', 'coc', 'ale']
  for client in clients
    if s:Available(client)
      let command = fixmarge#lsp#Command()
      if command ==# ''
        " 見つからない知らせは、fixmarge のファイルを開いたときに出す (plugin/fixmarge.vim)
        let s:client = 'none'
        return
      endif
      call s:register[client](command)
      let s:client = client
      return
    endif
  endfor
endfunction

function! s:Available(client) abort
  if a:client ==# 'yegappan'
    return exists('*LspAddServer')
  elseif a:client ==# 'vimlsp'
    return exists('g:lsp_loaded')
  elseif a:client ==# 'coc'
    return exists('g:did_coc_loaded')
  elseif a:client ==# 'ale'
    return exists(':ALEInfo') == 2
  endif
  return 0
endfunction

let s:register = {}

function! s:register.yegappan(command) abort
  call LspAddServer([{
        \ 'name': 'fixmarge',
        \ 'filetype': s:filetypes,
        \ 'path': a:command,
        \ 'args': ['trench'],
        \ 'rootSearch': ['fixmarge.mod', '.git'],
        \ }])
endfunction

function! s:register.vimlsp(command) abort
  call lsp#register_server({
        \ 'name': 'fixmarge',
        \ 'cmd': [a:command, 'trench'],
        \ 'allowlist': s:filetypes,
        \ 'root_uri': {server_info -> lsp#utils#path_to_uri(fixmarge#lsp#Root())},
        \ })
endfunction

function! s:register.coc(command) abort
  call coc#config('languageserver.fixmarge', {
        \ 'command': a:command,
        \ 'args': ['trench'],
        \ 'filetypes': s:filetypes,
        \ 'rootPatterns': ['fixmarge.mod', '.git'],
        \ })
endfunction

function! s:register.ale(command) abort
  for filetype in s:filetypes
    call ale#linter#Define(filetype, {
          \ 'name': 'fixmarge',
          \ 'lsp': 'stdio',
          \ 'executable': a:command,
          \ 'command': '%e trench',
          \ 'project_root': {buffer -> fixmarge#lsp#Root(bufname(buffer))},
          \ })
  endfor
endfunction
