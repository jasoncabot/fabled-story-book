package jabl

import (
	"fmt"
	"strconv"
	"strings"
)

type Character struct {
	Name string
}

type Result struct {
	Output     string    `json:"output"`
	Choices    []*Choice `json:"choices"`
	Transition SectionId `json:"transition"`
}

type Choice struct {
	Text string `json:"text"`
	Code string `json:"code"`
}

type SectionId string

type Interpreter struct {
	loader SectionLoader
}

type State interface {
	Get(key string) (float64, error)
	Set(key string, value float64) error
}

type SectionLoader interface {
	LoadSection(identifier SectionId, onLoad func(code string, err error))
}

func NewInterpreter(loader SectionLoader) *Interpreter {
	return &Interpreter{
		loader: loader,
	}
}

func (i *Interpreter) Execute(identifier SectionId, state State, callback func(result *Result, err error)) {
	i.loader.LoadSection(identifier, func(code string, err error) {
		if err != nil {
			callback(nil, err)
			return
		}
		i.Evaluate(code, state, callback)
	})
}

func (i *Interpreter) Evaluate(code string, state State, callback func(result *Result, err error)) {
	// parse code
	src := strings.NewReader(code)

	// gimme a good error
	yyErrorVerbose = true

	lexer := newLexer(src)
	parseResult := yyParse(lexer)
	if lexer.err != nil || parseResult != 0 {
		callback(nil, fmt.Errorf("error parsing code: %v", lexer.err))
		return
	}

	result := &Result{}
	if err := i.eval(lexer.ast, state, result); err != nil {
		callback(nil, fmt.Errorf("error evaluating expression: %v", err))
		return
	}

	// Congrats we've evaluated our code
	callback(result, nil)
}

func (i *Interpreter) eval(node any, state State, result *Result) error {
	// We can evaluate a nil statement, it just does nothing, it's not an error
	if node == nil {
		return nil
	}

	switch t := node.(type) {
	case *program:
		return i.eval(t.body, state, result)
	case *blockStmt:
		return i.eval(t.stmt, state, result)
	case *seqStmt:
		if err := i.eval(t.first, state, result); err != nil {
			return err
		}
		if err := i.eval(t.rest, state, result); err != nil {
			return err
		}
		return nil
	case *fnStmt:
		switch t.fn {
		case PRINT:
			expression, err := i.evalStr(t.expr)
			if err != nil {
				return err
			}
			result.Output += expression + "\n"
			return nil
		case GOTO:
			identifier, err := i.evalStr(t.expr)
			if err != nil {
				return err
			}
			result.Transition = SectionId(identifier)
			return nil
		case CHOICE:
			expression, err := i.evalStr(t.expr)
			if err != nil {
				return err
			}
			sb := &strings.Builder{}
			i.printCode(t.block, sb, 0)

			code := sb.String()
			result.Choices = append(result.Choices, &Choice{
				Text: expression,
				Code: code,
			})
			return nil
		case SET:
			key, err := i.evalStr(t.expr)
			if err != nil {
				return err
			}
			val, err := i.evalNum(t.expr2, state)
			if err != nil {
				return err
			}
			return state.Set(key, val)
		default:
			if x, ok := printableFn[t.fn]; ok {
				return fmt.Errorf("%s not supported in this context", x)
			}
			return fmt.Errorf("invalid function token: %d", t.fn)
		}
	case *ifStmt:
		val, err := i.evalBool(t.cond, state)
		if err != nil {
			return err
		}
		if val {
			return i.eval(t.block, state, result)
		} else if t.other != nil {
			return i.eval(t.other, state, result)
		}
		return nil
	default:
		return fmt.Errorf("invalid statement type: %s", t)
	}
}

func addIndent(sb *strings.Builder, indent uint8) {
	sb.WriteRune('\n')
	for j := uint8(0); j < indent; j++ {
		sb.WriteRune('\t')
	}
}

