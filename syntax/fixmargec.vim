" 語の表 (internal/lexicon) から internal/grammar が作る。手で直さず FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
scriptencoding utf-8
if exists('b:current_syntax')
  finish
endif

syntax iskeyword @,48-57,_

syntax match fixmargeOperator /\*\*\|\/\/\|<<\|>>\|<=\|>=\|\.\.=\=\|[-+*\/%<>&|^~]/
syntax match fixmargeNumber /\<\d\%(_\=\d\)*\%(\.\d\%(_\=\d\)*\)\=\%([eE][+-]\=\d\+\)\=\>/
syntax match fixmargeCall /\<\h\w*\ze\s*(/
syntax match fixmargeNamespace /\C\%([[:alnum:]_'’~·-]\)\@<!\%(i·zana-vited\|ievil JIO\|VariZaMa\|dEstiMay\|Coword\|GE-mon\|Memoid\|KeyGa\|MeMe\|Toim\)\w\@!/
syntax match fixmargeBuiltin /\C\%([[:alnum:]_'’~·-]\)\@<!\%(i·zana-vitedeXi'\|FinecloWtKuwPace\|madRazysAnAnd\|Iz'uKuwPace\|Iz'udeFuTai\|Lai NeiAlte\|MEifKuwPace\|BatteKolSu\|Sen-choose\|eATah·Guy\|fOreSeNaji\|Divil JIO\|FinecloWt\|PyusenS'e\|Vani·Ell\|Ligh-Luw\|SodaKast\|Sodagu-n\|conf-Soy\|rEcanRey\|Commity\|KuwPace\|Lai Nei\|RacDA-m\|Say mAl\|ShoJikA\|deFuTai\|madRazy\|mu·Yon\|sAiMoku\|tSuYoNg\|Matill\|SaLook\|YowEAk\|ieyWOO\|mOurce\|nOzuki\|sAgAin\|sAnAnd\|ChiDd\|GoKuu\|KolSu\|Litto\|Tarip\|Zelln\|nAddY\|oNkhe\|GoMi\|Kast\|SurC\|gDie\|gu-n\|nAng\|oRei\|Ado\|Eda\|Edo\|Hir\|LEA\)\w\@!/
syntax match fixmargeKeyword /\C\%([[:alnum:]_'’~·-]\)\@<!\%(AccUhEarts\|zanfession\|Tei-Tance\|pray·lee\|yOuLugiVe\|CaNoSiri\|En-NOte\|ReDogh'\|ShouBeK\|aZaGara\|sent-ky\|FineNd\|KetZen\|meKaMo\|Kei-P\|itsWA\|weRer\|SiVa\|aSte\|Mof\)\w\@!/
syntax match fixmargeModifier /\C\%([[:alnum:]_'’~·-]\)\@<!\%(AbiRyoku\|showence\|itsver\|Oqun\)\w\@!/
syntax match fixmargeOperatorWord /\C\%([[:alnum:]_'’~·-]\)\@<!\%(S\~ond'\|SeNaji\|Batte\|Iz'u\|eXi'\|N'o\|Nii\|Owk\)\w\@!/
syntax match fixmargeConstant /\C\%([[:alnum:]_'’~·-]\)\@<!\%(Vani·tA\|TadACt\|ULieSo\)\w\@!/
syntax match fixmargeDeclare /\C\%([[:alnum:]_'’~·-]\)\@<!\%(We-Ho\|cUrOi\)\w\@!/
syntax match fixmargeFunctionName /\%(\C\%([[:alnum:]_'’~·-]\)\@<!\%(We-Ho\|cUrOi\)\w\@!\s\+\)\@<=\h\w*/
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

let b:current_syntax = 'fixmargec'
