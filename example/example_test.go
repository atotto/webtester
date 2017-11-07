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
	d.SetPageLoadTimeout(10 * time.Second)

	d.VisitTo("https://tour.golang.org/")
	d.WaitFor("id:run").Element().Click()
	d.Expect("class:stdout", "Hello") // Deprecated

	d.MustFindElement("class:stdout").VerifyText(strings.Contains, "Hello")
	//d.MustFindElements("class:stdout").Verify(func(e *Element) {
	//	strings.Contains("Hello")
	//})

	d.Find("class:next-page").Click()
	d.ExpectTransitTo("/welcome/2").TakeScreenshot("page2.png")
}

var code = `
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, 世界!")
}
`

func TestPlayground(t *testing.T) {
	ts := webtester.Setup(t, chrome.DriverPath)
	defer ts.TearDown()

	b := ts.OpenBrowser()
	b.SetPageLoadTimeout(10 * time.Second)

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
}
