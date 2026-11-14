" Vim のプラグイン (リポジトリの直下) を確かめる。リポジトリの直下で動かす:
"   FIXMARGE=dist/vimtest/fixmarge vim -u NONE -N -Es -S scripts/check/vim.vim
" 通れば終了コード 0、落ちれば 1 で、落ちた所を標準出力に書く。
let s:root = getcwd()
execute 'set runtimepath^=' . fnameescape(s:root)
filetype plugin indent on
syntax on
runtime plugin/fixmarge.vim

" ファイルタイプと色分けと字下げの設定
for [s:file, s:type] in [['examples/fizzbuzz.fxma', 'fixmarge'], ['x.fxmac', 'fixmargec'], ['fixmarge.mod', 'fixmargemod'], ['fixmarge.sum', 'fixmargesum']]
  execute 'edit ' . s:file
  call assert_equal(s:type, &filetype, s:file)
endfor
edit examples/fizzbuzz.fxma
call assert_equal('fixmarge', get(b:, 'current_syntax', ''), 'syntax')
call assert_equal('fixmarge#indent#Expr()', &indentexpr, 'indentexpr')

" 字下げは fixmarge beaukshi と同じ
enew
setfiletype fixmarge
call setline(1, ['weho f(n)', 'mof n > 1', '"a" redogh', 'siva', '"b" redogh', 'finend', 'n itswa', 'finend'])
normal! gg=G
call assert_equal([0, 2, 4, 2, 4, 2, 2, 0], map(range(1, 8), 'indent(v:val)'), 'indent')
bwipeout!

" fixmarge を探す
let g:fixmarge_path = $FIXMARGE
call assert_equal(exepath($FIXMARGE), fixmarge#lsp#Command(), 'command from g:fixmarge_path')
let g:fixmarge_path = 'fixmarge-not-installed'
call assert_equal('', fixmarge#lsp#Command(), 'command not found')
let g:fixmarge_path = $FIXMARGE

" パッケージの根: fixmarge.mod がなければ Git の根
call assert_equal(s:root, fnamemodify(fixmarge#lsp#Root(s:root . '/examples/lib/stats.fxma'), ':p:h'), 'root')

" 入っているクライアントに登録する (ここでは yegappan/lsp の代わりの関数で受け取る)
let s:servers = []
function! LspAddServer(servers) abort
  call extend(s:servers, a:servers)
endfunction
call fixmarge#lsp#Register()
call assert_equal('yegappan', fixmarge#lsp#Client(), 'client')
call assert_equal(1, len(s:servers), 'registered once')
call assert_equal([exepath($FIXMARGE), ['trench'], ['fixmarge', 'fixmargec']], [s:servers[0].path, s:servers[0].args, s:servers[0].filetype], 'server')
" 2回目は登録しない (診断が二重に出ないように)
call fixmarge#lsp#Register()
call assert_equal(1, len(s:servers), 'not registered twice')

if empty(v:errors)
  call writefile(['PASS vim'], '/dev/stdout')
  qall!
endif
call writefile(v:errors, '/dev/stdout')
cquit!
