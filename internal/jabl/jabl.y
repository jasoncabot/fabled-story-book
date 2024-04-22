%{
package jabl
%}

%union{
String string
Number float64
Boolean bool
Statement stmt
Expression expr 
}

%token START END IF ELSE GET GETN GETB SET PRINT CHOICE GOTO
%token CMP_LT CMP_GT CMP_LTE CMP_GTE CMP_EQ CMP_NEQ CMP_AND CMP_OR CMP_NOT
%token<String> STRING
%token<Number> NUMBER DICE
%token<Boolean> BOOLEAN

%type <Statement> stmt block stmt_list
%type <Expression> strExpr boolExpr numExpr rollExpr

%left '+' '-'
%left '*' '/'
%left CMP_NOT CMP_LT CMP_GT CMP_LTE CMP_GTE CMP_EQ CMP_NEQ CMP_AND CMP_OR

%start program

%%
program:
       block                                              { yylex.(*lexer).ast = &program{body: $1} }
    ;

block:
       START stmt_list END                                { $$ = &blockStmt{stmt: $2} }
     ;

stmt_list:
        stmt_list stmt                                    { $$ = &seqStmt{first: $1, rest: $2} }
      | stmt                                              { $$ = &seqStmt{first: $1, rest: nil} }
      ;

stmt:
        PRINT '(' strExpr ')'                             { $$ = &fnStmt{fn: PRINT, expr: $3} }
      | PRINT '(' numExpr ')'                             { $$ = &fnStmt{fn: PRINT, expr: $3} }
      | PRINT '(' boolExpr ')'                            { $$ = &fnStmt{fn: PRINT, expr: $3} }
      | GOTO '(' strExpr ')'                              { $$ = &fnStmt{fn: GOTO, expr: $3} }
      | CHOICE '(' strExpr ',' block ')'                  { $$ = &fnStmt{fn: CHOICE, expr: $3, block: $5} }
      | IF '(' boolExpr ')' block                         { $$ = &ifStmt{cond: $3, block: $5} }
      | IF '(' boolExpr ')' block ELSE block              { $$ = &ifStmt{cond: $3, block: $5, other: $7} }
      | SET '(' strExpr ',' strExpr ')'                   { $$ = &fnStmt{fn: SET, expr: $3, expr2: $5} }
      | SET '(' strExpr ',' numExpr ')'                   { $$ = &fnStmt{fn: SET, expr: $3, expr2: $5} }
      | SET '(' strExpr ',' boolExpr ')'                  { $$ = &fnStmt{fn: SET, expr: $3, expr2: $5} }
      | '/' '/' STRING                                    { $$ = &commentStmt{comment: $3} }
      ;

strExpr:
  STRING                                                  { $$ = $1 }
      | '(' strExpr ')'                                   { $$ = &parenExpr{expr: $2} }
      | strExpr '+' strExpr                               { $$ = &mathExpr{op: '+', left: $1, right: $3} }
      | strExpr '+' numExpr                               { $$ = &mathExpr{op: '+', left: $1, right: $3} }
      | numExpr '+' strExpr                               { $$ = &mathExpr{op: '+', left: $1, right: $3} }
      | strExpr '+' boolExpr                              { $$ = &mathExpr{op: '+', left: $1, right: $3} }
      | boolExpr '+' strExpr                              { $$ = &mathExpr{op: '+', left: $1, right: $3} }
      | GET '(' strExpr ')'                               { $$ = &fnStmt{fn: GET, expr: $3} }
      | SET '(' strExpr ',' strExpr ')'                   { $$ = &fnStmt{fn: SET, expr: $3, expr2: $5} }
      ;

boolExpr:
        BOOLEAN                                           { $$ = $1 }
      | '(' boolExpr ')'                                  { $$ = &parenExpr{expr: $2} }
      | CMP_NOT boolExpr                                  { $$ = &notExpr{expr: $2} }
      | boolExpr CMP_AND boolExpr                         { $$ = &cmpExpr{op: CMP_AND, t: BOOLEAN, left: $1, right: $3} }
      | boolExpr CMP_OR boolExpr                          { $$ = &cmpExpr{op: CMP_OR, t: BOOLEAN, left: $1, right: $3} }
      | boolExpr CMP_EQ boolExpr                          { $$ = &cmpExpr{op: CMP_EQ, t: BOOLEAN, left: $1, right: $3} }
      | boolExpr CMP_NEQ boolExpr                         { $$ = &cmpExpr{op: CMP_NEQ, t: BOOLEAN, left: $1, right: $3} }
      | numExpr CMP_LT numExpr                            { $$ = &cmpExpr{op: CMP_LT, t: NUMBER, left: $1, right: $3} }
      | numExpr CMP_GT numExpr                            { $$ = &cmpExpr{op: CMP_GT, t: NUMBER, left: $1, right: $3} }
      | numExpr CMP_LTE numExpr                           { $$ = &cmpExpr{op: CMP_LTE, t: NUMBER, left: $1, right: $3} }
      | numExpr CMP_GTE numExpr                           { $$ = &cmpExpr{op: CMP_GTE, t: NUMBER, left: $1, right: $3} }
      | numExpr CMP_EQ numExpr                            { $$ = &cmpExpr{op: CMP_EQ, t: NUMBER, left: $1, right: $3} }
      | numExpr CMP_NEQ numExpr                           { $$ = &cmpExpr{op: CMP_NEQ, t: NUMBER, left: $1, right: $3} }
      | strExpr CMP_EQ strExpr                            { $$ = &cmpExpr{op: CMP_EQ, t: STRING, left: $1, right: $3} }
      | strExpr CMP_NEQ strExpr                           { $$ = &cmpExpr{op: CMP_NEQ, t: STRING, left: $1, right: $3} }
      | GETB '(' strExpr ')'                              { $$ = &fnStmt{fn: GETB, expr: $3} }
      | SET '(' strExpr ',' boolExpr ')'                  { $$ = &fnStmt{fn: SET, expr: $3, expr2: $5} }
     ;

numExpr:
        NUMBER                                            { $$ = $1 }
        | rollExpr                                        { $$ = $1 }
        | '(' numExpr ')'                                 { $$ = &parenExpr{expr: $2} }
        | numExpr '+' numExpr                             { $$ = &mathExpr{op: '+', left: $1, right: $3} }
        | numExpr '-' numExpr                             { $$ = &mathExpr{op: '-', left: $1, right: $3} }
        | numExpr '*' numExpr                             { $$ = &mathExpr{op: '*', left: $1, right: $3} }
        | numExpr '/' numExpr                             { $$ = &mathExpr{op: '/', left: $1, right: $3} }
        | GETN '(' strExpr ')'                            { $$ = &fnStmt{fn: GETN, expr: $3} }
        | SET '(' strExpr ',' numExpr ')'                 { $$ = &fnStmt{fn: SET, expr: $3, expr2: $5} }
      ;

rollExpr:
        DICE                                              { $$ = &rollExpr{num: 1, sides: $1} }
      | NUMBER DICE                                       { $$ = &rollExpr{num: $1, sides: $2} }
      ;

%%