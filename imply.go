package imply

import (
	"sync"

	js "github.com/dop251/goja"
)

// Core is the basic struct of gode
type Core struct {
	*js.Runtime
	*sync.RWMutex
	Pkg map[string]js.Value
}

// NewCore creates a pointer to a new Core.
// A Core is a type that is used to organize a collection of runtime javascript modules.
func NewCore() *Core {
	return &Core{
		Runtime: js.New(),
		RWMutex: &sync.RWMutex{},
		Pkg:     make(map[string]js.Value),
	}
}

// RegisterBuiltins register some build in modules to the runtime
func (c *Core) RegisterBuiltins() error {
	if err := c.RegisterConsole(basicLog); err != nil {
		return err
	}
	return c.RegisterLoader()
}

// RegisterLoader registers a simple commonjs style loader to runtime.
func (c *Core) RegisterLoader() error {
	err := c.Set("require", func(call js.FunctionCall) js.Value {
		var val js.Value
		var err error
		val, err = c.LoadJSPackageFile(call.Argument(0).String())
		if err != nil {
			return c.NewGoError(err)
		}
		return val
	})
	return err
}
