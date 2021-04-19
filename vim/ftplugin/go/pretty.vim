" Copyright 2021 momaek. All rights reserved.
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.
"
" This filetype plugin add a new commands for go buffers:
"
"   :PrettyTag
"
"       Run formattag for the current Go file.
"
if exists("b:did_ftplugin_go_pretty_tag")
    finish
endif

if !executable("formattag")
    finish
endif

command! -buffer PrettyTag call s:GoPrettyTag()

function! s:GoPrettyTag() abort
    cexpr system('formattag -file ' . shellescape(expand('%')))
endfunction

let b:did_ftplugin_go_pretty_tag = 1
" vim:ts=4:sw=4:et
