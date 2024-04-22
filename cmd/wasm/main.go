//go:build js && wasm

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/jasoncabot/fabled-story-book/internal/jabl"
)

var interpreter *jabl.Interpreter
var state jabl.State

func main() {
	// Keep this 'program' running and ready to receive function calls from Javascript
	c := make(chan struct{}, 0)

	state = &localStorageMapper{}
	loader := &jsLoader{fn: "loadSection"}

	// Create an instance of the interpreter we will use to execute the code
	// There is nothing shared between executions of the interpreter, so we can use a single instance
	interpreter = jabl.NewInterpreter(loader)

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
		if !err.IsUndefined() {
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
	interpreter.Execute(section, state, func(section *jabl.Result, err error) {
		if err != nil {
			callback.Invoke(js.Undefined(), err.Error())
		} else {
			jsonValueOfRes, err := json.Marshal(section)
			if err != nil {
				callback.Invoke(js.Undefined(), err.Error())
			} else {
				callback.Invoke(string(jsonValueOfRes), js.Undefined())
			}
		}
	})

	return nil
}

func evalCode(this js.Value, inputs []js.Value) any {
	callback := inputs[2]

	interpreter.Evaluate(inputs[0].String(), inputs[1].String(), state, func(section *jabl.Result, err error) {
		jsonValueOfRes, err := json.Marshal(section)
		if err != nil {
			callback.Invoke(js.Undefined(), err.Error())
		} else {
			callback.Invoke(string(jsonValueOfRes), js.Undefined())
		}
	})
	return nil
}

func registerCallbacks() {
	js.Global().Set("execSection", js.FuncOf(execSection))
	js.Global().Set("evalCode", js.FuncOf(evalCode))
}

// A state mapper that delegates back to Javascript for getting and setting state
type localStorageMapper struct{}

func (s *localStorageMapper) Get(key string) (any, error) {
	value := js.Global().Get("bookStorage").Call("getItem", key)
	if value.IsUndefined() {
		return 0, nil
	}

	switch value.Type() {
	case js.TypeNumber:
		return value.Float(), nil
	case js.TypeString:
		return value.String(), nil
	case js.TypeBoolean:
		return value.Bool(), nil
	}
	return nil, fmt.Errorf("unknown type %s", value.Type())
}

func (s *localStorageMapper) Set(key string, value any) error {
	if key != "system:source" {
		switch v := value.(type) {
		case float64:
			js.Global().Get("bookStorage").Call("setItem", key, "n", strconv.FormatFloat(v, 'f', -1, 64))
		case string:
			js.Global().Get("bookStorage").Call("setItem", key, "s", v)
		case bool:
			js.Global().Get("bookStorage").Call("setItem", key, "b", strconv.FormatBool(v))
		}
	}

	return nil
}
