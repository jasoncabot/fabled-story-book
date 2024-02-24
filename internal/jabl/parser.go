// Code generated by goyacc -l -o parser.go jabl.y. DO NOT EDIT.
package jabl

import __yyfmt__ "fmt"

type yySymType struct {
	yys        int
	String     string
	Number     float64
	Boolean    bool
	Statement  stmt
	Expression expr
}

const START = 57346
const END = 57347
const IF = 57348
const ELSE = 57349
const GET = 57350
const SET = 57351
const PRINT = 57352
const CHOICE = 57353
const GOTO = 57354
const CMP_LT = 57355
const CMP_GT = 57356
const CMP_LTE = 57357
const CMP_GTE = 57358
const CMP_EQ = 57359
const CMP_NEQ = 57360
const CMP_AND = 57361
const CMP_OR = 57362
const CMP_NOT = 57363
const STRING = 57364
const NUMBER = 57365
const DICE = 57366
const BOOLEAN = 57367

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"START",
	"END",
	"IF",
	"ELSE",
	"GET",
	"SET",
	"PRINT",
	"CHOICE",
	"GOTO",
	"CMP_LT",
	"CMP_GT",
	"CMP_LTE",
	"CMP_GTE",
	"CMP_EQ",
	"CMP_NEQ",
	"CMP_AND",
	"CMP_OR",
	"CMP_NOT",
	"STRING",
	"NUMBER",
	"DICE",
	"BOOLEAN",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'('",
	"')'",
	"','",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

var yyExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 227

var yyAct = [...]int8{
	21, 99, 59, 2, 49, 50, 51, 50, 51, 43,
	44, 45, 46, 47, 48, 58, 17, 37, 38, 16,
	33, 41, 42, 49, 50, 51, 36, 72, 95, 49,
	50, 51, 97, 105, 18, 30, 31, 66, 34, 15,
	14, 39, 13, 74, 75, 77, 78, 79, 80, 81,
	82, 83, 84, 22, 57, 37, 38, 65, 68, 69,
	37, 38, 100, 73, 36, 94, 92, 93, 3, 36,
	32, 26, 25, 85, 40, 64, 5, 96, 1, 90,
	91, 12, 60, 24, 28, 19, 23, 29, 27, 4,
	67, 0, 0, 20, 37, 38, 102, 0, 103, 0,
	95, 49, 50, 51, 104, 72, 0, 86, 87, 88,
	89, 43, 44, 45, 46, 47, 48, 43, 44, 45,
	46, 47, 48, 0, 42, 49, 50, 51, 0, 0,
	0, 49, 50, 51, 43, 44, 45, 46, 47, 48,
	0, 0, 55, 56, 53, 54, 0, 0, 0, 50,
	51, 52, 55, 56, 53, 54, 71, 95, 49, 50,
	51, 52, 101, 37, 38, 0, 63, 26, 25, 0,
	37, 38, 36, 55, 56, 53, 54, 0, 62, 36,
	0, 0, 23, 29, 98, 37, 38, 37, 38, 76,
	0, 37, 38, 0, 36, 0, 36, 0, 0, 70,
	36, 61, 0, 0, 0, 35, 55, 56, 53, 54,
	0, 0, 11, 9, 0, 52, 10, 6, 8, 7,
	9, 0, 0, 10, 6, 8, 7,
}

var yyPact = [...]int16{
	64, -1000, -1000, 214, 207, -1000, 12, 10, 9, -11,
	-14, -1000, -1000, 63, 63, 63, 63, 63, 174, -1000,
	63, 98, 189, 30, -1000, -15, -28, -1000, 63, -1000,
	170, 146, 135, 38, 43, -1000, 63, 63, 63, 168,
	125, -4, 63, 159, 159, 159, 159, 159, 159, 159,
	159, 159, 63, 63, 63, 63, 63, -1000, 63, 63,
	-1000, -1000, 64, 64, 159, 77, 104, 156, -1000, -1000,
	-1000, -1000, -1000, 77, 121, -23, 159, -23, -23, -23,
	-23, -23, -21, -1000, -1000, 77, -1000, -1000, -1000, -1000,
	0, 153, -30, 55, 131, 159, 74, 159, -1000, -1000,
	64, -1000, -21, 2, -1000, -1000,
}

var yyPgo = [...]int8{
	0, 76, 3, 89, 20, 53, 0, 83, 78,
}

var yyR1 = [...]int8{
	0, 8, 2, 3, 3, 1, 1, 1, 1, 1,
	1, 4, 4, 4, 4, 4, 4, 4, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 7, 7,
}

var yyR2 = [...]int8{
	0, 1, 3, 2, 1, 4, 4, 6, 5, 7,
	6, 1, 3, 3, 3, 3, 3, 3, 1, 3,
	2, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 1, 1, 3, 3, 3, 3, 3,
	6, 4, 1, 2,
}

var yyChk = [...]int16{
	-1000, -8, -2, 4, -3, -1, 10, 12, 11, 6,
	9, 5, -1, 30, 30, 30, 30, 30, -4, 22,
	30, -6, -5, 23, -7, 9, 8, 25, 21, 24,
	-4, -4, -5, -4, -4, 31, 26, 17, 18, -4,
	-5, -6, 26, 13, 14, 15, 16, 17, 18, 27,
	28, 29, 26, 19, 20, 17, 18, 24, 30, 30,
	-5, 31, 32, 31, 32, -4, -6, -5, -4, -4,
	31, 31, 31, -4, -6, -6, 30, -6, -6, -6,
	-6, -6, -6, -6, -6, -4, -5, -5, -5, -5,
	-4, -4, -2, -2, -6, 26, -6, 32, 31, 31,
	7, 31, -6, -6, -2, 31,
}

