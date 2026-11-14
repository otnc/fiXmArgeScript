" 語の表 (internal/lexicon) から internal/grammar が作る。手で直さず FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
scriptencoding utf-8

" ブロックを開く語・閉じる語・区切る語の正規表現。
function! fixmarge#words#Block() abort
  return {
  \ 'opener': '\C\<\%(persemon\|teitance\|azagara\|sentky\|curoi\|weho\|mof\)\>',
  \ 'continued': '\C\<\%(siva\)\>\s\+\C\<\%(mof\)\>',
  \ 'closer': '\C\<\%(finend\)\>',
  \ 'dedent': '\C\<\%(youlugive\|ennote\|finend\|siva\)\>',
  \ }
endfunction

" 行の頭に来たら1段浅く書く語。打ち終えたら字下げし直す (indentkeys)。
function! fixmarge#words#Dedent() abort
  return ['finend', 'siva', 'youlugive', 'ennote']
endfunction