func (i *Interpreter) printCode(stmt any, sb *strings.Builder, indent uint8) {
	switch t := stmt.(type) {
	case *program:
		i.printCode(t.body, sb, indent)
	case *blockStmt:
		sb.WriteRune('{')
		addIndent(sb, indent+1)
		i.printCode(t.stmt, sb, indent+1)
		addIndent(sb, indent)
		sb.WriteRune('}')
	case *seqStmt:
		i.printCode(t.first, sb, indent)
		if t.rest != nil {
			addIndent(sb, indent)
			i.printCode(t.rest, sb, indent)
		}
	case *fnStmt:
		sb.WriteString(printableFn[t.fn])
		sb.WriteRune('(')
		i.printCode(t.expr, sb, indent)
		if t.block != nil {
			sb.WriteRune(',')
			i.printCode(t.block, sb, indent)
		} else if t.expr2 != nil {
			sb.WriteRune(',')
			i.printCode(t.expr2, sb, indent)
		}
		sb.WriteRune(')')
	case *ifStmt:
		sb.WriteString("if (")
		i.printCode(t.cond, sb, indent)
		sb.WriteString(") ")
		i.printCode(t.block, sb, indent)
		if t.other != nil {
			sb.WriteString(" else ")
			i.printCode(t.other, sb, indent)
		}
	case *parenExpr:
		sb.WriteRune('(')
		i.printCode(t.expr, sb, indent)
		sb.WriteRune(')')
	case *notExpr:
		sb.WriteRune('!')
		i.printCode(t.expr, sb, indent)
	case *cmpExpr:
		i.printCode(t.left, sb, indent)
		sb.WriteString(printableOp[t.op])
		i.printCode(t.right, sb, indent)
	case *mathExpr:
		i.printCode(t.left, sb, indent)
		sb.WriteRune(rune(t.op))
		i.printCode(t.right, sb, indent)
	case float64:
		sb.WriteString(strconv.FormatFloat(t, 'f', -1, 64))
	case string:
		sb.WriteRune('"')
		sb.WriteString(t)
		sb.WriteRune('"')
	case bool:
		if t {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	}
}

var (
	printableOp = map[int]string{
		CMP_EQ:  " == ",
		CMP_NEQ: " != ",
		CMP_LT:  " < ",
		CMP_LTE: " <= ",
		CMP_GT:  " > ",
		CMP_GTE: " >= ",
	}

	printableFn = map[int]string{
		PRINT:  "print",
		GOTO:   "goto",
		CHOICE: "choice",
		SET:    "set",
		GET:    "get",
	}
)

func (i *Interpreter) evalBool(e expr, state State) (bool, error) {
	switch t := e.(type) {
	case bool:
		return t, nil
	case *parenExpr:
		return i.evalBool(t.expr, state)
	case *notExpr:
		val, err := i.evalBool(t.expr, state)
		if err != nil {
			return false, err
		}
		return !val, nil
	case *cmpExpr:
		switch t.t {
		case BOOLEAN:
			return i.compareBoolean(state, t.op, t.left, t.right)
		case STRING:
			return i.compareString(state, t.op, t.left, t.right)
		case NUMBER:
			return i.compareNumber(state, t.op, t.left, t.right)
		}
		return false, fmt.Errorf("invalid comparison operator")
	default:
		return false, fmt.Errorf("invalid node type")
	}
}

func (i *Interpreter) compareBoolean(state State, op int, left expr, right expr) (bool, error) {
	lval, err := i.evalBool(left, state)
	if err != nil {
		return false, err
	}
	rval, err := i.evalBool(right, state)
	if err != nil {
		return false, err
	}
	switch op {
	case CMP_AND:
		return lval && rval, nil
	case CMP_OR:
		return lval || rval, nil
	case CMP_EQ:
		return lval == rval, nil
	case CMP_NEQ:
		return lval != rval, nil
	}
	return false, fmt.Errorf("invalid boolean comparator")
}

func (i *Interpreter) compareString(state State, op int, left expr, right expr) (bool, error) {
	lval, err := i.evalBool(left, state)
	if err != nil {
		return false, err
	}
	rval, err := i.evalBool(right, state)
	if err != nil {
		return false, err
	}
	switch op {
	case CMP_EQ:
		return lval == rval, nil
	case CMP_NEQ:
		return lval != rval, nil
	}
	return false, fmt.Errorf("invalid string comparator")
}

func (i *Interpreter) compareNumber(state State, op int, left expr, right expr) (bool, error) {
	lval, err := i.evalNum(left, state)
	if err != nil {
		return false, err
	}
	rval, err := i.evalNum(right, state)
	if err != nil {
		return false, err
	}

	switch op {
	case CMP_EQ:
		return lval == rval, nil
	case CMP_NEQ:
		return lval != rval, nil
	case CMP_LT:
		return lval < rval, nil
	case CMP_LTE:
		return lval <= rval, nil
	case CMP_GT:
		return lval > rval, nil
	case CMP_GTE:
		return lval >= rval, nil
	}
	return false, fmt.Errorf("invalid number comparator")
}

func (i *Interpreter) evalNum(e expr, state State) (float64, error) {
	switch t := e.(type) {
	case float64:
		return t, nil
	case *parenExpr:
		return i.evalNum(t.expr, state)
	case *fnStmt:
		switch t.fn {
		case GET:
			key, err := i.evalStr(t.expr)
			if err != nil {
				return 0, err
			}
			return state.Get(key)
		case SET:
			key, err := i.evalStr(t.expr)
			if err != nil {
				return 0, err
			}
			val, err := i.evalNum(t.expr2, state)
			if err != nil {
				return 0, err
			}
			return val, state.Set(key, val)
		}
		return 0, fmt.Errorf("invalid numeric function")
	case *mathExpr:
		left, err := i.evalNum(t.left, state)
		if err != nil {
			return 0, err
		}
		right, err := i.evalNum(t.right, state)
		if err != nil {
			return 0, err
		}
		switch t.op {
		case '+':
			return left + right, nil
		case '-':
			return left - right, nil
		case '*':
			return left * right, nil
		case '/':
			return left / right, nil
		}
		return 0, fmt.Errorf("invalid math operator")

	default:
		return 0, fmt.Errorf("invalid node type")
	}
}

func (i *Interpreter) evalStr(e expr) (string, error) {
	switch t := e.(type) {
	case string:
		return t, nil
	case *parenExpr:
		return i.evalStr(t.expr)
	default:
		return "", fmt.Errorf("invalid node type")
	}
}
