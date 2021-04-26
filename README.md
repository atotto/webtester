# web test library for golang

example:

```example.go
func TestSimple(tt *testing.T) {
	t := webtester.Setup(tt, chrome.DriverPath)
	defer t.TearDown()

	d := t.OpenBrowser()
	d.SetPageLoadTimeout(4 * time.Second)

	d.VisitTo("https://tour.golang.org/welcome/1")

	d.WaitFor("id:run")
	d.MustFindElement("id:run").Click()

	d.WaitFor("class:stdout")
	d.MustFindElement("class:stdout").VerifyText(strings.Contains, "Hello")

	d.MustFindElement("class:next-page").Click()
	d.ExpectTransitTo("/welcome/2").TakeScreenshot("page2.png")
}
```