if exists('b:did_indent')
  finish
endif
let b:did_indent = 1

setlocal indentexpr=fxma#indent#Expr() autoindent
let &l:indentkeys = '0),0],0},!^F,o,O,e,' . join(map(fxma#words#Dedent(&filetype), '"=" . v:val'), ',')
let b:undo_indent = 'setlocal indentexpr< autoindent< indentkeys<'
