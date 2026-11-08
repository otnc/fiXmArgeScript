" 語の表 (internal/lexicon) から internal/grammar が作る。手で直さず FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
scriptencoding utf-8
if exists('b:current_syntax')
  finish
endif

syntax iskeyword @,48-57,_

syntax match fxmaOperator /\*\*\|\/\/\|<<\|>>\|<=\|>=\|\.\.=\=\|[-+*\/%<>&|^~]/
syntax match fxmaNumber /\<\d\%(_\=\d\)*\%(\.\d\%(_\=\d\)*\)\=\%([eE][+-]\=\d\+\)\=\>/
syntax match fxmaCall /\<\h\w*\ze\s*(/
syntax match fxmaNamespace /\c\<\%(izanavited\|destimay\|ieviljio\|varizama\|coword\|memoid\|gemon\|keyga\|meme\|toim\)\>/
syntax match fxmaBuiltin /\c\<\%(fineclowtKuwpace\|izanavitedExi\|madrazySanand\|meifKuwpace\|batteKolsu\|foreSenaji\|izuDefutai\|izuKuwpace\|laineiAlte\|fineclowt\|senchoose\|diviljio\|eatahguy\|pyusense\|recanrey\|sodaKast\|commity\|confsoy\|defutai\|kuwpace\|lighluw\|madrazy\|saimoku\|shojika\|sodaGun\|tsuyong\|vaniell\|ieywoo\|lainei\|matill\|mource\|nozuki\|racdam\|sagain\|salook\|sanand\|saymal\|yoweak\|chidd\|gokuu\|kolsu\|litto\|muyon\|naddy\|onkhe\|tarip\|zelln\|gdie\|gomi\|kast\|nang\|orei\|surc\|ado\|eda\|edo\|gun\|hir\|lea\)\>/
syntax match fxmaKeyword /\C\<\%(accuhearts\|zanfession\|youlugive\|canosiri\|teitance\|azagara\|praylee\|shoubek\|ennote\|finend\|ketzen\|mekamo\|redogh\|sentky\|itswa\|werer\|aste\|keip\|siva\|mof\)\>/
syntax match fxmaModifier /\C\<\%(abiryoku\|showence\|itsver\|oqun\)\>/
syntax match fxmaOperatorWord /\C\<\%(senaji\|batte\|sond\|exi\|izu\|nii\|owk\|no\)\>/
syntax match fxmaConstant /\C\<\%(tadact\|ulieso\|vanita\)\>/
syntax match fxmaDeclare /\C\<\%(curoi\|weho\)\>/
syntax match fxmaFunctionName /\%(\C\<\%(curoi\|weho\)\>\s\+\)\@<=\h\w*/
syntax match fxmaComment /#.*$/ contains=@Spell
syntax match fxmaEscape /\\\%(u{\x\+}\|.\)/ contained
syntax region fxmaInterpolation matchgroup=fxmaInterpolationDelimiter start=/#{/ end=/}/ contained contains=TOP
syntax region fxmaString matchgroup=fxmaStringDelimiter start=/"/ skip=/\\./ end=/"/ contains=fxmaEscape,fxmaInterpolation,@Spell
syntax region fxmaString matchgroup=fxmaStringDelimiter start=/"""/ end=/"""/ contains=fxmaEscape,fxmaInterpolation,@Spell

highlight default link fxmaOperator Operator
highlight default link fxmaNumber Number
highlight default link fxmaCall Function
highlight default link fxmaNamespace Type
highlight default link fxmaBuiltin Function
highlight default link fxmaKeyword Keyword
highlight default link fxmaModifier StorageClass
highlight default link fxmaOperatorWord Operator
highlight default link fxmaConstant Constant
highlight default link fxmaDeclare Keyword
highlight default link fxmaFunctionName Function
highlight default link fxmaComment Comment
highlight default link fxmaEscape SpecialChar
highlight default link fxmaInterpolationDelimiter Delimiter
highlight default link fxmaStringDelimiter String
highlight default link fxmaString String

let b:current_syntax = 'fxma'
