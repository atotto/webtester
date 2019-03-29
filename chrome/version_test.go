package chrome

import "testing"

func TestParseChromeVersion(t *testing.T) {
	tests := []struct {
		line    string
		version string
	}{
		{`Chromium 73.0.3683.75 Built on Ubuntu , running on Ubuntu 16.04`, "73"},
		{`Chromium 72.0.3626.122 built on Debian 9.8, running on Debian 9.4`, "72"},
	}

	for n, tt := range tests {
		actual := parseChromeVersion([]byte(tt.line))
		if tt.version != actual {
			t.Errorf("#%d want %v, got %v", n, tt.version, actual)
		}
	}
}
