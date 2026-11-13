" 語の表 (internal/lexicon) から internal/grammar が作る。手で直さず FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
scriptencoding utf-8
if exists('b:current_syntax')
  finish
endif

syntax match fixmargemodField /^\s*\zs\c\<\%(youlugive\|foryulvu\|mokpose\|lainei\|neyo\|ado\)\>/
syntax region fixmargemodString start=/"/ skip=/\\./ end=/"/ oneline

highlight default link fixmargemodField Keyword
highlight default link fixmargemodString String

let b:current_syntax = 'fixmargemod'
