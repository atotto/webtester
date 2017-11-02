package webtester

import (
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bborbe/webdriver"
)

type Driver struct {
	testing.TB
	webDriver webdriver.WebDriver
	sessions  []*webdriver.Session
}

func Setup(tb testing.TB, path string) *Driver {
	tb.Helper()

	webDriver := webdriver.NewChromeDriver(path)
	err := webDriver.Start()
	if err != nil {
		tb.Fatal(err)
	}

	return &Driver{
		TB:        tb,
		webDriver: webDriver,
	}
}

func (d *Driver) TearDown() {
	for _, session := range d.sessions {
		session.Delete()
	}
	d.webDriver.Stop()
}

type browser struct {
	testing.TB
	session *webdriver.Session
	element webdriver.WebElement
}

func (d *Driver) OpenBrowser() *browser {
	d.Helper()

	desired := webdriver.Capabilities{"Platform": "Linux"}
	required := webdriver.Capabilities{}
	session, err := d.webDriver.NewSession(desired, required)
	if err != nil {
		d.Fatal(err)
	}

	d.sessions = append(d.sessions, session)

	return &browser{
		TB:      d.TB,
		session: session,
	}
}

func (b *browser) Session() (session *webdriver.Session) {
	return b.session
}

func (b *browser) Element() (elem webdriver.WebElement) {
	return b.element
}

func toMillisecond(d time.Duration) int {
	return int(d / time.Millisecond)
}

func (b *browser) SetPageLoadTimeout(timeout time.Duration) {
	b.Helper()
	if err := b.session.SetTimeouts("page load", toMillisecond(timeout)); err != nil {
		b.Fatal(err)
	}
}

func (b *browser) VisitTo(rawurl string) *browser {
	b.Helper()
	if _, err := url.Parse(rawurl); err != nil {
		b.Fatal(err)
	}
	if err := b.session.Url(rawurl); err != nil {
		b.Fatal(err)
	}
	return b
}

func (b *browser) WaitFor(target string) *browser {
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

func (b *browser) Expect(target string, text string) {
	b.Helper()
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

func (b *browser) Find(target string) webdriver.WebElement {
	b.Helper()
	using, value := splitTarget(b.TB, target)

	elem, err := b.session.FindElement(using, value)
	if err != nil {
		b.Fatal(err)
	}
	b.element = elem
	return elem
}

func (b *browser) FindElements(target string) []webdriver.WebElement {
	b.Helper()
	using, value := splitTarget(b.TB, target)

	elems, err := b.session.FindElements(using, value)
	if err != nil {
		b.Fatal(err)
	}
	// b.element = elem
	return elems
}

func (b *browser) TakeScreenshot(name string) *browser {
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

func (b *browser) ExpectTransitTo(rawurl string) *browser {
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
