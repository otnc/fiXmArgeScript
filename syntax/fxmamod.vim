" 語の表 (internal/lexicon) から internal/grammar が作る。手で直さず FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
scriptencoding utf-8
if exists('b:current_syntax')
  finish
endif

syntax match fxmamodField /^\s*\zs\c\<\%(youlugive\|foryulvu\|mokpose\|lainei\|unneyo\|neyo\|ado\)\>/
syntax region fxmamodString start=/"/ skip=/\\./ end=/"/ oneline

highlight default link fxmamodField Keyword
highlight default link fxmamodString String

let b:current_syntax = 'fxmamod'