var yyDef = [...]int8{
	0, -2, 1, 0, 0, 4, 0, 0, 0, 0,
	0, 2, 3, 0, 0, 0, 0, 0, 0, 11,
	0, 0, 0, 33, 34, 0, 0, 18, 0, 42,
	0, 0, 0, 0, 0, 5, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 43, 0, 0,
	20, 6, 0, 0, 0, 13, 14, 16, 31, 32,
	12, 19, 35, 15, 36, 25, 0, 26, 27, 28,
	29, 30, 37, 38, 39, 17, 21, 22, 23, 24,
	0, 0, 0, 8, 0, 0, 0, 0, 41, 7,
	0, 10, 36, 0, 9, 40,
}

var yyTok1 = [...]int8{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	30, 31, 28, 26, 32, 27, 3, 29,
}

var yyTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25,
}

var yyTok3 = [...]int8{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := int(yyPact[state])
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && int(yyChk[int(yyAct[n])]) == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || int(yyExca[i+1]) != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := int(yyExca[i])
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = int(yyTok1[0])
		goto out
	}
	if char < len(yyTok1) {
		token = int(yyTok1[char])
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = int(yyTok2[char-yyPrivate])
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = int(yyTok3[i+0])
		if token == char {
			token = int(yyTok3[i+1])
			goto out
		}
	}

out:
	if token == 0 {
		token = int(yyTok2[1]) /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = int(yyPact[yystate])
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = int(yyAct[yyn])
	if int(yyChk[yyn]) == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = int(yyDef[yystate])
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && int(yyExca[xi+1]) == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = int(yyExca[xi+0])
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = int(yyExca[xi+1])
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = int(yyPact[yyS[yyp].yys]) + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = int(yyAct[yyn]) /* simulate a shift of "error" */
					if int(yyChk[yystate]) == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= int(yyR2[yyn])
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = int(yyR1[yyn])
	yyg := int(yyPgo[yyn])
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = int(yyAct[yyg])
	} else {
		yystate = int(yyAct[yyj])
		if int(yyChk[yystate]) != -yyn {
			yystate = int(yyAct[yyg])
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yylex.(*lexer).ast = &program{body: yyDollar[1].Statement}
		}
	case 2:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Statement = &blockStmt{stmt: yyDollar[2].Statement}
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.Statement = &seqStmt{first: yyDollar[1].Statement, rest: yyDollar[2].Statement}
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Statement = &seqStmt{first: yyDollar[1].Statement, rest: nil}
		}
	case 5:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.Statement = &fnStmt{fn: PRINT, expr: yyDollar[3].Expression}
		}
	case 6:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.Statement = &fnStmt{fn: GOTO, expr: yyDollar[3].Expression}
		}
	case 7:
		yyDollar = yyS[yypt-6 : yypt+1]
		{
			yyVAL.Statement = &fnStmt{fn: CHOICE, expr: yyDollar[3].Expression, block: yyDollar[5].Statement}
		}
	case 8:
		yyDollar = yyS[yypt-5 : yypt+1]
		{
			yyVAL.Statement = &ifStmt{cond: yyDollar[3].Expression, block: yyDollar[5].Statement}
		}
	case 9:
		yyDollar = yyS[yypt-7 : yypt+1]
		{
			yyVAL.Statement = &ifStmt{cond: yyDollar[3].Expression, block: yyDollar[5].Statement, other: yyDollar[7].Statement}
		}
	case 10:
		yyDollar = yyS[yypt-6 : yypt+1]
		{
			yyVAL.Statement = &fnStmt{fn: SET, expr: yyDollar[3].Expression, expr2: yyDollar[5].Expression}
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Expression = yyDollar[1].String
		}
	case 12:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &parenExpr{expr: yyDollar[2].Expression}
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &mathExpr{op: '+', left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &mathExpr{op: '+', left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &mathExpr{op: '+', left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &mathExpr{op: '+', left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &mathExpr{op: '+', left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Expression = yyDollar[1].Boolean
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &parenExpr{expr: yyDollar[2].Expression}
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.Expression = &notExpr{expr: yyDollar[2].Expression}
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_AND, t: BOOLEAN, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_OR, t: BOOLEAN, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_EQ, t: BOOLEAN, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_NEQ, t: BOOLEAN, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_LT, t: NUMBER, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_GT, t: NUMBER, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_LTE, t: NUMBER, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_GTE, t: NUMBER, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_EQ, t: NUMBER, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_NEQ, t: NUMBER, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_EQ, t: STRING, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &cmpExpr{op: CMP_NEQ, t: STRING, left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Expression = yyDollar[1].Number
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Expression = yyDollar[1].Expression
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &parenExpr{expr: yyDollar[2].Expression}
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &mathExpr{op: '+', left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &mathExpr{op: '-', left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &mathExpr{op: '*', left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.Expression = &mathExpr{op: '/', left: yyDollar[1].Expression, right: yyDollar[3].Expression}
		}
	case 40:
		yyDollar = yyS[yypt-6 : yypt+1]
		{
			yyVAL.Expression = &fnStmt{fn: SET, expr: yyDollar[3].Expression, expr2: yyDollar[5].Expression}
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.Expression = &fnStmt{fn: GET, expr: yyDollar[3].Expression}
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.Expression = &rollExpr{num: 1, sides: yyDollar[1].Number}
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.Expression = &rollExpr{num: yyDollar[1].Number, sides: yyDollar[2].Number}
		}
	}
	goto yystack /* stack new state and value */
}
