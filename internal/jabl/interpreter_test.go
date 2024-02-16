package jabl

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterpreter(t *testing.T) {

	loader := &testLoader{}
	interpreter := NewInterpreter(loader)

	err := filepath.Walk(filepath.Join("testdata"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() || path == "testdata" {
			return nil
		}

		// get the name of the file without the extension
		name := filepath.Base(path)

		// Run a test for every example found in the testdata folder
		t.Run(name, func(t *testing.T) {
			// Ensure that each test runs with a clean state
			state, err := NewTestState(filepath.Join("testdata", name, "state_before.json"))
			require.NoError(t, err)

			// execute the code in the file
			interpreter.Execute(SectionId(name), state, func(r *Result, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, r)

				// marshal result as JSON string and compare to the contents of the file print_string.json
				resultAsJSON, err := json.Marshal(r)
				assert.Nil(t, err)
				// load file print_string.json
				expectedResult := loadFile(t, filepath.Join("testdata", name, "result.json"))

				assert.JSONEq(t, expectedResult, string(resultAsJSON))

				// compare the state before and after
				stateAsJSON, err := json.Marshal(state.testValues)
				require.NoError(t, err)

				// load file .state.2
				expectedState := loadFile(t, filepath.Join("testdata", name, "state_after.json"))
				assert.JSONEq(t, expectedState, string(stateAsJSON))
			})

		})
		return nil
	})
	require.NoError(t, err)
}

type testLoader struct {
}

func (t *testLoader) LoadSection(identifier SectionId, onLoad func(code string, err error)) {
	file, err := os.Open(filepath.Join("testdata", string(identifier), "code.jabl"))
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
