package imply

import (
	"fmt"

	js "github.com/dop251/goja"
)

func basicLog(call js.FunctionCall) js.Value {
	str := call.Argument(0)
	fmt.Println(str.String())
	return str
}

// RegisterConsole register a console.basicLog to runtime
func (c *Core) RegisterConsole(log func(call js.FunctionCall) js.Value) error {
	o := c.NewObject()
	err := o.Set("log", log)
	if err != nil {
		return err
	}
	err = c.Set("console", o)
	return err
}
