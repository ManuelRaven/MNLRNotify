package executors

import (
	"github.com/dop251/goja"
)

type GojaExecutor struct{}

func init() {
	RegisterExecutor("goja", &GojaExecutor{})
}

func (e *GojaExecutor) Execute(message string, script string) string {
	vm := goja.New()

	// Set message as a global variable
	err := vm.Set("message", message)
	if err != nil {
		return "Error: Failed to set message variable"
	}

	// Execute the script
	val, err := vm.RunString(script)
	if err != nil {
		return "Error: " + err.Error()
	}

	// Convert result to string
	result := val.String()
	if result == "undefined" {
		return message
	}

	return result
}
