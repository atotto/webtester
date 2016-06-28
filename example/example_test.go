package example_test

import (
	"log"
	"os"
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

	d.VisitTo("https://tour.golang.org/").WaitFor("id:run").Element().Click()
	d.Expect("class:stdout", "Hello, 世界")

	d.Find("class:next-page").Click()
	d.ExpectTransitTo("/welcome/2").TakeScreenshot("page2.png")
}
