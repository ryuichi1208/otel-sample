package main

import (
	"strconv"
	"syscall/js"
)

func increment(this js.Value, args []js.Value) any {
	counter := js.Global().Get("document").Call("getElementById", "counter")
	counterValue, err := strconv.ParseInt(counter.Get("textContent").String(), 10, 64)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	counterValue += int64(args[0].Int())
	counter.Set("textContent", counterValue)
	return map[string]any{"message": counterValue}
}

func main() {
	js.Global().Set("goIncrement", js.FuncOf(increment))
	select {} // keep running
}
