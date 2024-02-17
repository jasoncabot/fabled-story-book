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

				// marshal result as JSON string and compare to the contents of the file print_string.json
				resultAsJSON, err := json.Marshal(r)
				assert.Nil(t, err)
				// load file print_string.json
				expectedResult := loadFile(t, filepath.Join("testdata", "examples", name, "result.json"))

				assert.JSONEq(t, expectedResult, string(resultAsJSON))

				// compare the state before and after
				stateAsJSON, err := json.Marshal(state.testValues)
				require.NoError(t, err)

				// load file .state.2
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
		loader := &testLoader{testErr: loadErr}
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
		lexer := newLexer(strings.NewReader(code))
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
	testValues map[string]float64
}

func NewEmptyState() *testState {
	return &testState{
		testValues: map[string]float64{},
	}
}

func NewTestState(before string) (*testState, error) {
	initialState := map[string]float64{}
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

func (t *testState) Get(key string) (float64, error) {
	if value, ok := t.testValues[key]; ok {
		return value, nil
	}
	return 0, nil
}

func (t *testState) Set(key string, value float64) error {
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
