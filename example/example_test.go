package example_test

import (
	"fmt"
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
	d.SetPageLoadTimeout(4 * time.Second)

	d.VisitTo("https://tour.golang.org/welcome/1")
	// ngのレンダリングを待たなければrunの結果が出てこない
	d.TakeSource("./before.html")
	time.Sleep(2 * time.Second)
	d.TakeSource("./after.html")

	d.WaitFor("id:run")
	d.MustFindElement("id:run").Click()

	d.WaitFor("class:stdout")
	d.MustFindElement("class:stdout").VerifyText(strings.Contains, "Hello")

	d.MustFindElement("class:next-page").Click()
	d.ExpectTransitTo("/welcome/2").TakeScreenshot("page2.png")
}

var code = `
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, go世界!")
}
`

func TestPlayground(t *testing.T) {
	ts := webtester.Setup(t, chrome.DriverPath)
	defer ts.TearDown()

	b := ts.OpenBrowser()
	b.SetPageLoadTimeout(2 * time.Second)

	b.VisitTo("https://play.golang.org/")
	b.WaitFor("id:code").Element().Clear().Input(code)
	b.MustFindElement("id:run").Click()

	b.WaitForText("class:stdout", "Hello")
	b.MustFindElement("class:stdout").VerifyText(strings.Contains, "Hello")

	es := b.MustFindElements("id:controls")
	for _, e := range es {
		text, _ := e.WebElement().Text()
		fmt.Printf("%+v\n", text)
	}

	b.TakeScreenshot("./screenshot.png")
	b.TakeSource("./source.html")
}
