package jabl

import (
	"fmt"
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
	state  StateMapper
	loader SectionLoader
}

type StateMapper interface {
	Get(key string) (float64, error)
	Set(key string, value float64) error
}

type SectionLoader interface {
	LoadSection(identifier SectionId, onLoad func(code string, err error))
}

func NewInterpreter(state StateMapper, loader SectionLoader) *Interpreter {
	return &Interpreter{
		state:  state,
		loader: loader,
	}
}

func (i *Interpreter) Execute(identifier SectionId, callback func(result *Result, err error)) {
	i.loader.LoadSection(identifier, func(code string, err error) {
		if err != nil {
			callback(nil, err)
			return
		}
		i.Evaluate(code, callback)
	})
}

func (i *Interpreter) Evaluate(code string, callback func(result *Result, err error)) {
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
	if err := i.eval(lexer.ast, result); err != nil {
		callback(nil, fmt.Errorf("error evaluating expression: %v", err))
		return
	}

	// Congrats we've evaluated our code
	callback(result, nil)
}

func (i *Interpreter) eval(node any, result *Result) error {
	// We can evaluate a nil statement, it just does nothing, it's not an error
	if node == nil {
		return nil
	}

	switch t := node.(type) {
	case *program:
		return i.eval(t.body, result)
	case *blockStmt:
		return i.eval(t.stmt, result)
	case *seqStmt:
		if err := i.eval(t.first, result); err != nil {
			return err
		}
		if err := i.eval(t.rest, result); err != nil {
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
			i.printCode(t.block, sb)

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
			val, err := i.evalNum(t.expr2)
			if err != nil {
				return err
			}
			return i.state.Set(key, val)
		default:
			switch t.fn {
			case GET:
				return fmt.Errorf("GET not supported in this context")
			case SET:
				return fmt.Errorf("SET not supported in this context")
			case CHOICE:
				return fmt.Errorf("CHOICE not supported in this context")
			case PRINT:
				return fmt.Errorf("PRINT not supported in this context")
			case GOTO:
				return fmt.Errorf("GOTO not supported in this context")
			default:
				return fmt.Errorf("invalid function: %d", t.fn)
			}
		}
	case *ifStmt:
		val, err := i.evalBool(t.cond)
		if err != nil {
			return err
		}
		if val {
			return i.eval(t.block, result)
		} else if t.other != nil {
			return i.eval(t.other, result)
		}
		return nil
	default:
		return fmt.Errorf("invalid statement type: %s", t)
	}
}

func (i *Interpreter) printCode(stmt any, sb *strings.Builder) {
	switch t := stmt.(type) {
	case *program:
		i.printCode(t.body, sb)
	case *blockStmt:
		sb.WriteRune('{')
		i.printCode(t.stmt, sb)
		sb.WriteRune('}')
	case *seqStmt:
		i.printCode(t.first, sb)
		i.printCode(t.rest, sb)
	case *fnStmt:
		switch t.fn {
		case PRINT:
			sb.WriteString("print")
		case GOTO:
			sb.WriteString("goto")
		case CHOICE:
			sb.WriteString("choice")
		case GET:
			sb.WriteString("get")
		case SET:
			sb.WriteString("set")
		}
		sb.WriteRune('(')
		i.printCode(t.expr, sb)
		if t.block != nil {
			sb.WriteRune(',')
			i.printCode(t.block, sb)
		} else if t.expr2 != nil {
			sb.WriteRune(',')
			i.printCode(t.expr2, sb)
		}
		sb.WriteRune(')')
	case *ifStmt:
		sb.WriteString("if(")
		i.printCode(t.cond, sb)
		sb.WriteRune(')')
		i.printCode(t.block, sb)
		if t.other != nil {
			sb.WriteString("else")
			i.printCode(t.other, sb)
		}
	case *parenExpr:
		sb.WriteRune('(')
		i.printCode(t.expr, sb)
		sb.WriteRune(')')
	case *notExpr:
		sb.WriteRune('!')
		i.printCode(t.expr, sb)
	case *cmpExpr:
		i.printCode(t.left, sb)
		switch t.op {
		case CMP_EQ:
			sb.WriteString("==")
		case CMP_NEQ:
			sb.WriteString("!=")
		case CMP_LT:
			sb.WriteString("<")
		case CMP_LTE:
			sb.WriteString("<=")
		case CMP_GT:
			sb.WriteString(">")
		case CMP_GTE:
			sb.WriteString(">=")
		}
		i.printCode(t.right, sb)
	case *mathExpr:
		i.printCode(t.left, sb)
		sb.WriteRune(rune(t.op))
		i.printCode(t.right, sb)
	case float64:
		sb.WriteString(fmt.Sprintf("%f", t))
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

func (i *Interpreter) evalBool(e expr) (bool, error) {
	switch t := e.(type) {
	case bool:
		return t, nil
	case *parenExpr:
		return i.evalBool(t.expr)
	case *notExpr:
		val, err := i.evalBool(t.expr)
		if err != nil {
			return false, err
		}
		return !val, nil
	case *cmpExpr:
		switch t.t {
		case BOOLEAN:
			return i.compareBoolean(t.op, t.left, t.right)
		case STRING:
			return i.compareString(t.op, t.left, t.right)
		case NUMBER:
			return i.compareNumber(t.op, t.left, t.right)
		}

		return false, fmt.Errorf("invalid comparison operator")
	default:
		return false, fmt.Errorf("invalid node type")
	}
}

func (i *Interpreter) compareBoolean(op int, left expr, right expr) (bool, error) {
	lval, err := i.evalBool(left)
	if err != nil {
		return false, err
	}
	rval, err := i.evalBool(right)
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

func (i *Interpreter) compareString(op int, left expr, right expr) (bool, error) {
	lval, err := i.evalBool(left)
	if err != nil {
		return false, err
	}
	rval, err := i.evalBool(right)
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

func (i *Interpreter) compareNumber(op int, left expr, right expr) (bool, error) {
	lval, err := i.evalNum(left)
	if err != nil {
		return false, err
	}
	rval, err := i.evalNum(right)
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

func (i *Interpreter) evalNum(e expr) (float64, error) {
	switch t := e.(type) {
	case float64:
		return t, nil
	case *parenExpr:
		return i.evalNum(t.expr)
	case *fnStmt:
		switch t.fn {
		case GET:
			key, err := i.evalStr(t.expr)
			if err != nil {
				return 0, err
			}
			return i.state.Get(key)
		case SET:
			key, err := i.evalStr(t.expr)
			if err != nil {
				return 0, err
			}
			val, err := i.evalNum(t.expr2)
			if err != nil {
				return 0, err
			}
			return val, i.state.Set(key, val)
		}
		return 0, fmt.Errorf("invalid numeric function")
	case *mathExpr:
		left, err := i.evalNum(t.left)
		if err != nil {
			return 0, err
		}
		right, err := i.evalNum(t.right)
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
