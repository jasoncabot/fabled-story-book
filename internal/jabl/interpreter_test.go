package jabl

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterpreter_QuickCodeCheck(t *testing.T) {
	t.Run("should print a single line of output properly", func(t *testing.T) {
		state := &testState{
			testValues: map[string]any{},
		}

		loader := &testLoader{}
		// use a known seed for generating random numbers in tests to make them deterministic
		interpreter := NewInterpreter(loader, WithRandomNumberGenerator(&fixedRandom{}))
		// execute the code in the file
		code := `{
			set("string", "1")
			set("bool", true)
			set("num", 1)
			print("str" + get("string"))
			print("str" + get("bool"))
			print("str" + get("num"))
			print(1 + getn("string"))
			print(1 + getn("bool"))
			print(1 + getn("num"))
			print(true && getb("string"))
			print(true && getb("bool"))
			print(true && getb("num"))
		}`
		interpreter.Evaluate("inline_test", code, state, func(r *Result, err error) {
			require.Nil(t, err)
			require.NotNil(t, r)
			require.NotNil(t, r.Output)

			// split output by newlines
			output := strings.Split(r.Output, "\n")
			require.Len(t, output, 10, "expected 10 lines of output")
			assert.Equal(t, "str1", output[0], "string")
			assert.Equal(t, "strtrue", output[1], "bool")
			assert.Equal(t, "str1", output[2], "num")
			assert.Equal(t, "2", output[3], "string")
			assert.Equal(t, "2", output[4], "bool")
			assert.Equal(t, "2", output[5], "num")
			assert.Equal(t, "true", output[6], "string")
			assert.Equal(t, "true", output[7], "bool")
			assert.Equal(t, "true", output[8], "num")
		})
	})
}

func TestInterpreter(t *testing.T) {

	loader := &testLoader{}
	// use a known seed for generating random numbers in tests to make them deterministic
	interpreter := NewInterpreter(loader, WithRandomNumberGenerator(&fixedRandom{}))

	err := filepath.Walk(filepath.Join("testdata", "examples"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() || path == "testdata/examples" {
			return nil
		}

		// get the name of the file without the extension
		name := filepath.Base(path)

		// Run a test for every example found in the testdata folder
		t.Run(name, func(t *testing.T) {
			// Ensure that each test runs with a clean state
			state, err := NewTestState(filepath.Join("testdata", "examples", name, "state_before.json"))
			require.NoError(t, err)

			// execute the code in the file
			interpreter.Execute(SectionId(name), state, func(r *Result, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, r)

				// marshal result as JSON string and compare to the contents of the file `result.json`
				resultAsJSON, err := json.Marshal(r)
				assert.Nil(t, err)
				// load file result.json
				expectedResult := loadFile(t, filepath.Join("testdata", "examples", name, "result.json"))

				assert.JSONEq(t, expectedResult, string(resultAsJSON))

				// compare the state before and after
				stateAsJSON, err := json.Marshal(state.testValues)
				require.NoError(t, err)

				// load file `state_after.json`
				expectedState := loadFile(t, filepath.Join("testdata", "examples", name, "state_after.json"))
				assert.JSONEq(t, expectedState, string(stateAsJSON))
			})

		})
		return nil
	})
	require.NoError(t, err)
}

func TestInterpreterError(t *testing.T) {
	t.Run("when the section loader returns an error, the interpreter should return it", func(t *testing.T) {
		loadErr := errors.New("failed to load due to some dodgy reason")
		loader := &testLoader{
			testErr: loadErr,
		}
		interpreter := NewInterpreter(loader)

		state := NewEmptyState()
		interpreter.Execute(SectionId("section_loader_error"), state, func(r *Result, err error) {
			assert.NotNil(t, err)
			assert.Nil(t, r)

			assert.ErrorIs(t, err, loadErr)
		})
	})
}

func TestInterpreterPrintString(t *testing.T) {
	loader := &testLoader{}
	interpreter := NewInterpreter(loader, WithRandomNumberGenerator(&fixedRandom{}))

	yyErrorVerbose = true

	err := filepath.Walk(filepath.Join("testdata", "pretty_printing"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".jabl") {
			return nil
		}

		code := loadFile(t, path)
		formattedCode := loadFile(t, path+".formatted")
		lexer := newLexer(path, strings.NewReader(code))
		parseResult := yyParse(lexer)
		require.NoError(t, lexer.err)
		require.Equal(t, parseResult, 0, "failed to parse code")

		sb := &strings.Builder{}
		interpreter.printCode(lexer.ast, sb, 0)

		assert.Equal(t, formattedCode, sb.String())

		// write formattedCode to an .out file at the same place to make it easy to debug
		if formattedCode != sb.String() {
			outFile, err := os.Create(path + ".out")
			require.NoError(t, err)
			defer outFile.Close()
			_, err = outFile.WriteString(sb.String())
			require.NoError(t, err)
		}

		return nil
	})
	assert.NoError(t, err)
}

type fixedRandom struct {
}

func (f *fixedRandom) Float64() float64 {
	return 0.99999999999 // chosen by a fair dice roll #221 - actually it's just to give the max value of any dice roll so 1d6 will always be 6
}

type testLoader struct {
	testErr error
}

func (t *testLoader) LoadSection(identifier SectionId, onLoad func(code string, err error)) {
	if t.testErr != nil {
		onLoad("", t.testErr)
		return
	}

	file, err := os.Open(filepath.Join("testdata", "examples", string(identifier), "code.jabl"))
	if err != nil {
		onLoad("", err)
		return
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		onLoad("", err)
		return
	}
	onLoad(string(bytes), nil)
}

type testState struct {
	testValues map[string]any
}

func NewEmptyState() *testState {
	return &testState{
		testValues: map[string]any{},
	}
}

func NewTestState(before string) (*testState, error) {
	initialState := map[string]any{}
	file, err := os.Open(before)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&initialState); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return &testState{
		testValues: initialState,
	}, nil
}

func (t *testState) Get(key string) (any, error) {
	if value, ok := t.testValues[key]; ok {
		return value, nil
	}
	return 0, nil
}

func (t *testState) Set(key string, value any) error {
	t.testValues[key] = value
	return nil
}

func loadFile(t *testing.T, filename string) string {
	file, err := os.Open(filename)
	require.NoError(t, err)
	defer file.Close()
	bytes, err := io.ReadAll(file)
	require.NoError(t, err)
	return string(bytes)
}
