package xpm

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// FuzzParsePanic runs a fuzzer to detect possible panics in the XPM
// parser.  To run, use:
//
//	go test -fuzz=FuzzParsePanic
func FuzzParsePanic(f *testing.F) {
	matches, err := filepath.Glob("testdata/*.xpm")
	if err != nil {
		f.Fatalf("filepath.Glob: %v", err)
	}
	for _, m := range matches {
		buf, err := os.ReadFile(m)
		if err != nil {
			f.Fatalf("os.ReadFile(%q): %v", m, err)
		}
		f.Add(buf)
	}
	f.Fuzz(func(t *testing.T, buf []byte) {
		// It is expected that this may return an error during
		// fuzzing, but we can at least make sure that it does
		// not panic.
		_, _ = Decode(bytes.NewReader(buf))
	})
}
