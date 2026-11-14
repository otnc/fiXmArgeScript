" fixmarge beaukshi と同じ字下げを、前の行だけを見て決める。
" ブロックを開く語が多ければ1段深く、finend が多ければ1段浅くする。区切りと終わりの語で始まる行は1段浅く書く。
function! fixmarge#indent#Expr() abort
  let words = fixmarge#words#Block()
  if s:InString(v:lnum)
    return -1
  endif
  let prev_lnum = prevnonblank(v:lnum - 1)
  if prev_lnum == 0
    return 0
  endif
  let prev = s:Code(getline(prev_lnum))
  let cur = s:Code(getline(v:lnum))

  let ind = indent(prev_lnum)
  " 浅く書いた行 (siva や finend や閉じ括弧で始まる行) の次は、その分を戻してから数える。
  if prev =~# '^\s*\%(' . words.dedent . '\)' || prev =~# '^\s*[)\]}]'
    let ind += shiftwidth()
  endif
  let blocks = s:Count(prev, words.opener) - s:Count(prev, words.continued) - s:Count(prev, words.closer)
  " 括弧は1行でいくつ開いても1段 (beaukshi と同じ)。
  let brackets = s:Count(prev, '[([{]') - s:Count(prev, '[)\]}]')
  let ind += shiftwidth() * (blocks + max([-1, min([1, brackets])]))
  if cur =~# '^\s*\%(' . words.dedent . '\)' || cur =~# '^\s*[)\]}]'
    let ind -= shiftwidth()
  endif
  return max([ind, 0])
endfunction

" s:Code は行から文字列の中身とコメントを除く。語や括弧を数え違えないため。
function! s:Code(line) abort
  let line = substitute(a:line, '"\%(\.\|[^"\]\)*"', '""', 'g')
  return substitute(line, '#.*$', '', '')
endfunction

function! s:Count(line, pattern) abort
  let n = 0
  while match(a:line, a:pattern, 0, n + 1) >= 0
    let n += 1
  endwhile
  return n
endfunction

" s:InString は """ の文字列の中の行 (閉じる行も含む) か。中身の字下げは書いたとおりに残す。
" 構文の色は字下げの途中で古くなるので使わず、前の行までの """ の数で決める。
function! s:InString(lnum) abort
  if search('"""', 'bnW') == 0
    return 0
  endif
  let n = 0
  for line in getline(1, a:lnum - 1)
    let n += count(line, '"""')
  endfor
  return n % 2 == 1
endfunction
