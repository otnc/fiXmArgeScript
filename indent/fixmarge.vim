if exists('b:did_indent')
  finish
endif
let b:did_indent = 1

setlocal indentexpr=fixmarge#indent#Expr() autoindent
let &l:indentkeys = '0),0],0},!^F,o,O,e,' . join(map(fixmarge#words#Dedent(), '"=" . v:val'), ',')
let b:undo_indent = 'setlocal indentexpr< autoindent< indentkeys<'
