//go:build js && wasm

package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"syscall/js"

	"github.com/jasoncabot/fabled-story-book/internal/jabl"
)

var interpreter *jabl.Interpreter

func main() {
	// Keep this 'program' running and ready to receive function calls from Javascript
	c := make(chan struct{}, 0)

	mapper := &jsState{}
	loader := &jsLoader{fn: "loadSection"}

	// Create an instance of the interpreter we will use to execute the code
	// There is nothing shared between executions of the interpreter, so we can use a single instance
	interpreter = jabl.NewInterpreter(mapper, loader)

	registerCallbacks()

	<-c
}

// A loader that delegates back to a Javascript function `loadSection` for loading the next block of code to execute
type jsLoader struct {
	fn string
}

func (l *jsLoader) LoadSection(identifier jabl.SectionId, onLoaded func(code string, err error)) {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) any {
		result := args[0].String()
		err := args[1]
		if !err.IsNull() {
			onLoaded("", errors.New(err.String()))
		} else {
			onLoaded(result, nil)
		}
		cb.Release()
		return nil
	})
	js.Global().Get("window").Call(l.fn, string(identifier), cb)
}

func execSection(this js.Value, inputs []js.Value) any {
	section := jabl.SectionId(inputs[0].String())
	callback := inputs[1]

	// The interpreter delegates back to the loader for getting the code to execute from an identifier
	interpreter.Execute(section, func(section *jabl.Result, err error) {
		if err != nil {
			callback.Invoke(js.Null(), err.Error())
		} else {
			jsonValueOfRes, err := json.Marshal(section)
			if err != nil {
				callback.Invoke(js.Null(), err.Error())
			} else {
				callback.Invoke(string(jsonValueOfRes), js.Null())
			}
		}
	})

	return nil
}

func registerCallbacks() {
	js.Global().Set("execSection", js.FuncOf(execSection))
}

// A state mapper that delegates back to Javascript for getting and setting state
type jsState struct{}

func (s *jsState) Get(key string) (float64, error) {
	value := js.Global().Get("localStorage").Call("getItem", key)
	if value.IsNull() {
		return 0, nil
	}
	return strconv.ParseFloat(value.String(), 64)
}

func (s *jsState) Set(key string, value float64) error {
	js.Global().Get("localStorage").Call("setItem", key, value)
	return nil
}
