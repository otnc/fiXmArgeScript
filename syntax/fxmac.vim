" 語の表 (internal/lexicon) から internal/grammar が作る。手で直さず FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
scriptencoding utf-8
if exists('b:current_syntax')
  finish
endif

syntax iskeyword @,48-57,_

syntax match fxmaOperator /\*\*\|\/\/\|<<\|>>\|<=\|>=\|\.\.=\=\|[-+*\/%<>&|^~]/
syntax match fxmaNumber /\<\d\%(_\=\d\)*\%(\.\d\%(_\=\d\)*\)\=\%([eE][+-]\=\d\+\)\=\>/
syntax match fxmaCall /\<\h\w*\ze\s*(/
syntax match fxmaNamespace /\C\%([[:alnum:]_'’~·-]\)\@<!\%(i·zana-vited\|ievil JIO\|VariZaMa\|dEstiMay\|Coword\|GE-mon\|Memoid\|KeyGa\|MeMe\|Toim\)\w\@!/
syntax match fxmaBuiltin /\C\%([[:alnum:]_'’~·-]\)\@<!\%(i·zana-vitedeXi'\|FinecloWtKuwPace\|madRazysAnAnd\|Iz'uKuwPace\|Iz'udeFuTai\|Lai NeiAlte\|MEifKuwPace\|BatteKolSu\|Sen-choose\|eATah·Guy\|fOreSeNaji\|Divil JIO\|FinecloWt\|PyusenS'e\|Vani·Ell\|Ligh-Luw\|SodaKast\|Sodagu-n\|conf-Soy\|rEcanRey\|Commity\|KuwPace\|Lai Nei\|RacDA-m\|Say mAl\|ShoJikA\|deFuTai\|madRazy\|mu·Yon\|sAiMoku\|tSuYoNg\|Matill\|SaLook\|YowEAk\|ieyWOO\|mOurce\|nOzuki\|sAgAin\|sAnAnd\|ChiDd\|GoKuu\|KolSu\|Litto\|Tarip\|Zelln\|nAddY\|oNkhe\|GoMi\|Kast\|SurC\|gDie\|gu-n\|nAng\|oRei\|Ado\|Eda\|Edo\|Hir\|LEA\)\w\@!/
syntax match fxmaKeyword /\C\%([[:alnum:]_'’~·-]\)\@<!\%(AccUhEarts\|zanfession\|Tei-Tance\|pray·lee\|yOuLugiVe\|CaNoSiri\|En-NOte\|ReDogh'\|ShouBeK\|aZaGara\|sent-ky\|FineNd\|KetZen\|meKaMo\|Kei-P\|itsWA\|weRer\|SiVa\|aSte\|Mof\)\w\@!/
syntax match fxmaModifier /\C\%([[:alnum:]_'’~·-]\)\@<!\%(AbiRyoku\|showence\|itsver\|Oqun\)\w\@!/
syntax match fxmaOperatorWord /\C\%([[:alnum:]_'’~·-]\)\@<!\%(S\~ond'\|SeNaji\|Batte\|Iz'u\|eXi'\|N'o\|Nii\|Owk\)\w\@!/
syntax match fxmaConstant /\C\%([[:alnum:]_'’~·-]\)\@<!\%(Vani·tA\|TadACt\|ULieSo\)\w\@!/
syntax match fxmaDeclare /\C\%([[:alnum:]_'’~·-]\)\@<!\%(We-Ho\|cUrOi\)\w\@!/
syntax match fxmaFunctionName /\%(\C\%([[:alnum:]_'’~·-]\)\@<!\%(We-Ho\|cUrOi\)\w\@!\s\+\)\@<=\h\w*/
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

let b:current_syntax = 'fxmac'
