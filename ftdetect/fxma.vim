" 下書きと清書は語の書き方が違うので、別のファイルタイプにする (VS Code 拡張と同じ)。
autocmd BufRead,BufNewFile *.fxma setfiletype fxma
autocmd BufRead,BufNewFile *.fxmac setfiletype fxmac
" パッケージのファイルとロックファイル。Vim は *.mod を modsim3、*.sum を hercules として見分けるので、setfiletype ではなく上書きする。
autocmd BufRead,BufNewFile fixmarge.mod set filetype=fxmamod
autocmd BufRead,BufNewFile fixmarge.sum set filetype=fxmasum
