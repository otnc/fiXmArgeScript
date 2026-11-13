" 語の表 (internal/lexicon) から internal/grammar が作る。手で直さず FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
scriptencoding utf-8
if exists('b:current_syntax')
  finish
endif

syntax iskeyword @,48-57,_

syntax match fixmargeOperator /\*\*\|\/\/\|<<\|>>\|<=\|>=\|\.\.=\=\|[-+*\/%<>&|^~]/
syntax match fixmargeNumber /\<\d\%(_\=\d\)*\%(\.\d\%(_\=\d\)*\)\=\%([eE][+-]\=\d\+\)\=\>/
syntax match fixmargeCall /\<\h\w*\ze\s*(/
syntax match fixmargeNamespace /\c\<\%(izanavited\|destimay\|ieviljio\|varizama\|coword\|memoid\|gemon\|keyga\|meme\|toim\)\>/
syntax match fixmargeBuiltin /\c\<\%(fineclowtKuwpace\|izanavitedExi\|madrazySanand\|meifKuwpace\|batteKolsu\|foreSenaji\|izuDefutai\|izuKuwpace\|laineiAlte\|fineclowt\|senchoose\|diviljio\|eatahguy\|pyusense\|recanrey\|sodaKast\|commity\|confsoy\|defutai\|kuwpace\|lighluw\|madrazy\|saimoku\|shojika\|sodaGun\|tsuyong\|vaniell\|ieywoo\|lainei\|matill\|mource\|nozuki\|racdam\|sagain\|salook\|sanand\|saymal\|yoweak\|chidd\|gokuu\|kolsu\|litto\|muyon\|naddy\|onkhe\|tarip\|zelln\|gdie\|gomi\|kast\|nang\|orei\|surc\|ado\|eda\|edo\|gun\|hir\|lea\)\>/
syntax match fixmargeKeyword /\C\<\%(accuhearts\|zanfession\|youlugive\|canosiri\|persemon\|teitance\|azagara\|praylee\|shoubek\|ennote\|finend\|ketzen\|mekamo\|redogh\|sentky\|itswa\|werer\|aste\|keip\|siva\|mof\)\>/
syntax match fixmargeModifier /\C\<\%(abiryoku\|showence\|itsver\|oqun\)\>/
syntax match fixmargeOperatorWord /\C\<\%(senaji\|batte\|sond\|exi\|izu\|nii\|owk\|no\)\>/
syntax match fixmargeConstant /\C\<\%(tadact\|ulieso\|vanita\)\>/
syntax match fixmargeDeclare /\C\<\%(curoi\|weho\)\>/
syntax match fixmargeFunctionName /\%(\C\<\%(curoi\|weho\)\>\s\+\)\@<=\h\w*/
syntax match fixmargeComment /#.*$/ contains=@Spell
syntax match fixmargeEscape /\\\%(u{\x\+}\|.\)/ contained
syntax region fixmargeInterpolation matchgroup=fixmargeInterpolationDelimiter start=/#{/ end=/}/ contained contains=TOP
syntax region fixmargeString matchgroup=fixmargeStringDelimiter start=/"/ skip=/\\./ end=/"/ contains=fixmargeEscape,fixmargeInterpolation,@Spell
syntax region fixmargeString matchgroup=fixmargeStringDelimiter start=/"""/ end=/"""/ contains=fixmargeEscape,fixmargeInterpolation,@Spell

highlight default link fixmargeOperator Operator
highlight default link fixmargeNumber Number
highlight default link fixmargeCall Function
highlight default link fixmargeNamespace Type
highlight default link fixmargeBuiltin Function
highlight default link fixmargeKeyword Keyword
highlight default link fixmargeModifier StorageClass
highlight default link fixmargeOperatorWord Operator
highlight default link fixmargeConstant Constant
highlight default link fixmargeDeclare Keyword
highlight default link fixmargeFunctionName Function
highlight default link fixmargeComment Comment
highlight default link fixmargeEscape SpecialChar
highlight default link fixmargeInterpolationDelimiter Delimiter
highlight default link fixmargeStringDelimiter String
highlight default link fixmargeString String

let b:current_syntax = 'fixmarge'
