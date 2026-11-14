" 清書は字下げしない (fixmarge sinruits が行頭の空白を消す)。前の行の字下げを引き継ぐ autoindent だけ残す
if exists('b:did_indent')
  finish
endif
let b:did_indent = 1

setlocal indentexpr= autoindent
let b:undo_indent = 'setlocal indentexpr< autoindent<'
