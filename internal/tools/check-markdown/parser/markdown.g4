/** Parser subset of Markdown copied from:
	http://daringfireball.net/projects/markdown/basics
	# Header 1
	## Header 2
	paragraphs are lines of text separated by blank lines. No indent.
	*   Candy.
	*   Gum.
	*   Booze. might
	    be 2 lines
	> This is a blockquote. No markdown allowed inside
	>
	> This is the second paragraph in the blockquote.
	This is an [example link](http://example.com/).
	Some of these words *are emphasized*.
    Some of these words _are emphasized also_.
*/
grammar markdown;

//@header {package org.antlr.md;}

//options {tokenVocab=CharVocab;}

file_
    : (header |list |line )+ EOF
;

header:
    HEAD TEXT ;

list:
    LIST TEXT ;

//para2:
//    line+? '\n';

line:
    ~(LIST|HEAD) TEXT? '\n'
//    LINE
    ;

LINE:
     (WS|OPTION|CODEONE|TEXT)+
    ;


LIST
    : '-'
    | '*'
    | '[]'
    | '[x]'
    ;

HEAD
    : '#'+;
//
//TEXT:
//    [a-zA-Z0-9_]+;

OPTION
    : '(' ('Required'|'Optional') ')' ;

CODEONE
    : '`' [a-zA-Z0-9_ ]+'`';

TEXT
//    :~('#'|'*'|'>'|'['|']'|'_'|'\n')+
    :~[#*>\\[\]_\r\n] ~[#*>\\[\]_\r\n]+
    ;

WS:
    [\t ]
;

LN
: ('\r'?'\n'|'\r')+ -> skip;

END:
EOF -> skip;

//
//LN:
//    [\r\n]{1} -> skip;

//
//para:
//    '\n'* paraContent '\n' (NL|EOF) ; // if \n\n, exists loop. if \n not \n, stays in loop.
//
//paraContent : (TEXT|bold|italics|link|astericks|underscore|'\n')+? ~'#';
//
//bold:	'*' ~('\n'|' ') TEXT'*' ;
//
//astericks :  WS '*' WS ;
//
//underscore : WS '_' WS ;
//
//italics : '_' ~('\n'|' ') TEXT '_' ;
//
//link : '[' TEXT ']' '(' ~')'* ')' ;
//
//quote : quoteElem+ NL ;
//
//quoteElem : '>' ~'\n'* '\n' ;
//
//list:	listElem+ NL NL ~'#' ;
//
////listElem : (' ' (' ' ' '?)?)? '*' WS paraContent ;
//listElem : (' ' (' ' ' '?)?)? '*' WS paraContent ;
//
////code: '`' text '`' ;
//
//TEXT:
////  ~('#'|'*'|'>'|'['|']'|'_'|'\n'|'`')+
//    [a-zA-Z0-9_ '"]+
//;
//
//WS
//   : [ \r\n]+ -> skip
//;
//
//NL:
//    '\r'? '\n'
//;