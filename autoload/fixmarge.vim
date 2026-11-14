" ファイルタイプ fixmarge (下書き) と fixmargec (清書) のバッファの設定とインデント。plugin/fixmarge.vim の FileType から呼ぶ。
" ftplugin と indent のフォルダに分けないのは、リポジトリの直下のフォルダを減らすため。
scriptencoding utf-8

function! fixmarge#Setup() abort
  setlocal comments=:# commentstring=#\ %s
  setlocal formatoptions-=t formatoptions+=croql
  " gf で werer の "./lib/util" を開けるように。拡張子は読み込む側にそろう (清書は清書を読む)。
  let &l:suffixesadd = '.' . &filetype
  let b:undo_ftplugin = 'setlocal comments< commentstring< formatoptions< suffixesadd<'

  " fixmarge beaukshi と同じインデント (空白2つ)。
  if get(g:, 'fixmarge_recommended_style', 1)
    setlocal expandtab shiftwidth=2 softtabstop=2
    let b:undo_ftplugin .= ' | setlocal expandtab< shiftwidth< softtabstop<'
  endif

  if &filetype ==# 'fixmarge'
    setlocal indentexpr=fixmarge#indent#Expr() autoindent
    let &l:indentkeys = '0),0],0},!^F,o,O,e,' . join(map(fixmarge#words#Dedent(), '"=" . v:val'), ',')
  else
    " 清書はインデントしない (fixmarge sinruits が行頭の空白を消す)。前の行のインデントを引き継ぐ autoindent だけ残す
    setlocal indentexpr= autoindent
  endif
  let b:undo_indent = 'setlocal indentexpr< autoindent< indentkeys<'
endfunction
