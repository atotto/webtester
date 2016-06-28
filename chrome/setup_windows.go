package chrome

import (
	"os"
	"path/filepath"
)

func init() {
	DriverPath = filepath.Join(os.TempDir(), "chromedriver.exe")
}
