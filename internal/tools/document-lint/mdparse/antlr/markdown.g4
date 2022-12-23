grammar markdown;


HEAD_EXAMPLE: 'Example Usage';

HEAD_ARGS: 'Argument'[s]?' Reference';

HEAD_ATTR: 'Attribute'[s]?' Reference';

CODE:
    '`'[a-zA-Z0-9_]+'`'
    ;

DELIMETER:
    '---'
    ;

HEAD_TIMEOUT: 'Timeout'[s]? ;

HEAD_IMPORT: 'Import'[s]? ;

LIST_MARK:
    '*'
    ;

REQUIRED:
    '(Required)'
    ;

OPTIONAL:
    '(Optional)'
    ;

TEXT:
    ~[\r\n]+
    ;

WS:
    [\r\n\t ] -> skip;

head:
    '##' (HEAD_EXAMPLE
    | HEAD_ARGS
    | HEAD_ATTR
    | HEAD_ATTR
    | HEAD_IMPORT)
    ;

list_item:
    LIST_MARK CODE '-' (REQUIRED|OPTIONAL)? TEXT
    ;

delimeter:
    DELIMETER
    ;

text:
    TEXT;

file_:
    (head | list_item | delimeter | text )+;
