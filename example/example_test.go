package example_test

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/atotto/webtester"
	"github.com/atotto/webtester/chrome"
)

func TestMain(m *testing.M) {
	if err := chrome.SetupDriver(); err != nil {
		log.Fatal(err)
	}
	ret := m.Run()
	os.Exit(ret)
}

func TestSimple(tt *testing.T) {
	t := webtester.Setup(tt, chrome.DriverPath)
	defer t.TearDown()

	d := t.OpenBrowser()
	d.SetPageLoadTimeout(10 * time.Second)

	d.VisitTo("https://tour.golang.org/")
	d.WaitFor("id:run").Element().Click()
	d.Expect("class:stdout", "Hello")
	d.MustFindElement("class:stdout").VerifyText(strings.Contains, "Hello")
	//d.MustFindElements("class:stdout").Verify(func(e *Element) {
	//	strings.Contains("Hello")
	//})

	d.Find("class:next-page").Click()
	d.ExpectTransitTo("/welcome/2").TakeScreenshot("page2.png")
}
