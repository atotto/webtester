package webtester

import (
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bborbe/webdriver"
)

type Browser struct {
	testing.TB
	session *webdriver.Session
	element webdriver.WebElement
}

func (b *Browser) Session() (session *webdriver.Session) {
	return b.session
}

func (b *Browser) WebElement() (elem webdriver.WebElement) {
	return b.element
}

func (b *Browser) Element() (elem *Element) {
	return &Element{
		TB:   b.TB,
		elem: b.element,
	}
}

func (b *Browser) SetPageLoadTimeout(timeout time.Duration) {
	b.Helper()
	if err := b.session.SetTimeouts("page load", toMillisecond(timeout)); err != nil {
		b.Fatal(err)
	}
}

func toMillisecond(d time.Duration) int {
	return int(d / time.Millisecond)
}

func (b *Browser) VisitTo(rawurl string) *Browser {
	b.Helper()
	if _, err := url.Parse(rawurl); err != nil {
		b.Fatal(err)
	}
	if err := b.session.Url(rawurl); err != nil {
		b.Fatal(err)
	}
	return b
}

func (b *Browser) WaitFor(target string) *Browser {
	b.Helper()
	using, value := splitTarget(b.TB, target)

	var elem webdriver.WebElement
	var err error
	ok := wait(func() bool {
		elem, err = b.session.FindElement(using, value)
		return err == nil
	})
	if !ok {
		b.Fatal(err)
	}
	b.element = elem
	return b
}

func (b *Browser) WaitForText(target string, text string) *Browser {
	b.Helper()
	using, value := splitTarget(b.TB, target)

	var elem webdriver.WebElement
	var err error
	ok := wait(func() bool {
		elem, err = b.session.FindElement(using, value)
		if err != nil {
			return false
		}
		content, err := elem.Text()
		if err != nil {
			return false
		}
		if !strings.Contains(content, text) {
			return false
		}
		return true
	})
	if !ok {
		b.Fatal(err)
	}
	b.element = elem
	return b
}

func splitTarget(tb testing.TB, target string) (using webdriver.FindElementStrategy, value string) {
	tb.Helper()
	tags := strings.SplitN(target, ":", 2)
	if len(tags) != 2 {
		tb.Fatal("expect target format `using:value`")
	}

	using, ok := toStrategy(tags[0])
	if !ok {
		tb.Fatalf("not supported: using=%s", using)
	}
	return using, tags[1]
}

func toStrategy(usingString string) (using webdriver.FindElementStrategy, ok bool) {
	u := webdriver.FindElementStrategy(usingString)
	switch u {
	case webdriver.ClassName, webdriver.CSS_Selector, webdriver.ID, webdriver.Name, webdriver.LinkText, webdriver.PartialLinkText, webdriver.TagName, webdriver.XPath:
		return u, true
	case "class":
		return webdriver.ClassName, true
	case "css":
		return webdriver.CSS_Selector, true
	case "tag":
		return webdriver.TagName, true
	default:
		return "", false
	}
}

func (b *Browser) Expect(target string, text string) {
	b.Helper()
	b.Log("Deprecated")
	using, value := splitTarget(b.TB, target)

	var elems []webdriver.WebElement
	var err error
	ok := wait(func() bool {
		elems, err = b.session.FindElements(using, value)
		if err != nil {
			return false
		}
		for _, elem := range elems {
			actual, err := elem.Text()
			if err != nil {
				return false
			}
			if strings.Contains(actual, text) {
				return true
			}
		}
		return false
	})
	if !ok {
		b.Log(err)
		b.Fatalf("not found: %s", text)
	}
}

func (b *Browser) MustFindElement(target string) *Element {
	b.Helper()
	using, value := splitTarget(b.TB, target)

	elem, err := b.session.FindElement(using, value)
	if err != nil {
		b.Fatalf(`element "%s" notfound: %+v`, target, err)
	}
	b.element = elem
	return b.Element()
}

func (b *Browser) MustFindElements(target string) []*Element {
	b.Helper()
	using, value := splitTarget(b.TB, target)

	elems, err := b.session.FindElements(using, value)
	if err != nil {
		b.Fatalf(`elements "%s" notfound: %+v`, target, err)
	}
	if len(elems) != 0 {
		b.element = elems[0]
	}
	es := make([]*Element, len(elems))
	for i, elem := range elems {
		es[i] = &Element{
			TB:   b.TB,
			elem: elem,
		}
	}
	return es
}

func (b *Browser) TakeScreenshot(name string) *Browser {
	b.Helper()
	buf, err := b.session.Screenshot()
	if err != nil {
		b.Fatal(err)
	}

	err = ioutil.WriteFile(name, buf, 0644)
	if err != nil {
		b.Fatal(err)
	}
	return b
}

func (b *Browser) TakeSource(name string) *Browser {
	b.Helper()
	str, err := b.session.Source()
	if err != nil {
		b.Fatal(err)
	}

	err = ioutil.WriteFile(name, []byte(str), 0644)
	if err != nil {
		b.Fatal(err)
	}
	return b
}

func (b *Browser) ExpectTransitTo(rawurl string) *Browser {
	b.Helper()
	expect, err := url.Parse(rawurl)
	if err != nil {
		b.Fatal(err)
	}
	ok := wait(func() bool {
		ru, err := b.session.GetUrl()
		if err != nil {
			b.Log(err)
		}
		u, err := url.Parse(ru)
		if err != nil {
			b.Log(err)
		}
		return u.Path == expect.Path
	})
	if !ok {
		b.Log(err)
		b.Fatalf("not found: %s", rawurl)
	}
	return b
}

func (b *Browser) SetWindowSize(width, height int) {
	b.session.GetCurrentWindowHandle().SetSize(webdriver.Size{Width: width, Height: height})
}

func (b *Browser) GetWindowSize() (width, height int, err error) {
	size, err := b.session.GetCurrentWindowHandle().GetSize()
	if err != nil {
		return
	}
	return size.Width, size.Height, nil
}
