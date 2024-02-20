package chrome

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var DriverPath = filepath.Join(os.TempDir(), "chromedriver")

var chromeBrowser [][]string = [][]string{
	{"chromium-browser", "--version"},
	{"chromium", "--version"},
	{"google-chrome", "--version"},
}

func parseChromeVersion(line []byte) string {
	re := regexp.MustCompile(`.* (\d+)\.(\d+)\.(\d+).*`)
	ss := re.FindSubmatch(line)

	fmt.Printf("%q\n", ss)
	if len(ss) != 4 {
		return ""
	}
	return string(ss[1])
}

func chromeVersion() (version string) {
	var line []byte
	var err error
	for _, chrome := range chromeBrowser {
		cmd := exec.Command(chrome[0], chrome[1:]...)
		line, err = cmd.CombinedOutput()
		if err == nil {
			break
		}
	}
	if err != nil {
		return "120"
	}
	return parseChromeVersion(line)
}

func latestRelease() (version string) {
	var url = "https://googlechromelabs.github.io/chrome-for-testing/LATEST_RELEASE_"
	chromeVersion := chromeVersion()

	res, err := http.Get(url + chromeVersion)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	return string(buf)
}

func targetArch() (target string, err error) {
	archTable := map[string]map[string]string{
		"linux": {
			"amd64": "linux64",
		},
		"darwin": {
			"arm64": "mac-arm64",
			"amd64": "mac-x64",
		},
		"windows": {
			"amd64": "win64",
			"386":   "win32",
		},
	}

	archs, osSupported := archTable[runtime.GOOS]
	if !osSupported {
		return "", fmt.Errorf("not supported: %s", runtime.GOOS)
	}

	target, archSupported := archs[runtime.GOARCH]
	if !archSupported {
		return "", fmt.Errorf("not supported on %s: %s", runtime.GOOS, runtime.GOARCH)
	}

	return target, nil
}

func SetupDriver() error {
	target, err := targetArch()
	if err != nil {
		return err
	}

	version := latestRelease()

	_, err = os.Stat(DriverPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !os.IsNotExist(err) {
		buf, err := exec.Command(DriverPath, "--version").CombinedOutput()
		if err != nil {
			return err
		}
		infos := bytes.Split(buf, []byte(" "))
		if len(infos) != 3 {
			return fmt.Errorf("unexpected version string: %s", string(buf))
		}
		current := string(infos[1])

		if strings.HasPrefix(current, version) {
			return nil
		}
	}

	url := fmt.Sprintf("https://storage.googleapis.com/chrome-for-testing-public/%s/%s/chromedriver-%s.zip", version, target, target)
	log.Printf("download from: %s", url)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	r, err := zip.NewReader(bytes.NewReader(body), res.ContentLength)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		savepath := filepath.Join(os.TempDir(), filepath.Base(f.Name))
		dst, err := os.OpenFile(savepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer dst.Close()
		src, err := f.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		io.Copy(dst, src)

		log.Printf("saved: %s", savepath)
	}
	return nil
}
