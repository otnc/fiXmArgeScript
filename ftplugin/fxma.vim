if exists('b:did_ftplugin')
  finish
endif
let b:did_ftplugin = 1

setlocal comments=:# commentstring=#\ %s
setlocal formatoptions-=t formatoptions+=croql
" gf で werer の "./lib/util" を開けるように。拡張子は読み込む側にそろう (清書は清書を読む)。
let &l:suffixesadd = '.' . &filetype
let b:undo_ftplugin = 'setlocal comments< commentstring< formatoptions< suffixesadd<'

" fxma beaukshi と同じ字下げ (空白2つ)。
if get(g:, 'fxma_recommended_style', 1)
  setlocal expandtab shiftwidth=2 softtabstop=2
  let b:undo_ftplugin .= ' | setlocal expandtab< shiftwidth< softtabstop<'
endif
