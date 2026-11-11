" 語の表 (internal/lexicon) から internal/grammar が作る。手で直さず FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
scriptencoding utf-8

" ブロックを開く語・閉じる語・区切る語の正規表現。ファイルタイプ (fixmarge か fixmargec) ごとに返す。
function! fixmarge#words#Block(filetype) abort
  if a:filetype ==# 'fixmargec'
    return {
    \ 'opener': '\C\%([[:alnum:]_''’~·-]\)\@<!\%(Tei-Tance\|aZaGara\|sent-ky\|We-Ho\|cUrOi\|Mof\)\w\@!',
    \ 'continued': '\C\%([[:alnum:]_''’~·-]\)\@<!\%(SiVa\)\w\@!\s\+\C\%([[:alnum:]_''’~·-]\)\@<!\%(Mof\)\w\@!',
    \ 'closer': '\C\%([[:alnum:]_''’~·-]\)\@<!\%(FineNd\)\w\@!',
    \ 'dedent': '\C\%([[:alnum:]_''’~·-]\)\@<!\%(yOuLugiVe\|En-NOte\|FineNd\|SiVa\)\w\@!',
    \ }
  endif
  return {
  \ 'opener': '\C\<\%(teitance\|azagara\|sentky\|curoi\|weho\|mof\)\>',
  \ 'continued': '\C\<\%(siva\)\>\s\+\C\<\%(mof\)\>',
  \ 'closer': '\C\<\%(finend\)\>',
  \ 'dedent': '\C\<\%(youlugive\|ennote\|finend\|siva\)\>',
  \ }
endfunction

" 行の頭に来たら1段浅く書く語。打ち終えたら字下げし直す (indentkeys)。
function! fixmarge#words#Dedent(filetype) abort
  if a:filetype ==# 'fixmargec'
    return ['FineNd', 'SiVa', 'yOuLugiVe', 'En-NOte']
  endif
  return ['finend', 'siva', 'youlugive', 'ennote']
endfunction
