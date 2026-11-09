" 下書きと清書は語の書き方が違うので、別のファイルタイプにする (VS Code 拡張と同じ)。
autocmd BufRead,BufNewFile *.fxma setfiletype fxma
autocmd BufRead,BufNewFile *.fxmac setfiletype fxmac
" パッケージのファイル。Vim は *.ado を Stata として見分けるので、fxma.ado だけは setfiletype ではなく上書きする。
autocmd BufRead,BufNewFile fxma.ado set filetype=fxmaado
