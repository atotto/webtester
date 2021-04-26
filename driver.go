package webtester

import (
	"os"
	"strings"
	"testing"

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

func (d *Driver) OpenBrowser() *Browser {
	d.Helper()

	desired := webdriver.Capabilities{"Platform": "Linux"}
	required := webdriver.Capabilities{}
	var args []string
	if os.Getenv("CI") != "" {
		args = []string{"--headless", "--no-sandbox"}
	}
	if v := os.Getenv("CHROME_OPTIONS_ARGS"); v != "" {
		args = append(args, strings.Split(v, " ")...)
	}
	if args != nil {
		desired["chromeOptions"] = webdriver.Capabilities{"args": args}
	}

	session, err := d.webDriver.NewSession(desired, required)
	if err != nil {
		d.Fatal(err)
	}

	d.sessions = append(d.sessions, session)

	return &Browser{
		TB:      d.TB,
		session: session,
	}
}
