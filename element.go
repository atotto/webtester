package webtester

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/bborbe/webdriver"
)

type Element struct {
	testing.TB
	elem webdriver.WebElement
}

func (e *Element) WebElement() (elem webdriver.WebElement) {
	return e.elem
}

func (e *Element) Click() *Element {
	e.Helper()
	if err := e.elem.Click(); err != nil {
		e.Fatalf("click failed: %+v", err)
	}
	return e
}

func (e *Element) Clear() *Element {
	if err := e.elem.Clear(); err != nil {
		e.Fatalf(`clear failed: %+v`, err)
	}
	return e
}

func (e *Element) Input(text string) *Element {
	if err := e.elem.SendKeys(text); err != nil {
		e.Fatalf(`input text "%s" failed: %+v`, text, err)
	}
	return e
}

func (e *Element) VerifyText(fn func(string, string) bool, expect string) *Element {
	e.Helper()
	actual, err := e.elem.Text()
	if err != nil {
		e.Fatal(err)
	}
	if !fn(actual, expect) {
		name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		ss := strings.Split(name, ".")
		if len(ss) == 2 {
			name = ss[1]
		}
		e.Fatalf("want %s %s, got %s", strings.ToLower(name), expect, actual)
	}
	return e
}

func (e *Element) WaitForEnabled() *Element {
	e.Helper()

	var enabled bool
	var err error
	ok := wait(func() bool {
		enabled, err = e.elem.IsEnabled()
		if err != nil {
			e.Fatal(err)
		}
		return enabled
	})
	if !ok {
		e.Fatal(err)
	}
	return e
}

func (e *Element) WaitForDisabled() *Element {
	e.Helper()

	var enabled bool
	var err error
	ok := wait(func() bool {
		enabled, err = e.elem.IsEnabled()
		if err != nil {
			e.Fatal(err)
		}
		return !enabled
	})
	if !ok {
		e.Fatal(err)
	}
	return e
}
