" 語の表 (internal/lexicon) から internal/grammar が作る。手で直さず FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
scriptencoding utf-8
if exists('b:current_syntax')
  finish
endif

syntax match fxmaadoField /^\s*\zs\c\<\%(youlugive\|foryulvu\|mokpose\|lainei\|unneyo\|neyo\|ado\)\>/
syntax region fxmaadoString start=/"/ skip=/\\./ end=/"/ oneline

highlight default link fxmaadoField Keyword
highlight default link fxmaadoString String

let b:current_syntax = 'fxmaado'
