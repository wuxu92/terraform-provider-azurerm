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
    : elem+
;

elem
//@init {System.err.println(_input.LT(1));} //With predicates, it seems debugging the grammar helps where is normally does not.
	:	header
	|	para
	|	quote
	|	list
	|	'\n'
	;

header : '#'+ ~'\n'* '\n' ;

para:
    '\n'* paraContent '\n' (nl|EOF) ; // if \n\n, exists loop. if \n not \n, stays in loop.

paraContent : (text|bold|italics|link|astericks|underscore|'\n')+ ;

bold:	'*' ~('\n'|' ') text'*' ;

astericks :  WS '*' WS ;

underscore : WS '_' WS ;

italics : '_' ~('\n'|' ') text '_' ;

link : '[' text ']' '(' ~')'* ')' ;

quote : quoteElem+ nl ;

quoteElem : '>' ~'\n'* '\n' ;

list:	listElem+ nl nl ;

listElem : (' ' (' ' ' '?)?)? '*' WS paraContent ;

//code: '`' text '`' ;

text: ~('#'|'*'|'>'|'['|']'|'_'|'\n'|'`')+ ;

WS
   : [ \r\n\t]+ -> skip
;

nl	:	'\r'? '\n' ;