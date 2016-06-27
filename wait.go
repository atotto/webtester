package webtest

import (
	"time"

	"github.com/bborbe/webdriver"
)

func wait(fn func() bool) bool {
	for _, w := range []time.Duration{1, 2, 3, 5, 7, 11, 13, 17} {
		if ok := fn(); ok {
			return true
		}
		time.Sleep(w * time.Second)
	}
	return false
}

func WaitElement(session *webdriver.Session, using webdriver.FindElementStrategy, value string) (elem webdriver.WebElement, err error) {
	ok := wait(func() bool {
		elem, err = session.FindElement(using, value)
		return err == nil
	})
	if !ok {
		return webdriver.WebElement{}, err
	}
	return elem, nil
}
