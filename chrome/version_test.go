package chrome

import "testing"

func TestParseChromeVersion(t *testing.T) {
	tests := []struct {
		line    string
		version string
	}{
		{`Chromium 83.0.4103.39 Built on Ubuntu , running on Ubuntu 16.04`, "83"},
		{`Chromium 81.0.4044.92 built on Debian 9.8, running on Debian 9.4`, "81"},
		{`Google Chrome 83.0.4103.39`, "83"},
	}

	for n, tt := range tests {
		actual := parseChromeVersion([]byte(tt.line))
		if tt.version != actual {
			t.Errorf("#%d want %v, got %v", n, tt.version, actual)
		}
	}
}
