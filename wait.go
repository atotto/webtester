package webtester

import "time"

func wait(fn func() bool) bool {
	for _, w := range []time.Duration{1, 2, 3, 5, 7, 11, 13, 17} {
		if ok := fn(); ok {
			return true
		}
		time.Sleep(w * time.Second)
	}
	return false
}
