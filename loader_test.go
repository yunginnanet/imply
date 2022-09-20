package imply

import (
	"strings"
	"testing"

	js "github.com/dop251/goja"
)

type testConsole struct {
	t *testing.T
}

func (tc *testConsole) log(call js.FunctionCall) js.Value {
	tc.t.Helper()
	ret := call.Argument(0)
	tc.t.Logf("[JSLOG]%v: [%T] %v", call.This, ret, ret)
	return ret
}

func newTestConsole(t *testing.T) *testConsole {
	t.Helper()
	return &testConsole{t}
}

func testCore(t *testing.T) *Core {
	t.Helper()
	core := NewCore()
	tc := newTestConsole(t)
	if err := core.RegisterConsole(tc.log); err != nil {
		t.Fatalf("RegisterConsole failed: %s", err)
	}
	if err := core.RegisterLoader(); err != nil {
		t.Fatalf("RegisterLoader failed: %s", err)
	}
	_, err := core.RunString("console.log('init');")
	if err != nil {
		t.Fatalf("failed sanity check: %v", err)
	}
	return core
}

func ExampleNewCore() {
	core := NewCore()
	if err := core.RegisterBuiltins(); err != nil {
		panic(err)
	}

}

func TestImply(t *testing.T) {
	const (
		one   = "let t = require('tests/"
		two   = "test.js');\n"
		three = "test-mustfail.js');\n"
	)

	var base *strings.Builder
	var want = "Hello Golang~"

	res := func(base *strings.Builder, js string, shouldFail bool) {
		core := testCore(t)
		jscode := base.String() + js
		t.Logf("Running:\n%s", jscode)
		v, err := core.RunString(jscode)
		switch {
		case err != nil && !shouldFail:
			t.Errorf("RunString failed: %s", err)
		case v == nil && !shouldFail:
			t.Errorf("RunString failed: got nil")
		case v != nil && v.String() != want && !shouldFail:
			t.Errorf("RunString failed.\nwanted: %s\ngot: %s", want, v.String())
		case err == nil && shouldFail:
			t.Errorf("RunString should have failed, got nil error")
		case v != nil && shouldFail:
			t.Errorf("RunString result should have been nil!\nGot: [%T] %v", v, v)
		case err != nil && shouldFail:
			// t.Logf("success, wanted err: [%T] %v", err, err)
		}
	}

	t.Run("Simple", func(t *testing.T) {
		base = &strings.Builder{}
		base.WriteString(one)
		base.WriteString(two)
		res(base, "t.test();", false)
		res(base, "t.foo();", false)
		res(base, "t.bar();", true)
		base.Reset()
	})

	t.Run("InvalidSyntax", func(t *testing.T) {
		base = &strings.Builder{}
		base.WriteString(one)
		base.WriteString(three)
		res(base, "t.test();", true)
		base.Reset()
	})
}
