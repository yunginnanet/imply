package imply

import (
	"fmt"
	"os"

	js "github.com/dop251/goja"
)

func moduleTemplate(c string) string {
	return "(function(module, exports) {\n" + c + "\n})"
}

func (c *Core) createPack() (*js.Object, error) {
	m := c.NewObject()
	e := c.NewObject()
	err := m.Set("exports", e)
	return m, err
}

func load(name string, code []byte) (*js.Program, error) {
	text := moduleTemplate(string(code))
	return js.Compile(name, text, false)
}

func (c *Core) validate(prg *js.Program) (jsExports js.Value, err error) {
	var f js.Value
	f, err = c.RunProgram(prg)
	if err != nil {
		return
	}
	g, ok := js.AssertFunction(f)
	if !ok {
		err = fmt.Errorf("[%T] %v is not a function", f, f)
		return
	}
	pkg, err := c.createPack()
	if err != nil {
		return
	}
	jsExports = pkg.Get("exports")
	_, err = g(jsExports, pkg, jsExports)
	if err != nil {
		return
	}
	// fmt.Printf("VALIDATED: [%T] %v\n", jsExports, spew.Sprint(jsExports))
	return
}

func (c *Core) checkExisting(p string) (js.Value, bool) {
	c.RLock()
	defer c.RUnlock()
	val, ok := c.Pkg[p]
	return val, ok
}

func (c *Core) registerModule(p string, jsExports js.Value) {
	c.Lock()
	c.Pkg[p] = jsExports
	c.Unlock()
}

func (c *Core) validateAndRegister(name string, prg *js.Program) (jsExports js.Value, err error) {
	if jsExports, err = c.validate(prg); err != nil {
		return nil, err
	}
	if jsExports == nil {
		return nil, fmt.Errorf("module %s does not export anything", name)
	}
	c.registerModule(name, jsExports)
	return jsExports, nil
}

// LoadJSPackageFile loads a javascript package from a file, validates it, and registers it.
// The given path will be used as the name of the module.
//
// Warning: If name is already registered, the existing module will be returned.
func (c *Core) LoadJSPackageFile(path string) (js.Value, error) {
	code, err := os.ReadFile(path)
	if err != nil {
		// println("!!ERR!! " + err.Error())
		return nil, err
	}
	// println("new: " + path)
	return c.LoadJSPackage(path, code)
}

// LoadJSPackage loads a javascript package from code, validates it, and registers it by name.
//
// Warning: If name is already registered, the existing module will be returned.
func (c *Core) LoadJSPackage(name string, code []byte) (js.Value, error) {
	if mod, ok := c.checkExisting(name); ok {
		return mod, nil
	}
	prg, err := load(name, code)
	if err != nil {
		return nil, err
	}
	return c.validateAndRegister(name, prg)
}
